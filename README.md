# go-pool

go-pool 是一个用 Golang 实现的**高性能**、**功能丰富**、**简单易用**的 goroutine 池库。它会管理和回收一组 goroutine 来并发完成任务，从而提高你的应用程序的效率和性能。

## 性能测试

这个表格展示了三个 Go 库 go-pool、[ants](https://github.com/panjf2000/ants) 和 [pond](https://github.com/alitto/pond)的性能测试结果。表格包括每个库处理 100 万个任务所需的时间和内存消耗（以 MB 为单位）。

| 项目                                   | 处理一百万任务耗时 (s) | 内存消耗 (MB) |
| -------------------------------------- | :--------------------: | :-----------: |
| go-pool                                |          1.13          |     2.11     |
| [ants](https://github.com/panjf2000/ants) |          1.43          |     8.94     |
| [pond](https://github.com/alitto/pond)    |          3.32          |     2.20     |

## 特性

- **任务队列**：go-pool 使用一个线程安全的任务队列来存储等待处理的任务。多个工作器可以同时从这个队列中获取任务。任务队列的大小可配置。

- **并发控制**：go-pool 可以控制并发任务的数量，防止系统过载。

- **动态工作器调整**：go-pool 可以根据任务数量和系统负载动态调整工作器的数量。

- **优雅关闭**：go-pool 可以优雅地关闭。当没有更多的任务或收到关闭信号时，它会停止接受新的任务，并等待所有进行中的任务完成后再关闭。

- **任务错误处理**：go-pool 可以处理任务执行过程中出现的错误。

- **任务超时处理**：go-pool 可以处理任务执行超时。如果一个任务在指定的超时期限内没有完成，该任务被认为失败，返回一个超时错误。

- **任务结果获取**：go-pool 提供了一种获取任务结果的方式。

- **任务重试**：go-pool 为失败的任务提供了重试机制。

- **锁定制**：go-pool 支持不同类型的锁。你可以使用内置的 `sync.Mutex`或自定义锁，如 `spinlock.SpinLock`。

## 安装

要安装go-pool，使用 `go get`：

```bash
go get -u github.com/morsuning/go-pool
```

## 使用

这是一个如何使用带有 `sync.Mutex` 的go-pool 的简单示例：

```go
package main

import (
    "sync"
    "time"

    "github.com/morsuning/go-pool"
)

func main() {
    pool := go-pool.NewGoPool(100)
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error){
            time.Sleep(10 * time.Millisecond)
            return nil, nil
        })
    }
    pool.Wait()
}
```

这是如何使用带有 `spinlock.SpinLock` 的 go-pool 的示例：

```go
package main

import (
    "time"

    "github.com/daniel-hutao/spinlock"
    "github.com/morsuning/go-pool"
)

func main() {
    pool := go-pool.NewGoPool(100, gopool.WithLock(new(spinlock.SpinLock)))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error){
            time.Sleep(10 * time.Millisecond)
            return nil, nil
        })
    }
    pool.Wait()
}
```

## 配置任务队列大小

go-pool 使用一个线程安全的任务队列来存储等待处理的任务。多个工作器可以同时从这个队列中获取任务。任务队列的大小可配置。可以通过在创建池时设置 `WithQueueSize` 选项来配置任务队列的大小。

这是一个如何配置 go-pool 任务队列大小的示例：

```go
package main

import (
    "time"

    "github.com/morsuning/go-pool"
)

func main() {
    pool := go-pool.NewGoPool(100, gopool.WithTaskQueueSize(5000))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error){
            time.Sleep(10 * time.Millisecond)
            return nil, nil
        })
    }
    pool.Wait()
}
```

## 动态工作器调整

go-pool 支持动态工作器调整。这意味着池中的工作器数量可以根据队列中的任务数量增加或减少。可以通过在创建池时设置 MinWorkers 选项来启用此功能。

这是如何使用动态工作器调整的 go-pool 的示例：

```go
package main

import (
    "time"

    "github.com/morsuning/go-pool"
)

func main() {
    pool := go-pool.NewGoPool(100, gopool.WithMinWorkers(50))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error){
            time.Sleep(10 * time.Millisecond)
            return nil, nil
        })
    }
    pool.Wait()
}
```

在这个示例中，池开始时有50个工作器。如果队列中的任务数量超过当前工作器数量的3/4，并且当前工作器数量小于 MaxWorkers，池将翻倍工作器数量，直到达到 MaxWorkers。如果队列中的任务数量为零，并且当前工作器数量大于 MinWorkers，池将把工作器数量减半，直到达到 MinWorkers。

## 任务超时处理

go-pool支持任务超时。如果一个任务花费的时间超过指定的超时时间，它将被取消。可以通过在创建池时设置 `WithTimeout` 选项来启用此功能。

这是如何使用任务超时的 go-pool 的示例：

```go
package main

import (
    "time"

    "github.com/morsuning/go-pool"
)

func main() {
    pool := go-pool.NewGoPool(100, gopool.WithTimeout(1*time.Second))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error) {
            time.Sleep(2 * time.Second)
            return nil, nil
        })
    }
    pool.Wait()
}
```

在这个示例中，如果任务花费的时间超过1秒，任务将被取消。

## 任务错误处理

go-pool 支持任务错误处理。如果一个任务返回一个错误，错误回调函数将被调用。可以通过在创建池时设置 `WithErrorCallback` 选项来启用此功能。

这是如何使用错误处理的 go-pool 的示例：

```go
package main

import (
    "errors"
    "fmt"

    "github.com/morsuning/go-pool"
)

func main() {
    pool := go-pool.NewGoPool(100, gopool.WithErrorCallback(func(err error) {
        fmt.Println("Task error:", err)
    }))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error) {
            return nil, errors.New("task error")
        })
    }
    pool.Wait()
}
```

在这个示例中，如果一个任务返回一个错误，错误将被打印到控制台。

## 任务结果获取

go-pool 支持任务结果获取。如果一个任务返回一个结果，结果回调函数将被调用。可以通过在创建池时设置 `WithResultCallback` 选项来启用此功能。

这是如何使用任务结果获取的 go-pool 的示例：

```go
package main

import (
    "fmt"

    "github.com/morsuning/go-pool"
)

func main() {
    pool := go-pool.NewGoPool(100, gopool.WithResultCallback(func(result interface{}) {
        fmt.Println("Task result:", result)
    }))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error) {
            return "task result", nil
        })
    }
    pool.Wait()
}
```

在这个示例中，如果一个任务返回一个结果，结果将被打印到控制台。

## 任务重试

go-pool 支持任务重试。如果任务失败，可以重试指定的次数。可以通过在创建池时设置 `WithRetryCount` 选项来启用此功能。

以下是如何使用带有任务重试的 go-pool 的示例：

```go
package main

import (
    "errors"
    "fmt"

    "github.com/morsuning/go-pool"
)

func main() {
    pool := go-pool.NewGoPool(100, gopool.WithRetryCount(3))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error) {
            return nil, errors.New("task error")
        })
    }
    pool.Wait()
}
```

在这个示例中，如果任务失败，它将重试最多3次。
