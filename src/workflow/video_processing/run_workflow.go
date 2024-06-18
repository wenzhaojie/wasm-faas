package main

import (
	"fmt"
	"github.com/second-state/WasmEdge-go/wasmedge"
	bindgen "github.com/second-state/wasmedge-bindgen/host/go"
	"os"
	"sync"
	"time"
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

func testWorkflow() {
	fmt.Println("Go: Args:", os.Args)

	// 将输入图片放入 Redis
	input_img_path := "input.jpg"
	input_obj_key := "input_img"
	executeWasmFunctionParallel(os.Args, "put_input_img_into_redis", input_img_path, input_obj_key, 1, 1)

	// 设置并行任务数量和最大并发工作者数
	start_t := time.Now()
	executeWasmFunctionParallel(os.Args, "handler", "input_img", "input_img_putlut", 20, 10)
	end_t := time.Now()
	fmt.Println("Go: 执行时间:", end_t.Sub(start_t))
	fmt.Println("Go: 执行成功")
}

func testWorkflowWithDiffParameters() {
	fmt.Println("Go: Args:", os.Args)

	// 将输入图片放入 Redis
	input_img_path := "input.jpg"
	input_obj_key := "input_img"
	executeWasmFunctionParallel(os.Args, "put_input_img_into_redis", input_img_path, input_obj_key, 1, 1)

	parallelNum_list := []int{10, 20, 40, 80}
	maxWorkers_list := []int{1, 2, 4, 8}
	result_dict_list := make([]map[string]interface{}, 0)

	for _, parallelNum := range parallelNum_list {
		for _, maxWorkers := range maxWorkers_list {
			// 设置并行任务数量和最大并发工作者数
			start_t := time.Now()
			executeWasmFunctionParallel(os.Args, "handler", "input_img", "input_img_putlut", parallelNum, maxWorkers)
			end_t := time.Now()
			// 执行时间，换算秒
			duration := end_t.Sub(start_t).Seconds()

			fmt.Println("parallelNum: 并行任务数量", parallelNum)
			fmt.Println("maxWorkers: 最大并发工作者数", maxWorkers)
			fmt.Println("Go: 执行时间:", duration)
			fmt.Println("Go: 执行成功")
			result_dict := map[string]interface{}{
				"parallelNum": parallelNum,
				"maxWorkers":  maxWorkers,
				"duration":    end_t.Sub(start_t),
			}
			result_dict_list = append(result_dict_list, result_dict)
		}
	}
	fmt.Println(result_dict_list)
	fmt.Println("Go: 执行完毕")
}

// 主函数
func main() {
	// testWorkflow()
	testWorkflowWithDiffParameters()
}

// parallelNum: 并行任务数量 100
// maxWorkers: 最大并发工作者数 5
// Go: 执行时间: 33.7446385s

// parallelNum: 并行任务数量 100
// maxWorkers: 最大并发工作者数 10
// Go: 执行时间: 20.028265167s

// parallelNum: 并行任务数量 20
// maxWorkers: 最大并发工作者数 10
// Go: 执行时间: 4.005883375s
