package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
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

func DebugLogPrintf(format string, a ...interface{}) {
	if false {
		return
	}
	_, f, l, _ := runtime.Caller(1)
	s := fmt.Sprintf("%s:%d ", f, l)
	s += fmt.Sprintf(format, a...)
	log.Print(s)
	return
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

func (pnode *Node) GetName() string {
	return pnode.name
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

func (graph *BKGraph) FormArcName(from string, to string) string {
	return fmt.Sprintf("%s->%s", from, to)
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
	aarc.SetName(graph.FormArcName(nodeiname, nodejname))
	arevarc.SetNext(pj.GetFirst())
	pj.SetFirst(arevarc)
	arevarc.SetHead(pi)
	arevarc.SetName(graph.FormArcName(nodejname, nodeiname))
	aarc.SetCap(caps)
	arevarc.SetCap(rev_caps)

	graph.arcs[aarc.GetName()] = aarc
	graph.arcs[arevarc.GetName()] = arevarc
	return
}

func (graph *BKGraph) SortNodeArcs(pnode *Node) {
	var arcarray []*Arc
	var arcnames []string
	var curarc *Arc
	var i, j int

	arcarray = []*Arc{}
	arcnames = []string{}

	curarc = pnode.GetFirst()
	for curarc != nil {
		arcarray = append(arcarray, curarc)
		arcnames = append(arcnames, curarc.GetName())
		curarc = curarc.GetNext()
	}

	if len(arcarray) <= 1 {
		return
	}

	for i = 0; i < len(arcarray); i++ {
		for j = (i + 1); j < len(arcarray); j++ {
			if strings.Compare(arcnames[i], arcnames[j]) > 0 {
				tmpnames := arcnames[i]
				arcnames[i] = arcnames[j]
				arcnames[j] = tmpnames
				tmparc := arcarray[i]
				arcarray[i] = arcarray[j]
				arcarray[j] = tmparc
			}
		}
	}

	pnode.SetFirst(arcarray[0])
	for i = 0; i < (len(arcarray) - 1); i++ {
		arcarray[i].SetNext(arcarray[i+1])
	}
	i = len(arcarray)
	arcarray[(i - 1)].SetNext(nil)
	return
}

func (graph *BKGraph) SortNodes() []*Node {
	var retnodes []*Node
	var retnodeidx []int

	retnodeidx = []int{}
	retnodes = []*Node{}

	for _, pnode := range graph.nodes {
		val, _ := strconv.Atoi(pnode.GetName())
		retnodeidx = append(retnodeidx, val)
		retnodes = append(retnodes, pnode)
	}

	for i := 0; i < len(retnodeidx); i++ {
		for j := (i + 1); j < len(retnodeidx); j++ {
			if retnodeidx[i] > retnodeidx[j] {
				tmpidx := retnodeidx[i]
				retnodeidx[i] = retnodeidx[j]
				retnodeidx[j] = tmpidx
				tmpnode := retnodes[i]
				retnodes[i] = retnodes[j]
				retnodes[j] = tmpnode
			}
		}
	}

	return retnodes
}

func (graph *BKGraph) maxflow_init() {
	for _, pnode := range graph.SortNodes() {
		pnode.SetNext(nil)
		pnode.SetTS(graph.TIME)
		graph.SortNodeArcs(pnode)
		if pnode.GetCap() > 0 {
			pnode.SetSink(false)
			pnode.SetParent(MAXFLOW_TERMINAL)
			graph.SetActive(pnode)
			pnode.SetDIST(1)
		} else if pnode.GetCap() < 0 {
			pnode.SetSink(true)
			pnode.SetParent(MAXFLOW_TERMINAL)
			graph.SetActive(pnode)
			pnode.SetDIST(1)
		} else {
			pnode.SetParent(nil)
		}

	}
	return
}

func (graph *BKGraph) InitGraph(caps *StringGraph, neighbour *Neigbour, source string, sink string) error {
	for _, iname := range neighbour.Iter() {
		for _, jname := range neighbour.GetValue(iname) {
			capto := caps.GetValue(iname, jname)
			caprev := caps.GetValue(jname, iname)
			if iname == source {
				DebugLogPrintf("add_tweights (%s,%d,0)", jname, capto)
				graph.add_tweights(jname, capto, 0)

			} else if iname == sink {
				/*nothing to do*/
			} else if jname == source {
				/*nothing to do*/
			} else if jname == sink {
				DebugLogPrintf("add_tweights (%s,0,%d)", iname, capto)
				graph.add_tweights(iname, 0, capto)
			} else {
				//fromarcname := fmt.Sprintf("%s -> %s", iname, jname)
				fromarcname := graph.FormArcName(iname, jname)
				//toarcname := fmt.Sprintf("%s -> %s", jname, iname)
				toarcname := graph.FormArcName(jname, iname)
				_, ok1 := graph.arcs[fromarcname]
				_, ok2 := graph.arcs[toarcname]
				if !ok1 && !ok2 {
					DebugLogPrintf("add_edge (%s,%s,%d,%d)", iname, jname, capto, caprev)
					graph.add_edge(iname, jname, capto, caprev)
				}
			}
		}
	}

	graph.maxflow_init()
	return nil
}

func (graph *BKGraph) GetNodeNames() []string {
	narr := []string{}
	for n, _ := range graph.nodes {
		narr = append(narr, n)
	}
	return narr
}

func (graph *BKGraph) GetArcNames() []string {
	narr := []string{}
	for n, _ := range graph.arcs {
		narr = append(narr, n)
	}
	return narr
}

func (graph *BKGraph) GetNextList(pnode *Node) string {
	s := "["
	pj := pnode.GetNext()
	i := 0
	for pj != nil {
		if i != 0 {
			s += ","
		}
		i++
		s += fmt.Sprintf("%s", pj.GetName())
		if pj == pj.GetNext() {
			break
		}
		pj = pj.GetNext()
	}

	s += fmt.Sprintf("]cnt(%d)", i)

	return s
}

func (graph *BKGraph) GetNodeName(pnode *Node) string {
	if pnode == nil {
		return "NULL"
	}
	return pnode.GetName()
}

func (graph *BKGraph) GetFirstList(pnode *Node) string {
	s := "["
	i := 0
	parc := pnode.GetFirst()
	for parc != nil {
		if i != 0 {
			s += ","
		}
		i++
		s += parc.GetName()
		parc = parc.GetNext()
	}
	s += fmt.Sprintf("]cnt(%d)", i)

	return s
}

func (graph *BKGraph) GetParentList(pnode *Node) string {
	s := "["
	i := 0
	parc := pnode.GetParent()
	for parc != nil {
		if i != 0 {
			s += ","
		}
		if parc == MAXFLOW_ORPHAN {
			s += "MAXFLOW_ORPHAN"
			break
		}
		if parc == MAXFLOW_TERMINAL {
			s += "MAXFLOW_TERMINAL"
			break
		}
		i++
		pj := parc.GetHead()
		s += fmt.Sprintf("%s(%s)", graph.GetNodeName(pj), parc.GetName())
		if pj == nil {
			break
		}
		parc = pj.GetParent()
	}
	s += fmt.Sprintf("]cnt(%d)", i)
	return s
}

func (graph *BKGraph) DebugNode(pnode *Node) {
	DebugLogPrintf("==============================")
	if pnode.IsSink() {
		DebugLogPrintf("node[%s].is_sink (True)", pnode.GetName())
	} else {
		DebugLogPrintf("node[%s].is_sink (False)", pnode.GetName())
	}
	DebugLogPrintf("node[%s].arc_first list(%s)", pnode.GetName(), graph.GetFirstList(pnode))
	DebugLogPrintf("node[%s].arc_parent list(%s)", pnode.GetName(), graph.GetParentList(pnode))
	DebugLogPrintf("node[%s].node_next list(%s)", pnode.GetName(), graph.GetNextList(pnode))
	DebugLogPrintf("node[%s].tr_cap (%d)", pnode.GetName(), pnode.GetCap())
	DebugLogPrintf("node[%s].TS (%d) node[%s].DIST (%d)", pnode.GetName(), pnode.GetTS(), pnode.GetName(), pnode.GetDIST())
	DebugLogPrintf("******************************")
	return
}

func (graph *BKGraph) GetArcName(parc *Arc) string {
	if parc == nil {
		return "NULL"
	}
	return parc.GetName()
}

func (graph *BKGraph) GetArcNextList(parc *Arc) string {
	s := "["
	i := 0
	pnext := parc.GetNext()
	for pnext != nil {
		if i != 0 {
			s += ","
		}
		i++
		s += pnext.GetName()
		pnext = pnext.GetNext()
	}
	s += fmt.Sprintf("]cnt(%d)", i)
	return s
}

func (graph *BKGraph) DebugArc(parc *Arc) {
	DebugLogPrintf("+++++++++++++++++++++++++++++++")
	DebugLogPrintf("arc[%s].node_head (%s)", parc.GetName(), graph.GetNodeName(parc.GetHead()))
	DebugLogPrintf("arc[%s].arc_next list(%s)", parc.GetName(), graph.GetArcNextList(parc))
	DebugLogPrintf("arc[%s].arc_sister (%s)", parc.GetName(), parc.GetSister().GetName())
	DebugLogPrintf("arc[%s].r_cap (%d)", parc.GetName(), parc.GetCap())
	DebugLogPrintf("-------------------------------")
	return
}

func (graph *BKGraph) GetQueueFirst() string {
	s := "["
	i := 0
	pnode := graph.queue_first
	for pnode != nil {
		if i != 0 {
			s += ","
		}
		i++
		s += pnode.GetName()
		if pnode == pnode.GetNext() {
			break
		}
		pnode = pnode.GetNext()
	}
	s += fmt.Sprintf("]cnt(%d)", i)
	return s
}

func (graph *BKGraph) GetOrphanString() string {
	s := fmt.Sprintf("cnt(%d)[", graph.orphanlist.Len())
	i := 0
	for curlist := graph.orphanlist.Front(); curlist != nil; curlist = curlist.Next() {
		pnode := curlist.Value.(*Node)
		if i != 0 {
			s += ","
		}
		i++
		s += pnode.GetName()
	}
	s += fmt.Sprintf("]")
	return s
}

func (graph *BKGraph) DebugState(notice string) {
	DebugLogPrintf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	DebugLogPrintf("%s", notice)

	for _, curn := range SortArrayString(graph.GetNodeNames()) {
		pnode, _ := graph.nodes[curn]
		graph.DebugNode(pnode)
	}

	for _, curn := range SortArrayString(graph.GetArcNames()) {
		parc, _ := graph.arcs[curn]
		graph.DebugArc(parc)
	}

	DebugLogPrintf("queue_first list(%s)", graph.GetQueueFirst())
	DebugLogPrintf("TIME (%d) flow (%d)", graph.TIME, graph.flow)
	DebugLogPrintf("orphan_list (%s)", graph.GetOrphanString())
	DebugLogPrintf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	return
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
		pnode.SetNext(pnode)
	}
	return
}

func (graph *BKGraph) GetActive() *Node {
	var pnode *Node
	pnode = nil
	for {
		pnode = graph.queue_first
		if pnode == nil {
			return nil
		}
		if pnode.GetNext() == pnode {
			graph.queue_first = nil
			graph.queue_last = nil
		} else {
			graph.queue_first = pnode.GetNext()
		}
		pnode.SetNext(nil)
		if pnode.GetParent() != nil {
			return pnode
		}
	}
	return nil
}

func (graph *BKGraph) Augment(parc *Arc) {
	var pi *Node
	bottlecap := parc.GetCap()
	DebugLogPrintf("Get arc (%s)", parc.GetName())

	/*for source side*/
	pi = parc.GetSister().GetHead()
	for {
		pcurarc := pi.GetParent()
		if pcurarc == MAXFLOW_TERMINAL {
			break
		}
		DebugLogPrintf("curarc (%s)", graph.GetArcName(pcurarc))
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
	graph.DebugState(fmt.Sprintf("debug state after init"))

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
		graph.DebugState(fmt.Sprintf("debug state after arcs handle(%d)", graph.TIME))

		if gotarc != nil {
			curnode.SetNext(curnode)
			curgetnode = curnode

			graph.Augment(gotarc)
			graph.DebugState(fmt.Sprintf("debug state after augment(%d)", graph.TIME))

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

			graph.DebugState(fmt.Sprintf("debug state after orphan handle (%d)", graph.TIME))

		} else {
			curgetnode = nil
		}

	}

	return graph.flow, nil
}

type CapPair struct {
	cap_source int
	cap_sink   int
}

func BKParseFile(infile string) *BKGraph {
	var source, sink string
	var tweights_pair map[string]*CapPair
	var sarr []string
	var caps int
	var keys []string

	tweights_pair = make(map[string]*CapPair)
	graph := NewBkGraph()
	file, e := os.Open(infile)
	if e != nil {
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		l := scanner.Text()
		l = strings.Trim(l, "\r\n")
		if strings.HasPrefix(l, "#") {
			continue
		}

		if strings.HasPrefix(l, "source=") {
			sarr = strings.Split(l, "=")
			if len(sarr) < 2 {
				continue
			}
			source = sarr[1]
			continue
		}

		if strings.HasPrefix(l, "sink=") {
			sarr = strings.Split(l, "=")
			if len(sarr) < 2 {
				continue
			}
			sink = sarr[1]
			continue
		}

		sarr = strings.Split(l, ",")
		if len(sarr) < 3 {
			continue
		}

		caps, _ = strconv.Atoi(sarr[2])
		curs := sarr[0]
		curt := sarr[1]
		if curs == source && curt != sink {
			_, ok := tweights_pair[curt]
			if !ok {
				tweights_pair[curt] = &CapPair{}
				tweights_pair[curt].cap_source = caps
				tweights_pair[curt].cap_sink = 0
			} else {
				tweights_pair[curt].cap_source = caps
				graph.add_tweights(curt, tweights_pair[curt].cap_source, tweights_pair[curt].cap_sink)
				DebugLogPrintf("g -> add_tweights(%s,%d,%d);", curt, tweights_pair[curt].cap_source, tweights_pair[curt].cap_sink)
				delete(tweights_pair, curt)
			}
		} else if curt == sink && curs != source {
			_, ok := tweights_pair[curs]
			if !ok {
				tweights_pair[curs] = &CapPair{}
				tweights_pair[curs].cap_source = 0
				tweights_pair[curs].cap_sink = caps

			} else {
				tweights_pair[curs].cap_sink = caps
				graph.add_tweights(curs, tweights_pair[curs].cap_source, tweights_pair[curs].cap_sink)
				DebugLogPrintf("g -> add_tweights(%s,%d,%d);", curs, tweights_pair[curs].cap_source, tweights_pair[curs].cap_sink)
				delete(tweights_pair, curs)
			}
		} else {
			graph.add_edge(curs, curt, caps, 0)
			DebugLogPrintf("g -> add_edge(%s,%s,%d,0);", curs, curt, caps)
		}
	}

	keys = []string{}
	for kk := range tweights_pair {
		keys = append(keys, kk)
	}

	for i := 0; i < len(keys); i++ {
		for j := (i + 1); j < len(keys); j++ {
			if strings.Compare(keys[i], keys[j]) > 0 {
				tmp := keys[i]
				keys[i] = keys[j]
				keys[j] = tmp
			}
		}
	}

	for _, kk := range keys {
		graph.add_tweights(kk, tweights_pair[kk].cap_source, tweights_pair[kk].cap_sink)
		DebugLogPrintf("g -> add_tweights(%s,%d,%d);", kk, tweights_pair[kk].cap_source, tweights_pair[kk].cap_sink)
	}

	graph.maxflow_init()
	return graph
}
