package main

import (
	"fmt"
	"github.com/second-state/WasmEdge-go/wasmedge"
	bindgen "github.com/second-state/wasmedge-bindgen/host/go"
	"math"
	"os"
	"strings"
	"time"
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

	// 准备输入数据
	// 准备输入数据
	const dataSizeMB = 1000 // 希望测试的字符串大小，以 MB 为单位

	// 计算数据大小，确保不小于 1 字节
	dataSize := int(math.Round(dataSizeMB * 1024 * 1024))
	if dataSize < 1 {
		dataSize = 1
	}

	// 生成输入数据
	inputData := strings.Repeat("x", dataSize)

	// 获取开始时间
	startTime := time.Now()

	// 调用 wasm 模块中的 helloworld 函数
	wasm_success, _, err := bg.Execute("helloworld", "inputData")
	// 打印wasm内部的耗时
	if len(wasm_success) > 0 {
		fmt.Println("状态:", wasm_success[0])
	} else {
		fmt.Println("未找到 wasm 内部状态")
	}

	// 调用 wasm 模块中的 bandwidth 函数
	wasm_duration, _, err := bg.Execute("bandwidth", inputData)
	// 打印wasm内部的耗时
	if len(wasm_duration) > 0 {
		fmt.Println("wasm内部打点计时:", wasm_duration[0])
	} else {
		fmt.Println("未找到 wasm 内部打点计时")
	}

	// 获取结束时间
	endTime := time.Now()

	if err != nil {
		fmt.Println("运行 bindgen -- 带宽 测试失败")
		return
	}

	// 计算耗时
	duration := endTime.Sub(startTime).Seconds()

	// 计算带宽，并将结果转换为 MB/s
	bandwidthMB := float64(dataSize) / (duration * 1024 * 1024)

	// 打印结果
	fmt.Printf("带宽: %f MB/s\n", bandwidthMB)
}
