# 变更日志

本项目的所有显著更改都将记录在此文件中。

## [v1.0.0] - 2025-12-17

### 重大变更
- **项目更名**：项目名称从 `go-pool` 变更为 `lifopool`，以更好地体现 LIFO 调度特性。
- **模块路径变更**：Go module 路径更改为 `github.com/morsuning/lifopool`。
- **API 变更**：
    - `NewGoPool` -> `New`
    - `GoPool` interface -> `Pool`
    - 包名 `gopool` -> `lifopool`

## [0.1.0]

### 优化
- **文档重构**：更新了 README.md（英文）和 README.md（中文），增加了详细的性能基准测试对比和特性说明。
- **项目结构**：增加了 `examples/` 目录提供示例代码，将基准测试移动到 `benchmarks/` 目录。
- **性能验证**：更新了与其他流行库（ants, pond）的最新性能对比数据。

### 修复
- **Worker 健壮性**：修复了一个严重问题，即任务内的 panic 会导致 worker goroutine 崩溃并永久丢失，可能导致池死锁。添加了 `recover()` 机制来捕获 panic，将其作为错误报告，并确保 worker 返回到池中。
