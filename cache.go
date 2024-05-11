package cache

import "fmt"

type item map[string]any

type InterfaceCache interface {
	Set(key string, value any)
	Get(key string) (any, bool)
	Delete(key string) string
}

type Cache struct {
	items item
}

var storage = Cache{
	items: make(item),
}

func New() *Cache {
	storage.items = make(map[string]any)
	return &storage
}

func (c *Cache) Set(key string, value any) {
	storage.items[key] = value
	fmt.Println(key, " added successfully")
}

func (c *Cache) Get(key string) (any, bool) {
	value, found := storage.items[key]
	return value, found
}

func (c *Cache) Delete(key string) string {
	result := "Cache not found!"
	for value, _ := range storage.items {
		if value == key {
			result = "Deleted successful"
		}
	}

	delete(storage.items, key)
	return result
}
