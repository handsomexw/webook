package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// 指定源文件和目标文件
	originalFile := "C:\\Users\\artrc\\Desktop\\测试.txt"
	copyFile := "C:\\Users\\artrc\\Desktop\\111.txt"

	// 以只读模式打开原始文件
	src, err := os.OpenFile(originalFile, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error opening original file:", err)
		return
	}
	defer src.Close()

	// 以写入模式创建新文件
	dst, err := os.Create(copyFile)
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}

	fmt.Println("File copied successfully.")
}
