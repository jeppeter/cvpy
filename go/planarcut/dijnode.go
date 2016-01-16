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

type DijVertice struct {
	name      string
	dist      int
	prev      *DijVertice
	visited   bool
	edges     []*DijEdge
	nedges    int
	link_next *DijVertice
}

type DijEdge struct {
	from   *DijVertice
	to     *DijVertice
	length int
	name   string
}

func NewDijVertice(name string) *DijVertice {
	p := &DijVertice{}
	p.name = name
	p.dist = MAXINT
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
		log.Fatalf("vert (%s) != j (%s)", vert.TypeName(), j.TypeName())
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
		log.Fatalf("vert (%s) != j (%s)", vert.TypeName(), j.TypeName())
	}
	jv = ((*DijVertice)(unsafe.Pointer((reflect.ValueOf(j).Pointer()))))
	if vert.dist < jv.dist {
		return true
	}
	return false
}

func (vert *DijVertice) GetDist() int {
	return vert.dist
}

func (vert *DijVertice) SetDist(dist int) {
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

func NewDijEdge(from, to *DijVertice, length int) *DijEdge {
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

func (e *DijEdge) GetLength() int {
	return e.length
}

func (e *DijEdge) GetName() string {
	return e.name
}

type DijGraph struct {
	edges      map[string]*DijEdge
	verts      map[string]*DijVertice
	vertnum    int
	source     string
	sink       string
	queue2     *RBTree
	queue      []*DijVertice
	queuestart int
	queueend   int
}

func NewDijGraph() *DijGraph {
	p := &DijGraph{}
	p.edges = make(map[string]*DijEdge)
	p.verts = make(map[string]*DijVertice)
	p.vertnum = 0
	p.source = ""
	p.sink = ""
	p.queue2 = NewRBTree()
	p.queue = nil
	p.queuestart = -1
	p.queueend = -1
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

func (g *DijGraph) AddEdge(from, to string, caps int) error {
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
	re = NewDijEdge(tvert, fvert, caps)
	g.edges[DijFormEdgeName(fvert, tvert)] = e
	g.edges[DijFormEdgeName(tvert, fvert)] = re

	/*now we should add edge for the */
	fvert.AddEdge(e)
	tvert.AddEdge(re)
	return nil
}

func (g *Graph) InsertQueue2(vert *Vertice) {
	g.queue2.Insert(vert)
	return
}
