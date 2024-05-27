use serde::{Serialize, Deserialize};
use serde_json;

#[derive(Serialize, Deserialize, Debug)]
struct MyData {
    name: String,
    age: u32,
}

fn main() {
    // 创建一个对象
    let my_object = MyData {
        name: String::from("Alice"),
        age: 30,
    };

    // 序列化对象为JSON字符串
    let serialized = serde_json::to_string(&my_object).unwrap();
    println!("Serialized: {}", serialized);

    // 反序列化JSON字符串为对象
    let deserialized: MyData = serde_json::from_str(&serialized).unwrap();
    println!("Deserialized: {:?}", deserialized);
}
