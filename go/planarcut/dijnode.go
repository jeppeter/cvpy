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

func FormEdgeName(from, to *DijVertice) string {
	return fmt.Sprint("%s->%s", from.GetName(), to.GetName())
}

func NewDijEdge(from, to *DijVertice, length int) *DijEdge {
	p := &DijEdge{}
	p.from = from
	p.to = to
	p.length = length
	p.name = FormEdgeName(from, to)
	return p
}
