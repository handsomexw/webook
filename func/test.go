package main

import "fmt"

func Types[T any](name T) {
	fmt.Printf("type %T\n", name)
}

func main() {
	var name string
	name = "猜猜看"
	var name2 int
	name2 = 2
	name3 := 3.14
	Types(name)
	Types(name2)
	Types(name3)
	//kk := make(chan string, 5)
	//ww := make([]int, 5)

}
