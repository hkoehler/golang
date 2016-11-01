package main

import (
	"fmt"
	"flag"
)

var (
	relayChainLen int = 10
)

func worker(c, done chan int) {
	v := <- c
	//fmt.Println(v)
	if v < relayChainLen {
		go worker(c, done)
		c <- v+1
	} else {
		fmt.Println("Done")
		done <- 1
	}
}

func main() {
	var num int
	var done []chan int

	flag.IntVar(&relayChainLen, "chain", 10, "length of relay chain")
	flag.IntVar(&num, "num", 10, "number of concurrent workers")
	flag.Parse()
	for i := 0; i < num; i++ {
		d := make(chan int)
		c := make(chan int)
		go worker(c, d)
		c <- 1
		done = append(done, d)
	}
	for _, d := range done {
		<- d
	}
}
