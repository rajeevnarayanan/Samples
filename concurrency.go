package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	str  string
	wait chan bool
}

func proc(msg string) <-chan Message {

	c := make(chan Message)
	wait := make(chan bool)
	go func() {

		for i := 0; ; i++ {
			c <- Message{fmt.Sprintf("%s: %d", msg, i), wait}
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
			bContinue := <-wait //waiting to receive a go ahead to reply

			if bContinue == false {
				fmt.Println("Condition ", bContinue)
				break
			}
		}
		fmt.Println(msg, "Quit")
	}()
	return c
}
func start(input1, input2 <-chan Message) <-chan Message {
	c := make(chan Message)

	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()
	return c
}

func main() {
	
	c := start(proc("First"), proc("Second"))

	var bContinue = true

	for i := 0; i < 5; i++ {
		msg1 := <-c
		fmt.Println(msg1.str)
		msg2 := <-c
		fmt.Println(msg2.str)
		if i == 4 {
			bContinue = false
			fmt.Println("Quit")
		}
		msg1.wait <- bContinue
		msg2.wait <- bContinue
	}
	time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond) //Just for everyone to quit gracefully
}
