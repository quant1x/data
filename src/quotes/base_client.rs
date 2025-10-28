use std::io::{self, Read, Write};
use std::net::ToSocketAddrs;
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::{Arc, Mutex, MutexGuard};
use std::thread::{self, JoinHandle};
use std::time::{Duration, Instant};

use anyhow::{anyhow, Context, Result};
use log::warn;
use mio::event::Event;
use mio::net::TcpStream;
use mio::{Events, Interest, Poll, Token};

use crate::quotes::cmd_security_count::{new_security_count, SecurityCountRequest};
use crate::quotes::message::{process, Message};
use crate::quotes::options::Options;
use crate::quotes::server::Server;

const ERR_TIMEOUT: &str = "connect timeout";
const CLIENT_TOKEN: Token = Token(0);

struct TcpClientInner {
    stream: Option<TcpStream>,
    poll: Poll,
    events: Events,
    server: Option<Server>,
    opt: Options,
    completed_at: Instant,
    closed: bool,
}

impl TcpClientInner {
    fn new(mut opt: Options) -> Result<Self> {
        if opt.max_retry_times == 0 {
            opt.max_retry_times = 3;
        }
        if opt.connection_timeout == Duration::ZERO {
            opt.connection_timeout = Duration::from_secs(10);
        }
        if opt.read_timeout == Duration::ZERO {
            opt.read_timeout = Duration::from_secs(5);
        }
        if opt.write_timeout == Duration::ZERO {
            opt.write_timeout = Duration::from_secs(5);
        }

        let poll = Poll::new().context("failed to create mio poll")?;
        let events = Events::with_capacity(16);

        Ok(Self {
            stream: None,
            poll,
            events,
            server: None,
            opt,
            completed_at: Instant::now(),
            closed: true,
        })
    }

    fn update_completed(&mut self) {
        self.completed_at = Instant::now();
    }

    fn wait_for<F>(&mut self, timeout: Duration, mut predicate: F) -> Result<()>
    where
        F: FnMut(&Event) -> bool,
    {
        let start = Instant::now();
        loop {
            let elapsed = start.elapsed();
            if elapsed >= timeout {
                return Err(anyhow!(ERR_TIMEOUT));
            }
            let remaining = timeout - elapsed;
            self.events.clear();
            self.poll
                .poll(&mut self.events, Some(remaining))
                .context("poll wait")?;

            for event in self.events.iter() {
                if event.token() != CLIENT_TOKEN {
                    continue;
                }
                if event.is_error() || event.is_read_closed() || event.is_write_closed() {
                    return Err(anyhow!("socket closed"));
                }
                if predicate(event) {
                    return Ok(());
                }
            }
        }
    }

    fn wait_for_writable(&mut self, timeout: Duration) -> Result<()> {
        self.wait_for(timeout, |event| event.is_writable())
    }

    fn wait_for_readable(&mut self, timeout: Duration) -> Result<()> {
        self.wait_for(timeout, |event| event.is_readable())
    }

    fn write_all(&mut self, mut data: &[u8]) -> Result<()> {
        while !data.is_empty() {
            let write_result = {
                let stream = self
                    .stream
                    .as_mut()
                    .ok_or_else(|| anyhow!("tcp client not connected"))?;
                stream.write(data)
            };

            match write_result {
                Ok(0) => return Err(anyhow!("socket closed while writing")),
                Ok(n) => {
                    data = &data[n..];
                    self.update_completed();
                }
                Err(err) => {
                    if err.kind() == io::ErrorKind::WouldBlock {
                        self.wait_for_writable(self.opt.write_timeout)?;
                    } else {
                        return Err(err.into());
                    }
                }
            }
        }
        Ok(())
    }

    fn read_exact(&mut self, buf: &mut [u8]) -> Result<()> {
        let mut offset = 0;

        while offset < buf.len() {
            let read_result = {
                let stream = self
                    .stream
                    .as_mut()
                    .ok_or_else(|| anyhow!("tcp client not connected"))?;
                stream.read(&mut buf[offset..])
            };

            match read_result {
                Ok(0) => return Err(anyhow!("socket closed while reading")),
                Ok(n) => {
                    offset += n;
                    self.update_completed();
                }
                Err(err) => {
                    if err.kind() == io::ErrorKind::WouldBlock {
                        self.wait_for_readable(self.opt.read_timeout)?;
                    } else {
                        return Err(err.into());
                    }
                }
            }
        }
        Ok(())
    }

    fn has_timed_out(&self) -> bool {
        self.completed_at.elapsed() >= self.opt.connection_timeout
    }

    fn cross_time(&self) -> f64 {
        self.completed_at.elapsed().as_secs_f64()
    }
}

struct HeartbeatControl {
    stop: AtomicBool,
    handle: Mutex<Option<JoinHandle<()>>>,
}

impl HeartbeatControl {
    fn new() -> Self {
        Self {
            stop: AtomicBool::new(true),
            handle: Mutex::new(None),
        }
    }
}

/// 对应 Go 版的 `TcpClient`，使用 mio 提供的非阻塞网络能力。
#[derive(Clone)]
pub struct TcpClient {
    inner: Arc<Mutex<TcpClientInner>>,
    heartbeat: Arc<HeartbeatControl>,
}

