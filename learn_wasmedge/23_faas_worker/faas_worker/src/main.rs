// 这是一个handler，用于使用wasm环境沙箱来运行FaaS函数
use serde_json;
use anyhow::Result;
use redis::Commands;
use std::time::Instant;
use std::collections::HashMap;


fn get_url() -> String {
    if let Ok(url) = std::env::var("REDIS_URL") {
        url
    } else {
        "redis://127.0.0.1/".into()
    }
}


pub fn test_serde() -> Result<()> {
    // 连接 to Redis
    // 计时

    let client = redis::Client::open(&*get_url())?;
    let mut con = client.get_connection()?;

    // 测试 字符串
    let input_str = "wenzhaojie".to_string();
    // 序列化字符串 input_str 为输入对象 input_obj_serde_str
    let input_obj_serde_str = serde_json::to_string(&input_str).unwrap();

    // 存入redis，key为phd
    let input_obj_key = "phd".to_string();
    let output_obj_key = "phd_graduate".to_string();
    let _ : () = con.set(&input_obj_key, &input_obj_serde_str).unwrap();

    // 调用 handler 函数
    let stat = handler(input_obj_key, output_obj_key);

    // 打印输出结果
    println!("Output JSON string: {}", stat);
    Ok(())
}


pub fn handler(input_obj_key: String, output_obj_key: String) -> String {
    // input_obj_key 是redis的key
    // 从redis中获取输入字符串 input_obj_serde_str
    // 计时
    let conn_redis_t = Instant::now();
    let client = redis::Client::open(&*get_url()).unwrap();
    let mut con = client.get_connection().unwrap();
    let conn_redis_t = conn_redis_t.elapsed();
    // 计时
    let start_get_input_t = Instant::now();
    let input_obj_serde_str: String = con.get(input_obj_key).unwrap();
    let get_input_t = start_get_input_t.elapsed();
    // 反序列化输入字符串 input_obj_serde_str 为输入对象 input_obj
    let input_obj = serde_json::from_str(&input_obj_serde_str).unwrap();
    // 调用函数（假设这里调用了名为 `invoke_function` 的函数）
    // 计时
    let start_compute_t = Instant::now();
    let output_obj = invoke_function(input_obj);
    let compute_t = start_compute_t.elapsed();
    // 重新序列化得到输出字符串 output_obj_serde_str
    // 计时
    let start_set_output_t = Instant::now();
    let output_obj_serde_str = serde_json::to_string(&output_obj).unwrap();
    let set_output_t = start_set_output_t.elapsed();
    // 需要将输出字符串 output_obj_serde_str 存入redis，key为output_obj_key
    let _ : () = con.set(&output_obj_key, &output_obj_serde_str).unwrap();

    // 返回一些统计信息，一个hashmap
    let mut stats_dict = HashMap::new();
    stats_dict.insert("conn_redis_t", conn_redis_t.as_millis());
    stats_dict.insert("get_input_t", get_input_t.as_millis());
    stats_dict.insert("compute_t", compute_t.as_millis());
    stats_dict.insert("set_output_t", set_output_t.as_millis());
    // 返回统计信息, 以json字符串的形式
    serde_json::to_string(&stats_dict).unwrap()
}


// invoke_function 是一个示例函数，打印hello world并返回输入对象加上一个字符串
fn invoke_function(input_str: String) -> String {
    println!("Hello, world!");
    let result = format!("{} - from invoke_function", input_str);
    result
}


fn main() {
    test_serde().unwrap();
}