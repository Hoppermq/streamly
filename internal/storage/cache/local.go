package cache

import (
	"sync"
	"time"
)

const defaultMaxSize = 100

type Item[T comparable] struct {
	value   T
	expires time.Time
}

type LocalStorage[T comparable, K ~string] struct {
	items map[K]*Item[T]
	mu    sync.RWMutex
}

func NewLocalStorage[T comparable, K ~string]() *LocalStorage[T, K] {
	return &LocalStorage[T, K]{
		items: make(map[K]*Item[T], defaultMaxSize),
	}
}

func (l *LocalStorage[T, K]) Get(key string) (T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	item, ok := l.items[K(key)]
	if !ok {
		return item.value, false
	}

	if item.expires.Before(time.Now()) {
		return item.value, false
	}

	return item.value, true
}

func (l *LocalStorage[T, K]) Set(key string, value T, ttl time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	item := &Item[T]{
		value:   value,
		expires: time.Now().Add(ttl),
	}

	l.items[K(key)] = item
}

func (l *LocalStorage[T, K]) Delete(key K) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.items, key)
}

func (l *LocalStorage[T, K]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.items = make(map[K]*Item[T], defaultMaxSize)
}
