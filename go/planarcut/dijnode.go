package main

import (
	"fmt"
	"os"
	"reflect"
	"unsafe"
)

type DijVertice struct {
	name      string
	dist      float64
	prev      *DijVertice
	visited   bool
	edges     []*DijEdge
	nedges    int
	link_next *DijVertice
}

type DijEdge struct {
	from   *DijVertice
	to     *DijVertice
	length float64
	name   string
}

func NewDijVertice(name string) *DijVertice {
	p := &DijVertice{}
	p.name = name
	p.dist = CAP_INF
	p.prev = nil
	p.visited = false
	p.edges = []*DijEdge{}
	p.nedges = 0
	p.link_next = nil
	return p
}

func (vert *DijVertice) Stringer() string {
	return fmt.Sprintf("(%s) dist (%d)", vert.name, vert.dist)
}
func (vert *DijVertice) TypeName() string {
	return "DijVertice"
}
func (vert *DijVertice) Equal(j RBTreeData) bool {
	var jv *DijVertice
	if vert.TypeName() != j.TypeName() {
		Error("vert (%s) != j (%s)", vert.TypeName(), j.TypeName())
		os.Exit(5)
	}
	jv = ((*DijVertice)(unsafe.Pointer((reflect.ValueOf(j).Pointer()))))
	if jv == vert {
		return true
	}
	return false
}

func (vert *DijVertice) Less(j RBTreeData) bool {
	var jv *DijVertice
	if vert.TypeName() != j.TypeName() {
		Error("vert (%s) != j (%s)", vert.TypeName(), j.TypeName())
		os.Exit(5)
	}
	jv = ((*DijVertice)(unsafe.Pointer((reflect.ValueOf(j).Pointer()))))
	if vert.dist < jv.dist {
		return true
	}
	return false
}

func (vert *DijVertice) GetDist() float64 {
	return vert.dist
}

func (vert *DijVertice) SetDist(dist float64) {
	vert.dist = dist
	return
}

func (vert *DijVertice) SetPrev(p *DijVertice) {
	vert.prev = p
	return
}

func (vert *DijVertice) GetPrev() *DijVertice {
	return vert.prev
}

func (vert *DijVertice) IsVisited() bool {
	return vert.visited
}

func (vert *DijVertice) Visit() {
	vert.visited = true
	return
}

func (vert *DijVertice) UnVisit() {
	vert.visited = false
	return
}

func (vert *DijVertice) AddEdge(pe *DijEdge) {
	vert.edges = append(vert.edges, pe)
	vert.nedges++
	return
}

func (vert *DijVertice) GetEdges() []*DijEdge {
	return vert.edges
}

func (vert *DijVertice) GetName() string {
	return vert.name
}

func (vert *DijVertice) GetNext() *DijVertice {
	return vert.link_next
}

func (vert *DijVertice) SetNext(pnext *DijVertice) {
	vert.link_next = pnext
}

func DijFormEdgeName(from, to *DijVertice) string {
	return fmt.Sprint("%s->%s", from.GetName(), to.GetName())
}

func NewDijEdge(from, to *DijVertice, length float64) *DijEdge {
	p := &DijEdge{}
	p.from = from
	p.to = to
	p.length = length
	p.name = DijFormEdgeName(from, to)
	return p
}

func (e *DijEdge) GetFrom() *DijVertice {
	return e.from
}

func (e *DijEdge) GetTo() *DijVertice {
	return e.to
}

func (e *DijEdge) GetLength() float64 {
	return e.length
}

func (e *DijEdge) GetName() string {
	return e.name
}

type DijGraph struct {
	edges   map[string]*DijEdge
	verts   map[string]*DijVertice
	vertnum int
	source  string
	sink    string
	queue2  *RBTree
}

