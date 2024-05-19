extern "C" {
    // 定义外部函数 host_add，用于将两个整数相加，并返回结果
	fn host_add(a: i32, b: i32) -> i32;
    // 定义外部函数 host_get_datetime，用于获取当前日期和时间
	fn host_get_datetime() -> i32;
}

#[no_mangle]
// 定义不进行名称修饰的 run 函数，可在其他语言中直接调用
pub unsafe extern fn run() -> i32 {

    // 测试 add 函数 两个数相加，并打印结果
    let result = host_add(10, 20);
    println!("Addition result: {}", result);

    // 测试 get_datetime 函数，并打印结果
    let datetime = host_get_datetime();
    println!("Current datetime: {}", datetime);

    // result
    result as i32
}

