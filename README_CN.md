# LifoPool

<p align="center">
  <img src="https://socialify.git.ci/morsuning/lifopool/image?description=1&font=Inter&language=1&name=1&owner=1&pattern=Circuit%20Board&theme=Auto" alt="go-pool" width="640" height="320" />
</p>

<p align="center">
    <a href="https://pkg.go.dev/github.com/morsuning/lifopool"><img src="https://pkg.go.dev/badge/github.com/morsuning/lifopool.svg" alt="GoDoc"></a>
    <a href="https://goreportcard.com/report/github.com/morsuning/lifopool"><img src="https://goreportcard.com/badge/github.com/morsuning/lifopool" alt="Go Report Card"></a>
    <a href="LICENSE"><img src="https://img.shields.io/github/license/morsuning/lifopool" alt="License"></a>
</p>

<p align="center">
  <strong>Go 语言实现的高性能、功能丰富且生产就绪的 Goroutine 池。</strong>
</p>

<p align="center">
  <a href="README.md">English</a> | <a href="README_CN.md">中文文档</a>
</p>

---

## 🚀 简介

`lifopool` 是一个用于管理一组特定 goroutine 工作者来并发处理任务的库。它通过复用 goroutine 来限制并发数，降低资源消耗，并提高应用程序的稳定性和性能。

核心优势：
*   **LIFO 调度**：通过优先使用最近活跃的工作者来提高 CPU 缓存亲和性。
*   **动态扩缩容**：根据流量负载自动增加或减少工作者数量。
*   **健壮性**：内置 Panic 恢复、超时处理和重试机制。

## 📊 性能测试

`lifopool` 专为速度和效率进行了优化。以下是与热门库（`ants`, `pond`）及原生 goroutine 的对比基准测试结果。

| 库 | 优化项 | 耗时 (ns/op) | 内存 (B/op) | 分配次数 (allocs/op) |
| :--- | :--- | :--- | :--- | :--- |
| **lifopool** | 默认 | **1,114,703,667** | **2,757,944** | **14,906** |
| [ants](https://github.com/panjf2000/ants) | - | 1,141,786,333 | 4,533,200 | 59,463 |
| [pond](https://github.com/alitto/pond) | - | 1,479,714,792 | 1,035,432 | 10,788 |
| *原生 Goroutine* | *无* | *336,763,680* | *128,759,930* | *3,007,120* |

**亮点**：
- 在处理 100 万个任务时，速度比 `ants` 和 `pond` **更快**。
- 相比 `ants`，内存使用量减少了 **50%**。
- 相比原生 goroutine 和 `ants`，内存分配次数显著减少。

*测试环境：Apple M3，处理 100 万任务。*

## ✨ 特性

- **🚀 高性能**：低开销，基于栈的 LIFO 工作者管理。
- **⚖️ 动态伸缩**：基于队列监控自动扩展工作者数量。
- **🛡️ Panic 恢复**：自动捕获工作者 Panic，防止协程泄漏。
- **⏱️ 超时支持**：集成 `context` 的任务执行超时控制。
- **🔄 重试机制**：支持失败任务的可配置重试。
- **🔒 自定义锁**：支持标准 `sync.Mutex` 或高性能 `SpinLock`。

## 📦 安装

```bash
go get -u github.com/morsuning/lifopool
```

## ⚡ 快速开始

### 简单示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/morsuning/lifopool"
)

func main() {
    pool := lifopool.New(100)
    defer pool.Release()

    pool.AddTask(func() (any, error) {
        time.Sleep(10 * time.Millisecond)
        fmt.Println("Hello, lifopool!")
        return nil, nil
    })
    
    pool.Wait()
}
```

查看 [examples/](examples/) 目录获取更多用法，包括自定义锁和超时处理。

## ⚙️ 配置

`lifopool.New` 接受函数式选项来定制行为：

| 选项 | 说明 | 默认值 |
| :--- | :--- | :--- |
| `WithMinWorkers(n)` | 保持的最小空闲工作者数量。 | `maxWorkers` |
| `WithTaskQueueSize(n)` | 任务缓冲通道的大小。 | `1e6` |
| `WithTimeout(d)` | 任务执行超时时间。 | `0` (无超时) |
| `WithRetryCount(n)` | 失败重试次数。 | `0` |
| `WithErrorCallback(fn)` | 任务错误/Panic 回调函数。 | `nil` |
| `WithLock(l)` | 自定义锁实现。 | `sync.Mutex` |

## 🤝 贡献代码

欢迎提交 Pull Request 或 [Issues](https://github.com/morsuning/lifopool/issues)！

## 📄 许可证

MIT © [morsuning](https://github.com/morsuning)
