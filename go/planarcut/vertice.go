package main

import (
	"fmt"
	//"log"
)

type Vertice struct {
	edgesccw []*Edge
	faces    []*Face
	numedges int
	x, y     int
	idx      int
	name     string
}

func NewVertice(name string) *Vertice {
	p := &Vertice{}
	p.edgesccw = []*Edge{}
	p.numedges = 0
	p.x = 0
	p.y = 0
	p.idx = 0
	p.name = name
	return p
}

func (vert *Vertice) GetEdge(idx int) *Edge {
	if vert.numedges == 0 {
		return nil
	}

	idx %= vert.numedges
	return vert.edgesccw[idx]
}

func (vert *Vertice) SetCCWEdges(num int, edges []*Edge) {
	vert.numedges = num
	vert.edgesccw = edges

	for i := 0; i < num; i++ {
		if edges[i].GetTail() == vert {
			edges[i].SetTailEdgeId(i)
		} else if edges[i].GetHead() == vert {
			edges[i].SetHeadEdgeId(i)
		}
	}
	return
}

func (vert *Vertice) SetFaces(faces []*Face) {
	vert.faces = faces
	return
}

func (vert *Vertice) GetFaces() []*Face {
	return vert.faces
}

func (vert *Vertice) SetXY(x, y int) {
	vert.x = x
	vert.y = y
	return
}

func (vert *Vertice) SetIdx(idx int) {
	vert.idx = idx
	return
}

func (vert *Vertice) GetIdx() int {
	return vert.idx
}

func (vert *Vertice) GetX() int {
	return vert.x
}

func (vert *Vertice) GetName() string {
	return vert.name
}

func (vert *Vertice) GetY() int {
	return vert.y
}

func (vert *Vertice) GetEdgeNum() int {
	return vert.numedges
}

func (vert *Vertice) GetEdgeId(e *Edge) int {
	if e.GetTail() == vert {
		return e.GetTailEdgeId()
	} else if e.GetHead() == vert {
		return e.GetHeadEdgeId()
	} else {
		return -1
	}
}

func (vert *Vertice) PushEdge(e *Edge) error {
	var err error
	if e == nil {
		err = fmt.Errorf("nil edge error")
		return err
	}

	for _, ce := range vert.edgesccw {
		if ce == e {
			return nil
		}
	}

	vert.edgesccw = append(vert.edgesccw, e)
	vert.numedges++
	return nil
}

func (vert *Vertice) get_east_edge() *Edge {
	var overt *Vertice
	var e *Edge

	for i := 0; i < vert.numedges; i++ {
		e = vert.edgesccw[i]
		if vert == e.GetTail() {
			overt = e.GetHead()
		} else {
			overt = e.GetTail()
		}

		if vert.GetX() < overt.GetX() {
			return e
		}
	}
	return nil
}

func (vert *Vertice) get_west_edge() *Edge {
	var overt *Vertice
	var e *Edge

	for i := 0; i < vert.numedges; i++ {
		e = vert.edgesccw[i]
		if vert == e.GetTail() {
			overt = e.GetHead()
		} else {
			overt = e.GetTail()
		}

		if vert.GetX() > overt.GetX() {
			return e
		}
	}
	return nil
}

func (vert *Vertice) get_south_edge() *Edge {
	var overt *Vertice
	var e *Edge

	for i := 0; i < vert.numedges; i++ {
		e = vert.edgesccw[i]
		if vert == e.GetTail() {
			overt = e.GetHead()
		} else {
			overt = e.GetTail()
		}

		if vert.GetY() < overt.GetY() {
			return e
		}
	}
	return nil
}

func (vert *Vertice) get_north_edge() *Edge {
	var overt *Vertice
	var e *Edge

	for i := 0; i < vert.numedges; i++ {
		e = vert.edgesccw[i]
		if vert == e.GetTail() {
			overt = e.GetHead()
		} else {
			overt = e.GetTail()
		}

		if vert.GetY() > overt.GetY() {
			return e
		}
	}
	return nil
}

func (vert *Vertice) CounterClockWise() {
	var ccw []*Edge
	var e *Edge
	if vert.numedges == 0 || vert.numedges == 1 {
		return
	}

	ccw = []*Edge{}
	e = vert.get_east_edge()
	if e != nil {
		ccw = append(ccw, e)
	}

	e = vert.get_north_edge()
	if e != nil {
		ccw = append(ccw, e)
	}

	e = vert.get_west_edge()
	if e != nil {
		ccw = append(ccw, e)
	}

	e = vert.get_south_edge()
	if e != nil {
		ccw = append(ccw, e)
	}

	vert.edgesccw = ccw
	return
}
