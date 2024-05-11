// package main

// import "fmt"

// var ch8 = make(chan int, 6)

// func mm1() {
// 	for i := 0; i < 10; i++ {
// 		ch8 <- 8 * i
// 	}
// 	// close(ch8)
// }
// func main() {

// 	go mm1()
// 	for {
// 		for data := range ch8 {
// 			fmt.Print(data,"\t")
// 		}
// 		break
// 	}
// }

package main

import (
	"fmt"
	"strings"
	"time"
)

func main1() {
	ch := make(chan int) // 创建一个无缓冲的channel

	// 启动一个goroutine发送数据到channel
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(time.Second) // 模拟发送延迟
		}
		close(ch) // 发送完所有数据后关闭channel
	}()

	// 在主goroutine中接收数据
	// aaa:
	for {
		// 使用两个变量来接收值，第二个变量是一个bool值，表示channel是否已关闭
		value, ok := <-ch
		if !ok {
			// 如果channel已关闭，并且没有更多的数据可接收，则退出循环
			fmt.Println("Channel closed")
			// break aaa
			break
		}
		fmt.Println("Received:", value)
	}
}

// 当一个channel被关闭后，尝试向它发送数据会引发panic，但可以继续从中接收数据，直到channel为空。
// 当从已关闭的channel中接收数据时，接收操作会立即返回【零值】，并且第二个返回值（通常是一个布尔值）会指示channel是否已关闭。

func main2() {
	type aa struct {
		a1 int
		a2 string
	}

	d := aa{1, "mlee"}
	for index, num := range strings.Repeat("mlee", 100) {
		go func(index int, num rune) {
			d.a2 += string(num)
			d.a1 += index
		}(index, num)
	}
	fmt.Println(d)

	time.Sleep(10 * time.Second)

}
