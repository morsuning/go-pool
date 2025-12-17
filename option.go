package gopool

import (
	"sync"
	"time"
)

// Option 代表池的一个选项。
type Option func(*goPool)

// WithLock 设置池的锁。
func WithLock(lock sync.Locker) Option {
	return func(p *goPool) {
		p.lock = lock
		p.cond = sync.NewCond(p.lock)
	}
}

// WithMinWorkers 设置池的最小工作者数量。
func WithMinWorkers(minWorkers int) Option {
	return func(p *goPool) {
		p.minWorkers = minWorkers
	}
}

// WithTimeout 设置池的超时时间。
func WithTimeout(timeout time.Duration) Option {
	return func(p *goPool) {
		p.timeout = timeout
	}
}

// WithResultCallback 设置池的结果回调。
func WithResultCallback(callback func(any)) Option {
	return func(p *goPool) {
		p.resultCallback = callback
	}
}

// WithErrorCallback 设置池的错误回调。
func WithErrorCallback(callback func(error)) Option {
	return func(p *goPool) {
		p.errorCallback = callback
	}
}

// WithRetryCount 设置池的重试次数。
func WithRetryCount(retryCount int) Option {
	return func(p *goPool) {
		p.retryCount = retryCount
	}
}

// WithTaskQueueSize 设置池的任务队列大小。
func WithTaskQueueSize(size int) Option {
	return func(p *goPool) {
		p.taskQueueSize = size
	}
}
