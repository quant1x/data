use std::any::Any;
use std::collections::VecDeque;
use std::sync::{Arc, Mutex};

use anyhow::{anyhow, Result};

use crate::quotes::base_consts::QuotesError;

/// 简易连接池，占位实现，后续将依据 Go 版本补全功能。
pub struct ConnPool {
    max_idle: usize,
    state: Arc<Mutex<PoolState>>,
    factory: Arc<dyn Fn() -> Result<Box<dyn Any + Send>> + Send + Sync>,
    close: Arc<dyn Fn(&mut (dyn Any + Send)) -> Result<()> + Send + Sync>,
    #[allow(dead_code)]
    ping: Arc<dyn Fn(&mut (dyn Any + Send)) -> Result<()> + Send + Sync>,
}

struct PoolState {
    idle: VecDeque<Box<dyn Any + Send>>,
    closed: bool,
}

impl ConnPool {
    pub fn new<F, C, P>(
        max_cap: usize,
        max_idle: usize,
        factory: F,
        close: C,
        ping: P,
    ) -> Result<Self>
    where
        F: Fn() -> Result<Box<dyn Any + Send>> + Send + Sync + 'static,
        C: Fn(&mut (dyn Any + Send)) -> Result<()> + Send + Sync + 'static,
        P: Fn(&mut (dyn Any + Send)) -> Result<()> + Send + Sync + 'static,
    {
        if max_cap == 0 {
            return Err(anyhow!(QuotesError::Unimplemented(
                "connection pool requires capacity"
            )));
        }
        Ok(Self {
            max_idle,
            state: Arc::new(Mutex::new(PoolState {
                idle: VecDeque::new(),
                closed: false,
            })),
            factory: Arc::new(factory),
            close: Arc::new(close),
            ping: Arc::new(ping),
        })
    }

    pub fn get_max_idle_count(&self) -> usize {
        self.max_idle
    }

    pub fn get_conn(&self) -> Result<Box<dyn Any + Send>> {
        let mut guard = self.state.lock().expect("connection pool poisoned");
        if guard.closed {
            return Err(anyhow!("connection pool closed"));
        }
        if let Some(conn) = guard.idle.pop_front() {
            return Ok(conn);
        }
        drop(guard);
        (self.factory)()
    }

    pub fn return_conn(&self, conn: Box<dyn Any + Send>) -> Result<()> {
        let mut guard = self.state.lock().expect("connection pool poisoned");
        if guard.closed {
            let mut conn = conn;
            (self.close)(conn.as_mut())?;
            return Ok(());
        }
        if guard.idle.len() < self.max_idle {
            guard.idle.push_back(conn);
        } else {
            let mut value = conn;
            (self.close)(value.as_mut())?;
        }
        Ok(())
    }

    pub fn close_conn(&self, mut conn: Box<dyn Any + Send>) -> Result<()> {
        (self.close)(conn.as_mut())
    }

    pub fn close_all(&self) {
        let mut guard = self.state.lock().expect("connection pool poisoned");
        guard.closed = true;
        for mut conn in guard.idle.drain(..) {
            if let Err(err) = (self.close)(conn.as_mut()) {
                log::warn!("failed to close pooled connection: {err}");
            }
        }
    }
}
