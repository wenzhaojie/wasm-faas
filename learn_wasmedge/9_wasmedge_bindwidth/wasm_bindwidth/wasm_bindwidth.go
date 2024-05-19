package main

import (
	"fmt"
	"github.com/second-state/WasmEdge-go/wasmedge"
	"math"
	"os"
	"strings"
	"time"
	"unsafe"
)

func main() {
	// 期望的 Args[0]: 程序名称 (./bindgen_funcs)
	// 期望的 Args[1]: wasm 文件 (rust_bindgen_funcs_lib.wasm))

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

	// 准备输入数据
	const dataSizeMB = 10 // 希望测试的字符串大小，以 MB 为单位

	// 计算数据大小，确保不小于 1 字节
	dataSize := int(math.Round(dataSizeMB * 1024 * 1024))
	if dataSize < 1 {
		dataSize = 1
	}

	// 生成输入数据
	inputData := strings.Repeat("x", dataSize)

	// 转换输入数据为字节切片
	inputBytes := []byte(inputData)

	// 获取开始时间
	startTime := time.Now()

	// 获取输入数据的指针和长度
	dataPointer := unsafe.Pointer(&inputBytes[0])
	dataLength := len(inputBytes)

	// 将指针和长度转换为 int32 类型
	ptrInt32 := int32(uintptr(dataPointer))
	lengthInt32 := int32(dataLength)

	// 加载 Wasm 文件并执行，传入数据指针和长度
	_, err := vm.RunWasmFile(os.Args[1], "bandwidth", ptrInt32, lengthInt32)

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
