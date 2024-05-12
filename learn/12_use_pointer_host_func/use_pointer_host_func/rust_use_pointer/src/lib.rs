extern "C" {
    // 定义外部函数 host_get_data_from_cache
    fn host_get_data_from_cache(key_pointer: *const u8, key_length: i32) -> *const u8, i32;
    // 定义外部函数 host_put_data_to_cache
    fn host_set_data_to_cache(key_pointer: *const u8, key_length: i32, value_pointer: *const u8, value_length: i32);
}

#[no_mangle]
// 定义不进行名称修饰的 run 函数，可在其他语言中直接调用
pub unsafe extern "C" fn run() -> i32 {
    // 准备存放一个key，key="abc", value="wenzhaojie"
    let key = "abc";
    let value = "wenzhaojie";

    // 将 key 和 value 转换为 C 风格的字符串
    let key_bytes = key.as_bytes();
    let value_bytes = value.as_bytes();

    // 记录函数开始时间
    let start_time = std::time::Instant::now();

    // 调用外部函数将数据放入缓存
    host_set_data_to_cache(key_bytes.as_ptr(), key_bytes.len() as i32, value_bytes.as_ptr(), value_bytes.len() as i32);

    // 读取这个 key，返回这个 value
    let _retrieved_value = host_get_data_from_cache(key_bytes.as_ptr(), key_bytes.len() as i32);

    // 计算函数运行时间
    let elapsed_time_ms = start_time.elapsed().as_millis() as i32;

    // 返回函数运行时间
    elapsed_time_ms
}
