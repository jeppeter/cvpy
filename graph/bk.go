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

const MAXFLOW_INFINITE_D = (1 << 31)

func init() {
	MAXFLOW_ORPHAN = NewNode("MAXFLOW_ORPHAN")
	MAXFLOW_TERMINAL = NewNode("MAXFLOW_TERMINAL")
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
	p.caps = nil
	p.neigh = nil
	p.source = ""
	p.sink = ""
	p.TIME = 0
	p.flow = 0
	return p
}

func (graph *BKGraph) InitSource(pnode *Node) error {
	pnode.SetSink(false)
	pnode.SetDist(0)
	pnode.SetParent(MAXFLOW_TERMINAL)
	graph.SetActive(pnode)
	return nil
}
func (graph *BKGraph) AddSourceNode(pnode *Node, parent *Node) error {
	pnode.SetSink(false)
	pnode.SetDist(parent.GetDist() + 1)
	pnode.SetParent(parent)
	graph.SetActive(pnode)
	return nil
}

func (graph *BKGraph) InitSink(pnode *Node) error {
	pnode.SetSink(true)
	pnode.SetDist(0)
	pnode.SetParent(MAXFLOW_TERMINAL)
	graph.SetActive(pnode)
	return nil
}

func (graph *BKGraph) AddSinkNode(pnode *Node, parent *Node) error {
	pnode.SetSink(true)
	pnode.SetDist(parent.GetDist() + 1)
	pnode.SetParent(parent)
	graph.SetActive(pnode)
	return nil
}

func (graph *BKGraph) SetActive(pnode *Node) {
	if pnode.GetNext() != nil {
		return
	}
	log.Printf("Set (%s) active", pnode.GetName())
	if graph.queue_first == nil {
		pnode.SetNext(nil)
		graph.queue_first = pnode
		graph.queue_last = pnode
		return
	}

	/*queue_last should be has*/
	graph.queue_last.SetNext(pnode)
	pnode.SetNext(nil)
	graph.queue_last = pnode
	return
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
	var pnode *Node
	/*because the two dimensions are not all the same*/
	for _, k1 := range neighbour.Iter() {
		graph.nodemap.AddNode_NoError(k1)
		for _, k2 := range neighbour.GetValue(k1) {
			graph.nodemap.AddNode_NoError(k2)
		}
	}

	/*now search for the source and sink*/
	psrcnode := graph.nodemap.GetNode(source)
	if psrcnode == nil {
		return fmt.Errorf("can not find (%s) source in graph", source)
	}
	graph.InitSource(psrcnode)
	for _, nname := range neighbour.GetValue(source) {
		if nname != sink {
			pnode = graph.nodemap.GetNode(nname)
			graph.AddSourceNode(pnode, psrcnode)
		}
	}

	psinknode := graph.nodemap.GetNode(sink)
	if psinknode == nil {
		return fmt.Errorf("can not find (%s) sink in graph", sink)
	}
	graph.InitSink(psinknode)
	for _, nname := range neighbour.GetValue(sink) {
		if nname != source {
			pnode = graph.nodemap.GetNode(nname)
			graph.AddSinkNode(pnode, psinknode)
		}
	}

	graph.caps = caps
	graph.neigh = neighbour
	graph.source = source
	graph.sink = sink
	return nil
}

func (graph *BKGraph) GetOrphan() *Node {
	var lv *list.Element
	for {
		if graph.orphans.Len() == 0 {
			return nil
		}
		lv = graph.orphans.Front()
		graph.orphans.Remove(lv)
		if lv.Value.(*Node).GetParent() == MAXFLOW_ORPHAN {
			/*it is the orphan we have pushed*/
			break
		}
	}
	return lv.Value.(*Node)
}

func (graph *BKGraph) PushOrphanFront(pnode *Node) int {
	if pnode.GetParent() == MAXFLOW_ORPHAN {
		/*it is already pushed in*/
		return graph.orphans.Len()
	}
	log.Printf("push (%s) front", pnode.GetName())
	pnode.SetParent(MAXFLOW_ORPHAN)
	graph.orphans.PushFront(pnode)
	return graph.orphans.Len()
}

func (graph *BKGraph) PushOrphanBack(pnode *Node) int {
	if pnode.GetParent() == MAXFLOW_ORPHAN {
		/*it is already pushed in*/
		return graph.orphans.Len()
	}
	pnode.SetParent(MAXFLOW_ORPHAN)
	graph.orphans.PushBack(pnode)
	return graph.orphans.Len()
}

func (graph *BKGraph) CanFlow(from string, to string) bool {
	if (graph.caps.GetValue(from, to) - graph.flows.GetValue(from, to)) > 0 {
		return true
	}
	return false
}

func GetNodeName(pnode *Node) string {
	if pnode == nil {
		return "NULL"
	}
	return pnode.GetName()
}

