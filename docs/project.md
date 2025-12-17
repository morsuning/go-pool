# Project Documentation: lifopool

`lifopool` 是一个高性能、线程安全的 Golang goroutine 池，旨在通过复用 goroutine 来限制并发数，优化资源使用并提高系统稳定性。

## 核心特性

*   **动态扩缩容**：根据流量负载自动增加或减少 Worker 数量。
*   **任务队列**：缓冲任务，削峰填谷。
*   **超时控制**：支持任务超时自动取消（Fallback）。
*   **异常恢复**：**[新增]** 自动捕获任务中的 Panic，防止 Worker 崩溃导致 Pool 不可用。
*   **任务重试**：支持失败任务自动重试。
*   **回调机制**：支持结果回调和错误回调。
*   **锁策略**：支持自定义锁（如 SpinLock）以优化性能。

## 快速开始

### 安装

```bash
go get -u github.com/morsuning/lifopool
```

### 基础使用

```go
package main

import (
    "fmt"
    "time"
    "github.com/morsuning/lifopool"
)

func main() {
    // 创建一个最大容量为 100 的池
    pool := lifopool.New(100)
    defer pool.Release()

    // 添加任务
    pool.AddTask(func() (any, error) {
        fmt.Println("Hello, World!")
        return nil, nil
    })

    pool.Wait()
}
```

## 配置选项

创建 Pool 时支持多种 Option 配置：

| 选项 | 说明 | 默认值 |
| :--- | :--- | :--- |
| `WithMinWorkers(n)` | 最小 Worker 数量，支持缩容到底线 | `maxWorkers` |
| `WithTaskQueueSize(n)` | 任务缓冲队列大小 | `1e6` |
| `WithTimeout(d)` | 任务执行超时时间 | 0 (无超时) |
| `WithRetryCount(n)` | 任务失败重试次数 | 0 |
| `WithErrorCallback(fn)` | 错误回调函数 | `nil` |
| `WithResultCallback(fn)` | 结果回调函数 | `nil` |

## 变更日志

请参考 [changelog.md](changelog.md)。
