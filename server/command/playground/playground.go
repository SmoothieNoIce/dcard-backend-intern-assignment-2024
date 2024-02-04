package main

import "fmt"

// 這個是用來隨便測試程式的
func main() {
	ch := make(chan int, 3)

	go func() {
		for {
			select {
			case v := <-ch:
				fmt.Println(v)
			default:
				fmt.Println("chan no data")
			}
		}
	}()

	for i := 0; i < 10; i++ {
		ch <- i
	}

}
