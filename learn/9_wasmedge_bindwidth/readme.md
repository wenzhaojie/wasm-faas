# 目标
测试wasm函数的带宽，不使用bindgen的方法，使用传递指针加长度的方式；

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

```

