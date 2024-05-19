#[allow(unused_imports)]
use wasmedge_bindgen::*;
use wasmedge_bindgen_macro::*;
use std::hint::black_box;
use std::fs::File;
use std::io::Write;

#[wasmedge_bindgen] // 使用 wasmedge_bindgen 宏
pub fn bandwidth(data: String) -> String {
    let start_time = std::time::Instant::now(); // 记录开始时间点
    // 使用 black_box 防止编译器优化
    black_box(&data); // 假设这里是你的逻辑处理

    // 写入数据到文件
    if let Ok(mut file) = File::create("output.txt") {
        if let Err(_) = file.write_all(data.as_bytes()) {
            return "Failed to write to file".to_string();
        }
    } else {
        return "Failed to create file".to_string();
    }

    let end_time = std::time::Instant::now(); // 记录结束时间点
    let elapsed_time = end_time.duration_since(start_time); // 计算时间差
    let elapsed_seconds = elapsed_time.as_secs(); // 转换为秒数
    elapsed_seconds.to_string() // 返回时间（秒）
}
