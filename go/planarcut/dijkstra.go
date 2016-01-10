package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const MAXINT int = (1 << 31)

type Vertice struct {
	name      string
	dist      int
	prev      *Vertice
	visited   bool
	edges     []*Edge
	nedges    int
	link_next *Vertice
}

func NewVertice(name string) *Vertice {
	p := &Vertice{}
	p.name = name
	p.dist = MAXINT
	p.prev = nil
	p.visited = false
	p.edges = []*Edge{}
	p.nedges = 0
	p.link_next = nil
	return p
}

func (vert *Vertice) Stringer() string {
	return fmt.Sprintf("(%s) dist (%d)", vert.name, vert.dist)
}
func (vert *Vertice) TypeName() string {
	return "Vertice"
}
func (vert *Vertice) Equal(j RBTreeData) bool {
	var jv *Vertice
	if vert.TypeName() != j.TypeName() {
		log.Fatalf("vert (%s) != j (%s)", vert.TypeName(), j.TypeName())
	}
	jv = ((*Vertice)(unsafe.Pointer((reflect.ValueOf(j).Pointer()))))
	if jv == vert {
		return true
	}
	return false
}

func (vert *Vertice) Less(j RBTreeData) bool {
	var jv *Vertice
	if vert.TypeName() != j.TypeName() {
		log.Fatalf("vert (%s) != j (%s)", vert.TypeName(), j.TypeName())
	}
	jv = ((*Vertice)(unsafe.Pointer((reflect.ValueOf(j).Pointer()))))
	if vert.dist < jv.dist {
		return true
	}
	return false
}

func (vert *Vertice) GetDist() int {
	return vert.dist
}

func (vert *Vertice) SetDist(dist int) {
	vert.dist = dist
	return
}

func (vert *Vertice) SetPrev(p *Vertice) {
	vert.prev = p
	return
}

func (vert *Vertice) GetPrev() *Vertice {
	return vert.prev
}

func (vert *Vertice) IsVisited() bool {
	return vert.visited
}

func (vert *Vertice) Visit() {
	vert.visited = true
	return
}

func (vert *Vertice) UnVisit() {
	vert.visited = false
	return
}

func (vert *Vertice) AddEdge(pe *Edge) {
	vert.edges = append(vert.edges, pe)
	vert.nedges++
	return
}

func (vert *Vertice) GetEdges() []*Edge {
	return vert.edges
}

func (vert *Vertice) GetName() string {
	return vert.name
}

func (vert *Vertice) GetNext() *Vertice {
	return vert.link_next
}

func (vert *Vertice) SetNext(pnext *Vertice) {
	vert.link_next = pnext
}

func FormEdgeName(from, to *Vertice) string {
	return fmt.Sprint("%s->%s", from.GetName(), to.GetName())
}

type Edge struct {
	from   *Vertice
	to     *Vertice
	length int
	name   string
}

func NewEdge(from, to *Vertice, length int) *Edge {
	p := &Edge{}
	p.from = from
	p.to = to
	p.length = length
	p.name = FormEdgeName(from, to)
	return p
}

func (e *Edge) GetFrom() *Vertice {
	return e.from
}

func (e *Edge) GetTo() *Vertice {
	return e.to
}

func (e *Edge) GetLength() int {
	return e.length
}

func (e *Edge) GetName() string {
	return e.name
}

type Graph struct {
	edges  map[string]*Edge
	verts  map[string]*Vertice
	source string
	sink   string
	queue  *RBTree
}

func NewGraph() *Graph {
	p := &Graph{}
	p.edges = make(map[string]*Edge)
	p.verts = make(map[string]*Vertice)
	p.source = ""
	p.sink = ""
	p.queue = NewRBTree()
	return p
}

func (g *Graph) SetSource(source string) {
	g.source = source
	return
}

func (g *Graph) SetSink(sink string) {
	g.sink = sink
	return
}

