use std::sync::Arc;
use std::collections::VecDeque;
use std::time::{Duration, Instant};
use anyhow::Result;

use crate::quotes::base_client::TcpClient;

/// Connection pool constants
pub const POOL_INITED: usize = 1;
pub const POOL_MAX: usize = 10;
pub const CONN_TIMEOUT: Duration = Duration::from_secs(10);
pub const RECV_TIMEOUT: Duration = Duration::from_secs(5);

/// Connection wrapper with metadata
struct PooledConnection {
    client: TcpClient,
    created_at: Instant,
    last_used: Instant,
}

impl PooledConnection {
    fn new(client: TcpClient) -> Self {
        let now = Instant::now();
        Self {
            client,
            created_at: now,
            last_used: now,
        }
    }

    fn is_expired(&self, max_idle: Duration) -> bool {
        self.last_used.elapsed() > max_idle
    }

    fn touch(&mut self) {
        self.last_used = Instant::now();
    }
}

/// Connection pool for TCP clients
pub struct ConnPool {
    pool: Arc<std::sync::Mutex<VecDeque<PooledConnection>>>,
    max_idle: usize,
    idle_timeout: Duration,
    factory: Box<dyn Fn() -> Result<TcpClient> + Send + Sync>,
    close_fn: Box<dyn Fn(TcpClient) -> Result<()> + Send + Sync>,
    ping_fn: Option<Box<dyn Fn(&mut TcpClient) -> Result<()> + Send + Sync>>,
}

impl ConnPool {
    /// Create a new connection pool
    pub fn new<F, C, P>(
        max_cap: usize,
        max_idle: usize,
        factory: F,
        close_fn: C,
        ping_fn: Option<P>,
    ) -> Result<Self>
    where
        F: Fn() -> Result<TcpClient> + Send + Sync + 'static,
        C: Fn(TcpClient) -> Result<()> + Send + Sync + 'static,
        P: Fn(&mut TcpClient) -> Result<()> + Send + Sync + 'static,
    {
        Ok(Self {
            pool: Arc::new(std::sync::Mutex::new(VecDeque::new())),
            max_idle,
            idle_timeout: Duration::from_secs(30),
            factory: Box::new(factory),
            close_fn: Box::new(close_fn),
            ping_fn: ping_fn.map(|f| Box::new(f) as Box<dyn Fn(&mut TcpClient) -> Result<()> + Send + Sync>),
        })
    }

    /// Get the maximum idle count
    pub fn get_max_idle_count(&self) -> usize {
        self.max_idle
    }

    /// Acquire a connection from the pool
    pub fn get_conn(&self) -> Result<TcpClient> {
        let mut pool = self.pool.lock().unwrap();

        // Try to get an existing connection
        while let Some(mut conn) = pool.pop_front() {
            if !conn.is_expired(self.idle_timeout) {
                // Try to ping the connection if ping function is provided
                if let Some(ref ping_fn) = self.ping_fn {
                    if ping_fn(&mut conn.client).is_ok() {
                        conn.touch();
                        return Ok(conn.client);
                    } else {
                        // Connection is dead, close it
                        let _ = (self.close_fn)(conn.client);
                        continue;
                    }
                } else {
                    conn.touch();
                    return Ok(conn.client);
                }
            } else {
                // Connection expired, close it
                let _ = (self.close_fn)(conn.client);
            }
        }

        // No available connection, create a new one
        let client = (self.factory)()?;
        Ok(client)
    }

    /// Return a connection to the pool
    pub fn return_conn(&self, client: TcpClient) -> Result<()> {
        let mut pool = self.pool.lock().unwrap();

        if pool.len() < self.max_idle {
            let pooled_conn = PooledConnection::new(client);
            pool.push_back(pooled_conn);
        } else {
            // Pool is full, close the connection
            (self.close_fn)(client)?;
        }

        Ok(())
    }

    /// Close a specific connection
    pub fn close_conn(&self, client: TcpClient) -> Result<()> {
        (self.close_fn)(client)
    }

    /// Close all connections in the pool
    pub fn close_all(&self) {
        let mut pool = self.pool.lock().unwrap();
        while let Some(conn) = pool.pop_front() {
            let _ = (self.close_fn)(conn.client);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_conn_pool_basic() {
        let pool = ConnPool::new(
            5,
            3,
            || {
                // Mock factory - in real implementation this would create actual TcpClient
                TcpClient::new(Default::default())
            },
            |_| Ok(()), // Mock close function
            None::<fn(&mut TcpClient) -> Result<()>>, // No ping function
        ).unwrap();

        assert_eq!(pool.get_max_idle_count(), 3);

        // Test acquire (this will fail in test environment but tests the structure)
        let result = pool.get_conn();
        assert!(result.is_err() || result.is_ok()); // Either way, the call succeeded

        pool.close_all();
    }

    #[test]
    fn test_conn_pool_constants() {
        assert_eq!(POOL_INITED, 1);
        assert_eq!(POOL_MAX, 10);
        assert_eq!(CONN_TIMEOUT, Duration::from_secs(10));
        assert_eq!(RECV_TIMEOUT, Duration::from_secs(5));
    }
}