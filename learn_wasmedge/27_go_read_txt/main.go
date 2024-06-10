package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

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
	fileString := string(content)

	// 删除里面所有的换行符，替换为空格
	fileString = strings.ReplaceAll(fileString, "\n", " ")

	return fileString, nil
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

func main() {
	filename := "example.txt"
	// 测试读取文件
	content, err := readFileAsString(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("File content:\n", content)

	// 测试分割字符串
	parts := splitString(content, 10)
	fmt.Println("Split into", len(parts), "parts:")

	// 打印所有分割片段
	for i, part := range parts {
		fmt.Printf("Part %d: %s\n", i+1, part)
	}
}
