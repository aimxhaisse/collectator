package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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
	memActive float64 // current active memory in kB
}

func (c *Collectator) Run() {
	for {
		c.refreshMetrics()
		fmt.Printf("collectator> load_avg:%f mem_active:%f\n", c.loadAvg, c.memActive)
		time.Sleep(time.Duration(c.periodSec) * time.Second)
	}
}

func (c *Collectator) refreshMetrics() {
	c.refreshLoadAvg()
	c.refreshMemory()
}

var (
	re_mem_active = regexp.MustCompile("^Active:\\s+([0-9]+)\\s+kB$")
)

func (c *Collectator) refreshMemory() {
	file, err := os.Open("/proc/meminfo")
	if err == nil {
		defer file.Close()
		scan := bufio.NewScanner(file)
		for scan.Scan() {
			line := scan.Text()

			if m := re_mem_active.FindStringSubmatch(line); m != nil {
				fmt.Printf("match=%s\n", m[1])
				c.memActive, err = strconv.ParseFloat(m[1], 64)
			}
		}
	}
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
		0,
	}
}

func main() {
	flag.Parse()
	c := NewCollectator()
	c.Run()
}
