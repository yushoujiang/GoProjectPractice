package Pratice

import (
	"fmt"
	"runtime"
	"sync"
)

var counter int = 0

func CountNum(lock *sync.Mutex, index int) {

	if lock == nil {
		return
	}

	defer func() {
		lock.Unlock()
	}()

	lock.Lock()
	counter++
	fmt.Println(index, counter)
	// lock.Unlock()

}
func StartCount() {

	lock := &sync.Mutex{}

	for i := 0; i < 10; i++ {
		go CountNum(lock, i)
	}

	for {
		lock.Lock()

		c := counter

		lock.Unlock()

		runtime.Gosched()

		if c >= 10 {
			break
		}

	}

}
