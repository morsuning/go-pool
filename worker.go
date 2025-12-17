package lifopool

import (
	"context"
	"fmt"
)

// worker 代表池中的一个工作者。
type worker struct {
	taskQueue chan task
}

func newWorker() *worker {
	return &worker{
		taskQueue: make(chan task, 1),
	}
}

// start 在单独的 goroutine 中启动工作者。
// 工作者将从其 taskQueue 运行任务，直到 taskQueue 关闭。
// 因为 taskQueue 的长度为 1，工作者在执行完 1 个任务后将被推回池中。
func (w *worker) start(pool *lifoPool, workerIndex int) {
	go func() {
		for t := range w.taskQueue {
			if t != nil {
				func() {
					defer func() {
						if r := recover(); r != nil {
							var err error
							if e, ok := r.(error); ok {
								err = e
							} else {
								err = fmt.Errorf("工作者发生 panic: %v", r)
							}
							w.handleResult(nil, err, pool)
						}
					}()
					result, err := w.executeTask(t, pool)
					w.handleResult(result, err, pool)
				}()
			}
			pool.pushWorker(workerIndex)
		}
	}()
}

// executeTask 执行任务并返回结果和错误。
// 如果任务失败，它将根据池的 retryCount 进行重试。
func (w *worker) executeTask(t task, pool *lifoPool) (result any, err error) {
	for i := 0; i <= pool.retryCount; i++ {
		if pool.timeout > 0 {
			result, err = w.executeTaskWithTimeout(t, pool)
		} else {
			result, err = w.executeTaskWithoutTimeout(t)
		}
		if err == nil || i == pool.retryCount {
			return result, err
		}
	}
	return
}

// executeTaskWithTimeout 执行带有超时的任务并返回结果和错误。
func (w *worker) executeTaskWithTimeout(t task, pool *lifoPool) (result any, err error) {
	// 创建一个带有超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), pool.timeout)
	defer cancel()
	// 创建一个通道来接收任务的结果
	resultChan := make(chan any)
	errChan := make(chan error)

	// 在单独的 goroutine 中运行任务
	go func() {
		res, err := t()
		select {
		case resultChan <- res:
		case errChan <- err:
		case <-ctx.Done():
			// context 被取消，停止任务
			return
		}
	}()

	// 等待任务完成或 context 超时
	select {
	case result = <-resultChan:
		err = <-errChan
		// 任务成功完成
		return result, err
	case <-ctx.Done():
		// context 超时，任务耗时太长
		return nil, fmt.Errorf("任务超时")
	}
}

func (w *worker) executeTaskWithoutTimeout(t task) (result any, err error) {
	// 如果未设置超时或为 0，则直接运行任务
	return t()
}

// handleResult 处理任务的结果。
func (w *worker) handleResult(result any, err error, pool *lifoPool) {
	if err != nil && pool.errorCallback != nil {
		pool.errorCallback(err)
	} else if pool.resultCallback != nil {
		pool.resultCallback(result)
	}
}
