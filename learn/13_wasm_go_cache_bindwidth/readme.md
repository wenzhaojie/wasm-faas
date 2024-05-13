# 目标
在Go语言中，实现了在wasm函数，调用访问宿主go程序中的go-cache。测量带宽。
# 步骤
```bash
rustup target add wasm32-wasi
cd use_pointer_host_func
cd rust_use_pointer
cargo build --target wasm32-wasi --release

cp target/wasm32-wasi/release/rust_use_pointer.wasm ../
cd ..
go build use_pointer_host_func.go
./use_pointer_host_func rust_use_pointer.wasm
```
# 结果
```shell
wzj@ZhaojiedeMacBook-Pro use_pointer_host_func % ./use_pointer_host_func rust_use_pointer.wasm 
Go: Args: [./use_pointer_host_func rust_use_pointer.wasm]
Set value wenzhaojie for key abc in cache
Type of value: string
Retrieved value 123123 for key hust from cache
writeMem: tmpDataStr.DataStr is 123123
The result from wasm module is 0 
wzj@ZhaojiedeMacBook-Pro use_pointer_host_func % 
```