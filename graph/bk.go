package main

import (
	"fmt"
)

type Node struct {
	name    string /*name of the node ,it is the key to search node*/
	parent  *Node  /*this is the route recursive pointer ,for example ,if the new route find ,it will recursive find the route by this */
	next    *Node  /*in the queue */
	TS      int    /*the handle cycle time */
	DIST    int    /*if is_sink == False ,it means the distance to source ,if is_sink == True it means the distance to sink*/
	is_sink bool   /*it means it is in the sink or source map*/
}

type NodeMap struct {
	inner map[string]*Node
}

type BKGraph struct {
	nodemap     *NodeMap
	queue_first *Node
	queue_last  *Node
	flows       *StringGraph
	TIME        int
	flow        int
}

func NewNode(name string) *Node {
	p := &Node{}
	p.name = name
	p.parent = NULL
	p.next = NULL
	p.TS = 0
	p.DIST = 0
	p.is_sink = False
	return p
}

func (pnode *Node) GetTS() int {
	return pnode.TS
}

func (pnode *Node) SetTS(ts int) {
	pnode.TS = ts
}

func (pnode *Node) GetParent() *Node {
	return pnode.parent
}

func (pnode *Node) SetParent(parent *Node) {
	pnode.parent = parent
	return
}

func (pnode *Node) GetDist() int {
	return pnode.DIST
}
func (pnode *Node) SetDist(dist int) {
	pnode.DIST = dist
	return
}

func (pnode *Node) GetNext() *Node {
	return pnode.next
}

func (pnode *Node) SetNext(next *Node) {
	pnode.next = next
	return
}

func (pnode *Node)SetSink( val bool) {
	pnode.is_sink = val
	return
}
func (pnode *Node)GetSink() bool {
	return pnode.is_sink
}

func (pnode *Node) GetName() string {
	return pnode.name
}

func NewNodeMap() *NodeMap {
	p := &NodeMap{}
	p.inner = make(map[string]*Node)
	return p
}

func (pmap *NodeMap) GetNode(name string) *Node {
	pnode, ok := pmap.inner[name]
	if !ok {
		return nil
	}
	return pnode
}

func (pmap *NodeMap) AddNode(name string) error {
	pnode := pmap.GetNode(name)
	if pnode != nil {
		return fmt.Errorf("already exists (%s) node", name)
	}

	pmap.inner[name] = NewNode(name)
	return nil
}

func (pmap *NodeMap)AddNode_NoError(name string) {
	pnode := pmap.GetNode(name)
	if pnode == nil{
		pmap.AddNode(name)
	}
	return
	
}

func NewBkGraph() *BKGraph {
	p := &BKGraph{}
	p.nodemap = NewNodeMap()
	p.queue_first = nil
	p.queue_last = nil
	p.flows = NewStringGraph()
	p.TIME = 0
	p.flow = 0
	return p
}


func (graph *BKGraph)InitSource(pnode *Node) error {
	pnode.SetSink(false)
	pnode.SetDist(1)
	pnode.SetParent(((Node*)))
	
}

/**********************************************
*function :
*         to init the BKGraph inner structure
*         by caps and neighbour
**********************************************/
func (graph *BKGraph) InitGraph(caps *StringGraph, neighbour *Neigbour, source string, sink string) error {
	/*because the two dimensions are not all the same*/
	for _,k1 := caps.Iter(){
		graph.nodemap.AddNode_NoError(k1)
		for _,k2 := caps.IterIdx(k1){
			graph.nodemap.AddNode_NoError(k2)
		}
	}

	/*now search for the source and sink*/
	pnode := graph.nodemap.GetNode(source)
	if pnode == nil{
		return fmt.Errorf("can not find (%s) source in graph",source)
	}

}
