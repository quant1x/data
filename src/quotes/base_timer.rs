use std::sync::mpsc::{self, Receiver, Sender};
use std::thread::{self, JoinHandle};
use std::time::Duration;
use std::io::Write;

/// Heartbeat timer control structure
pub struct HeartbeatTimer {
    stop_tx: Sender<()>,
    handle: Option<JoinHandle<()>>,
}

impl HeartbeatTimer {
    /// Create a new heartbeat timer
    pub fn new() -> Self {
        let (stop_tx, _stop_rx) = mpsc::channel();
        Self {
            stop_tx,
            handle: None,
        }
    }

    /// Start the heartbeat timer with the given writer and interval
    pub fn start<W: Write + Send + 'static>(
        &mut self,
        writer: W,
        interval: Duration,
        reset_rx: Receiver<Duration>,
    ) {
        // Create a new channel for this specific timer instance
        let (_stop_tx_inner, stop_rx_inner) = mpsc::channel();
        let handle = thread::spawn(move || {
            pinger(writer, interval, reset_rx, stop_rx_inner);
        });
        self.handle = Some(handle);
    }

    /// Stop the heartbeat timer
    pub fn stop(&mut self) {
        let _ = self.stop_tx.send(());
        if let Some(handle) = self.handle.take() {
            let _ = handle.join();
        }
    }
}

impl Drop for HeartbeatTimer {
    fn drop(&mut self) {
        self.stop();
    }
}

/// Pinger function that sends periodic "ping" messages
fn pinger<W: Write>(
    mut writer: W,
    mut interval: Duration,
    reset_rx: Receiver<Duration>,
    stop_rx: Receiver<()>,
) {
    loop {
        // Wait for either timeout, reset signal, or stop signal
        match reset_rx.recv_timeout(interval) {
            Ok(new_interval) => {
                // Received reset signal
                if new_interval < Duration::from_millis(1) {
                    interval = Duration::from_secs(10); // default_ping_interval
                } else {
                    interval = new_interval;
                }
                continue;
            }
            Err(mpsc::RecvTimeoutError::Timeout) => {
                // Timeout - send ping
                if let Err(_) = writer.write_all(b"ping") {
                    // Write failed, exit
                    return;
                }
                if let Err(_) = writer.flush() {
                    // Flush failed, exit
                    return;
                }
            }
            Err(mpsc::RecvTimeoutError::Disconnected) => {
                // Channel disconnected, exit
                return;
            }
        }

        // Check for stop signal (non-blocking)
        match stop_rx.try_recv() {
            Ok(_) => return, // Stop signal received
            Err(mpsc::TryRecvError::Empty) => {} // No stop signal, continue
            Err(mpsc::TryRecvError::Disconnected) => return, // Stop channel disconnected
        }
    }
}

/// Example pinger function demonstrating usage
pub fn example_pinger() {
    use std::io::Cursor;

    let (reset_tx, reset_rx) = mpsc::channel();
    let (stop_tx, stop_rx) = mpsc::channel();

    let handle = thread::spawn(move || {
        let mut buffer = Vec::new();
        let mut writer = Cursor::new(&mut buffer);
        pinger(&mut writer, Duration::from_secs(1), reset_rx, stop_rx);
        buffer
    });

    // Simulate the example sequence
    let test_sequence = [
        Duration::from_millis(0),    // reset to 0
        Duration::from_millis(200),  // reset to 200ms
        Duration::from_millis(300),  // reset to 300ms
        Duration::from_millis(0),    // reset to 0
        Duration::from_millis(0),    // reset to 0 (negative case)
        Duration::from_millis(0),    // reset to 0 (negative case)
        Duration::from_millis(0),    // reset to 0 (negative case)
    ];

    for (i, &delay) in test_sequence.iter().enumerate() {
        println!("Run {}", i + 1);

        if delay.as_millis() > 0 {
            println!("resetting time ({})", delay.as_millis());
            let _ = reset_tx.send(delay);
        }

        // Wait a bit to allow ping
        thread::sleep(Duration::from_millis(50));
    }

    // Stop the pinger
    let _ = stop_tx.send(());
    let buffer = handle.join().unwrap();
    println!("Final buffer contents: {:?}", String::from_utf8_lossy(&buffer));
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::io::Cursor;

    #[test]
    fn test_heartbeat_timer_creation() {
        let timer = HeartbeatTimer::new();
        // Timer should be created successfully
        assert!(timer.handle.is_none());
    }

    #[test]
    fn test_pinger_basic() {
        let (_reset_tx, reset_rx) = mpsc::channel();
        let (stop_tx, stop_rx) = mpsc::channel();

        let handle = thread::spawn(move || {
            let mut buffer = Vec::new();
            let mut writer = Cursor::new(&mut buffer);
            pinger(&mut writer, Duration::from_millis(50), reset_rx, stop_rx);
            buffer
        });

        // Wait for a ping to be sent
        thread::sleep(Duration::from_millis(100));

        // Stop the pinger
        let _ = stop_tx.send(());
        let buffer = handle.join().unwrap();

        // Check that some data was written (ping messages)
        assert!(!buffer.is_empty());
    }
}