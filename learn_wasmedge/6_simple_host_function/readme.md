# 目标
自己写两个hostfunc，在wasm中，调用go的函数。

# 步骤
```bash
cd rust
cargo build --target wasm32-wasi --release

cp target/wasm32-wasi/release/rust.wasm ../
cd ..
go build hostfunc.go
./hostfunc rust.wasm
```

# 结果
```bash
wzj@ZhaojiedeMacBook-Pro hostfunc % ./hostfunc rust.wasm
Go: Args: [./hostfunc rust.wasm]
host_add res: 30
host_get_datetime year: 2024
add result: 30
```