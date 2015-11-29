package main

import (
	"fmt"
)

type Arc struct {
	name   string
	head   *Node
	next   *Arc
	sister *Arc
	r_cap  int
}

func NewArc() *Arc {
	p := &Arc{}
	p.name = ""
	p.head = nil
	p.next = nil
	p.sister = nil
	p.r_cap = 0
	return p
}

func (parc *Arc) SetName(name string) {
	parc.name = name
	return
}
func (parc *Arc) GetName() string {
	return parc.name
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
	name    string
	first   *Arc
	parent  *Arc
	next    *Node
	TS      int
	DIST    int
	is_sink bool
	tr_cap  int
}

func NewNode(name string) *Node {
	p := &Node{}
	p.name = name
	p.first = nil
	p.parent = nil
	p.next = nil
	p.TS = 0
	p.DIST = 0
	p.is_sink = false
	p.tr_cap = 0
	return p
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

func (pnode *Node) GetDIST() int {
	return pnode.DIST
}

func (pnode *Node) SetDIST(dist int) {
	pnode.DIST = dist
	return
}

func (pnode *Node) IsSink() bool {
	return pnode.is_sink
}

func (pnode *Node) SetSink(sink bool) {
	pnode.is_sink = sink
	return
}

func (pnode *Node) GetCap() int {
	return pnode.tr_cap
}

func (pnode *Node) SetCap(caps int) {
	pnode.tr_cap = caps
	return
}

type BKGraph struct {
	nodes       map[string]*Node
	arcs        map[string]*Arc
	flow        int
	orphanlist  []*Node
	queue_first *Node
	queue_last  *Node
}

func NewBKGraph() *BKGraph {
	p := &BKGraph{}
	p.nodes = make(map[string]*Node)
	p.arcs = make(map[string]*Arc)
	p.flow = 0
	p.orphanlist = []*Node{}
	p.queue_first = nil
	p.queue_last = nil
	return p
}

func (graph *BKGraph) add_tweights(nodename string, cap_source, cap_sink int) {
	var pnode *Node
	var ok bool
	pnode, ok = graph.nodes[nodename]
	if !ok {
		pnode = NewNode(nodename)
		graph.nodes[nodename] = pnode
	}

}
