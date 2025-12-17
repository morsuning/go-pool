package lifopool

import (
	"context"
	"sort"
	"sync"
	"time"
)

type Pool interface {
	// AddTask 向池中添加一个任务。
	AddTask(t task)
	// Wait 等待所有任务被分发并完成。
	Wait()
	// Release 释放池及其所有工作者。
	Release()
	// GetRunning 返回正在运行的工作者数量。
	Running() int
	// GetWorkerCount 返回工作者的总数量。
	GetWorkerCount() int
	// GetTaskQueueSize 返回任务队列的大小。
	GetTaskQueueSize() int
}

// task 代表将被工作者执行的函数。
// 它返回一个结果和一个错误。
type task func() (any, error)

// goPool 代表一个工作者池。
type lifoPool struct {
	workers     []*worker
	workerStack []int
	maxWorkers  int
	// 由 WithMinWorkers() 设置，用于调整工作者数量。默认等于 maxWorkers。
	minWorkers int
	// 任务首先被添加到这个通道，然后分发给工作者。默认缓冲区大小为 100 万。
	taskQueue chan task
	// 由 WithTaskQueueSize() 设置，用于设置任务队列的大小。默认为 1e6。
	taskQueueSize int
	// 由 WithRetryCount() 设置，用于在任务失败时重试。默认为 0。
	retryCount int
	lock       sync.Locker
	cond       *sync.Cond
	// 由 WithTimeout() 设置，用于设置任务的超时时间。默认为 0，意味着没有超时。
	timeout time.Duration
	// 由 WithResultCallback() 设置，用于处理任务的结果。默认为 nil。
	resultCallback func(any)
	// 由 WithErrorCallback() 设置，用于处理任务的错误。默认为 nil。
	errorCallback func(error)
	// adjustInterval 是调整工作者数量的时间间隔。默认为 1 秒。
	adjustInterval time.Duration
	ctx            context.Context
	// cancel 用于取消上下文。当调用 Release() 时会被调用。
	cancel context.CancelFunc
}

// NewGoPool 创建一个新的工作者池。
func New(maxWorkers int, opts ...Option) Pool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &lifoPool{
		maxWorkers: maxWorkers,
		// 默认将 minWorkers 设置为 maxWorkers
		minWorkers: maxWorkers,
		// workers 和 workerStack 应该在调用 WithMinWorkers() 之后初始化
		workers:        nil,
		workerStack:    nil,
		taskQueue:      nil,
		taskQueueSize:  1e6,
		retryCount:     0,
		lock:           new(sync.Mutex),
		timeout:        0,
		adjustInterval: 1 * time.Second,
		ctx:            ctx,
		cancel:         cancel,
	}
	// 应用选项
	for _, opt := range opts {
		opt(pool)
	}

	pool.taskQueue = make(chan task, pool.taskQueueSize)
	pool.workers = make([]*worker, pool.minWorkers)
	pool.workerStack = make([]int, pool.minWorkers)

	if pool.cond == nil {
		pool.cond = sync.NewCond(pool.lock)
	}
	// 使用最小数量创建工作者。这里不要使用 pushWorker()。
	for i := 0; i < pool.minWorkers; i++ {
		worker := newWorker()
		pool.workers[i] = worker
		pool.workerStack[i] = i
		worker.start(pool, i)
	}
	go pool.adjustWorkers()
	go pool.dispatch()
	return pool
}

// AddTask 向池中添加一个任务。
func (p *lifoPool) AddTask(t task) {
	p.taskQueue <- t
}

// Wait 等待所有任务被分发并完成。
func (p *lifoPool) Wait() {
	for {
		p.lock.Lock()
		workerStackLen := len(p.workerStack)
		p.lock.Unlock()

		if len(p.taskQueue) == 0 && workerStackLen == len(p.workers) {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}
}

// Release 停止所有工作者并释放资源。
func (p *lifoPool) Release() {
	close(p.taskQueue)
	p.cancel()
	p.cond.L.Lock()
	for len(p.workerStack) != p.minWorkers {
		p.cond.Wait()
	}
	p.cond.L.Unlock()
	for _, worker := range p.workers {
		close(worker.taskQueue)
	}
	p.workers = nil
	p.workerStack = nil
}

func (p *lifoPool) popWorker() int {
	p.lock.Lock()
	workerIndex := p.workerStack[len(p.workerStack)-1]
	p.workerStack = p.workerStack[:len(p.workerStack)-1]
	p.lock.Unlock()
	return workerIndex
}

func (p *lifoPool) pushWorker(workerIndex int) {
	p.lock.Lock()
	p.workerStack = append(p.workerStack, workerIndex)
	p.lock.Unlock()
	p.cond.Signal()
}

// adjustWorkers 根据队列中的任务数量调整工作者数量。
func (p *lifoPool) adjustWorkers() {
	ticker := time.NewTicker(p.adjustInterval)
	defer ticker.Stop()

	var adjustFlag bool

	for {
		adjustFlag = false
		select {
		case <-ticker.C:
			p.cond.L.Lock()
			if len(p.taskQueue) > len(p.workers)*3/4 && len(p.workers) < p.maxWorkers {
				adjustFlag = true
				// 将工作者数量翻倍，直到达到最大值
				newWorkers := min(len(p.workers)*2, p.maxWorkers) - len(p.workers)
				for i := 0; i < newWorkers; i++ {
					worker := newWorker()
					p.workers = append(p.workers, worker)
					// 这里不要使用 len(p.workerStack)-1，因为当池繁忙时它会小于 len(p.workers)-1
					p.workerStack = append(p.workerStack, len(p.workers)-1)
					worker.start(p, len(p.workers)-1)
				}
			} else if len(p.taskQueue) == 0 && len(p.workerStack) == len(p.workers) && len(p.workers) > p.minWorkers {
				adjustFlag = true
				// 将工作者数量减半，直到达到最小值
				removeWorkers := (len(p.workers) - p.minWorkers + 1) / 2
				// 在移除工作者之前对 workerStack 进行排序。
				// [1,2,3,4,5] -工作-> [1,2,3] -扩容-> [1,2,3,6,7] -空闲-> [1,2,3,6,7,4,5]
				sort.Ints(p.workerStack)
				p.workers = p.workers[:len(p.workers)-removeWorkers]
				p.workerStack = p.workerStack[:len(p.workerStack)-removeWorkers]
			}
			p.cond.L.Unlock()
			if adjustFlag {
				p.cond.Broadcast()
			}
		case <-p.ctx.Done():
			return
		}
	}
}

// dispatch 将任务分发给工作者。
func (p *lifoPool) dispatch() {
	for t := range p.taskQueue {
		p.cond.L.Lock()
		for len(p.workerStack) == 0 {
			p.cond.Wait()
		}
		p.cond.L.Unlock()
		workerIndex := p.popWorker()
		p.workers[workerIndex].taskQueue <- t
	}
}

// Running 返回当前正在工作的工作者数量。
func (p *lifoPool) Running() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return len(p.workers) - len(p.workerStack)
}

// GetWorkerCount 返回池中的工作者数量。
func (p *lifoPool) GetWorkerCount() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return len(p.workers)
}

// GetTaskQueueSize 返回任务队列的大小。
func (p *lifoPool) GetTaskQueueSize() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.taskQueueSize
}
