# 目标
在go程序中，有的时候希望在多个协程间，共享内存数据。因此，需要在go程序内部，创建一个本地cache对象，提供key value的put，get接口。

# 步骤
```bash
cd usage
go build
```


# 结果
```bash
wzj@ZhaojiedeMacBook-Pro usage % ./usage 
Put value 'value1' for key 'key1' into cache
Put value 'value0' for key 'key0' into cache
Put value 'value4' for key 'key4' into cache
Put value 'value3' for key 'key3' into cache
Put value 'value2' for key 'key2' into cache
Dictionary Info:
key1: value1
key0: value0
key4: value4
key3: value3
key2: value2

```