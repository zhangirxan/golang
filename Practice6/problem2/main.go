package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

//Version 1 Fix using sync.Mutex
func withMutex() {
	var counter int
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("Mutex result:", counter)
}

//Version 2 Fix using atomic
func withAtomic() {
	var counter atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}

	wg.Wait()
	fmt.Println("Atomic result:", counter.Load())
}

func main() {
	withMutex()
	withAtomic()
}
