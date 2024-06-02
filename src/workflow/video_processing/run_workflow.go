package main

import (
	"fmt"
	"github.com/second-state/WasmEdge-go/wasmedge"
	"github.com/second-state/wasmedge-bindgen/host/go"
	"os"
	"sync"
)

// 初始化和配置 VM 和 bindgen
func initVM(args []string) (*wasmedge.VM, *bindgen.Bindgen, error) {
	wasmedge.SetLogErrorLevel()
	var conf = wasmedge.NewConfigure(wasmedge.WASI)
	var vm = wasmedge.NewVMWithConfig(conf)

	var wasi = vm.GetImportModule(wasmedge.WASI)
	wasi.InitWasi([]string{}, os.Environ(), []string{".:."})

	err := vm.LoadWasmFile(args[1])
	if err != nil {
		return nil, nil, err
	}
	err = vm.Validate()
	if err != nil {
		return nil, nil, err
	}
	bg := bindgen.New(vm)
	bg.Instantiate()

	return vm, bg, nil
}

// 并行初始化 VM 并执行 Wasm 模块的特定函数
func executeWasmFunctionParallel(args []string, functionName string, p1 string, p2 string, parallelNum int, maxWorkers int) {
	var wg sync.WaitGroup
	taskChan := make(chan struct{}, maxWorkers)
	errors := make(chan error, parallelNum)

	for i := 0; i < parallelNum; i++ {
		wg.Add(1)
		taskChan <- struct{}{}
		go func() {
			defer wg.Done()
			vm, bg, err := initVM(args)
			if err != nil {
				errors <- err
				return
			}
			defer vm.Release() // 确保每个 VM 都被释放

			wasmSuccess, _, err := bg.Execute(functionName, p1, p2)
			if err != nil {
				errors <- err
				return
			}
			if len(wasmSuccess) > 0 {
				fmt.Println("状态:", wasmSuccess[0])
			} else {
				fmt.Println("未找到 wasm 内部状态")
			}
			<-taskChan
		}()
	}

	wg.Wait()
	close(errors)
	close(taskChan)

	// 检查错误
	for err := range errors {
		if err != nil {
			fmt.Println("执行失败:", err)
			return
		}
	}

	fmt.Println("所有任务执行完毕")
}

// 主函数
func main() {
	fmt.Println("Go: Args:", os.Args)

	// 将输入图片放入 Redis
	input_img_path := "input.jpg"
	input_obj_key := "input_img"
	executeWasmFunctionParallel(os.Args, "put_input_img_into_redis", input_img_path, input_obj_key, 1, 1)

	// 设置并行任务数量和最大并发工作者数
	executeWasmFunctionParallel(os.Args, "handler", "input_img", "input_img_putlut", 10, 5)

	fmt.Println("Go: 执行成功")
}
