package main


type Vertice struct {
	edgesccw []*Edge
	numedges int
	x, y     int
}

func NewVertice() *Vertice {
	p := &Vertice{}
	p.edgesccw = []*Edge{}
	p.numedges = 0
	p.x = p.y = 0
	return p
}

func (vert *Vertice) GetEdge(idx int) *Edge {
	if vert.numedges == 0 {
		return nil
	}

	idx %= vert.numedges
	return vert.edgesccw[idx]
}

func (vert *Vertice)SetCCWEdges(num int ,edges []*Edge) {
	vert.numedges = num
	vert.edgesccw = edges

	for i := 0;i <num ;i ++{
		if edges[i].GetTail() == vert{
			edges[i].SetTailEdgeId(i)
		} else if edges[i].GetHead() == vert {
			edges[i].SetHeadEdgeId(i)
		}
	}
	return
}

func (vert *Vertice)SetXY( x, y int) {
	vert.x = x
	vert.y = y
	return 
}

func (vert *Vertice)GetX() int {
	return vert.x
}

func (vert *Vertice)GetY() int{
	return vert.y
}

func (vert *Vertice)GetEdgeNum()int {
	return vert.numedges
}
