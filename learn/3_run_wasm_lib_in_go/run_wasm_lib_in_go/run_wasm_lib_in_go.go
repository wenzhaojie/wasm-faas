package main

import (
	"fmt"
	"log"
	"os"

	"github.com/second-state/WasmEdge-go/wasmedge"
)

func main() {
	// 打印程序名称和 Wasm 文件名称
	fmt.Println("Go: Args:", os.Args)

	// 检查是否提供了 Wasm 文件路径作为参数
	if len(os.Args) != 2 {
		log.Fatalf("请提供 Wasm 文件路径作为参数！")
	}

	// 设置日志级别
	wasmedge.SetLogErrorLevel()

	// 创建配置上下文并添加 WASI 支持
	// 除非你需要 WASI 支持，否则此步骤不是必需的。
	conf := wasmedge.NewConfigure(wasmedge.WASI)

	// 使用配置创建虚拟机
	vm := wasmedge.NewVMWithConfig(conf)

	// 加载 Wasm 文件并执行
	res, err := vm.RunWasmFile(os.Args[1], "fibonacci", uint64(21))
	if err == nil {
		fmt.Println("获取第 21 项斐波那契数:", res[0].(int64))
	} else {
		fmt.Println("运行失败:", err.Error())
	}

	// 释放资源
	vm.Release()
	conf.Release()
}
