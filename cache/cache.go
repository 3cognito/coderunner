package cache

import (
	"3cognito/coderunner/types"
	"errors"
	"sync"
)

type CacheInterface interface {
	Get(key string) (types.ExecutionOutput, error)
	Set(key string, output types.ExecutionOutput)
}

type LRUCache struct {
	capacity int
	list     *DoublyLinkedList
	entries  map[string]*Node
	mu       sync.Mutex
}

func NewLRUCache(cap int) CacheInterface {
	if cap <= 0 {
		panic("cache capacity must be greater than zero")
	}
	return &LRUCache{
		capacity: cap,
		list:     &DoublyLinkedList{},
		entries:  make(map[string]*Node),
	}
}

func (c *LRUCache) Get(key string) (types.ExecutionOutput, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.entries[key]
	if exists {
		c.list.MovetoFront(node)
		return node.Value, nil
	}

	return types.ExecutionOutput{}, errors.New("key not found")
}

func (c *LRUCache) Set(key string, value types.ExecutionOutput) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.entries[key]; exists {
		node.Value = value
		c.list.MovetoFront(node)
		return
	}
	//I decided to panic cache initialization if set to cap zero (because what's the use of a zero capacity cache???)
	//however if the cache initialization behaviour is not to panic, the commented out code below is needed to prevent setting an entry
	// if c.capacity == 0 {
	// 	return
	// }

	if len(c.entries) >= c.capacity {
		node := c.list.RemoveTail()
		if node != nil {
			delete(c.entries, node.Key)
		}
	}

	node := &Node{
		Key:   key,
		Value: value,
	}

	c.entries[key] = node
	c.list.AddtoFront(node)
}
