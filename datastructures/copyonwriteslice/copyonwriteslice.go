package copyonwriteslice

import (
	"sync"
)

type CopyOnWriteSlice[T comparable] struct {
	data []T
	mu   sync.RWMutex
}

func NewCopyOnWriteSlice[T comparable]() *CopyOnWriteSlice[T] {
	return &CopyOnWriteSlice[T]{}
}

func (c *CopyOnWriteSlice[T]) Add(value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	length := len(c.data)
	new := make([]T, length+1)
	copy(new, c.data)
	new[length] = value
	c.data = new
}

func (c *CopyOnWriteSlice[T]) List() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	copied := make([]T, len(c.data))
	copy(copied, c.data)
	return copied
}
