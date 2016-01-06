package main

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
