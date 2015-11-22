package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

type IntGraph struct {
	inner map[int]map[int]int
}

func (p *IntGraph) SetValue(k1, k2, val int) {
	_, ok := p.inner[k1][k2]
	if !ok {
		_, ok = p.inner[k1]
		if !ok {
			p.inner[k1] = make(map[int]int)
		}
	}
	p.inner[k1][k2] = val
	return
}

func (p *IntGraph) GetValue(k1, k2 int) int {
	val, ok := p.inner[k1][k2]
	if !ok {
		return 0
	}
	return val
}

func (p *IntGraph) Iter() []int {
	q := []int{}
	for k, _ := range p.inner {
		q = append(q, k)
	}
	return q
}

func (p *IntGraph) IterIdx(k1 int) []int {
	q := []int{}
	_, ok := p.inner[k1]
	if !ok {
		return q
	}

	for k, _ := range p.inner[k1] {
		q = append(q, k)
	}

	return q
}

func NewIntGraph() *IntGraph {
	p := &IntGraph{}
	p.inner = make(map[int]map[int]int)
	return p
}

func Debug(format string, a ...interface{}) {
	_, fname, lineno, _ := runtime.Caller(1)
	s := fmt.Sprintf("%s:%d\t", fname, lineno)
	s += fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, s)
	return
}

func SetGraphs(graph *IntGraph, k1 int, k2 int, cap32 int, source int, sink int) bool {

	/*we can not set
	self-cycle k1 == k2
	sink as from k1 == sink
	source as to k2 == source
	value zero cap32 == 0
	*/
	if k1 == k2 || k1 == sink || k2 == source || cap32 == 0 {
		return false
	}

	/* we can not set twice value*/
	if graph.GetValue(k1, k2) != 0 {
		return false
	}

	/*we can not accept the reserver way in the value*/
	if graph.GetValue(k2, k1) != 0 {
		return false
	}
	if k2 == source {
		Debug("can not be source k2\n")
		os.Exit(4)
	}

	if k1 == sink {
		Debug("can not be sink k1\n")
		os.Exit(4)
	}

	graph.SetValue(k1, k2, cap32)
	return true
}

func RandMake(pntcnt int, edgecnt int, maxcap int) *IntGraph {

	var p32, m32 int32
	var ok bool

	p32 = int32(pntcnt)
	m32 = int32(maxcap)
	graph := NewIntGraph()
	i := 0
	for i < edgecnt {
		pnt1 := int(rand.Int31n(p32))
		pnt2 := int(rand.Int31n(p32))
		cap32 := int(rand.Int31n(m32))

		ok = SetGraphs(graph, pnt1, pnt2, cap32, 0, pntcnt-1)
		if ok {
			i += 1
		}
	}
	return graph
}

func OutPutGraph(graph *IntGraph, source int, sink int) {
	fmt.Fprintf(os.Stdout, "source=%d\n", source)
	fmt.Fprintf(os.Stdout, "sink=%d\n", sink)
	w := int(math.Sqrt(float64(sink + 1)))
	h := w + 1
	fmt.Fprintf(os.Stdout, "width=%d\n", w)
	fmt.Fprintf(os.Stdout, "height=%d\n", h)
	for _, from := range graph.Iter() {
		for _, to := range graph.IterIdx(from) {
			fmt.Fprintf(os.Stdout, "%d,%d,%d\n", from, to, graph.GetValue(from, to))
		}
	}
	return

}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "%s pointcnt edgecnt [maxcap]\n", os.Args[0])
		os.Exit(4)
	}
	rand.Seed(int64(time.Now().Nanosecond()))
	pntcnt, err := strconv.Atoi(os.Args[1])
	if err != nil || pntcnt == 0 {
		fmt.Fprintf(os.Stderr, "%s not valid number\n", os.Args[1])
		os.Exit(4)
	}
	edgecnt, err := strconv.Atoi(os.Args[2])
	if err != nil || edgecnt == 0 || edgecnt <= pntcnt || edgecnt >= (pntcnt*(pntcnt-1))/2 {
		fmt.Fprintf(os.Stderr, "%s not valid number\n", os.Args[2])
		os.Exit(4)
	}
	maxcap := 100
	if len(os.Args) > 3 {
		maxcap, _ = strconv.Atoi(os.Args[3])
	}
	graph := RandMake(pntcnt, edgecnt, maxcap)
	OutPutGraph(graph, 0, pntcnt-1)
	return
}
