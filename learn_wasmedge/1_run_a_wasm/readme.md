# 目标
在go程序中，调用wasm来执行函数。

# 步骤
## Rust

```bash
cargo build --target wasm32-wasi --release
wasmedge target/wasm32-wasi/release/helloworld.wasm
````

# 结果
```bash
wzj@ZhaojiedeMacBook-Pro helloworld % wasmedge target/wasm32-wasi/release/helloworld.wasm
Hello, world!
```

