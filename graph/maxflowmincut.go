package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Edge struct {
	source  string
	sink    string
	capcity int
}

type FlowNetwork struct {
	adj  map[string][]*Edge
	flow map[string]int
}

func NewEdge(source, sink string, capcity int) *Edge {
	e := &Edge{}
	e.source = source
	e.sink = sink
	e.capcity = capcity
	return e
}

func (e *Edge) Stringer() string {
	s := fmt.Sprintf("%s->%s:%d", e.source, e.sink, e.capcity)
	return s
}

func NewFlowNetwork() *FlowNetwork {
	f := &FlowNetwork{}
	f.adj = make(map[string][]*Edge)
	f.flow = make(map[string]int)
	return f
}

func (f *FlowNetwork) Add_Vertex(vertex string) {
	if _, ok := f.adj[vertex]; ok {
		return
	}

	f.adj[vertex] = nil
	return
}

func (f *FlowNetwork) Get_Edges(vertex string) []*Edge {
	if _, ok := f.adj[vertex]; !ok {
		return nil
	}
	return f.adj[vertex]
}

func (f *FlowNetwork) Add_Edge(source, sink string, capcity int) error {
	f.Add_Vertex(source)
	f.Add_Vertex(sink)
	edge := NewEdge(source, sink, capcity)
	redge := NewEdge(sink, source, 0)
	f.adj[source] = append(f.adj[source], edge)
	f.adj[sink] = append(f.adj[sink], redge)
	f.flow[source] = 0
	f.flow[sink] = 0
	return nil
}

func SetDictDefValue(caps map[string]map[string]int, fk string, sk string, defvalue int) map[string]map[string]int {
	mm, ok := caps[fk]
	if !ok {
		mm = make(map[string]int)
		caps[fk] = mm
		mm[sk] = defvalue
	}
	if _, ok = mm[sk]; !ok {
		mm[sk] = defvalue
	}
	return caps
}

func IsInArray(arr []string, key string) int {
	for _, k := range arr {
		if key == k {
			return 1
		}
	}

	return 0
}

func Debug(format string, a ...interface{}) {
	_, fname, lineno, _ := runtime.Caller(1)
	s := fmt.Sprintf("%s:%d\t", fname, lineno)
	s += fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stdout, s)
	return
}

func DebugMapString(caps map[string]map[string]int, format string, a ...interface{}) {
	var sortkeys []string
	var longestkey int
	longestkey = 4
	sortkeys = MakeSortKeys(caps)
	if format != "" {
		s := fmt.Sprintf(format, a...)
		fmt.Fprintf(os.Stdout, s)
		fmt.Fprintf(os.Stdout, "\n")
	}
	for _, k1 := range sortkeys {
		if longestkey < len(k1) {
			longestkey = len(k1)
		}
	}

	fmt.Fprintf(os.Stdout, "%*s[", longestkey, "tags")
	for _, k1 := range sortkeys {
		fmt.Fprintf(os.Stdout, "%*s", longestkey, k1)
	}
	fmt.Fprintf(os.Stdout, "]\n")

	for _, k1 := range sortkeys {
		fmt.Fprintf(os.Stdout, "%*s[", longestkey, k1)
		for _, k2 := range sortkeys {
			val := 0
			if _, ok := caps[k1][k2]; ok {
				val = caps[k1][k2]
			}
			fmt.Fprintf(os.Stdout, "%*d", longestkey, val)

		}
		fmt.Fprintf(os.Stdout, "]\n")
	}
}

func (f *FlowNetwork) Get_Cap_Neighbour() (capcities map[string]map[string]int,
	neighbours map[string][]string) {
	var sortkeys []string
	caps := make(map[string]map[string]int)
	neigh := make(map[string][]string)
	for k := range f.adj {
		sortkeys = append(sortkeys, k)
	}

	for k, ev := range f.adj {
		for _, edge := range ev {
			caps = SetDictDefValue(caps, edge.source, edge.sink, edge.capcity)
			if IsInArray(neigh[k], edge.sink) == 0 {
				neigh[k] = append(neigh[k], edge.sink)
			}
			if IsInArray(neigh[edge.sink], edge.source) == 0 {
				neigh[edge.sink] = append(neigh[edge.sink], edge.source)
			}
		}
	}

	return caps, neigh
}

func GetGraphFromFile(infile string) (f *FlowNetwork, source string, sink string, err error) {
	var sarr []string
	var caps int
	file, e := os.Open(infile)
	if e != nil {
		return nil, "", "", e
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	sink = ""
	source = ""
	f = NewFlowNetwork()
	for scanner.Scan() {
		l := scanner.Text()
		l = strings.Trim(l, "\r\n")
		if strings.HasPrefix(l, "#") {
			continue
		}

		if strings.HasPrefix(l, "source=") {
			sarr = strings.Split(l, "=")
			if len(sarr) < 2 {
				continue
			}
			source = sarr[1]
			continue
		}

		if strings.HasPrefix(l, "sink=") {
			sarr = strings.Split(l, "=")
			if len(sarr) < 2 {
				continue
			}
			sink = sarr[1]
			continue
		}

		sarr = strings.Split(l, ",")
		if len(sarr) < 3 {
			continue
		}

		caps, _ = strconv.Atoi(sarr[2])
		f.Add_Edge(sarr[0], sarr[1], caps)
	}

	if sink == "" {
		err = fmt.Errorf("no sink specify")
		return nil, "", "", err
	}

	if source == "" {
		err = fmt.Errorf("no source specify")
		return nil, "", "", err
	}
	return f, source, sink, nil
}

func MakeSortKeys(caps map[string]map[string]int) []string {
	var retstr []string
	for k := range caps {
		retstr = append(retstr, k)
	}

	for i := 0; i < len(retstr); i++ {
		for j := (i + 1); j < len(retstr); j++ {
			if retstr[j] < retstr[i] {
				tmp := retstr[j]
				retstr[j] = retstr[i]
				retstr[i] = tmp
			}
		}
	}

	return retstr

}

func main() {
	var flow int
	var flows map[string]map[string]int
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stdout, "%s [ed|gt] infile\n", os.Args[0])
		fmt.Fprintf(os.Stdout, "\ted for Edmonds-Karp algorithm\n")
		fmt.Fprintf(os.Stdout, "\tgt for Goldberg-Tarjan algorithm\n")
		os.Exit(4)
	}

	f, s, t, e := GetGraphFromFile(os.Args[2])
	if e != nil {
		fmt.Fprintf(os.Stderr, "error %v\n", e.Error())
		os.Exit(4)
	}
	caps, neighs := f.Get_Cap_Neighbour()
	if os.Args[1] == "edmonds" {
		flow, flows = EdmondsWarp(caps, neighs, s, t)
	} else {
		flow, flows = GoldbergTarjan(caps, neighs, s, t)
	}
	fmt.Fprintf(os.Stdout, "flow %d\n", flow)
	DebugMapString(caps, "caps ")
	DebugMapString(flows, "flows ")
	return
}
