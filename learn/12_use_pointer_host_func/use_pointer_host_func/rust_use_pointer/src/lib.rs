extern "C" {
    // 定义外部函数 host_get_data_from_cache
    fn host_get_data_from_cache(key_pointer: *const u8, key_length: i32) -> i32;
    // 定义外部函数 host_put_data_to_cache
    fn host_set_data_to_cache(key_pointer: *const u8, key_length: i32, value_pointer: *const u8, value_length: i32);
    // 定义外部函数 write_mem，用于将数据写入到内存中
    fn write_mem(pointer: *const u8);
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

    // 测试读取key
    // 准备读取一个key，key="hust", value="123123"
    let new_key = "hust";
    // 将 key 转换为 C 风格的字符串
    let new_key_bytes = new_key.as_bytes();
    // 读取这个 key，返回这个 value的空间长度
    let value_length = host_get_data_from_cache(new_key_bytes.as_ptr(), new_key_bytes.len() as i32) as usize;

    // 分配内存空间准备写入value
    // 分配内存空间
    let mut buffer = Vec::with_capacity(value_length);
    let pointer = buffer.as_mut_ptr();

    // 调用宿主函数 write_mem 将源代码写入内存
    write_mem(pointer);
    // 设置缓冲区长度
    buffer.set_len(value_length);
    // 将字节切片转换为字符串
    let num_str = std::str::from_utf8(&buffer);

    // 为了证明 num_str == "123123"
     assert_eq!(num_str, Ok("123123"));

    // 计算函数运行时间
    let elapsed_time_ms = start_time.elapsed().as_millis() as i32;

    //
    elapsed_time_ms
}
