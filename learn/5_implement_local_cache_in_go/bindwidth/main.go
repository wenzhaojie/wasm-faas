package main

import (
	"fmt"
	"time"

	"github.com/wenzhaojie/localcache" // 请替换为实际的包路径
)

func main() {
	cache := localcache.NewCache()

	// 设置数据大小
	dataSizeMB := 100 // 假设为 100 MB
	dataSize := dataSizeMB * 1024 * 1024

	// 生成数据
	data := make([]byte, dataSize)
	for i := 0; i < dataSize; i++ {
		data[i] = byte(i % 256)
	}

	// 存储测试
	start := time.Now()
	cache.Put("testKey", data)
	storageTime := time.Since(start)

	// 读取测试
	start = time.Now()
	_, ok := cache.Get("testKey")
	if !ok {
		fmt.Println("Key not found in cache")
		return
	}
	readTime := time.Since(start)

	// 计算存储和读取带宽
	storageBandwidth := float64(dataSize) / storageTime.Seconds()
	readBandwidth := float64(dataSize) / readTime.Seconds()

	// 转换为 MB/s
	storageBandwidth /= 1024 * 1024
	readBandwidth /= 1024 * 1024

	fmt.Printf("存储带宽: %.2f MB/秒\n", storageBandwidth)
	fmt.Printf("读取带宽: %.2f MB/秒\n", readBandwidth)

}