func NewDijGraph() *DijGraph {
	p := &DijGraph{}
	p.edges = make(map[string]*DijEdge)
	p.verts = make(map[string]*DijVertice)
	p.vertnum = 0
	p.source = ""
	p.sink = ""
	p.queue2 = NewRBTree()
	return p
}

func (g *DijGraph) SetSource(source string) {
	g.source = source
	return
}

func (g *DijGraph) SetSink(sink string) {
	g.sink = sink
	return
}

func (g *DijGraph) AddEdge(from, to string, caps, rcaps float64) error {
	fvert, fok := g.verts[from]
	tvert, tok := g.verts[to]
	if !fok {
		fvert = NewDijVertice(from)
		g.verts[from] = fvert
		g.vertnum++
	}

	if !tok {
		tvert = NewDijVertice(to)
		g.verts[to] = tvert
		g.vertnum++
	}

	e, eok := g.edges[DijFormEdgeName(fvert, tvert)]
	re, reok := g.edges[DijFormEdgeName(tvert, fvert)]
	if eok {
		return fmt.Errorf("%s has already in", DijFormEdgeName(fvert, tvert))
	}

	if reok {
		return fmt.Errorf("%s has already in", DijFormEdgeName(tvert, fvert))
	}

	e = NewDijEdge(fvert, tvert, caps)
	re = NewDijEdge(tvert, fvert, rcaps)
	g.edges[DijFormEdgeName(fvert, tvert)] = e
	g.edges[DijFormEdgeName(tvert, fvert)] = re

	/*now we should add edge for the */
	fvert.AddEdge(e)
	tvert.AddEdge(re)
	return nil
}

func (g *DijGraph) InsertQueue2(vert *DijVertice) {
	g.queue2.Insert(vert)
	return
}

func (g *DijGraph) GetQueue2() *DijVertice {
	var rbdata RBTreeData
	var pvert *DijVertice
	rbdata = g.queue2.GetMin()
	if rbdata == nil {
		return nil
	}

	pvert = ((*DijVertice)(unsafe.Pointer((reflect.ValueOf(rbdata).Pointer()))))
	return pvert
}

func (g *DijGraph) Dijkstra() (dist float64, err error) {
	var cvert, svert, dstvert *DijVertice

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
		cvert.SetDist(CAP_INF)
		cvert.UnVisit()
	}
	svert.SetDist(CAP_ZERO)
	cvert = svert
	svert.Visit()
	g.InsertQueue2(svert)

	for {
		cvert = g.GetQueue2()
		if cvert == nil || cvert == dstvert {
			//if cvert == nil {
			break
		}
		//Debug("get (%s) dist (%d)", cvert.GetName(), cvert.GetDist())
		for _, e := range cvert.GetEdges() {
			tvert := e.GetTo()
			alt := cvert.GetDist() + e.GetLength()
			if alt < tvert.GetDist() {
				//Debug("set (%s) (%d -> %d)", tvert.GetName(), tvert.GetDist(), alt)
				/*we delete it and reinsert it into */
				_, err := g.queue2.Delete(tvert)
				tvert.SetDist(alt)
				tvert.SetPrev(cvert)
				if err == nil {
					g.queue2.Insert(tvert)
				}
			}
			if !tvert.IsVisited() {
				//Debug("push into (%s) dist(%d)", tvert.GetName(), tvert.GetDist())
				tvert.Visit()
				g.InsertQueue2(tvert)
			}
		}
	}

	if dstvert.GetPrev() == nil {
		return 0, fmt.Errorf("(%s->%s) not connected", g.source, g.sink)
	}

	return dstvert.GetDist(), nil
}

func (g *DijGraph) GetPath() []string {
	var rs, s []string
	var sinkvert, sourcevert, curvert *DijVertice
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

func (g *DijGraph) GetWeigth(name string) float64 {
	v, ok := g.verts[name]
	if !ok {
		Error("can not find %s verts", name)
		os.Exit(5)
		return -1
	}

	return v.GetDist()

}
