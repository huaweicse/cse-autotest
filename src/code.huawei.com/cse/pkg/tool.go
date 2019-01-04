package pkg

import (
	"errors"
	"sync"
)

type Counter struct {
	sync.Mutex
	limit int
	n     int
}

func (c *Counter) Increase() error {
	c.Lock()
	defer c.Unlock()
	c.n++
	if c.n == c.limit {
		c.n = 0
		return errors.New("reach the limit")
	}
	return nil
}

func (c *Counter) Value() int {
	c.Lock()
	defer c.Unlock()
	return c.n
}

func NewCounter(n int) *Counter {
	return &Counter{
		limit: n,
	}
}
