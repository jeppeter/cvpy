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
	edges     []*Edge
	nedges    int
	link_next *DijVertice
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
