package main // 主程序包

import (
	"fmt"
	"os"

	"github.com/second-state/WasmEdge-go/wasmedge"
)

type host struct { // 定义一个结构体 host
	addResult int32 // 存储结果
}

func (h *host) host_add(data interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	// add: i32, i32 -> i32
	res := params[0].(int32) + params[1].(int32)

	// print res
	fmt.Printf("host_add res: %d\n", res)

	// Set the returns
	returns := make([]interface{}, 1)
	returns[0] = res

	// 将结果存储到结构体中
	h.addResult = res

	// Return
	return returns, wasmedge.Result_Success
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

	// Create a function type: {i32, i32} -> {i32}.
	funcAddType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{wasmedge.ValType_I32, wasmedge.ValType_I32},
		[]wasmedge.ValType{wasmedge.ValType_I32},
	)
	func_add := wasmedge.NewFunction(funcAddType, h.host_add, nil, 0)
	obj.AddFunction("host_add", func_add)

	vm.RegisterModule(obj)      // 注册模块
	vm.LoadWasmFile(os.Args[1]) // 加载指定的 Wasm 文件
	vm.Validate()               // 验证模块
	vm.Instantiate()            // 实例化模块

	r, _ := vm.Execute("run") // 执行 run 函数
	fmt.Printf("add result from return: %d\n", r[0])

	// 从host struct中提取结果
	fmt.Printf("add result from struct: %d\n", h.addResult)

	obj.Release()  // 释放模块资源
	vm.Release()   // 释放虚拟机资源
	conf.Release() // 释放配置资源

}