impl TcpClient {
    pub fn new(opt: Options) -> Result<Self> {
        let inner = TcpClientInner::new(opt)?;
        Ok(Self {
            inner: Arc::new(Mutex::new(inner)),
            heartbeat: Arc::new(HeartbeatControl::new()),
        })
    }

    fn lock_inner(&self) -> MutexGuard<'_, TcpClientInner> {
        self.inner.lock().expect("tcp client lock poisoned")
    }

    pub fn connect(&self, server: &Server) -> Result<()> {
        if server.host.is_empty() || server.port == 0 {
            return Err(anyhow!("invalid server address"));
        }

        // If an existing connection is active, shut it down first.
        let _ = self.close();

        let addr = (server.host.as_str(), server.port)
            .to_socket_addrs()
            .context("resolve server address")?
            .next()
            .ok_or_else(|| anyhow!("no socket address resolved"))?;

        let mut stream = TcpStream::connect(addr).context("tcp connect failed")?;

        {
            let mut inner = self.lock_inner();
            inner
                .poll
                .registry()
                .register(
                    &mut stream,
                    CLIENT_TOKEN,
                    Interest::WRITABLE | Interest::READABLE,
                )
                .context("register tcp stream")?;

            let connect_timeout = inner.opt.connection_timeout;
            inner.wait_for_writable(connect_timeout)?;

            if let Some(err) = stream.take_error().context("take socket error")? {
                let _ = inner.poll.registry().deregister(&mut stream);
                return Err(anyhow!("connect error after poll: {err}"));
            }
            stream.set_nodelay(true).ok();

            inner.stream = Some(stream);
            inner.server = Some(server.clone());
            inner.closed = false;
            inner.update_completed();
        }

        self.start_heartbeat();
        Ok(())
    }

    pub fn close(&self) -> Result<()> {
        self.close_internal(true)
    }

    fn close_internal(&self, join_heartbeat: bool) -> Result<()> {
        self.heartbeat.stop.store(true, Ordering::SeqCst);

        let (release, server) = {
            let mut inner = self.lock_inner();
            if inner.closed {
                return Ok(());
            }

            let release = inner.opt.release_address.clone();
            let server = inner.server.take();

            if let Some(mut stream) = inner.stream.take() {
                if let Err(err) = inner.poll.registry().deregister(&mut stream) {
                    warn!("failed to deregister tcp stream: {err}");
                }
            }

            inner.closed = true;
            (release, server)
        };

        if let (Some(callback), Some(server)) = (release, server) {
            callback(&server);
        }

        if join_heartbeat {
            if let Some(handle) = self.heartbeat.handle.lock().unwrap().take() {
                let _ = handle.join();
            }
        } else {
            let mut guard = self.heartbeat.handle.lock().unwrap();
            if let Some(handle) = guard.take() {
                drop(handle);
            }
        }

        Ok(())
    }

    pub fn write_all(&self, data: &[u8]) -> Result<()> {
        let mut inner = self.lock_inner();
        inner.write_all(data)
    }

    pub fn read_exact(&self, buf: &mut [u8]) -> Result<()> {
        let mut inner = self.lock_inner();
        inner.read_exact(buf)
    }

    pub fn command(&self, msg: &mut dyn Message) -> Result<()> {
        process(self, msg)
    }

    pub fn has_timed_out(&self) -> bool {
        let inner = self.lock_inner();
        inner.has_timed_out()
    }

    pub fn cross_time(&self) -> f64 {
        let inner = self.lock_inner();
        inner.cross_time()
    }

    fn start_heartbeat(&self) {
        self.heartbeat.stop.store(false, Ordering::SeqCst);
        let mut handle_guard = self.heartbeat.handle.lock().unwrap();
        if handle_guard.is_some() {
            return;
        }

        let client = self.clone();
        *handle_guard = Some(thread::spawn(move || client.run_heartbeat_loop()));
    }

    fn run_heartbeat_loop(self) {
        while !self.heartbeat.stop.load(Ordering::SeqCst) {
            thread::sleep(Duration::from_secs(1));
            if self.heartbeat.stop.load(Ordering::SeqCst) {
                break;
            }

            let should_ping = {
                let inner = self.lock_inner();
                if inner.closed {
                    return;
                }
                inner.has_timed_out()
            };

            if !should_ping {
                continue;
            }

            let server_label = {
                let inner = self.lock_inner();
                inner.server.as_ref().map(|srv| srv.addr())
            };

            match new_security_count() {
                Ok(mut pkg) => {
                    pkg.set_params(SecurityCountRequest { market: 1 });
                    if let Err(err) = self.command(&mut pkg) {
                        if let Some(label) = server_label.as_deref() {
                            warn!("client -> server[{label}]: heartbeat failed: {err}");
                        } else {
                            warn!("client heartbeat failed: {err}");
                        }
                        let _ = self.close_internal(false);
                        break;
                    } else {
                        if let Some(label) = server_label.as_deref() {
                            warn!("client -> server[{label}]: heartbeat");
                        } else {
                            warn!("client heartbeat");
                        }
                    }
                }
                Err(err) => {
                    warn!("failed to construct security count heartbeat: {err}");
                }
            }
        }

        let mut guard = self.heartbeat.handle.lock().unwrap();
        if let Some(handle) = guard.take() {
            if thread::current().id() != handle.thread().id() {
                let _ = handle.join();
            } else {
                drop(handle);
            }
        }
    }
}
