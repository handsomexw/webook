package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// 假设 GetData 是一个返回字节切片的函数
func GetData() []byte {
	// 这里返回一个 JSON 字符串，但实际上可以是任何格式
	return []byte(`{"name":"John","age":30}`)
}

func processBytes(data []byte) (interface{}, error) {
	// 尝试将数据解析为 JSON
	var jsonData interface{}

	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		// 如果不是 JSON 格式，这里可以添加其他逻辑来处理其他格式
		return nil, fmt.Errorf("invalid data format: %w", err)
	}

	return jsonData, nil
}

func Gorting(cnt int) {
	go func(cnt int) {
		for i := 0; i < cnt; i++ {
			fmt.Println(i)
			time.Sleep(time.Second)
		}
	}(cnt)
	fmt.Println("goroutine完成")

}

func main() {
	//data := GetData()
	//
	//result, err := processBytes(data)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//
	//// 现在 result 是原始数据的反序列化版本
	//fmt.Printf("Result: %v\n", result)
	cnt := 5
	Gorting(cnt)
	time.Sleep(time.Second * 10)

}
