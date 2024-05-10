# 目标
用rust写一个lib，编译成wasm模块，然后在go的代码中加载wasm模块，然后执行

# 步骤
```bash
rustup target add wasm32-wasi
cd fibonacci_lib
cargo build --target wasm32-wasi --release

cp target/wasm32-wasi/release/fibonacci_lib.wasm ../
cd ..
go build run_wasm_lib_in_go.go
./run_wasm_lib_in_go fibonacci_lib.wasm
```

# 结果
```bash
wzj@ZhaojiedeMacBook-Pro run_wasm_lib_in_go % ./run_wasm_lib_in_go fibonacci_lib.wasm
Go: Args: [./run_wasm_lib_in_go fibonacci_lib.wasm]
获取第 21 项斐波那契数: 10946
```