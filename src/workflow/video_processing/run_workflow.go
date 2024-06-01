package main

import (
	"fmt"
	bindgen "github.com/second-state/wasmedge-bindgen/host/go"
	"os"

	"github.com/second-state/WasmEdge-go/wasmedge"
)

// 初始化和配置 VM 和 bindgen
func initVM(args []string) (*wasmedge.VM, *bindgen.Bindgen, error) {
	wasmedge.SetLogErrorLevel()
	var conf = wasmedge.NewConfigure(wasmedge.WASI)
	var vm = wasmedge.NewVMWithConfig(conf)

	var wasi = vm.GetImportModule(wasmedge.WASI)
	wasi.InitWasi(args, os.Environ(), []string{".:."})

	vm.LoadWasmFile(args[1])
	vm.Validate()
	bg := bindgen.New(vm)
	bg.Instantiate()

	return vm, bg, nil
}

// 调用 Wasm 模块的特定函数
func executeWasmFunction(bg *bindgen.Bindgen, functionName string, params ...interface{}) ([]interface{}, error) {
	wasmSuccess, _, err := bg.Execute(functionName, params...)
	if len(wasmSuccess) > 0 {
		fmt.Println("状态:", wasmSuccess[0])
	} else {
		fmt.Println("未找到 wasm 内部状态")
	}
	return wasmSuccess, err
}

// 主函数
func main() {
	fmt.Println("Go: Args:", os.Args)

	vm, bg, err := initVM(os.Args[1:])
	if err != nil {
		fmt.Println("VM 初始化失败:", err)
		return
	}
	defer vm.Release()

	// 执行 Wasm 函数
	if _, err := executeWasmFunction(bg, "helloworld", "Hello, World!"); err != nil {
		fmt.Println("执行 helloworld 失败:", err)
		return
	}

	if _, err := executeWasmFunction(bg, "put_input_img_into_redis", "input.jpg", "input_img"); err != nil {
		fmt.Println("执行 put_input_img_into_redis 失败:", err)
		return
	}

	if _, err := executeWasmFunction(bg, "handler", "input_img", "input_img_putlut"); err != nil {
		fmt.Println("执行 handler 失败:", err)
		return
	}

	if _, err := executeWasmFunction(bg, "get_output_img_from_redis", "putlut_output.jpg", "input_img_putlut"); err != nil {
		fmt.Println("执行 get_output_img_from_redis 失败:", err)
		return
	}
}
