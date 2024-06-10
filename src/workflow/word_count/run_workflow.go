package main

import (
	"fmt"
	"github.com/second-state/WasmEdge-go/wasmedge"
	"github.com/second-state/wasmedge-bindgen/host/go"
	"io"
	"os"
	"strings"
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

// executeWasmFunctionParallel 并行初始化 VM 并执行 Wasm 模块的特定函数
// args 是用于加载 wasm 文件的参数
// functionName 是要执行的函数名
// paraList 是参数列表
// parallelNum 是并行任务数量
// maxWorkers 是最大并发工作者数
func executeWasmFunctionParallel(args []string, functionName string, paraList []interface{}, parallelNum int, maxWorkers int) {
	var wg sync.WaitGroup
	taskChan := make(chan struct{}, maxWorkers)
	errors := make(chan error, parallelNum)

	for i := 0; i < parallelNum; i++ {
		wg.Add(1)
		taskChan <- struct{}{}
		go func(para interface{}) {
			defer wg.Done()
			vm, bg, err := initVM(args)
			if err != nil {
				errors <- err
				<-taskChan
				return
			}
			defer vm.Release() // 确保每个 VM 都被释放

			wasmSuccess, _, err := bg.Execute(functionName, para)
			if err != nil {
				errors <- err
				<-taskChan
				return
			}
			if len(wasmSuccess) > 0 {
				fmt.Println("状态:", wasmSuccess[0])
			} else {
				fmt.Println("未找到 wasm 内部状态")
			}
			<-taskChan
		}(paraList[i])
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

func readFileAsString(filename string) (string, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(file)

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// 将内容转换为字符串
	file_string := string(content)

	// 删除里面所有的换行符，替换为空格
	file_string = strings.ReplaceAll(file_string, "\n", " ")

	return file_string, nil

}

// 将字符串平均分割成 n 份，不要切断单词
func splitString(s string, n int) []string {
	// 计算每份的长度
	l := len(s)
	if l <= n {
		return []string{s}
	}
	per := l / n

	// 逐个字符判断，找到合适的分割位置
	ret := make([]string, 0)
	start := 0

	for i := 0; i < n-1; i++ {
		end := start + per
		if end >= len(s) {
			break
		}

		// 尽量在单词边界分割
		for end < len(s) && s[end] != ' ' {
			end++
		}

		// 处理特殊情况：没有找到空格
		if end == len(s) {
			break
		}

		ret = append(ret, s[start:end])
		start = end + 1
	}
	// 添加最后一段
	ret = append(ret, s[start:])

	return ret
}

// 主函数
func main() {
	fmt.Println("Go: Args:", os.Args)

	// 先读取本地txt文件，input.txt, 变成字符串
	fileString, err := readFileAsString("input.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 将字符串平均分割成 10 份，不要切断单词
	splitStrings := splitString(fileString, 10)

	// 将 []string 转换为 []interface{}
	interfaces := make([]interface{}, len(splitStrings))
	for i, v := range splitStrings {
		interfaces[i] = v
	}

	// 设置并行任务数量和最大并发工作者数
	executeWasmFunctionParallel(os.Args, "word_count", interfaces, 10, 5)

	fmt.Println("Go: 执行成功")
}
