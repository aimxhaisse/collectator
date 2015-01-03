package main

import (
	"flag"
	"time"
	"fmt"
)

var (
	periodSec = flag.Int("period", 1, "refresh period in seconds")
)

type Collectator struct {
	periodSec int // period in second between each collect
}

func (c *Collectator) Run() {
	for {
		fmt.Printf("Hello World\n")
		time.Sleep(time.Duration(c.periodSec) * time.Second)
	}
}

func NewCollectator() *Collectator {
	return &Collectator{
		*periodSec,
	}
}

func main() {
	flag.Parse()
	c := NewCollectator()
	c.Run()
}