func (g *Graph) AddEdge(from, to string, caps int) error {
	fvert, fok := g.verts[from]
	tvert, tok := g.verts[to]
	if !fok {
		fvert = NewVertice(from)
		g.verts[from] = fvert
	}

	if !tok {
		tvert = NewVertice(to)
		g.verts[to] = tvert
	}

	e, eok := g.edges[FormEdgeName(fvert, tvert)]
	re, reok := g.edges[FormEdgeName(tvert, fvert)]
	if eok {
		return fmt.Errorf("%s has already in", FormEdgeName(fvert, tvert))
	}

	if reok {
		return fmt.Errorf("%s has already in", FormEdgeName(tvert, fvert))
	}

	e = NewEdge(fvert, tvert, caps)
	re = NewEdge(tvert, fvert, caps)
	g.edges[FormEdgeName(fvert, tvert)] = e
	g.edges[FormEdgeName(tvert, fvert)] = re

	/*now we should add edge for the */
	fvert.AddEdge(e)
	tvert.AddEdge(re)
	return nil
}

func (g *Graph) InsertQueue(vert *Vertice) {
	g.queue.Insert(vert)
	return
}

func (g *Graph) ReinsertQueue(vert *Vertice) {
	if vert.IsVisited() {
		_, err := g.queue.Delete(vert)
		if err == nil {
			g.queue.Insert(vert)
		}
	}
}

func (g *Graph) GetQueue() *Vertice {
	var rbdata RBTreeData
	var pvert *Vertice
	rbdata = g.queue.GetMin()
	if rbdata == nil {
		return nil
	}

	pvert = ((*Vertice)(unsafe.Pointer((reflect.ValueOf(rbdata).Pointer()))))
	return pvert
}

func (g *Graph) Dijkstra1() (dist int, err error) {
	var tvert, cvert, svert, dstvert *Vertice

	svert, ok := g.verts[g.source]
	if !ok {
		return 0, fmt.Errorf("source (%s) not found", g.source)
	}
	dstvert, ok = g.verts[g.sink]
	if !ok {
		return 0, fmt.Errorf("sink (%s) not found", g.sink)
	}

	/*init for the */
	svert.SetDist(0)

	for _, tvert = range g.verts {
		g.InsertQueue(tvert)
	}
	cvert = svert

	for {
		cvert = g.GetQueue()
		if cvert == nil {
			break
		}
		log.Printf("get (%s)", cvert.GetName())

		alt := MAXINT
		for _, e := range cvert.GetEdges() {
			tvert := e.GetTo()
			if tvert.IsVisited() {
				continue
			}
			alt = cvert.GetDist() + e.GetLength()
			if alt < tvert.GetDist() {
				tvert.SetDist(alt)
				log.Printf("set (%s) parent (%s)", tvert.GetName(), cvert.GetName())
				tvert.SetPrev(cvert)
			}
		}
	}

	if dstvert.GetPrev() == nil {
		return 0, fmt.Errorf("(%s->%s) not connected", g.source, g.sink)
	}

	return dstvert.GetDist(), nil
}

func (g *Graph) Dijkstra2() (dist int, err error) {
	var tvert, cvert, svert, dstvert *Vertice

	svert, ok := g.verts[g.source]
	if !ok {
		return 0, fmt.Errorf("source (%s) not found", g.source)
	}
	dstvert, ok = g.verts[g.sink]
	if !ok {
		return 0, fmt.Errorf("sink (%s) not found", g.sink)
	}

	/*init for the */
	svert.SetDist(0)

	for _, tvert = range g.verts {
		g.InsertQueue(tvert)
	}
	cvert = svert

	for {
		cvert = g.GetQueue()
		if cvert == nil || cvert == dstvert {
			break
		}
		//log.Printf("get (%s)", cvert.GetName())

		alt := MAXINT
		for _, e := range cvert.GetEdges() {
			tvert := e.GetTo()
			if tvert.IsVisited() {
				continue
			}
			alt = cvert.GetDist() + e.GetLength()
			if alt < tvert.GetDist() {
				tvert.SetDist(alt)
				//log.Printf("set (%s) parent (%s)", tvert.GetName(), cvert.GetName())
				tvert.SetPrev(cvert)
			}
		}
	}

	if dstvert.GetPrev() == nil {
		return 0, fmt.Errorf("(%s->%s) not connected", g.source, g.sink)
	}

	return dstvert.GetDist(), nil
}

