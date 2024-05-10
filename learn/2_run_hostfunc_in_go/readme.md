# 目标
在Go语言中，调用wasm函数，完成word count。
# 步骤
```bash
rustup target add wasm32-wasi
cd rust_host_func
cargo build --target wasm32-wasi --release

cp target/wasm32-wasi/release/rust_host_func.wasm ../
cd ..
go build hostfunc.go
./hostfunc rust_host_func.wasm
```
# 结果
```shell
Go: Args: [./hostfunc rust_host_func.wasm]
There are 78 'baidu' in source code of baidu.com
```