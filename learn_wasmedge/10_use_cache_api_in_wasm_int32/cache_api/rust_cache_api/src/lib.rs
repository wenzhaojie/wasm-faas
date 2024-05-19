extern "C" {
    // 定义外部函数 host_get_data_from_cache
    fn host_get_data_from_cache(key: i32) -> i32;
    // 定义外部函数 host_put_data_to_cache
    fn host_put_data_to_cache(key: i32, value: i32);
}

#[no_mangle]
// 定义不进行名称修饰的 run 函数，可在其他语言中直接调用
pub unsafe extern fn run() -> i32 {
    // 准备存放一个key，key=123,value=456
    let key = 123;
    let value = 456;

    // 调用外部函数将数据放入缓存
    host_put_data_to_cache(key, value);

    // 读取这个key，返回这个value
    let retrieved_value = host_get_data_from_cache(key);

    retrieved_value
}

