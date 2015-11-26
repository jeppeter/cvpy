package main

import (
	"container/list"
	"fmt"
	"log"
)

type Node struct {
	name    string /*name of the node ,it is the key to search node*/
	parent  *Node  /*this is the route recursive pointer ,for example ,if the new route find ,it will recursive find the route by this */
	next    *Node  /*in the queue */
	TS      int    /*the handle cycle time */
	DIST    int    /*if is_sink == False ,it means the distance to source ,if is_sink == True it means the distance to sink*/
	is_sink bool   /*it means it is in the sink or source map*/
}

var MAXFLOW_TERMINAL *Node
var MAXFLOW_ORPHAN *Node

func init() {
	MAXFLOW_ORPHAN = &Node{}
	MAXFLOW_TERMINAL = &Node{}
	return
}

type NodeMap struct {
	inner map[string]*Node
}

type BKGraph struct {
	nodemap     *NodeMap
	queue_first *Node
	queue_last  *Node
	flows       *StringGraph
	caps        *StringGraph
	neigh       *Neigbour
	source      string
	sink        string
	orphans     *list.List
	TIME        int
	flow        int
}

func NewNode(name string) *Node {
	p := &Node{}
	p.name = name
	p.parent = nil
	p.next = nil
	p.TS = 0
	p.DIST = 0
	p.is_sink = false
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

func (pnode *Node) SetSink(val bool) {
	pnode.is_sink = val
	return
}
func (pnode *Node) GetSink() bool {
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

func (pmap *NodeMap) AddNode_NoError(name string) {
	pnode := pmap.GetNode(name)
	if pnode == nil {
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
	p.orphans = list.New()
	p.TIME = 0
	p.flow = 0
	return p
}

func (graph *BKGraph) InitSource(pnode *Node) error {
	pnode.SetSink(false)
	pnode.SetDist(1)
	pnode.SetParent(MAXFLOW_TERMINAL)
	graph.SetActive(pnode)
	return nil
}

func (graph *BKGraph) InitSink(pnode *Node) error {
	pnode.SetSink(true)
	pnode.SetDist(1)
	pnode.SetParent(MAXFLOW_TERMINAL)
	graph.SetActive(pnode)
	return nil
}

func (graph *BKGraph) SetActive(pnode *Node) error {
	if pnode.GetNext() != nil {
		return nil
	}
	if graph.queue_first == nil {
		pnode.SetNext(nil)
		graph.queue_first = pnode
		graph.queue_last = pnode
		return nil
	}

	/*queue_last should be has*/
	graph.queue_last.SetNext(pnode)
	pnode.SetNext(nil)
	graph.queue_last = pnode
	return nil
}

func (graph *BKGraph) GetActive() *Node {
	if graph.queue_first == nil {
		return nil
	}
	pnode := graph.queue_first
	graph.queue_first = pnode.GetNext()
	pnode.SetNext(nil)
	/*if we have remove all the nodes ,just remove the last*/
	if graph.queue_first == nil {
		graph.queue_last = nil
	}
	return pnode
}

/**********************************************
*function :
*         to init the BKGraph inner structure
*         by caps and neighbour
**********************************************/
func (graph *BKGraph) InitGraph(caps *StringGraph, neighbour *Neigbour, source string, sink string) error {
	/*because the two dimensions are not all the same*/
	for _, k1 := range caps.Iter() {
		graph.nodemap.AddNode_NoError(k1)
		for _, k2 := range caps.IterIdx(k1) {
			graph.nodemap.AddNode_NoError(k2)
		}
	}

	/*now search for the source and sink*/
	pnode := graph.nodemap.GetNode(source)
	if pnode == nil {
		return fmt.Errorf("can not find (%s) source in graph", source)
	}
	graph.InitSource(pnode)

	pnode = graph.nodemap.GetNode(sink)
	if pnode == nil {
		return fmt.Errorf("can not find (%s) sink in graph", sink)
	}
	graph.InitSink(pnode)
	graph.caps = caps
	graph.neigh = neighbour
	graph.source = source
	graph.sink = sink
	return nil
}

func (graph *BKGraph) GetOrphan() *Node {
	if graph.orphans.Len() == 0 {
		return nil
	}
	lv := graph.orphans.Front()
	graph.orphans.Remove(lv)
	return lv.Value.(*Node)
}

func (graph *BKGraph) PushOrphanFront(pnode *Node) int {
	graph.orphans.PushFront(pnode)
	return graph.orphans.Len()
}

func (graph *BKGraph) PushOrphanBack(pnode *Node) int {
	graph.orphans.PushBack(pnode)
	return graph.orphans.Len()
}

func (graph *BKGraph) ProcessSourceOrphan(orphan *Node) {
	return
}

func (graph *BKGraph) ProcessSinkOrphan(orphan *Node) {
	return
}

func (graph *BKGraph) Augment(srcnode *Node, sinknode *Node) int {
	var orphans int
	var bottlecap, curval int
	var curparent, curchld *Node
	orphans = 0

	bottlecap = (graph.caps.GetValue(srcnode.GetName(), sinknode.GetName()) - graph.flows.GetValue(srcnode.GetName(), sinknode.GetName()))

	/*now we should check source side*/
	curparent = srcnode.GetParent()
	curchld = srcnode
	for {
		if curparent == MAXFLOW_TERMINAL {
			break
		} else if curparent == nil {
			log.Printf("%s node parent is nil", curchld.GetName())
		} else if curparent == MAXFLOW_ORPHAN {
			log.Printf("%s node parent is orphan", curchld.GetName())
		}

		curval = graph.caps.GetValue(curparent.GetName(), curchld.GetName()) - graph.flows.GetValue(curparent.GetName(), curchld.GetName())
		if curval < bottlecap {
			bottlecap = curval
		}
		curchld = curparent
		curparent = curchld.GetParent()
	}

	/*now we should check sink side*/
	curchld = sinknode
	curparent = curchld.GetParent()
	for {
		if curparent == MAXFLOW_TERMINAL {
			break
		} else if curparent == nil {
			log.Printf("%s node parent is nil", curchld.GetName())
		} else if curparent == MAXFLOW_ORPHAN {
			log.Printf("%s node parent is orphan", curchld.GetName())
		}

		curval = graph.caps.GetValue(curchld.GetName(), curparent.GetName()) - graph.flows.GetValue(curchld.GetName(), curparent.GetName())
		if curval < bottlecap {
			bottlecap = curval
		}
		curchld = curparent
		curparent = curchld.GetParent()
	}

	/*now we get the bottle cap ,and add it to the flow*/
	curval = graph.flows.GetValue(srcnode.GetName(), sinknode.GetName())
	graph.flows.SetValue(srcnode.GetName(), sinknode.GetName(), curval+bottlecap)
	if graph.flows.GetValue(srcnode.GetName(), sinknode.GetName()) == graph.caps.GetValue(srcnode.GetName(), sinknode.GetName()) {
		graph.PushOrphanFront(srcnode)
		graph.PushOrphanFront(sinknode)
	}

	/*for source side add flow*/
	curchld = srcnode
	curparent = curchld.GetParent()
	for {
		if curparent == MAXFLOW_TERMINAL {
			break
		} else if curparent == nil {
			log.Printf("%s node parent is nil", curchld.GetName())
		} else if curparent == MAXFLOW_ORPHAN {
			log.Printf("%s node parent is orphan", curchld.GetName())
		}

		curval = graph.flows.GetValue(curparent.GetName(), curchld.GetName())
		graph.flows.SetValue(curparent.GetName(), curchld.GetName(), curval+bottlecap)
		if graph.flows.GetValue(curparent.GetName(), curchld.GetName()) == graph.caps.GetValue(curparent.GetName(), curchld.GetName()) {
			graph.PushOrphanFront(curchld)
		}

		curchld = curparent
		curparent = curchld.GetParent()
	}

	/*for sink side add flow*/
	curchld = sinknode
	curparent = curchld.GetParent()
	for {
		if curparent == MAXFLOW_TERMINAL {
			break
		} else if curparent == nil {
			log.Printf("%s node parent is nil", curchld.GetName())
		} else if curparent == MAXFLOW_ORPHAN {
			log.Printf("%s node parent is orphan", curchld.GetName())
		}

		curval = graph.flows.GetValue(curchld.GetName(), curparent.GetName())
		graph.flows.SetValue(curchld.GetName(), curparent.GetName(), curval+bottlecap)

		if graph.flows.GetValue(curchld.GetName(), curparent.GetName()) == graph.caps.GetValue(curchld.GetName(), curparent.GetName()) {
			graph.PushOrphanFront(curchld)
		}

		curchld = curparent
		curparent = curchld.GetParent()
	}

	/*now flows to add*/
	graph.flow += bottlecap

	return orphans
}

func (graph *BKGraph) MaxFlow() (flow int, err error) {
	var curnode *Node
	var curgetnode *Node
	var srcnode *Node
	var sinknode *Node
	var lnode *Node /*link node*/

	curnode = nil
	curgetnode = nil
	for {
		srcnode = nil
		sinknode = nil
		curnode = curgetnode
		if curnode != nil {
			curnode.SetNext(nil)
			if curnode.GetParent() == nil {
				/*if we do not has any upstream to regress ,so it make nothing to do*/
				curnode = nil
			}
		}

		if curnode == nil {
			curnode = graph.GetActive()
			if curnode == nil {
				/*no node is active ,so it is over*/
				break
			}
		}
		if !curnode.GetSink() {
			/*it is source code*/
			for _, lname := range graph.neigh.GetValue(curnode.GetName()) {
				/*to search for the neighbour node to satisfy connection from source to sink*/
				lnode = graph.nodemap.GetNode(lname)
				if lnode == nil {
					log.Printf("can not find lnode (%s) for (%s)", lname, curnode.GetName())
					continue
				}
				curname := curnode.GetName()
				/*it is some flow for the over*/
				if (graph.caps.GetValue(curname, lname)-graph.flows.GetValue(curname, lname)) > 0 ||
					(graph.caps.GetValue(lname, curname)-graph.flows.GetValue(lname, curname)) > 0 {
					if lnode.GetParent() == nil {
						/*it means it is new node in handling ,or it is the last orphan node ,so make it as source node*/
						lnode.SetSink(false)
						lnode.SetParent(curnode)
						lnode.SetTS(curnode.GetTS())
						lnode.SetDist(curnode.GetDist() + 1)
						graph.SetActive(lnode)
					} else if lnode.GetSink() {
						/*we find a route from source to sink ,so break to handle this*/
						srcnode = curnode
						sinknode = lnode
						break
					} else if lnode.GetTS() <= curnode.GetTS() && lnode.GetDist() > curnode.GetDist() {
						/*********************************************
						  it means it handles early before curnode and
						  distance to source is more than cur node
						  so we make it parent to curnode
						  we do not put it in active ,because for it may be not used
						**********************************************/
						lnode.SetParent(curnode)
						lnode.SetTS(curnode.GetTS())
						lnode.SetDist(curnode.GetDist())
					}
				}
			}

		} else {
			/*it is sink code*/
			for _, lname := range graph.neigh.GetValue(curnode.GetName()) {
				lnode = graph.nodemap.GetNode(lname)
				if lnode == nil {
					log.Printf("can not get lnode(%s) for %s", lname, curnode.GetName())
					continue
				}
				curname := curnode.GetName()
				if (graph.caps.GetValue(curname, lname)-graph.flows.GetValue(curname, lname)) > 0 ||
					(graph.caps.GetValue(lname, curname)-graph.flows.GetValue(lname, curname)) > 0 {
					if lnode.GetParent() == nil {
						/*new node just to add for the active sink side*/
						lnode.SetSink(true)
						lnode.SetParent(curnode)
						lnode.SetDist(curnode.GetDist() + 1)
						lnode.SetTS(curnode.GetTS())
						graph.SetActive(lnode)
					} else if !lnode.GetSink() {
						/*it is source side node*/
						srcnode = lnode
						sinknode = curnode
						break
					} else if lnode.GetTS() <= curnode.GetTS() &&
						lnode.GetDist() > curnode.GetDist() {
						lnode.SetParent(curnode)
						lnode.SetTS(curnode.GetTS())
						lnode.SetDist(curnode.GetDist() + 1)
					}
				}
			}

		}

		graph.TIME += 1

		if srcnode != nil && sinknode != nil {
			curnode.SetNext(curnode)
			curgetnode = curnode

			orph := graph.Augment(srcnode, sinknode)
			if orph > 0 {
				for {
					orphan := graph.GetOrphan()
					if orphan == nil {
						break
					}
					if orphan.GetSink() {
						graph.ProcessSinkOrphan(orphan)
					} else {
						graph.ProcessSourceOrphan(orphan)
					}
				}

			}
		} else {
			curgetnode = nil
		}

	}

	return graph.flow, nil
}
