package main

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/second-state/WasmEdge-go/wasmedge"
	"os"
)

// 定义主机结构体
type Host struct {
	cache *cache.Cache // 注意：缓存对象应该是指针类型
}

// 定义从缓存中Get数据的主机函数
func (h *Host) hostGetDataFromCache(_ interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	keyPointer := params[0].(int32)
	keyLength := params[1].(int32)

	// 获得 mem
	mem := callframe.GetMemoryByIndex(0)

	// 从内存中读取 key
	keyBytes, _ := mem.GetData(uint(keyPointer), uint(keyLength))
	key := string(keyBytes)

	// 从缓存中获取值
	value, exists := h.cache.Get(key)
	if !exists {
		fmt.Printf("Key %s not found in cache\n", key)
		return []interface{}{int32(0), int32(0)}, wasmedge.Result_Success
	}
	fmt.Printf("Retrieved value %v for key %s from cache\n", value, key)

	// 将值转换为字节切片
	valueBytes := []byte(value.(string))

	// 返回值的指针和长度
	return []interface{}{int32(valuePointer), int32(len(valueBytes))}, wasmedge.Result_Success
}

// 定义从缓存中Put数据的主机函数
func (h *Host) hostSetDataToCache(_ interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	keyPointer := params[0].(int32)
	keyLength := params[1].(int32)

	valuePointer := params[2].(int32)
	valueLength := params[3].(int32)

	// 获得 mem
	mem := callframe.GetMemoryByIndex(0)

	// 从内存中读取 key 和 value
	keyBytes, _ := mem.GetData(uint(keyPointer), uint(keyLength))
	key := string(keyBytes)

	valueBytes, _ := mem.GetData(uint(valuePointer), uint(valueLength))
	value := string(valueBytes)

	// 将值写入缓存
	h.cache.Set(key, value, cache.DefaultExpiration)

	fmt.Printf("Set value %v for key %s in cache\n", value, key)
	return nil, wasmedge.Result_Success
}

func main() {
	fmt.Println("Go: Args:", os.Args) // 输出命令行参数
	// 预期 Args[0]: 程序名称 (./externref)
	// 预期 Args[1]: Wasm 文件 (funcs.wasm)

	// 设置不打印调试信息
	wasmedge.SetLogErrorLevel()

	conf := wasmedge.NewConfigure(wasmedge.WASI) // 创建配置
	vm := wasmedge.NewVMWithConfig(conf)         // 创建 Wasm 虚拟机
	obj := wasmedge.NewModule("env")             // 创建 Wasm 模块

	h := Host{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration), // 使用 cache.New() 创建缓存对象并赋值给 cache 字段
	}

	// 将宿主函数添加到模块实例中
	hostGetDataFromCacheType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
		},
		[]wasmedge.ValType{
			wasmedge.ValType_I32,
		})

	hostGetDataFromCache := wasmedge.NewFunction(hostGetDataFromCacheType, h.hostGetDataFromCache, nil, 0)
	obj.AddFunction("host_get_data_from_cache", hostGetDataFromCache)

	hostSetDataToCacheType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
		},
		[]wasmedge.ValType{})

	hostSetDataToCache := wasmedge.NewFunction(hostSetDataToCacheType, h.hostSetDataToCache, nil, 0)
	obj.AddFunction("host_set_data_to_cache", hostSetDataToCache)

	vm.RegisterModule(obj) // 注册模块

	vm.LoadWasmFile(os.Args[1]) // 加载指定的 Wasm 文件
	vm.Validate()               // 验证模块
	vm.Instantiate()            // 实例化模块

	r, _ := vm.Execute("run")                                              // 执行 run 函数
	fmt.Printf("There are %d 'baidu' in source code of baidu.com\n", r[0]) // 输出在百度首页源代码中出现 "baidu" 的次数

	obj.Release()  // 释放模块资源
	vm.Release()   // 释放虚拟机资源
	conf.Release() // 释放配置资源
}
