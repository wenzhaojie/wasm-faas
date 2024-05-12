package main

import (
	"fmt"
	"github.com/second-state/WasmEdge-go/wasmedge"
	"os"
)

// 定义主机结构体
type Host struct {
	cache map[int32]int32
}

// 定义从缓存中获取数据的主机函数
func (h *Host) hostGetDataFromCache(data interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	key := params[0].(int32)
	value, exists := h.cache[key]
	if !exists {
		fmt.Printf("Key %s not found in cache\n", key)
		return []interface{}{""}, wasmedge.Result_Success
	}
	fmt.Printf("Retrieved value %d for key %s from cache\n", value, key)
	return []interface{}{value}, wasmedge.Result_Success
}

// 定义将数据存入缓存的主机函数
func (h *Host) hostPutDataToCache(data interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	key := params[0].(int32)
	value := params[1].(int32)
	h.cache[key] = value
	fmt.Printf("Stored value %d for key %s into cache\n", value, key)
	return []interface{}{}, wasmedge.Result_Success
}

func main() {
	// 初始化主机结构体
	h := Host{
		cache: make(map[int32]int32),
	}

	// 创建配置
	conf := wasmedge.NewConfigure(wasmedge.WASI)

	// 创建虚拟机
	vm := wasmedge.NewVMWithConfig(conf)

	// 创建模块
	obj := wasmedge.NewModule("env")

	// 创建获取数据的主机函数
	getDataFuncType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{wasmedge.ValType_I32},
		[]wasmedge.ValType{wasmedge.ValType_I32},
	)
	getDataFunc := wasmedge.NewFunction(getDataFuncType, h.hostGetDataFromCache, nil, 0)
	obj.AddFunction("host_get_data_from_cache", getDataFunc)

	// 创建存储数据的主机函数
	putDataFuncType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{wasmedge.ValType_I32, wasmedge.ValType_I32},
		nil,
	)
	putDataFunc := wasmedge.NewFunction(putDataFuncType, h.hostPutDataToCache, nil, 0)
	obj.AddFunction("host_put_data_to_cache", putDataFunc)

	// 注册模块
	vm.RegisterModule(obj)

	// 加载 Wasm 文件
	vm.LoadWasmFile(os.Args[1])

	// 验证并实例化模块
	vm.Validate()
	vm.Instantiate()

	// 执行 run 函数
	r, _ := vm.Execute("run")
	fmt.Printf("Run function returned: %d\n", r[0])

	// 查看 host里面的内容
	// 打印缓存
	fmt.Println("Go cache contents:")
	for key, value := range h.cache {
		fmt.Printf("Key: %d, Value: %d\n", key, value)
	}

	// 释放资源
	obj.Release()
	vm.Release()
	conf.Release()
}
