package Pratice

import (
	"fmt"
	"sync"
)

var channel chan int
var myMap map[string]chan bool
var once sync.Once

func CallTestFunc() {
	ch := make(chan int)

	go ttt(ch)

	value := <-ch

	fmt.Println("in here", value)
}

func Dododo() {
	once.Do(gggg)
	once.Do(gggg)
	once.Do(gggg)
}
func gggg() {
	fmt.Println("show show show")
}

func ttt(ch chan int) {
	ch <- 1
}

func CallTestSelectFunc() {
	ch := make(chan int, 1)

	for {
		select {
		case ch <- 0:
		case ch <- 1:
		}

		i := <-ch
		fmt.Println("value received:", i)
	}
}
