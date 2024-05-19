# 目标
测试https://github.com/patrickmn/go-cache

# 步骤
```bash
cd use_go_cache
go build use_cache.go
./use_cache
```

# 结果
```bash
Stored value 456 for key %!s(int32=123) into cache
Retrieved value 456 for key %!s(int32=123) from cache
Run function returned: 456
Go cache contents:
Key: 123, Value: 456

```

