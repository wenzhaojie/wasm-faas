# 目标
自己写两个hostfunc，在wasm中，调用go的函数。hostfunc是结构体的方法，结果保存在结构体内。

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
add result from return: 30
add result from struct: 30
```