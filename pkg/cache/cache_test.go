package cache

import (
	"bytes"
	"log"
	"testing"
	"time"
)

func TestInsertCacheHit(t *testing.T) {
	c := NewCache(10 * time.Second)
	key, value := "somekey", []byte("somevalue")

	c.Insert(key, value)

	cachedValue, found := c.Get(key)
	if !found {
		t.Fatal("value not found in cache")
	}
	if !bytes.Equal(value, cachedValue) {
		t.Fatalf("cached value was different from original")
	}
}

func TestInsertCacheMissing(t *testing.T) {
	c := NewCache(time.Millisecond)
	key, value := "somekey", []byte("somevalue")

	c.Insert("otherkey", value)
	time.Sleep(100 * time.Millisecond)

	cachedValue, found := c.Get(key)
	if found {
		t.Fatal("value should be missing")
	}
	if cachedValue != nil {
		t.Fatal("value from cache should be nil")
	}
}

func TestInsertCacheMiss(t *testing.T) {
	c := NewCache(time.Millisecond)
	key, value := "somekey", []byte("somevalue")

	c.Insert(key, value)
	time.Sleep(100 * time.Millisecond)

	cachedValue, found := c.Get(key)
	if found {
		t.Fatal("value should be missing")
	}
	if cachedValue != nil {
		t.Fatal("value from cache should be nil")
	}
}

func TestInsertOverwrite(t *testing.T) {
	c := NewCache(time.Millisecond)
	key, value := "somekey", []byte("somevalue")

	c.Insert(key, value)
	c.Insert(key, value)

	internalCache := c.(*cache)
	var seen []byte
	internalCache.m.Range(func(k, v interface{}) bool {
		seen = v.([]byte)
		return true
	})
	if !bytes.Equal(seen, value) {
		log.Fatalf("stored value different than inserted %s != %s", string(seen), string(value))
	}

}
