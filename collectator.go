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
	procStats ProcStatCollector
}

// man 5 proc
type ProcStat struct {
	user       float64
	nice       float64
	system     float64
	idle       float64
	iowait     float64
	irq        float64
	softirq    float64
	guest      float64
	guest_nice float64
}

type ProcStatCollector struct {
	last  *ProcStat // last values found in /proc/stat
	stats ProcStat  // in %
}

func (c *Collectator) Run() {
	for {
		c.refreshMetrics()
		fmt.Printf("collectator> load_avg:%f mem_active:%f cpu_user:%f cpu_nice:%f cpu_system:%f cpu_idle:%f cpu_iowait:%f cpu_irq:%f cpu_softirq:%f cpu_guest:%f cpu_guest_nice:%f\n",
			c.loadAvg,
			c.memActive,
			c.procStats.stats.user,
			c.procStats.stats.nice,
			c.procStats.stats.system,
			c.procStats.stats.idle,
			c.procStats.stats.iowait,
			c.procStats.stats.irq,
			c.procStats.stats.softirq,
			c.procStats.stats.guest,
			c.procStats.stats.guest_nice,
		)
		time.Sleep(time.Duration(c.periodSec) * time.Second)
	}
}

func (c *Collectator) refreshMetrics() {
	c.refreshLoadAvg()
	c.refreshMemory()
	c.refreshCPU()
}

func procStatDiff(left *ProcStat, right *ProcStat) *ProcStat {
	return &ProcStat{
		left.user - right.user,
		left.nice - right.nice,
		left.system - right.system,
		left.idle - right.idle,
		left.iowait - right.iowait,
		left.irq - right.irq,
		left.softirq - right.softirq,
		left.guest - right.guest,
		left.guest_nice - right.guest_nice,
	}
}

var (
	// Active:           546972 kB
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
				c.memActive, err = strconv.ParseFloat(m[1], 64)
			}
		}
	}
}

var (
	// cpu  4956172 1444 642144 166558262 27581 0 56040 0 0 0
	re_cpu_stat = regexp.MustCompile("^cpu\\s+(\\d+) (\\d+) (\\d+) (\\d+) (\\d+) (\\d+) (\\d+) (\\d+) (\\d+) (\\d+)$")
)

func (c *Collectator) refreshCPU() {
	file, err := os.Open("/proc/stat")
	if err == nil {
		defer file.Close()
		scan := bufio.NewScanner(file)
		for scan.Scan() {
			line := scan.Text()

			if m := re_cpu_stat.FindStringSubmatch(line); m != nil {
				var s ProcStat

				s.user, _ = strconv.ParseFloat(m[1], 64)
				s.nice, _ = strconv.ParseFloat(m[2], 64)
				s.system, _ = strconv.ParseFloat(m[3], 64)
				s.idle, _ = strconv.ParseFloat(m[4], 64)
				s.iowait, _ = strconv.ParseFloat(m[5], 64)
				s.irq, _ = strconv.ParseFloat(m[6], 64)
				s.softirq, _ = strconv.ParseFloat(m[7], 64)
				s.guest, _ = strconv.ParseFloat(m[8], 64)
				s.guest_nice, _ = strconv.ParseFloat(m[9], 64)

				if c.procStats.last != nil && *c.procStats.last != s {
					diff := procStatDiff(&s, c.procStats.last)
					total := diff.user + diff.nice + diff.system + diff.idle + diff.iowait + diff.irq + diff.softirq + diff.guest + diff.guest_nice
					if total != 0 {
						c.procStats.stats.user = (diff.user / total) * 100.0
						c.procStats.stats.nice = (diff.nice / total) * 100.0
						c.procStats.stats.system = (diff.system / total) * 100.0
						c.procStats.stats.idle = (diff.idle / total) * 100.0
						c.procStats.stats.iowait = (diff.iowait / total) * 100.0
						c.procStats.stats.irq = (diff.irq / total) * 100.0
						c.procStats.stats.softirq = (diff.softirq / total) * 100.0
						c.procStats.stats.guest = (diff.guest / total) * 100.0
						c.procStats.stats.guest_nice = (diff.guest_nice / total) * 100.0
					}
				}

				if c.procStats.last == nil || *c.procStats.last != s {
					c.procStats.last = &s
				}
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
			// 0.00 0.01 0.05 1/75 3930
			c.loadAvg, err = strconv.ParseFloat(line[:i], 64)
		}
	}
}

func NewCollectator() *Collectator {
	return &Collectator{
		*periodSec,
		0,
		0,
		ProcStatCollector{},
	}
}

func main() {
	flag.Parse()
	c := NewCollectator()
	c.Run()
}
