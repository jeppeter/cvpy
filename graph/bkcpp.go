package main

import (
	"container/list"
	"fmt"
	//"log"
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

var MAXFLOW_ORPHAN, MAXFLOW_TERMINAL *Arc
var MAXFLOW_INFINITE_D int

func init() {
	MAXFLOW_TERMINAL = NewArc()
	MAXFLOW_TERMINAL.SetName("MAXFLOW_TERMINAL")
	MAXFLOW_ORPHAN = NewArc()
	MAXFLOW_ORPHAN.SetName("MAXFLOW_ORPHAN")
	MAXFLOW_INFINITE_D = (1 << 31)
}

type BKGraph struct {
	nodes       map[string]*Node
	arcs        map[string]*Arc
	flow        int
	TIME        int
	orphanlist  *list.List
	queue_first *Node
	queue_last  *Node
}

func NewBkGraph() *BKGraph {
	p := &BKGraph{}
	p.nodes = make(map[string]*Node)
	p.arcs = make(map[string]*Arc)
	p.flow = 0
	p.TIME = 0
	p.orphanlist = list.New()
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

	delta := pnode.GetCap()
	if delta > 0 {
		cap_source += delta
	} else {
		cap_sink -= delta
	}

	if cap_source < cap_sink {
		graph.flow += cap_source
	} else {
		graph.flow += cap_sink
	}

	pnode.SetCap(cap_source - cap_sink)
	return
}

func (graph *BKGraph) add_edge(nodeiname, nodejname string, caps, rev_caps int) {
	var aarc, arevarc *Arc
	var pi, pj *Node
	var ok bool
	aarc = NewArc()
	arevarc = NewArc()
	pi, ok = graph.nodes[nodeiname]
	if !ok {
		pi = NewNode(nodeiname)
		graph.nodes[nodeiname] = pi
	}

	pj, ok = graph.nodes[nodejname]
	if !ok {
		pj = NewNode(nodejname)
		graph.nodes[nodejname] = pj
	}

	aarc.SetSister(arevarc)
	arevarc.SetSister(aarc)
	aarc.SetNext(pi.GetFirst())
	pi.SetFirst(aarc)
	aarc.SetHead(pj)
	aarc.SetName(fmt.Sprintf("%s -> %s", nodejname, nodeiname))
	pj.SetFirst(arevarc)
	arevarc.SetHead(pi)
	arevarc.SetName(fmt.Sprintf("%s -> %s", nodeiname, nodejname))
	aarc.SetCap(caps)
	arevarc.SetCap(rev_caps)
	graph.arcs[aarc.GetName()] = aarc
	graph.arcs[arevarc.GetName()] = arevarc
	return
}

func (graph *BKGraph) InitGraph(caps *StringGraph, neighbour *Neigbour, source string, sink string) error {
	for _, iname := range neighbour.Iter() {
		for _, jname := range neighbour.GetValue(iname) {
			capto := caps.GetValue(iname, jname)
			caprev := caps.GetValue(jname, iname)
			if iname == source {
				graph.add_tweights(jname, capto, 0)

			} else if jname == sink {
				graph.add_tweights(iname, 0, capto)
			} else {
				graph.add_edge(iname, jname, capto, caprev)
			}
		}
	}

	for _, pnode := range graph.nodes {
		pnode.SetNext(nil)
		pnode.SetTS(graph.TIME)
		if pnode.GetCap() > 0 {
			pnode.SetSink(false)
			pnode.SetParent(MAXFLOW_TERMINAL)
			graph.SetActive(pnode)
		} else if pnode.GetCap() < 0 {
			pnode.SetSink(true)
			pnode.SetParent(MAXFLOW_TERMINAL)
			graph.SetActive(pnode)
		} else {
			pnode.SetParent(nil)
		}

	}

	return nil
}

func (graph *BKGraph) SetActive(pnode *Node) {
	if pnode.GetNext() == nil {
		/*not in the queue or ,just used for in queue,just insert into it*/
		if graph.queue_first == nil {
			graph.queue_first = pnode
			graph.queue_last = pnode
		} else {
			graph.queue_last.SetNext(pnode)
			graph.queue_last = pnode
		}
	}
	return
}

func (graph *BKGraph) GetActive() *Node {
	var pnode *Node
	pnode = nil
	if graph.queue_first != nil {
		pnode = graph.queue_first
		graph.queue_first = pnode.GetNext()
		/*set next for will insert queue again*/
		pnode.SetNext(nil)
		if graph.queue_first == nil {
			/*nothing in the queue ,so we should set nil*/
			graph.queue_last = nil
		}
	}
	return pnode
}

func (graph *BKGraph) Augment(parc *Arc) {
	var pi *Node
	bottlecap := parc.GetCap()

	/*for source side*/
	pi = parc.GetSister().GetHead()
	for {
		pcurarc := pi.GetParent()
		if pcurarc == MAXFLOW_TERMINAL {
			break
		}
		pcursis := pcurarc.GetSister()
		if bottlecap > pcursis.GetCap() {
			bottlecap = pcursis.GetCap()
		}
		pi = pcurarc.GetHead()
	}

	/*for sink side*/
	pi = parc.GetHead()
	for {
		pcurarc := pi.GetParent()
		if pcurarc == MAXFLOW_TERMINAL {
			break
		}

		if bottlecap > pcurarc.GetCap() {
			bottlecap = pcurarc.GetCap()
		}
		pi = pcurarc.GetHead()
	}

	if bottlecap > -pi.GetCap() {
		bottlecap = -pi.GetCap()
	}

	psister := parc.GetSister()
	psister.SetCap(psister.GetCap() + bottlecap)
	parc.SetCap(parc.GetCap() - bottlecap)

	pi = psister.GetHead()
	for {
		pcurarc := pi.GetParent()
		if pcurarc == MAXFLOW_TERMINAL {
			break
		}

		pcursister := pcurarc.GetSister()
		pcurarc.SetCap(pcurarc.GetCap() + bottlecap)
		pcursister.SetCap(pcursister.GetCap() - bottlecap)

		if pcursister.GetCap() == 0 {
			graph.PushOrphanFront(pi)
		}
		pi = pcurarc.GetHead()
	}

	pi.SetCap(pi.GetCap() - bottlecap)
	if pi.GetCap() == 0 {
		graph.PushOrphanFront(pi)
	}

	pi = parc.GetHead()

	for {
		pcurarc := pi.GetParent()
		if pcurarc == MAXFLOW_TERMINAL {
			break
		}
		pcursister := pcurarc.GetSister()
		pcursister.SetCap(pcursister.GetCap() + bottlecap)
		pcurarc.SetCap(pcurarc.GetCap() - bottlecap)

		if pcurarc.GetCap() == 0 {
			graph.PushOrphanFront(pi)
		}

		pi = pcurarc.GetHead()
	}

	pi.SetCap(pi.GetCap() + bottlecap)
	if pi.GetCap() == 0 {
		graph.PushOrphanFront(pi)
	}
	graph.flow += bottlecap
	return
}

func (graph *BKGraph) PushOrphanBack(pnode *Node) {
	if pnode.GetParent() == MAXFLOW_ORPHAN {
		return
	}

	pnode.SetParent(MAXFLOW_ORPHAN)
	graph.orphanlist.PushBack(pnode)
	return
}

func (graph *BKGraph) PushOrphanFront(pnode *Node) {
	if pnode.GetParent() == MAXFLOW_ORPHAN {
		return
	}

	pnode.SetParent(MAXFLOW_ORPHAN)
	graph.orphanlist.PushFront(pnode)
	return
}

func (graph *BKGraph) GetOrphan() *Node {
	var pnode *Node
	pnode = nil
	for {
		if graph.orphanlist.Len() == 0 {
			return nil
		}

		front := graph.orphanlist.Front()
		graph.orphanlist.Remove(front)
		pnode = front.Value.(*Node)
		if pnode.GetParent() == MAXFLOW_ORPHAN {
			return pnode
		}
	}

	return nil
}

func (graph *BKGraph) ProcessSinkOrphan(pnode *Node) {
	var arc0_min *Arc
	dmin := MAXFLOW_INFINITE_D
	arc0_min = nil

	for arc0 := pnode.GetFirst(); arc0 != nil; arc0 = arc0.GetNext() {
		if arc0.GetCap() != 0 {
			pj := arc0.GetHead()
			if pj.IsSink() {
				arca := pj.GetParent()
				if arca != nil {
					d := 0
					for {
						if pj.GetTS() == graph.TIME {
							d += pj.GetDIST()
							break
						}
						arca = pj.GetParent()
						d++
						if arca == MAXFLOW_TERMINAL {
							pj.SetTS(graph.TIME)
							pj.SetDIST(1)
							break
						}
						if arca == MAXFLOW_ORPHAN {
							d = MAXFLOW_INFINITE_D
							break
						}
						pj = arca.GetHead()
					}

					if d < MAXFLOW_INFINITE_D {
						if d < dmin {
							dmin = d
							arc0_min = arc0
						}

						pj = arc0.GetHead()
						for pj.GetTS() != graph.TIME {
							pj.SetTS(graph.TIME)
							pj.SetDIST(d)
							d--
							pj = pj.GetParent().GetHead()
						}
					}
				}
			}
		}
	}

	pnode.SetParent(arc0_min)

	if arc0_min != nil {
		pnode.SetTS(graph.TIME)
		pnode.SetDIST(dmin + 1)
	} else {
		for arc0 := pnode.GetFirst(); arc0 != nil; arc0 = arc0.GetNext() {
			pj := arc0.GetHead()
			if pj.IsSink() {
				arca := pj.GetParent()
				if arca != nil {
					if arc0.GetCap() != 0 {
						graph.SetActive(pj)
					}

					if arca != MAXFLOW_ORPHAN && arca != MAXFLOW_TERMINAL && arca.GetHead() == pnode {
						graph.PushOrphanBack(pj)
					}
				}
			}
		}
	}
}

func (graph *BKGraph) ProcessSourceOrphan(pnode *Node) {
	var arc0_min *Arc
	dmin := MAXFLOW_INFINITE_D
	arc0_min = nil

	for arc0 := pnode.GetFirst(); arc0 != nil; arc0 = arc0.GetNext() {
		arc0sis := arc0.GetSister()
		if arc0sis.GetCap() != 0 {
			pj := arc0.GetHead()
			if !pj.IsSink() {
				arca := pj.GetParent()
				if arca != nil {
					d := 0
					for {
						if pj.GetTS() == graph.TIME {
							d += pj.GetDIST()
							break
						}
						arca = pj.GetParent()
						d++
						if arca == MAXFLOW_TERMINAL {
							pj.SetTS(graph.TIME)
							pj.SetDIST(1)
							break
						}
						if arca == MAXFLOW_ORPHAN {
							d = MAXFLOW_INFINITE_D
							break
						}

						pj = arca.GetHead()
					}
					if d < MAXFLOW_INFINITE_D {
						if d < dmin {
							dmin = d
							arc0_min = arc0
						}
						pj = arc0.GetHead()
						for pj.GetTS() != graph.TIME {
							pj.SetTS(graph.TIME)
							pj.SetDIST(d)
							d--
							pj = pj.GetParent().GetHead()
						}
					}
				}
			}
		}
	}

	pnode.SetParent(arc0_min)
	if arc0_min != nil {
		pnode.SetTS(graph.TIME)
		pnode.SetDIST(dmin + 1)
	} else {
		for arc0 := pnode.GetFirst(); arc0 != nil; arc0 = arc0.GetNext() {
			pj := arc0.GetHead()
			if !pj.IsSink() {
				arca := pj.GetParent()
				if arca != nil {
					psister := arc0.GetSister()
					if psister.GetCap() != 0 {
						graph.SetActive(pj)
					}

					if arca != MAXFLOW_ORPHAN && arca != MAXFLOW_TERMINAL && arca.GetHead() == pnode {
						graph.PushOrphanBack(pj)
					}
				}
			}
		}
	}
	return
}

func (graph *BKGraph) MaxFlow() (flow int, err error) {
	var curnode, curgetnode *Node
	var gotarc *Arc
	curnode = nil
	curgetnode = nil

	for {
		gotarc = nil
		curnode = curgetnode
		if curnode != nil {
			curnode.SetNext(nil)
			if curnode.GetParent() == nil {
				curnode = nil
			}
		}

		if curnode == nil {
			curnode = graph.GetActive()
			if curnode == nil {
				break
			}
		}

		if !curnode.IsSink() {
			/*if not */
			for arc := curnode.GetFirst(); arc != nil; arc = arc.GetNext() {
				if arc.GetCap() != 0 {
					pj := arc.GetHead()
					if pj.GetParent() == nil {
						/*to make for the node as the source side */
						pj.SetSink(false)
						pj.SetParent(arc.GetSister())
						pj.SetTS(curnode.GetTS())
						pj.SetDIST(curnode.GetDIST() + 1)
						graph.SetActive(pj)
					} else if pj.IsSink() {
						gotarc = arc
						break
					} else if pj.GetTS() <= curnode.GetTS() && pj.GetDIST() > curnode.GetDIST() {
						pj.SetParent(arc.GetSister())
						pj.SetTS(curnode.GetTS())
						pj.SetDIST(curnode.GetDIST() + 1)
					}
				}
			}

		} else {
			for arc := curnode.GetFirst(); arc != nil; arc = arc.GetNext() {
				if arc.GetCap() != 0 {
					pj := arc.GetHead()
					if pj.GetParent() == nil {
						/*set for the sink side*/
						pj.SetSink(true)
						pj.SetParent(arc.GetSister())
						pj.SetTS(curnode.GetTS())
						pj.SetDIST(curnode.GetDIST() + 1)
						graph.SetActive(pj)
					} else if !pj.IsSink() {
						gotarc = arc
						break
					} else if pj.GetTS() <= curnode.GetTS() && pj.GetDIST() > curnode.GetDIST() {
						pj.SetParent(arc.GetSister())
						pj.SetTS(curnode.GetTS())
						pj.SetDIST(curnode.GetDIST() + 1)
					}
				}
			}
		}

		graph.TIME++

		if gotarc != nil {
			curnode.SetNext(curnode)
			curgetnode = curnode

			graph.Augment(gotarc)

			for {
				orphan := graph.GetOrphan()
				if orphan == nil {
					break
				}

				if orphan.IsSink() {
					graph.ProcessSinkOrphan(orphan)
				} else {
					graph.ProcessSourceOrphan(orphan)
				}
			}

		} else {
			curgetnode = nil
		}

	}

	return graph.flow, nil
}
