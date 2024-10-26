package main

import (
	"fmt"
)

func producer(ch chan<- int) {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch)
}

func consumer(ch <-chan int) {
    for v := range ch {
        fmt.Println(v)
    }
}

func main() {
    ch := make(chan int) // 創建雙向通道

    // 將雙向通道轉換為單向通道，分別傳遞給 producer 和 consumer
    go producer(ch)      // 傳遞為發送專用通道 (chan<- int)
    go consumer(ch)      // 傳遞為接收專用通道 (<-chan int)

    // 為了防止 main 函數過早退出，加入簡單的等待
    var input string
    fmt.Scanln(&input)
}
