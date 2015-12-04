package main

const REV_MASK uint8 = 1
const TMP_MASK uint8 = 2
const MAP_MASK uint8 = 4
const CAP_INF float64 = float64(1e+308)

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

func (dyn *DynNode) GetNetmin(b bool) float64 {
	if b {
		return dyn.netminR
	}
	return dyn.netmin
}

func (dyn *DynNode) GetNetcost(b bool) float64 {
	if b {
		return dyn.netcostR
	}
	return dyn.netcost
}

func (dyn *DynNode) SetAsRChild(pn *DynNode, rstate bool) {
	var newtail *DynNode
	dyn.right = pn
	pn.parent = dyn

	if pn.IsLeaf() {
		newtail = pn
	} else {
		if pn.GetReserved() == rstate {
			newtail = pn.tail
		} else {
			newtail = pn.head
		}
	}

	if rstate {
		dyn.head = newtail
	} else {
		dyn.tail = newtail
	}
	return
}

func (dyn *DynNode) SetAsLChild(pn *DynNode, rstate bool) {
	var newhead *DynNode
	dyn.left = pn
	pn.parent = dyn

	if pn.IsLeaf() {
		newhead = pn
	} else {
		if pn.GetReserved() == rstate {
			newhead = pn.head
		} else {
			newhead = pn.tail
		}
	}

	if rstate {
		dyn.tail = newhead
	} else {
		dyn.head = newhead
	}
}

func MMin64(u, v, w float64) float64 {
	min := u
	if min > v {
		min = v
	}

	if min > w {
		min = w
	}
	return min
}

func (dyn *DynNode) GetNetMinPtr(pNetMin, pNetMinR **float64, rstate bool) {
	rstate ^= dyn.GetReserved()

	if rstate {
		*pNetMin = &(dyn.netminR)
		*pNetMinR = &(dyn.netmin)
	} else {
		*pNetMin = &(dyn.netmin)
		*pNetMinR = &(dyn.netminR)
	}
	return
}

func (dyn *DynNode) RotateRight(gross, grossR float64) {
	var u, v, w *DynNode
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
	w = v.left

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

	minVl := CAP_INF
	minVlR := CAP_INF

	minVr := CAP_INF
	minVrR := CAP_INF

	minUr := CAP_INF
	minUrR := CAP_INF

	if !v.left.IsLeaf() {
		/*here is DynPath.cpp:136*/
		v.left.GetNetMinPtr(&pnetmin, &pnetminR, false)
		minVl = *pnetmin + minvold
		minVlR = *pnetminR + minvoldR
	}

	if !v.right.IsLeaf() {
		v.right.GetNetMinPtr(&pnetmin, &pnetminR, false)
		minVr = *pnetmin + minvold
		minVrR = *pnetminR + minvoldR
	}

	if !u.right.IsLeaf() {
		u.right.GetNetMinPtr(&pnetmin, &pnetminR, false)
		minUr = *pnetmin + minU
		minUrR = *pnetminR + minUR
	}

	vnew := u
	unew := v

	vnew.SetAsLChild(vold.left, false)
	unew.SetAsLChild(vold.right, false)
	unew.SetAsRChild(uold.right, false)
	vnew.SetAsRChild(unew, false)

	minVNew := minU
	minVNewR := minUR

	minUNew := MMin64(costU, minUr, minVr)
	minUNewR := MMin64(costUR, minUrR, minVrR)

	/*for DynPath.cpp:170*/
}
