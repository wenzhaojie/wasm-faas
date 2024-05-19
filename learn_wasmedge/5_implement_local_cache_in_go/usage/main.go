package main

import (
	"fmt"
	"sync"

	"github.com/wenzhaojie/localcache"
)

func main() {
	cache := localcache.NewCache()

	// 创建一个字典用于存放信息
	info := make(map[string]string)

	// 在多个协程中操作字典
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", index)
			value := fmt.Sprintf("value%d", index)
			cache.Put(key, value)
			fmt.Printf("Put value '%s' for key '%s' into cache\n", value, key)

			// 将信息存入字典
			info[key] = value
		}(i)
	}

	// 等待协程执行完毕
	wg.Wait()

	// 打印字典的所有信息
	fmt.Println("Dictionary Info:")
	for key, value := range info {
		fmt.Printf("%s: %s\n", key, value)
	}
}
