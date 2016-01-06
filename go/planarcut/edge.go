package main

type Edge struct {
	caps       float64
	revcap     float64
	tail       *Vertice
	head       *Vertice
	taildual   *Face
	headdual   *Face
	tailedgeid int
	headedgeid int
	flags      uint32
	idx        int
	name       string
}

func NewEdge() *Edge {
	p := &Edge{}
	p.caps = 0.0
	p.revcap = 0.0
	p.tail = nil
	p.head = nil
	p.taildual = nil
	p.headdual = nil
	p.tailedgeid = -1
	p.headedgeid = -1
	p.flags = uint32(0)
	p.idx = 0
	p.name = ""
	return p
}

func (e *Edge) SetCap(caps float64) {
	e.caps = caps
	return
}

func (e *Edge) GetCap() float64 {
	return e.caps
}

func (e *Edge) SetRevCap(revcap float64) {
	e.revcap = revcap
	return
}

func (e *Edge) GetRevCap() float64 {
	return e.revcap
}

func (e *Edge) SetTail(tail *Vertice) {
	e.tail = tail
	return
}

func (e *Edge) SetHead(head *Vertice) {
	e.head = head
	return
}
func (e *Edge) SetHeadDual(headdual *Face) {
	e.headdual = headdual
	return
}

func (e *Edge) SetTailDual(taildual *Face) {
	e.taildual = taildual
	return
}

func (e *Edge) SetName(name string) {
	e.name = name
	return
}

func (e *Edge) GetName() string {
	return e.name
}

func (e *Edge) SetIdx(idx int) {
	e.idx = idx
	return
}

func (e *Edge) GetIdx() int {
	return e.idx
}

func (e *Edge) SetEdge(tail, head *Vertice, taildual, headdual *Face, caps, rcaps float64) {
	e.head = head
	e.tail = tail
	e.headdual = headdual
	e.taildual = taildual
	e.caps = caps
	e.revcap = rcaps

	return
}

func (e *Edge) SetFlags(flag uint32) {
	e.flags = flag
}

func (e *Edge) GetFlags() uint32 {
	return e.flags
}

func (e *Edge) GetHead() *Vertice {
	return e.head
}

func (e *Edge) GetTail() *Vertice {
	return e.tail
}

func (e *Edge) GetHeadDual() *Face {
	return e.headdual
}

func (e *Edge) GetTailDual() *Face {
	return e.taildual
}

func (e *Edge) GetTailEdgeId() int {
	return e.tailedgeid
}
func (e *Edge) GetHeadEdgeId() int {
	return e.headedgeid
}

func (e *Edge) SetTailEdgeId(id int) {
	e.tailedgeid = id
	return
}

func (e *Edge) SetHeadEdgeId(id int) {
	e.headedgeid = id
	return
}
