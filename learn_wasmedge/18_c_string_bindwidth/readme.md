# 目标
尝试使用c api来调用wasm文件

# 步骤
```bash
cd c_string
cd rust_bindgen
cargo build --target wasm32-wasi --release

cp target/wasm32-wasi/release/rust_bindgen.wasm  ../
cd ..

gcc bindwidth.c -lwasmedge
./a.out
./a.out rust_bindgen.wasm 
```


# 结果
```bash
(wasmtime) wzj@ZhaojiedeMacBook-Pro c_string % ./a.out rust_bindgen.wasm 
[2024-05-22 16:42:22.876] [error] instantiation failed: unknown import, Code: 0x62
[2024-05-22 16:42:22.876] [error]     When linking module: "wasi_snapshot_preview1" , function name: "clock_time_get"
[2024-05-22 16:42:22.876] [error]     At AST node: import description
[2024-05-22 16:42:22.876] [error]     This is a WASI related import. Please ensure that you've turned on the WASI configuration.
[2024-05-22 16:42:22.876] [error]     At AST node: import section
[2024-05-22 16:42:22.876] [error]     At AST node: module
Error message: unknown import
```
