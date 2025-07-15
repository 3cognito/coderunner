package cache

import "errors"

type ExecutionOutput struct {
	Stdout string
	Stderr string
	Err    string
}

type CacheInterface interface {
	Get(key string) (ExecutionOutput, error)
	Set(key string, output ExecutionOutput)
}

type LRUCache struct {
	capacity int
	list     *DoublyLinkedList
	entries  map[string]*Node
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

func (c *LRUCache) Get(key string) (ExecutionOutput, error) {
	if node, exists := c.entries[key]; exists {
		c.list.MovetoFront(node)
		return node.Value, nil
	}
	return ExecutionOutput{}, errors.New("key not found")
}

func (c *LRUCache) Set(key string, value ExecutionOutput) {
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
