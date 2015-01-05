package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var (
	periodSec = flag.Int("period", 1, "refresh period in seconds")
)

type Collectator struct {
	periodSec int     // period in second between each collect
	loadAvg   float64 // current load average
}

func (c *Collectator) Run() {
	for {
		c.refreshMetrics()
		fmt.Printf("collectator> load_avg:%f\n", c.loadAvg)
		time.Sleep(time.Duration(c.periodSec) * time.Second)
	}
}

func (c *Collectator) refreshMetrics() {
	c.refreshLoadAvg()
}

func (c *Collectator) refreshLoadAvg() {
	buf, err := ioutil.ReadFile("/proc/loadavg")
	if err == nil {
		line := string(buf)
		i := strings.IndexByte(line, ' ')
		if i > 0 {
			c.loadAvg, err = strconv.ParseFloat(line[:i], 64)
		}
	}
}

func NewCollectator() *Collectator {
	return &Collectator{
		*periodSec,
		0,
	}
}

func main() {
	flag.Parse()
	c := NewCollectator()
	c.Run()
}
