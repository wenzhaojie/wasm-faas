extern "C" {
	// 定义外部函数 fetch，用于从指定 URL 获取数据
	fn fetch(url_pointer: *const u8, url_length: i32) -> i32;
	// 定义外部函数 write_mem，用于将数据写入到内存中
	fn write_mem(pointer: *const u8);
}

#[no_mangle]
// 定义不进行名称修饰的 run 函数，可在其他语言中直接调用
pub unsafe extern fn run() -> i32 {
	// 定义要请求的 URL
	let url = "https://www.baidu.com";
	// 获取 URL 字节切片的指针
	let pointer = url.as_bytes().as_ptr();

	// 调用宿主函数 fetch 获取源代码，返回结果的长度
	let res_len = fetch(pointer, url.len() as i32) as usize;

	// 分配内存空间
	let mut buffer = Vec::with_capacity(res_len);
	let pointer = buffer.as_mut_ptr();

	// 调用宿主函数 write_mem 将源代码写入内存
	write_mem(pointer);

	// 设置缓冲区长度
	buffer.set_len(res_len);
	// 将字节切片转换为字符串，并匹配 "baidu" 出现的次数
	let str = std::str::from_utf8(&buffer).unwrap();
	str.matches("baidu").count() as i32
}
