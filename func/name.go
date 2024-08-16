package main

import (
	"fmt"
	"os"
)

func main() {
	SecretKey := os.Getenv("SMSSecretID")
	SecretId := os.Getenv("SMSSecretKey")
	fmt.Println(SecretId, SecretKey)
	//key := os.Getenv("TestName")
	//name := make([]string, 5, 10)
	//value := make(chan int, 5)
	//for i, _ := range name {
	//	name[i] = strconv.Itoa(i)
	//}
	//go func() {
	//	defer close(value)
	//	for i := 0; i < 5; i++ {
	//		value <- i * 2
	//	}
	//}()
	//
	//time.Sleep(time.Second * 2)
	//
	//for i := range value {
	//	fmt.Println(i)
	//}
	//
	//fmt.Printf("key=%s\n", key)
	//age1 := 1
	//var age = (*int)(nil)
	//var age2 int = int(6)
	//fmt.Println(age1, age2)
	//
	//age = &age1
	//fmt.Println(age)
}