func (graph *BKGraph) ProcessSourceOrphan(orphan *Node) {
	var newparent *Node
	var curnode, curparent, nearparent *Node
	var curd, dmin int
	dmin = MAXFLOW_INFINITE_D
	newparent = nil
	for _, curname := range graph.neigh.GetValue(orphan.GetName()) {
		curnode = graph.nodemap.GetNode(curname)
		if curnode == nil {
			log.Fatalf("%s not find node", curname)
		}
		log.Printf("orphan (%s) curnode (%s)", orphan.GetName(), curname)
		/*we search for all the flow to this orphan node*/
		if graph.CanFlow(curname, orphan.GetName()) {
			if !curnode.GetSink() {
				curparent = curnode.GetParent()
				nearparent = curnode
				curd = 0
				for curparent != nil {
					log.Printf("node[%s].TS (%d) TIME (%d)", curnode.GetName(), curnode.GetTS(), graph.TIME)
					if curnode.GetTS() == graph.TIME {
						curd += curnode.GetDist()
					}
					curd++
					if curparent == MAXFLOW_TERMINAL {
						/*we have find the terminal ,it is the source*/
						curparent.SetTS(graph.TIME)
						curparent.SetDist(1)
						break
					} else if curparent == MAXFLOW_ORPHAN {
						/*it is orphan for parent ,so set it as the not reachable*/
						curd = MAXFLOW_INFINITE_D
						break
					}
					curnode = curparent
					curparent = curnode.GetParent()
				}

				if curd < MAXFLOW_INFINITE_D {
					if curd < dmin {
						/*if we find a new part path from terminal of source to this orphan code ,just get it*/
						newparent = nearparent
						dmin = curd
					}
					curparent = nearparent
					log.Printf("curparent %s curd %d", GetNodeName(curparent), curd)
					for curparent.GetTS() != graph.TIME {
						/*it is not the current cycle do this*/
						curparent.SetTS(graph.TIME)
						curparent.SetDist(curd)
						curd--
						curnode = curparent
						curparent = curnode.GetParent()
						log.Printf("curnode (%s)curparent %s", curnode.GetName(), GetNodeName(curparent))

						if curparent == nil || curparent == MAXFLOW_TERMINAL || curparent == MAXFLOW_ORPHAN {
							break
						}
					}
				}
			}
		}
	}

	log.Printf("set (%s) Parent (%s) ", orphan.GetName(), GetNodeName(newparent))
	orphan.SetParent(newparent)
	if newparent != nil {
		/*we find the new parent for the orphan ,so we should give the TS and DIST*/
		orphan.SetTS(graph.TIME)
		orphan.SetDist(dmin + 1)
	} else {
		for _, curname := range graph.neigh.GetValue(orphan.GetName()) {
			/*get neighbour node */
			curnode = graph.nodemap.GetNode(curname)
			if curnode == nil {
				log.Fatalf("can not find (%s) node", curname)
			}
			if !curnode.GetSink() {
				/*it is source node*/
				curparent = curnode.GetParent()
				if curparent != nil && curparent != MAXFLOW_TERMINAL && curparent != MAXFLOW_ORPHAN {
					if graph.CanFlow(curparent.GetName(), curnode.GetName()) {
						/*can make flow of the parent ,so add it to the active */
						graph.SetActive(curparent)
					}

					if curparent != MAXFLOW_ORPHAN && curparent != MAXFLOW_TERMINAL && curparent == orphan {
						/*to push the child node into orphan when curnode*/
						graph.PushOrphanBack(curnode)
					}
				}
			}
		}
	}
	return
}

