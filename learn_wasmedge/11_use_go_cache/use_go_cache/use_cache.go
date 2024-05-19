package main

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

func main() {
	// 创建一个默认过期时间为5分钟、每10分钟清理过期项的缓存
	c := cache.New(5*time.Minute, 10*time.Minute)

	// 将键"foo"的值设置为"bar"，使用默认的过期时间
	c.Set("foo", "bar", cache.DefaultExpiration)

	// 将键"baz"的值设置为42，没有过期时间
	// （该项直到被重新设置或使用c.Delete("baz")删除时才会被移除）
	c.Set("baz", 42, cache.NoExpiration)

	// 从缓存中获取与键"foo"相关联的字符串
	foo, found := c.Get("foo")
	if found {
		fmt.Println("foo:", foo)
	}

	// 假设有一个 MyFunction 函数可以处理字符串参数
	MyFunction := func(s string) {
		fmt.Println("MyFunction called with:", s)
	}

	// 由于 Go 是静态类型的，而缓存值可以是任何类型，当将值传递给不接受任意类型的函数（即接口{}）时，需要进行类型断言。
	// 对于只会使用一次的值，例如传递给另一个函数，最简单的方法是：
	foo, found = c.Get("foo")
	if found {
		MyFunction(foo.(string))
	}

	// 如果在同一个函数中多次使用该值，这样做会变得很麻烦。
	// 您可以选择以下任一方式：
	if x, found := c.Get("foo"); found {
		foo := x.(string)
		fmt.Println("foo used multiple times:", foo)
	}
	// 或者
	var foo2 string
	if x, found := c.Get("foo"); found {
		foo2 = x.(string)
		fmt.Println("foo used multiple times (alternative):", foo2)
	}

	// 希望更高性能？存储指针！
	type MyStruct struct {
		Field1 string
		Field2 int
	}

	// 假设有一个 MyStruct 结构体
	myStruct := MyStruct{"Hello", 42}

	// 存储 MyStruct 的指针
	c.Set("myStruct", &myStruct, cache.DefaultExpiration)
	if x, found := c.Get("myStruct"); found {
		foo := x.(*MyStruct)
		fmt.Println("myStruct:", foo)
	}
}
