package main // 主程序包

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/second-state/WasmEdge-go/wasmedge"
)

type host struct { // 定义一个结构体 host
	fetchResult []byte // 存储从网络获取的结果
}

// 执行 HTTP 请求
func fetch(url string) []byte {
	resp, err := http.Get(string(url)) // 发起 HTTP GET 请求
	if err != nil {                    // 如果发生错误
		return nil // 返回空值
	}
	defer resp.Body.Close()            // 在函数返回前关闭响应体
	body, err := io.ReadAll(resp.Body) // 读取响应体数据
	if err != nil {                    // 如果发生错误
		return nil // 返回空值
	}

	return body // 返回响应体数据
}

// 用于获取数据的宿主函数
func (h *host) fetch(_ interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	// 从内存中获取 URL
	pointer := params[0].(int32)                      // 获取 URL 的指针
	size := params[1].(int32)                         // 获取 URL 的长度
	mem := callframe.GetMemoryByIndex(0)              // 获取内存
	data, _ := mem.GetData(uint(pointer), uint(size)) // 从内存中获取数据
	url := make([]byte, size)                         // 创建与 URL 长度相等的字节切片

	copy(url, data) // 将数据复制到字节切片中

	respBody := fetch(string(url)) // 执行 HTTP 请求获取数据

	if respBody == nil { // 如果获取数据失败
		return nil, wasmedge.Result_Fail // 返回失败
	}

	// 存储源代码
	h.fetchResult = respBody // 将获取的数据存储到 fetchResult 中

	return []interface{}{interface{}(len(respBody))}, wasmedge.Result_Success // 返回获取数据的长度和成功标志
}

// 用于写入内存的宿主函数
func (h *host) writeMem(_ interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	// 将源代码写入内存
	pointer := params[0].(int32)                                        // 获取内存指针
	mem := callframe.GetMemoryByIndex(0)                                // 获取内存
	mem.SetData(h.fetchResult, uint(pointer), uint(len(h.fetchResult))) // 将数据写入到内存中

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

	h := host{} // 实例化 host 结构体
	// 将宿主函数添加到模块实例中
	funcFetchType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{
			wasmedge.ValType_I32,
			wasmedge.ValType_I32,
		},
		[]wasmedge.ValType{
			wasmedge.ValType_I32,
		})

	hostFetch := wasmedge.NewFunction(funcFetchType, h.fetch, nil, 0)
	obj.AddFunction("fetch", hostFetch)

	funcWriteType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{
			wasmedge.ValType_I32,
		},
		[]wasmedge.ValType{})
	hostWrite := wasmedge.NewFunction(funcWriteType, h.writeMem, nil, 0)
	obj.AddFunction("write_mem", hostWrite)

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