func (graph *BKGraph) ProcessSinkOrphan(orphan *Node) {
	var newparent, nearparent, curparent, curnode *Node
	var dmin, curd int
	newparent = nil
	nearparent = nil
	dmin = MAXFLOW_INFINITE_D

	for _, curname := range graph.neigh.GetValue(orphan.GetName()) {
		curnode = graph.nodemap.GetNode(curname)
		if curnode == nil {
			log.Fatalf("can not find %s neighbour %s node", orphan.GetName(), curname)
		}
		if graph.CanFlow(orphan.GetName(), curnode.GetName()) {
			nearparent = curnode
			if curnode.GetSink() {
				/*it is sink node ,so we can do this*/
				curparent = curnode.GetParent()
				curd = 0
				for curparent != nil {
					if curnode.GetTS() == graph.TIME {
						curd += curnode.GetDist()
					}

					curd++
					if curparent == MAXFLOW_TERMINAL {
						curnode.SetTS(graph.TIME)
						curnode.SetDist(1)
						break
					} else if curparent == MAXFLOW_ORPHAN {
						/*it is orphan also*/
						curd = MAXFLOW_INFINITE_D
						break
					} else if curparent == nil {
						/*this is because the sink side not connected by sink*/
						curd = MAXFLOW_INFINITE_D
						break
					}
					curnode = curparent
					curparent = curnode.GetParent()
				}

				if curd < MAXFLOW_INFINITE_D {
					if curd < dmin {
						newparent = nearparent
						dmin = curd
					}
					curparent = curnode.GetParent()
					for curparent.GetTS() != graph.TIME {
						curparent.SetTS(graph.TIME)
						curparent.SetDist(curd)
						curd--
						curnode = curparent
						curparent = curnode.GetParent()
						if curparent == MAXFLOW_TERMINAL || curparent == MAXFLOW_ORPHAN || curparent == nil {
							break
						}
					}
				}
			}
		}
	}

	orphan.SetParent(newparent)
	if newparent != nil {
		orphan.SetTS(graph.TIME)
		orphan.SetDist(dmin + 1)
	} else {
		for _, curname := range graph.neigh.GetValue(orphan.GetName()) {
			curnode = graph.nodemap.GetNode(curname)
			if curnode == nil {
				log.Fatalf("can not find (%s) neighbour (%s) node", orphan.GetName(), curname)
			}

			if curnode.GetSink() {
				/*on the sink side to scan*/
				curparent = curnode.GetParent()
				if curparent != nil && curparent != MAXFLOW_ORPHAN && curparent != MAXFLOW_TERMINAL {
					if graph.CanFlow(curnode.GetName(), curparent.GetName()) {
						/*if have something to flow on the node ,just add it to the active*/
						graph.SetActive(curnode)
					}

					if curparent == orphan {
						graph.PushOrphanBack(curnode)
					}
				}
			}

		}
	}

	return
}

func (graph *BKGraph) GetFlow(from string, to string) int {
	return graph.caps.GetValue(from, to) - graph.flows.GetValue(from, to)
}

func (graph *BKGraph) AddFlow(from string, to string, addval int) bool {
	curval := graph.flows.GetValue(from, to)
	graph.flows.SetValue(from, to, curval+addval)
	if graph.flows.GetValue(from, to) > graph.caps.GetValue(from, to) {
		log.Fatalf("%s -> %s flow add %d error", from, to)
	}
	return graph.CanFlow(from, to)
}

func (graph *BKGraph) Augment(srcnode *Node, sinknode *Node) int {
	var orphans int
	var bottlecap, curval int
	var curparent, curchld *Node
	orphans = 0

	bottlecap = graph.GetFlow(srcnode.GetName(), sinknode.GetName())

	/*now we should check source side*/
	curparent = srcnode.GetParent()
	curchld = srcnode
	for {
		log.Printf("curchld (%s) curparent (%s)", curchld.GetName(), GetNodeName(curparent))
		if curparent == MAXFLOW_TERMINAL {
			break
		} else if curparent == nil {
			log.Fatalf("%s node parent is nil", curchld.GetName())
		} else if curparent == MAXFLOW_ORPHAN {
			log.Fatalf("%s node parent is orphan", curchld.GetName())
		}

		curval = graph.GetFlow(curparent.GetName(), curchld.GetName())
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
			log.Fatalf("%s node parent is nil", curchld.GetName())
		} else if curparent == MAXFLOW_ORPHAN {
			log.Fatalf("%s node parent is orphan", curchld.GetName())
		}

		curval = graph.GetFlow(curchld.GetName(), curparent.GetName())
		if curval < bottlecap {
			bottlecap = curval
		}
		curchld = curparent
		curparent = curchld.GetParent()
	}

	/*now we get the bottle cap ,and add it to the flow*/
	graph.AddFlow(srcnode.GetName(), sinknode.GetName(), bottlecap)

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

		if !graph.AddFlow(curparent.GetName(), curchld.GetName(), bottlecap) {
			graph.PushOrphanFront(curchld)
			orphans++
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

		if !graph.AddFlow(curchld.GetName(), curparent.GetName(), bottlecap) {
			graph.PushOrphanFront(curchld)
			orphans++
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
		log.Printf("curnode (%s)", GetNodeName(curnode))
		if !curnode.GetSink() {
			/*it is source node side*/
			for _, lname := range graph.neigh.GetValue(curnode.GetName()) {
				/*to search for the neighbour node to satisfy connection from source to sink*/
				lnode = graph.nodemap.GetNode(lname)
				if lnode == nil {
					log.Fatalf("can not find lnode (%s) for (%s)", lname, curnode.GetName())
				}
				curname := curnode.GetName()
				log.Printf("curname(%s) lname(%s)", curname, lname)
				/*it is some flow for the over*/
				if graph.CanFlow(curname, lname) {
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
						lnode.SetDist(curnode.GetDist() + 1)
					}
				}
			}

		} else {
			/*it is sink code*/
			for _, lname := range graph.neigh.GetValue(curnode.GetName()) {
				lnode = graph.nodemap.GetNode(lname)
				if lnode == nil {
					log.Fatalf("can not get lnode(%s) for %s", lname, curnode.GetName())
				}
				curname := curnode.GetName()
				if graph.CanFlow(lname, curname) {
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
