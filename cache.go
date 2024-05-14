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
	items map[string]item
	mu    sync.Mutex
	wg    sync.WaitGroup
}

var (
	storage = Cache{
		items: make(map[string]item),
	}
	wg sync.WaitGroup
)

func New() *Cache {
	return &Cache{
		items: make(map[string]item),
	}
}

func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	c.items[key] = item{
		value:     value,
		timestamp: time.Now(),
	}
	c.mu.Unlock()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
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
