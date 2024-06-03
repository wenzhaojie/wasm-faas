# 目标
测试序列化和反序列化

# 步骤
cargo run

# 结果
```bash
(wasmtime) wzj@ZhaojiedeMacBook-Pro serialize % cargo run
    Finished `dev` profile [unoptimized + debuginfo] target(s) in 0.01s
     Running `target/debug/serialize`
Serialized: {"name":"Alice","age":30}
Deserialized: MyData { name: "Alice", age: 30 }
```