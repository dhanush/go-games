// pingpong
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func hit(ch chan int) {
	i := <-ch
	if i%7 == 0 {
		fmt.Println("Ping lost")
		return
	}
	i = (i*3 + 2) / 5
	fmt.Printf("Ping Hit %d \n", i)

	go receive(ch)
	ch <- i
	time.Sleep(100 * time.Millisecond)
}

func receive(ch chan int) {
	i := <-ch
	if i%7 == 0 {
		fmt.Println("Pong lost")
		return
	}
	fmt.Printf("Pong Received %d \n", i)
	go hit(ch)
	time.Sleep(100 * time.Millisecond)
	ch <- i + 1

}

func main() {
	rand.Seed(20)
	r := rand.Intn(100)
	var ch = make(chan int)
	fmt.Printf("Serving with %d \n", r)
	go hit(ch)
	ch <- r
	time.Sleep(2 * time.Minute)

}
