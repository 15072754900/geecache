package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		// wg.Add(1)
		// go printOnce(10, &wg)
		go printOnce(100)
	}
	// wg.Wait() // 使用全局的同步信息，等待协程完毕
	time.Sleep(2 * time.Second) // 使用时间使协程运行完
}

// 定义一个全局并发访问的变量，进行打印

var set = make(map[int]bool, 0)
var m sync.Mutex

func printOnce(num int) {
	m.Lock()
	defer m.Unlock()
	if _, ok := set[num]; !ok {
		fmt.Println(num)
	}
	set[num] = true
}

// 第二种方法：使用sync.waitGroup
// 第三种方法：使用互斥锁sync.Mutex

// 前两种方法是使所有协程跑完，第三种才是使并发修改变为安全的
