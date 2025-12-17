package lifopool_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/daniel-hutao/spinlock"
	"github.com/morsuning/lifopool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LifoPool", func() {
	Describe("使用互斥锁", func() {
		It("应该正确工作", func() {
			pool := lifopool.New(100, lifopool.WithLock(new(sync.Mutex)))
			defer pool.Release()
			for i := 0; i < 1000; i++ {
				pool.AddTask(func() (any, error) {
					time.Sleep(10 * time.Millisecond)
					return nil, nil
				})
			}
			pool.Wait()
		})
	})

	Describe("使用自旋锁", func() {
		It("应该正确工作", func() {
			pool := lifopool.New(100, lifopool.WithLock(new(spinlock.SpinLock)))
			defer pool.Release()
			for i := 0; i < 1000; i++ {
				pool.AddTask(func() (any, error) {
					time.Sleep(10 * time.Millisecond)
					return nil, nil
				})
			}
			pool.Wait()
		})
	})

	Describe("处理错误", func() {
		It("应该正确工作", func() {
			var errTaskError = errors.New("task error")
			pool := lifopool.New(100, lifopool.WithErrorCallback(func(err error) {
				Expect(err).To(Equal(errTaskError))
			}))
			defer pool.Release()

			for i := 0; i < 1000; i++ {
				pool.AddTask(func() (any, error) {
					return nil, errTaskError
				})
			}
			pool.Wait()
		})
	})

	Describe("处理结果", func() {
		It("应该正确工作", func() {
			var expectedResult = "task result"
			pool := lifopool.New(100, lifopool.WithResultCallback(func(result any) {
				Expect(result).To(Equal(expectedResult))
			}))
			defer pool.Release()

			for i := 0; i < 1000; i++ {
				pool.AddTask(func() (any, error) {
					return expectedResult, nil
				})
			}
			pool.Wait()
		})
	})

	Describe("带有重试", func() {
		It("应该正确工作", func() {
			var retryCount = int32(3)
			var taskError = errors.New("task error")
			var taskRunCount int32 = 0

			pool := lifopool.New(100, lifopool.WithRetryCount(int(retryCount)))
			defer pool.Release()

			pool.AddTask(func() (any, error) {
				atomic.AddInt32(&taskRunCount, 1)
				if taskRunCount <= retryCount {
					return nil, taskError
				}
				return nil, nil
			})

			pool.Wait()

			Expect(atomic.LoadInt32(&taskRunCount)).To(Equal(retryCount + 1))
		})
	})

	Describe("带有超时", func() {
		It("应该正确工作", func() {
			var taskRun int32

			pool := lifopool.New(100, lifopool.WithTimeout(100*time.Millisecond), lifopool.WithErrorCallback(func(err error) {
				Expect(err.Error()).To(Equal("任务超时"))
				atomic.StoreInt32(&taskRun, 1)
			}))
			defer pool.Release()

			pool.AddTask(func() (any, error) {
				time.Sleep(200 * time.Millisecond)
				return nil, nil
			})

			pool.Wait()

			Expect(atomic.LoadInt32(&taskRun)).To(Equal(int32(1)))
		})
	})

	Describe("设置最小工作者数量", func() {
		It("应该正确工作", func() {
			var minWorkers = 50

			pool := lifopool.New(100, lifopool.WithMinWorkers(minWorkers))
			defer pool.Release()

			Expect(pool.GetWorkerCount()).To(Equal(minWorkers))
		})
	})

	Describe("设置任务队列大小", func() {
		It("应该正确工作", func() {
			size := 5000
			pool := lifopool.New(100, lifopool.WithTaskQueueSize(size))
			defer pool.Release()

			Expect(pool.GetTaskQueueSize()).To(Equal(size))
		})
	})
})
