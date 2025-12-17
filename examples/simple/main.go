package main

import (
	"fmt"
	"time"

	"github.com/morsuning/gopool"
)

func main() {
	// Create a pool with 100 workers
	// 创建一个包含 100 个工作者的池
	pool := gopool.NewGoPool(100)
	defer pool.Release()

	for i := 0; i < 1000; i++ {
		pool.AddTask(func() (any, error) {
			time.Sleep(10 * time.Millisecond)
			fmt.Println("Task executed")
			return nil, nil
		})
	}
	pool.Wait()
}
