package main

import (
	"fmt"
	"time"

	"github.com/daniel-hutao/spinlock"
	"github.com/morsuning/lifopool"
)

func main() {
	// Create a pool with 100 workers and a custom SpinLock
	// 创建一个包含 100 个工作者的池，并使用自定义的自旋锁
	pool := lifopool.New(100, lifopool.WithLock(new(spinlock.SpinLock)))
	defer pool.Release()

	for i := 0; i < 1000; i++ {
		pool.AddTask(func() (any, error) {
			time.Sleep(10 * time.Millisecond)
			fmt.Println("Task executed with spinlock")
			return nil, nil
		})
	}
	pool.Wait()
}
