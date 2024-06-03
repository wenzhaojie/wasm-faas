# 目标
faas_worker 用于在wasm中实现函数，以便在wasm中调用。

# 步骤
rust_build.sh
run.sh

# 结果
```bash
(base) wzj@ZhaojiedeMacBook-Pro faas_worker % sh run.sh 
Hello, world!
Output JSON string: {"compute_t":0,"conn_redis_t":0,"get_input_t":0,"set_output_t":0}
```