package cache

import (
	"fmt"
	"sync"
	"time"
)

type item struct {
	value     any
	timestamp time.Time
}

type InterfaceCache interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Delete(key string) string
}

type Cache struct {
	items      map[string]item
	mu         sync.Mutex
	wg         sync.WaitGroup
	workerPool *WorkerPool
}

type Task struct {
	key   string
	value any
}

type WorkerPool struct {
	tasks    chan Task
	workerWg sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	wp := &WorkerPool{
		tasks: make(chan Task, numWorkers),
	}
	for i := 0; i < numWorkers; i++ {
		wp.workerWg.Add(1)
	}
	return wp
}

func (wp *WorkerPool) AddTask(task Task) {
	wp.tasks <- task
}

func (wp *WorkerPool) Stop() {
	close(wp.tasks)
	wp.workerWg.Wait()
}

func New() *Cache {
	const numWorkers int = 5
	return &Cache{
		items:      make(map[string]item),
		workerPool: NewWorkerPool(numWorkers),
	}
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	c.items[key] = item{
		value:     value,
		timestamp: time.Now(),
	}
	c.mu.Unlock()

	c.workerPool.AddTask(Task{key: key, value: value})
	go func() {
		select {
		case <-time.After(5 * time.Second):
			c.mu.Lock()
			if itm, found := c.items[key]; found {
				if time.Since(itm.timestamp) >= 5*time.Second {
					delete(c.items, key)
					fmt.Println(key, "ttl is end")
				}
			}
			c.mu.Unlock()
		}
	}()
	fmt.Println(key, "added successfully")
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	itm, found := c.items[key]
	if found {
		return itm.value, true
	}
	return nil, false
}

func (c *Cache) Delete(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, found := c.items[key]; found {
		delete(c.items, key)
		return "Deleted successfully"
	}
	return "Cache not found!"
}

func (c *Cache) Stop() {
	c.workerPool.Stop()
}
