package main

import (
	"fmt"
)

type Arc struct {
	head   *Node
	next   *Arc
	sister *Arc
	r_cap  int
}

func NewArc() *Arc {
	p := &Arc{}
	p.head = nil
	p.next = nil
	p.sister = nil
	p.r_cap = 0
	return p
}

func (parc *Arc) SetHead(pnode *Node) {
	parc.head = pnode
	return
}

func (parc *Arc) GetHead() *Node {
	return parc.head
}

func (parc *Arc) SetNext(pnext *Arc) {
	parc.next = pnext
	return
}

func (parc *Arc) GetNext() *Arc {
	return parc.next
}

func (parc *Arc) SetSister(psister *Arc) {
	parc.sister = psister
	return
}

func (parc *Arc) GetSister() *Arc {
	return parc.sister
}

func (parc *Arc) GetCap() int {
	return parc.r_cap
}

func (parc *Arc) SetCap(caps int) {
	parc.r_cap = caps
	return
}

type Node struct {
	first   *Arc
	parent  *Arc
	next    *Node
	TS      int
	DIST    int
	is_sink bool
	tr_cap  int
}

func (pnode *Node) GetFirst() *Arc {
	return pnode.first
}

func (pnode *Node) SetFirst(pfirst *Arc) {
	pnode.first = pfirst
	return
}

func (pnode *Node) GetParent() *Arc {
	return pnode.parent
}

func (pnode *Node) SetParent(pparent *Arc) {
	pnode.parent = pparent
	return
}

func (pnode *Node) GetNext() *Node {
	return pnode.next
}

func (pnode *Node) SetNext(pnext *Node) {
	pnode.next = pnext
	return
}

func (pnode *Node) GetTS() int {
	return pnode.TS
}

func (pnode *Node) SetTS(ts int) {
	pnode.TS = ts
	return
}

type BKGraph struct {
	nodes             []*Node
	arcs              []*Arc
	flow              int
	maxflow_iteration int
	orphanlist        []*Node
	queue_first       *Node
	queue_last        *Node
}
