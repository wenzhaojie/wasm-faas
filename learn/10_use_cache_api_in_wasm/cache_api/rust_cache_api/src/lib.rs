use std::ffi::CString;
use std::os::raw::c_char;

// Rust 中的缓存结构体
struct Cache {
    data: std::collections::HashMap<String, i32>,
}

// 全局静态变量，用于存储缓存
static mut CACHE: Option<Cache> = None;

// 从缓存中获取数据的外部函数
#[no_mangle]
pub unsafe extern "C" fn host_get_data_from_cache(key_ptr: *const c_char) -> i32 {
    let key = {
        assert!(!key_ptr.is_null());
        CString::from_raw(key_ptr as *mut _)
    };
    
    let cache = match &mut CACHE {
        Some(cache) => cache,
        None => return 0,
    };
    
    match cache.data.get(&key.to_string_lossy()) {
        Some(value) => *value,
        None => 0,
    }
}

// 将数据存入缓存的外部函数
#[no_mangle]
pub unsafe extern "C" fn host_put_data_to_cache(key_ptr: *const c_char, value: i32) {
    let key = {
        assert!(!key_ptr.is_null());
        CString::from_raw(key_ptr as *mut _)
    };

    let cache = match &mut CACHE {
        Some(cache) => cache,
        None => {
            CACHE = Some(Cache { data: std::collections::HashMap::new() });
            CACHE.as_mut().unwrap()
        }
    };

    cache.data.insert(key.to_string_lossy().into_owned(), value);
}

// 运行函数，调用外部函数
#[no_mangle]
pub unsafe extern "C" fn run() -> i32 {
    let result = host_get_data_from_cache(CString::new("test_key").unwrap().into_raw());
    println!("Data retrieved from cache: {}", result);
    
    host_put_data_to_cache(CString::new("test_key").unwrap().into_raw(), 42);
    println!("Data stored into cache.");
    
    0
}
