# 目标
测试wasm函数的带宽

# 步骤
```bash
cd rust_bindgen
cargo build --target wasm32-wasi --release

cp target/wasm32-wasi/release/rust_bindgen.wasm  ../
cd ..
go build wasm_bindwidth.go
./wasm_bindwidth rust_bindgen.wasm 

```

# 结果
```bash
wzj@ZhaojiedeMacBook-Pro rust_bindgen % cp target/wasm32-wasi/release/rust_bindgen.wasm  ../
cd ..
go build wasm_bindwidth.go
./wasm_bindwidth rust_bindgen.wasm
wasm内部打点计时: [%!f(string=0)] s
带宽: 96.251719 MB/s
```

