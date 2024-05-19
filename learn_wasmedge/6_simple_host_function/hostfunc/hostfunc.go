package main // 主程序包

import (
	"fmt"
	"os"
	"time"

	"github.com/second-state/WasmEdge-go/wasmedge"
)

func host_add(data interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	// add: i32, i32 -> i32
	res := params[0].(int32) + params[1].(int32)

	// print res
	fmt.Printf("host_add res: %d\n", res)

	// Set the returns
	returns := make([]interface{}, 1)
	returns[0] = res

	// Return
	return returns, wasmedge.Result_Success
}

func host_get_datetime(data interface{}, callframe *wasmedge.CallingFrame, params []interface{}) ([]interface{}, wasmedge.Result) {
	// host_get_datetime: None -> i32

	// 获得当前年份的数字
	year := time.Now().Year()

	// print year
	fmt.Printf("host_get_datetime year: %d\n", year)

	// 构建返回结果
	returns := []interface{}{year}

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

	// Create a function type: {i32, i32} -> {i32}.
	funcAddType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{wasmedge.ValType_I32, wasmedge.ValType_I32},
		[]wasmedge.ValType{wasmedge.ValType_I32},
	)
	func_add := wasmedge.NewFunction(funcAddType, host_add, nil, 0)
	obj.AddFunction("host_add", func_add)

	// Create a function type: {} -> {i32}.
	funcGetDateTimeType := wasmedge.NewFunctionType(
		[]wasmedge.ValType{},
		[]wasmedge.ValType{wasmedge.ValType_I32},
	)
	func_get_datetime := wasmedge.NewFunction(funcGetDateTimeType, host_get_datetime, nil, 0)
	obj.AddFunction("host_get_datetime", func_get_datetime)

	vm.RegisterModule(obj)      // 注册模块
	vm.LoadWasmFile(os.Args[1]) // 加载指定的 Wasm 文件
	vm.Validate()               // 验证模块
	vm.Instantiate()            // 实例化模块

	r, _ := vm.Execute("run") // 执行 run 函数
	fmt.Printf("add result: %d\n", r[0])

	obj.Release()  // 释放模块资源
	vm.Release()   // 释放虚拟机资源
	conf.Release() // 释放配置资源

}
