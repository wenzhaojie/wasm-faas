# 目标
hostfunc操作宿主go cache，但是接口都是i32的数据类型。

# 步骤
```bash
cd rust_cache_api
cargo build --target wasm32-wasi --release

cp target/wasm32-wasi/release/rust_cache_api.wasm  ../
cd ..
go build hostfunc.go 
./hostfunc rust_cache_api.wasm

```

# 结果
```bash
Stored value 456 for key %!s(int32=123) into cache
Retrieved value 456 for key %!s(int32=123) from cache
Run function returned: 456
Go cache contents:
Key: 123, Value: 456

```