func (g *Graph) Dijkstra() (dist int, err error) {
	var cvert, svert, dstvert *Vertice

	svert, ok := g.verts[g.source]
	if !ok {
		return 0, fmt.Errorf("source (%s) not found", g.source)
	}
	dstvert, ok = g.verts[g.sink]
	if !ok {
		return 0, fmt.Errorf("sink (%s) not found", g.sink)
	}

	/*init for the */
	for _, cvert = range g.verts {
		cvert.SetDist(MAXINT)
		cvert.UnVisit()
	}
	svert.SetDist(0)
	cvert = svert
	svert.Visit()
	g.InsertQueue(cvert)

	for {
		cvert = g.GetQueue()
		//if cvert == nil || cvert == dstvert {
		if cvert == nil {
			break
		}
		//log.Printf("get (%s) dist (%d)", cvert.GetName(), cvert.GetDist())
		for _, e := range cvert.GetEdges() {
			tvert := e.GetTo()
			alt := cvert.GetDist() + e.GetLength()
			if alt < tvert.GetDist() {
				//log.Printf("set (%s) (%d -> %d)", tvert.GetName(), tvert.GetDist(), alt)
				tvert.SetDist(alt)
				tvert.SetPrev(cvert)
			}
			if !tvert.IsVisited() {
				//log.Printf("push into (%s) dist(%d)", tvert.GetName(), tvert.GetDist())
				tvert.Visit()
				g.InsertQueue(tvert)
			}
		}
	}

	if dstvert.GetPrev() == nil {
		return 0, fmt.Errorf("(%s->%s) not connected", g.source, g.sink)
	}

	return dstvert.GetDist(), nil
}

func (g *Graph) GetPath() []string {
	var rs, s []string
	var sinkvert, sourcevert, curvert *Vertice
	rs = []string{}
	sinkvert = g.verts[g.sink]
	sourcevert = g.verts[g.source]
	curvert = sinkvert
	for curvert != sourcevert {
		if curvert == nil {
			return rs
		}
		rs = append(rs, curvert.GetName())

		curvert = curvert.GetPrev()
	}
	rs = append(rs, curvert.GetName())

	s = []string{}
	for i := len(rs) - 1; i >= 0; i-- {
		s = append(s, rs[i])
	}
	return s

}

func ParseFile(infile string) *Graph {
	var words []string
	var caps int
	g := NewGraph()
	f, err := os.Open(infile)
	if err != nil {
		return nil
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		l = strings.Trim(l, "\r\n")
		if strings.HasPrefix(l, "#") {
			continue
		}
		if strings.HasPrefix(l, "source=") {
			words = strings.Split(l, "=")
			if len(words) < 2 {
				continue
			}
			g.SetSource(words[1])
			continue
		} else if strings.HasPrefix(l, "sink=") {
			words = strings.Split(l, "=")
			if len(words) < 2 {
				continue
			}
			g.SetSink(words[1])
			continue
		}
		words = strings.Split(l, ",")
		if len(words) < 3 {
			continue
		}
		caps, _ = strconv.Atoi(words[2])
		g.AddEdge(words[0], words[1], caps)
	}

	return g
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s infile\n", os.Args[0])
		os.Exit(4)
	}
	log.SetFlags(log.Lshortfile)

	g := ParseFile(os.Args[1])
	stime := time.Now()
	f, e := g.Dijkstra()
	if e != nil {
		fmt.Fprintf(os.Stderr, "error %s\n", e.Error())
		os.Exit(4)
	}
	etime := time.Now()

	path := g.GetPath()
	fmt.Fprintf(os.Stdout, "(%s) dist (%d) path (%v) (%s)\n", os.Args[1], f, path, etime.Sub(stime))

}
