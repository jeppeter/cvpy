package main

const REV_MASK uint8 = 1
const TMP_MASK uint8 = 2
const MAP_MASK uint8 = 4
const CAP_INF float64 = -1.0

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

func (dyn *DynNode) GetReserved() bool {
	v := dyn.reserved & REV_MASK
	if v != 0 {
		return true
	}
	return false
}

func (dyn *DynNode) GetMapping() bool {
	v := dyn.reserved & MAP_MASK
	if v != 0 {
		return true
	}
	return false
}

func (dyn *DynNode) SetMapping(b bool) {
	if dyn.GetMapping() != b {
		dyn.reserved ^= MAP_MASK
	}
	return
}

func (dyn *DynNode) SetReserved(b bool) {
	if dyn.GetReserved() != b {
		dyn.reserved ^= REV_MASK
	}
	return

}

func (dyn *DynNode) NormalizeReserveState() {
	if !dyn.GetReserved() {
		return
	}
	/*do not normalize for twice*/
	dyn.SetReserved(false)
	/*to reverse dyn*/
	dyn.SetMapping(!dyn.GetMapping())
	/*swap the left and right*/
	pn := dyn.left
	dyn.left = dyn.right
	dyn.right = pn

	/*swap tail and head*/
	pn = dyn.head
	dyn.head = dyn.tail
	dyn.tail = pn

	/*swap for net min*/
	c := dyn.netmin
	dyn.netmin = dyn.netminR
	dyn.netminR = c

	/*swap for net cost*/
	c = dyn.netcost
	dyn.netcost = dyn.netcostR
	dyn.netcostR = c

	if dyn.right && !dyn.right.IsLeaf() {
		/*to set the opposite*/
		dyn.right.SetReserved(!dyn.right.GetReserved())

	}

	if dyn.left && !dyn.left.IsLeaf() {
		dyn.left.SetReserved(!dyn.left.GetReserved())
	}
}

func (dyn *DynNode) GetData() interface{} {
	return dyn.data
}

func (dyn *DynNode) RotateRight(gross, grossR float64) {
	var u, v *DynNode
	var pnetmin, pnetminR *float64
	var rstate bool

	if dyn.IsLeaf() {
		/*it is in leaf mode*/
		return
	}
	dyn.NormalizeReserveState()
	u = dyn
	v = dyn.left

	if v.IsLeaf() {
		return
	}

	v.NormalizeReserveState()

	uold := *dyn
	vold := *v
	umapping := uold.GetMapping()
	vmapping := vold.GetMapping()

	udata := uold.GetData()
	vdata := vold.GetData()
	minU := gross
	minUR := grossR
	minvold := v.netmin + minU
	minvoldR := v.netminR + minUR

	costU := u.netcost + minU
	costUR := u.netcostR + minUR
	costV := v.netcost + minvold
	costVR := v.netcostR + minvoldR
	minvl := 
}
