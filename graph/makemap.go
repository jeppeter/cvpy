package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

func Debug(format string, a ...interface{}) {
	_, fname, lineno, _ := runtime.Caller(1)
	s := fmt.Sprintf("%s:%d\t", fname, lineno)
	s += fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, s)
	return
}

func SetGraphs(graph map[int]map[int]int, k1 int, k2 int, cap32 int) (gr map[int]map[int]int, bret bool) {
	var ok bool
	if k1 == k2 || cap32 == 0 {
		return graph, false
	}
	_, ok = graph[k1][k2]
	if ok {
		return graph, false
	}
	_, ok = graph[k1][k2]
	if ok {
		return graph, false
	}

	_, ok = graph[k1]
	if !ok {
		graph[k1] = make(map[int]int)
	}
	graph[k1][k2] = cap32
	return graph, true
}

func RandMake(pntcnt int, edgecnt int, maxcap int) map[int]map[int]int {

	var p32, m32 int32
	var ok bool
	p32 = int32(pntcnt)
	m32 = int32(maxcap)
	graph := make(map[int]map[int]int)
	graph = graph
	i := 0
	for i < edgecnt {
		pnt1 := int(rand.Int31n(p32))
		pnt2 := int(rand.Int31n(p32))
		cap32 := int(rand.Int31n(m32))

		graph, ok = SetGraphs(graph, pnt1, pnt2, cap32)
		if ok {
			i += 1
		}
	}
	return graph
}

func OutPutGraph(graph map[int]map[int]int, source int, sink int) {
	fmt.Fprintf(os.Stdout, "source=%d\n", source)
	fmt.Fprintf(os.Stdout, "sink=%d\n", sink)
	for from := range graph {
		for to, caps := range graph[from] {
			fmt.Fprintf(os.Stdout, "%d,%d,%d\n", from, to, caps)
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
	pntcnt, _ := strconv.Atoi(os.Args[1])
	edgecnt, _ := strconv.Atoi(os.Args[2])
	maxcap := 100
	if len(os.Args) > 3 {
		maxcap, _ = strconv.Atoi(os.Args[3])
	}
	graph := RandMake(pntcnt, edgecnt, maxcap)
	OutPutGraph(graph, 0, pntcnt-1)
	return
}
