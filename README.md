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
  <strong>High-Performance, Feature-Rich, and Production-Ready Goroutine Pool for Go.</strong>
</p>

<p align="center">
  <a href="README.md">English</a> | <a href="README_CN.md">中文文档</a>
</p>

---

## 🚀 Introduction

`lifopool` is a library designed to manage specific goroutine workers to handle tasks concurrently. It reuses goroutines to limit concurrency, reduce resource consumption, and improve the stability and performance of your applications.

Key differentiators:
*   **LIFO Scheduling**: Improves CPU cache locality by prioritizing recently used workers.
*   **Dynamic Resizing**: Automatically scales workers up or down based on traffic load.
*   **Robustness**: Built-in panic recovery, timeout handling, and retry mechanisms.

## 📊 Benchmarks

`lifopool` is optimized for speed and efficiency. Below are benchmark results comparing it with other popular libraries (`ants`, `pond`) and raw goroutines.

| Library | Optimization | Time (ns/op) | Memory (B/op) | Allocations (allocs/op) |
| :--- | :--- | :--- | :--- | :--- |
| **lifopool** | Default | **1,114,703,667** | **2,757,944** | **14,906** |
| [ants](https://github.com/panjf2000/ants) | - | 1,141,786,333 | 4,533,200 | 59,463 |
| [pond](https://github.com/alitto/pond) | - | 1,479,714,792 | 1,035,432 | 10,788 |
| *Raw Goroutines* | *None* | *336,763,680* | *128,759,930* | *3,007,120* |

**Highlights**:
- **Faster** than `ants` and `pond` in handling 1 million tasks.
- **50% Less Memory** usage compared to `ants`.
- **Significant** reduction in allocations compared to raw goroutines and `ants`.

*Benchmarks run on Apple M3, processing 1M tasks.*

## ✨ Features

- **🚀 High Performance**: Low overhead, LIFO stack-based worker management.
- **⚖️ Dynamic Scaling**: Auto-scale worker count based on queue monitoring.
- **🛡️ Panic Recovery**: Automatically recovers from worker panics to prevent leaks.
- **⏱️ Timeout Support**: Task execution timeouts with `context` integration.
- **🔄 Retry Mechanism**: Configurable retries for failed tasks.
- **🔒 Custom Lock**: Standard `sync.Mutex` or high-performance `SpinLock` support.

## 📦 Installation

```bash
go get -u github.com/morsuning/lifopool
```

## ⚡ Quick Start

### Simple Usage

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

See [examples/](examples/) for more usage patterns, including Custom Locks and Timeout handling.

## ⚙️ Configuration

`lifopool.New` accepts functional options to tailor behavior:

| Option | Description | Default |
| :--- | :--- | :--- |
| `WithMinWorkers(n)` | Minimum idle workers to keep. | `maxWorkers` |
| `WithTaskQueueSize(n)` | Size of the buffered task channel. | `1e6` |
| `WithTimeout(d)` | Execution timeout for tasks. | `0` (None) |
| `WithRetryCount(n)` | Retries upon failure. | `0` |
| `WithErrorCallback(fn)` | Callback for task errors/panics. | `nil` |
| `WithLock(l)` | Custom lock implementation. | `sync.Mutex` |

## 🤝 Contributing

Contributions are welcome! Please check out the [issues](https://github.com/morsuning/go-pool/issues) or submit a PR.

## 📄 License

MIT © [morsuning](https://github.com/morsuning)
