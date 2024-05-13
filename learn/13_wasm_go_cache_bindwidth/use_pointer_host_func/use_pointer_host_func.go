package main

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/second-state/WasmEdge-go/wasmedge"
	"os"
	"strings"
)

// 声明一个全局的变量来保存缓存对象
var globalCache *cache.Cache

// 定义 tmpHostDataStr 结构体
type tmpHostDataStr struct {
	DataStr string
}

// 定义从缓存中Get数据的主机函数
func (tmpDataStr *tmpHostDataStr) hostGetDataFromCache(_ interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	keyPointer := params[0].(int32)
	keyLength := params[1].(int32)

	// 获得 mem
	mem := callframe.GetMemoryByIndex(0)

	// 从内存中读取 key
	keyBytes, _ := mem.GetData(uint(keyPointer), uint(keyLength))
	key := string(keyBytes)

	// 从缓存中获取值
	value, exists := globalCache.Get(key)
	if !exists {
		fmt.Printf("Key %s not found in cache\n", key)
		return []interface{}{int32(123456)}, wasmedge.Result_Success
	}
	// 打印 value的类型
	// fmt.Printf("Type of value: %T\n", value)

	// 打印 value
	// fmt.Printf("Retrieved value %v for key %s from cache\n", value, key)

	// 将值转换为字节切片
	valueBytes := []byte(value.(string))

	// 写入tmpHostDataStr
	tmpDataStr.DataStr = value.(string)

	// 返回获取数据的长度和成功标志
	return []interface{}{interface{}(len(valueBytes))}, wasmedge.Result_Success // 返回获取数据的长度和成功标志
}

// 定义从缓存中Put数据的主机函数
func (tmpDataStr *tmpHostDataStr) hostSetDataToCache(_ interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
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
	globalCache.Set(key, value, cache.DefaultExpiration)

	// fmt.Printf("Set value %v for key %s in cache\n", value, key)
	return nil, wasmedge.Result_Success
}

// 用于写入内存的宿主函数
func (tmpDataStr *tmpHostDataStr) writeMem(_ interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	// 将源代码写入内存
	pointer := params[0].(int32) // 获取内存指针
	mem := callframe.GetMemoryByIndex(0)

	// fmt.Printf("writeMem: tmpDataStr.DataStr is %s\n", tmpDataStr.DataStr)

	// 转换成Bytes
	dataStrBytes := []byte(tmpDataStr.DataStr)
	// 写入内存
	mem.SetData(dataStrBytes, uint(pointer), uint(len(dataStrBytes))) // 将数据写入到内存中

	return nil, wasmedge.Result_Success // 返回成功标志
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

	globalCache = cache.New(cache.NoExpiration, cache.NoExpiration) // 使用 cache.New() 创建缓存对象
	// 提前存放 key="hust" , value="123123"
	globalCache.Set("hust", "123123", cache.NoExpiration)
	// 创建一个指定大小的字符串
	dataSizeMB := 1000 // MB
	dataSizeBytes := dataSizeMB * 1024 * 1024
	data := strings.Repeat("a", dataSizeBytes)

	// 将数据存储在缓存中，以"data"作为键
	globalCache.Set("data", data, cache.NoExpiration)

	tmp := tmpHostDataStr{}

	// 将宿主函数添加到模块实例中
	hostGetDataFromCacheType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
		},
		[]wasmedge.ValType{
			wasmedge.ValType_I32,
		})

	hostGetDataFromCache := wasmedge.NewFunction(hostGetDataFromCacheType, tmp.hostGetDataFromCache, nil, 0)
	obj.AddFunction("host_get_data_from_cache", hostGetDataFromCache)

	hostSetDataToCacheType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
		},
		[]wasmedge.ValType{})

	hostSetDataToCache := wasmedge.NewFunction(hostSetDataToCacheType, tmp.hostSetDataToCache, nil, 0)
	obj.AddFunction("host_set_data_to_cache", hostSetDataToCache)

	funcWriteType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{
			wasmedge.ValType_I32,
		},
		[]wasmedge.ValType{})
	hostWrite := wasmedge.NewFunction(funcWriteType, tmp.writeMem, nil, 0)
	obj.AddFunction("write_mem", hostWrite)

	vm.RegisterModule(obj) // 注册模块

	vm.LoadWasmFile(os.Args[1]) // 加载指定的 Wasm 文件
	vm.Validate()               // 验证模块
	vm.Instantiate()            // 实例化模块

	r, _ := vm.Execute("run") // 执行 run 函数

	fmt.Printf("The 用时 from wasm module is %d \n", r[0]) // 输出

	// 将毫秒数转换为秒
	transferTimeSeconds := float64(r[0].(int32)) / 1000.0

	// 计算传输速度（单位：MB/s）
	transferSpeedMBps := float64(dataSizeMB) / transferTimeSeconds

	// 计算带宽（单位：Mbps）
	bandwidthMbps := transferSpeedMBps * 8

	fmt.Printf("传输带宽为: %.2f Mbps\n", bandwidthMbps)
	fmt.Printf("传输带宽为: %.2f MB/s\n", transferSpeedMBps)

	obj.Release()  // 释放模块资源
	vm.Release()   // 释放虚拟机资源
	conf.Release() // 释放配置资源
}
