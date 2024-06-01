package main

import (
	"fmt"
	bindgen "github.com/second-state/wasmedge-bindgen/host/go"
	"os"

	"github.com/second-state/WasmEdge-go/wasmedge"
)

func main() {
	// putwatermark.wasm
	// putlut.wasm
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

	// 调用 wasm 模块中的 helloworld 函数
	helloworld_str := "Hello, World!"
	wasmSuccess, _, err := bg.Execute("helloworld", helloworld_str)
	// 打印wasm内部状态
	if len(wasmSuccess) > 0 {
		fmt.Println("状态:", wasmSuccess[0])
	} else {
		fmt.Println("未找到 wasm 内部状态")
	}

	if err != nil {
		fmt.Println("运行 bindgen -- 失败")
		return
	}

	// 调用 Wasm 模块中的 put_input_img_into_redis
	input_img_path := "input.jpg"
	input_obj_key := "input_img"

	// put_input_img_into_redis
	wasmSuccess, _, err = bg.Execute("put_input_img_into_redis", input_img_path, input_obj_key)
	// 打印wasm内部状态
	if len(wasmSuccess) > 0 {
		fmt.Println("状态:", wasmSuccess[0])
	} else {
		fmt.Println("未找到 wasm 内部状态")
	}

	if err != nil {
		fmt.Println("运行 bindgen -- 失败")
		return
	}

	// 调用 handler
	output_obj_key := "input_img_putlut"
	wasmSuccess, _, err = bg.Execute("handler", input_obj_key, output_obj_key)
	// 打印wasm内部状态
	if len(wasmSuccess) > 0 {
		fmt.Println("状态:", wasmSuccess[0])
	} else {
		fmt.Println("未找到 wasm 内部状态")
	}

	if err != nil {
		fmt.Println("运行 bindgen -- 失败")
		return
	}

	// get_output_img_from_redis
	output_img_path := "putlut_output.jpg"

	wasmSuccess, _, err = bg.Execute("get_output_img_from_redis", output_img_path, output_obj_key)
	// 打印wasm内部状态
	if len(wasmSuccess) > 0 {
		fmt.Println("状态:", wasmSuccess[0])
	} else {
		fmt.Println("未找到 wasm 内部状态")
	}

	if err != nil {
		fmt.Println("运行 bindgen -- 失败")
		return
	}

	// 释放资源
	vm.Release()
	conf.Release()
}
