package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/second-state/WasmEdge-go/wasmedge"
)

// 定义一个简单的缓存结构
type Cache struct {
	data map[string]string
	mu   sync.RWMutex
}

// 创建一个新的缓存实例
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

// 存储键值对到缓存中
func (c *Cache) Put(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// 从缓存中读取键对应的值
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

type host struct { // 定义一个结构体 host
	cache *Cache // 存储缓存实例
}

// 定义一个宿主函数，用于在 WebAssembly 模块中调用缓存的读取和存储接口
func (h *host) CacheHostFunction(data interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	if len(params) < 2 {
		return nil, wasmedge.Result_Fail
	}

	cache, ok := data.(*Cache)
	if !ok {
		return nil, wasmedge.Result_Fail
	}

	method, ok := params[0].(string)
	if !ok {
		return nil, wasmedge.Result_Fail
	}

	switch method {
	case "put":
		if len(params) < 3 {
			return nil, wasmedge.Result_Fail
		}
		key, ok := params[1].(string)
		if !ok {
			return nil, wasmedge.Result_Fail
		}
		value, ok := params[2].(string)
		if !ok {
			return nil, wasmedge.Result_Fail
		}
		cache.Put(key, value)
		return nil, wasmedge.Result_Fail

	case "get":
		if len(params) < 2 {
			return nil, wasmedge.Result_Fail
		}
		key, ok := params[1].(string)
		if !ok {
			return nil, wasmedge.Result_Fail
		}
		value, ok := cache.Get(key)
		if !ok {
			return nil, wasmedge.Result_Fail
		}
		return []interface{}{value}, wasmedge.Result_Fail

	default:
		return nil, wasmedge.Result_Fail
	}
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

	h := host{} // 实例化 host 结构体
	// 将宿主函数添加到模块实例中
	// 将宿主函数添加到模块实例中
	funcCacheType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{
			wasmedge.ValType_ExternRef, // host 结构体
			wasmedge.ValType_I32,       // key
			wasmedge.ValType_I32,       // value
		},
		[]wasmedge.ValType{
			wasmedge.ValType_I32, // 返回值
		})
	hostCache := wasmedge.NewFunction(
		funcCacheType,
		h.CacheHostFunction,
		nil,
		0,
	)
	obj.AddFunction("cache", hostCache)

	// 创建一个缓存实例
	h.cache = NewCache()

	res, err := vm.RunWasmFile(os.Args[1], "cache")
	if err == nil {
		fmt.Println("获取结果:", res[0].(int64))
	} else {
		fmt.Println("运行失败:", err.Error())
	}

}
