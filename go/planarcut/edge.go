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
	flags      uint8
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
	p.flags = 0
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
func (e *Edge) SetEdge(tail, head *Vertice, taildual, headdual *Face, caps, rcaps float64) {
	e.head = head
	e.tail = tail
	e.headdual = headdual
	e.taildual = taildual
	e.caps = caps
	e.revcap = rcaps

	if e.tailedgeid < 0 {
		for i := 0; i < tail.GetEdgeNum(); i++ {
			if tail.GetEdge(i) == e {
				e.tailedgeid = i
				break
			}
		}
	}

	if e.headedgeid < 0 {
		for i := 0; i < head.GetEdgeNum(); i++ {
			if head.GetEdge(i) == e {
				e.headedgeid = i
				break
			}
		}
	}

	return
}

func (e *Edge) SetFlags(flag uint8) {
	e.flags = flag
}

func (e *Edge) GetFlags() uint8 {
	return e.flags
}

func (e *Edge) GetHead() *Vertice {
	return e.head
}

func (e *Edge) GetTail() *Vertice {
	return e.tail
}

func (e *Edge) GetHeadDaul() *Face {
	return e.headdual
}

func (e *Edge) GetTailDaul() *Face {
	return e.taildual
}

func (e *Edge) SetTailEdgeId(id int) {
	e.tailedgeid = id
	return
}
func (e *Edge) SetHeadEdgeId(id int) {
	e.headedgeid = id
	return
}

func (e *Edge) GetTailEdgeId() int {
	return e.tailedgeid
}

func (e *Edge) GetHeadEdgeId() int {
	return e.headedgeid
}
