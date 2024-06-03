package main

import (
	"fmt"
	"github.com/second-state/WasmEdge-go/wasmedge"
	bindgen "github.com/second-state/wasmedge-bindgen/host/go"
	"os"
)

func main() {
	// 期望的 Args[0]: 程序名称 (./bindgen_funcs)
	// 期望的 Args[1]: wasm 文件 (rust_bindgen_funcs_lib.wasm))
	fmt.Println("Go: Args:", os.Args) // 输出命令行参数

	// 设置不打印调试信息
	wasmedge.SetLogErrorLevel()

	// 创建配置
	var conf = wasmedge.NewConfigure(wasmedge.WASI)

	// 使用配置创建虚拟机
	var vm = wasmedge.NewVMWithConfig(conf)

	// 初始化 WASI
	var wasi = vm.GetImportModule(wasmedge.WASI)
	wasi.InitWasi(
		os.Args[1:],     // 参数
		os.Environ(),    // 环境变量
		[]string{".:."}, // 映射的路径
	)

	// 加载并验证 wasm
	vm.LoadWasmFile(os.Args[1])
	vm.Validate()

	// 实例化 bindgen 和虚拟机
	bg := bindgen.New(vm)
	bg.Instantiate()

	// 调用 wasm 模块中的 split_text 函数
	inputData := "Hello, World!"

	wasm_result, _, err := bg.Execute("split_text", inputData)
	// 打印wasm内部的耗时
	if len(wasm_result) > 0 {
		fmt.Println("wasm_result:", wasm_result[0])
	} else {
		fmt.Println("未找到 wasm_result")
	}

	if err != nil {
		fmt.Println("运行 bindgen -- 测试失败")
		return
	}
}
