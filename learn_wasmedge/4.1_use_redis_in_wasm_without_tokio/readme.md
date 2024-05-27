# 目标
在wasm程序里，访问redis，并测量带宽。
# 步骤
```bash
cd redis_bindwidth
sh rust_build.sh
sh run.sh
```
# 结果
```bash
wzj@ZhaojiedeMacBook-Pro redis_bindwidth % sh run.sh 
Time taken to set 300.00 MB data in Redis: 90.308ms
Bandwidth: 3321.96 MB/s
```