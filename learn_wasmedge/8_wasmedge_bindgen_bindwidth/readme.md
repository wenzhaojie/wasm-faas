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
wzj@ZhaojiedeMacBook-Pro wasm_bindwidth % ./wasm_bindwidth rust_bindgen.wasm 
wasm内部打点计时: 0
带宽: 83.205145 MB/s
wzj@ZhaojiedeMacBook-Pro wasm_bindwidth % 

```

