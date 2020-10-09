package cache

import (
	"log"
	"sync"
	"time"
)

type cache struct {
	m           sync.Map
	duration    time.Duration
	keyEviction chan string
}

func NewCache(t time.Duration) Cacher {
	keyEviction := make(chan (string), 10)
	c := &cache{
		keyEviction: keyEviction,
		duration:    t,
		m:           sync.Map{},
	}
	go c.handleEvictions()

	return c
}

func (c *cache) handleEvictions() {
	for key := range c.keyEviction {
		log.Println("deleting", key)
		c.m.Delete(key)
	}
}

func (c *cache) Get(key string) ([]byte, bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return nil, false
	}
	return v.([]byte), true
}

func (c *cache) Insert(key string, value []byte) {
	if _, found := c.m.Load(key); found {
		return
	}
	go func() {
		time.Sleep(c.duration)
		c.keyEviction <- key
	}()
	c.m.Store(key, value)
}

type Cacher interface {
	// Insert adds a value to the cache
	Insert(keyName string, keyValue []byte)
	// Get returns the element from the cache
	Get(keyName string) (keyValue []byte, found bool)
}
