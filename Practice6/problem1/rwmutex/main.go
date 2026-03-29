package main

import (
	"fmt"
	"sync"
)

type SafeMap struct {
	mu sync.RWMutex
	m  map[string]int
}

func (s *SafeMap) Store(key string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func (s *SafeMap) Load(key string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.m[key]
	return val, ok
}

func main() {
	safeMap := &SafeMap{m: make(map[string]int)}
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			safeMap.Store("key", key)
		}(i)
	}

	wg.Wait()

	value, ok := safeMap.Load("key")
	if ok {
		fmt.Printf("Value: %d\n", value)
	}
}
