use std::sync::{Arc, Mutex, OnceLock};

use crate::quotes::{self, StdApi};

// 全局共享的 StdApi 实例，使用互斥锁保护，确保多线程安全复用。
static STD_API: OnceLock<Mutex<Option<Arc<StdApi>>>> = OnceLock::new();

fn init_std_api(slot: &mut Option<Arc<StdApi>>) {
    if slot.is_none() {
        if let Ok(api) = quotes::new_std_api() {
            *slot = Some(Arc::new(api));
        }
    }
}

pub fn get_api() -> Option<Arc<StdApi>> {
    let mutex = STD_API.get_or_init(|| Mutex::new(None));
    let mut guard = mutex.lock().expect("std api mutex poisoned");
    init_std_api(&mut guard);
    guard.clone()
}

pub fn reopen() {
    if let Some(mutex) = STD_API.get() {
        let mut guard = mutex.lock().expect("std api mutex poisoned");
        if let Some(api) = guard.take() {
            api.close();
        }
    }
}
