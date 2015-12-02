package main

type DynNode struct {
	parent   *DynNode
	head     *DynNode
	tail     *DynNode
	left     *DynNode
	right    *DynNode
	reserved uint8
	netcost  float64
	netcostR float64
	netmin   float64
	netminR  float64
	height   int
	data     interface{}
}

func NewDynNode() *DynNode {
	p := &DynNode{}
	p.parent = nil
	p.head = nil
	p.tail = nil
	p.left = nil
	p.right = nil
	p.reserved = 0
	p.netcost = float64(0.0)
	p.netcostR = float64(0.0)
	p.netmin = float64(0.0)
	p.netminR = float64(0.0)
	p.height = 0
	p.data = nil
	return p
}

func (dyn *DynNode) IsLeaf() bool {
	if dyn.left == nil {
		return true
	}
	return false
}

func (dyn *DynNode) RotateRight(gross, grossR float64) {
	var u, v *DynNode
	var pnetmin, pnetminR *float64
	var rstate bool

	if dyn.IsLeaf() {
		/*it is */
		return
	}
}
