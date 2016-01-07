package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	CAP_INF     = float64(1.79769313e308)
	CAP_ZERO    = float64(0.0)
	CAP_EPSILON = float64(0.000001)
)

type FaceArr struct {
	faces    []*Face
	ncols    int
	nrows    int
	totalnum int
}

func NewFaceArr(w, h int) *FaceArr {
	p := &FaceArr{}
	p.ncols = w - 1
	p.nrows = h - 1
	p.totalnum = p.ncols*p.nrows + 1
	p.faces = []*Face{}
	for i := 0; i < p.totalnum; i++ {
		f := NewFace()
		f.SetIdx(i)
		p.faces = append(p.faces, f)
	}
	return p
}

func (fa *FaceArr) GetFaces(x, y int) (retfaces []*Face, err error) {
	retfaces = []*Face{}

	if x < 0 || x > fa.ncols {
		err = fmt.Errorf("x(%d) not valid", x)
		log.Print(err.Error())
		return
	}

	if y < 0 || y > fa.nrows {
		err = fmt.Errorf("y(%d) not valid", y)
		log.Print(err.Error())
		return
	}

	//log.Printf("y:%d,x:%d", y, x)

	if x == 0 && y == 0 {
		retfaces = append(retfaces, fa.faces[0])
	} else if x == fa.ncols && y == fa.nrows {
		retfaces = append(retfaces, fa.faces[fa.totalnum-2])
	} else if x == 0 && y == fa.nrows {
		retfaces = append(retfaces, fa.faces[fa.ncols*(y-1)])
	} else if x == fa.ncols && y == 0 {
		retfaces = append(retfaces, fa.faces[fa.ncols-1])
	} else if x == 0 {
		retfaces = append(retfaces, fa.faces[fa.ncols*(y-1)])
		retfaces = append(retfaces, fa.faces[fa.ncols*(y)])
	} else if y == 0 {
		retfaces = append(retfaces, fa.faces[x-1])
		retfaces = append(retfaces, fa.faces[x])
	} else if x == fa.ncols {
		retfaces = append(retfaces, fa.faces[fa.ncols*y-1])
		retfaces = append(retfaces, fa.faces[fa.ncols*(y+1)-1])
	} else if y == fa.nrows {
		retfaces = append(retfaces, fa.faces[fa.ncols*(y-1)+x-1])
		retfaces = append(retfaces, fa.faces[fa.ncols*(y-1)+x])
	} else {
		/*we should add four faces*/
		retfaces = append(retfaces, fa.faces[fa.ncols*(y-1)+(x-1)])
		retfaces = append(retfaces, fa.faces[fa.ncols*(y-1)+x])
		retfaces = append(retfaces, fa.faces[fa.ncols*y+x-1])
		retfaces = append(retfaces, fa.faces[fa.ncols*y+x])
	}
	if x == 0 || y == 0 || x == fa.ncols || y == fa.nrows {
		/*put the infinite into the faces*/
		retfaces = append(retfaces, fa.faces[fa.totalnum-1])
	}

	err = nil
	return

}

type VertsHash struct {
	verts   map[string]*Vertice
	vertarr []*Vertice
}

func NewVertsHash() *VertsHash {
	p := &VertsHash{}
	p.verts = make(map[string]*Vertice)
	p.vertarr = []*Vertice{}
	return p
}

func (p *VertsHash) AddVerts(name string, w, h int, fa *FaceArr) error {
	_, ok := p.verts[name]
	if ok {
		return nil
	}

	p.verts[name] = NewVertice(name)
	idx, err := strconv.Atoi(name)
	if err != nil {
		return err
	}
	/*to increase the vertex*/
	vert := p.verts[name]
	x := idx % w
	y := idx / w
	vert.SetXY(x, y)
	vert.SetIdx(idx)
	faces, err := fa.GetFaces(x, y)
	if err != nil {
		return err
	}
	vert.SetFaces(faces)
	p.vertarr = append(p.vertarr, p.verts[name])
	return nil
}

func (p *VertsHash) GetVert(name string, w, h int, fa *FaceArr) (vert *Vertice, err error) {
	err = p.AddVerts(name, w, h, fa)
	if err != nil {

		return
	}
	vert = p.verts[name]
	err = nil
	return
}

type EdgeHash struct {
	edgemap map[string]*Edge
	edgearr []*Edge
}

func NewEdgeHash() *EdgeHash {
	p := &EdgeHash{}
	p.edgemap = make(map[string]*Edge)
	p.edgearr = []*Edge{}
	return p
}

func (eh *EdgeHash) get_name(from, to string) string {
	return fmt.Sprintf("%s -> %s", from, to)
}

func (eh *EdgeHash) get_edge(from, to string) *Edge {
	name := eh.get_name(from, to)
	e, ok := eh.edgemap[name]
	if ok {
		return e
	}
	return nil
}

func (eh *EdgeHash) compare_get_face(fromfaces []*Face, tofaces []*Face) []*Face {
	retfaces := []*Face{}
	for i := 0; i < len(fromfaces); i++ {
		for j := 0; j < len(tofaces); j++ {
			if fromfaces[i] == tofaces[j] {
				retfaces = append(retfaces, fromfaces[i])
			}
		}
	}
	return retfaces
}

func (eh *EdgeHash) is_upsidedown(fromv, tov *Vertice) bool {
	if fromv.GetX() != tov.GetX() {
		/*if we from the vertical mode return false*/
		if fromv.GetY() == 0 {
			return false
		}
		return true
	}

	if fromv.GetY() != tov.GetY() {
		if fromv.GetX() == 0 {
			return true
		}
		return false
	}

	log.Fatalf("can not reach here")
	return false
}

func (eh *EdgeHash) AddEdge(fromv, tov *Vertice, caps float64) error {
	var ed *Edge

	ed = eh.get_edge(fromv.GetName(), tov.GetName())
	if ed != nil {
		err := fmt.Errorf("set (%s) twice", eh.get_name(fromv.GetName(), tov.GetName()))
		log.Print(err.Error())
		return err
	}

	ed = eh.get_edge(tov.GetName(), fromv.GetName())
	if ed != nil {
		ed.SetCap(caps)
		fromv.PushEdge(ed)
		return nil
	}

	ed = NewEdge()
	ed.SetHead(tov)
	ed.SetTail(fromv)
	ed.SetRevCap(caps)
	eh.edgemap[eh.get_name(fromv.GetName(), tov.GetName())] = ed
	ed.SetName(eh.get_name(fromv.GetName(), tov.GetName()))
	ed.SetIdx(len(eh.edgearr))
	retfaces := eh.compare_get_face(fromv.GetFaces(), tov.GetFaces())
	if len(retfaces) != 2 {
		log.Printf("vert[%s] faces", fromv.GetName())
		for i, face := range fromv.GetFaces() {
			log.Printf("[%d].face %d", i, face.GetIdx())
		}

		log.Printf("vert[%s] faces", tov.GetName())
		for i, face := range tov.GetFaces() {
			log.Printf("[%d].face %d", i, face.GetIdx())
		}
		err := fmt.Errorf("%s not valid faces", ed.GetName())
		log.Print(err.Error())
		return err
	}
	if eh.is_upsidedown(fromv, tov) {
		//log.Printf("[%d][%d] -> [%d][%d] upsidedown Head %d Tail %d",
		//	fromv.GetY(), fromv.GetX(), tov.GetY(), tov.GetX(),
		//	retfaces[1].GetIdx(), retfaces[0].GetIdx())
		ed.SetHeadDual(retfaces[1])
		ed.SetTailDual(retfaces[0])
	} else {
		//log.Printf("[%d][%d] -> [%d][%d] not upsidedown Head %d Tail %d",
		//	fromv.GetY(), fromv.GetX(), tov.GetY(), tov.GetX(),
		//	retfaces[0].GetIdx(), retfaces[1].GetIdx())
		ed.SetHeadDual(retfaces[0])
		ed.SetTailDual(retfaces[1])
	}

	fromv.PushEdge(ed)
	eh.edgearr = append(eh.edgearr, ed)
	return nil
}

type PlanarGraph struct {
	verts     []*Vertice
	edges     []*Edge
	faces     []*Face
	sinkid    int
	sourceid  int
	preflowed int
}

/*we get the */
func NewPlanarGraph() *PlanarGraph {
	p := &PlanarGraph{}
	p.verts = []*Vertice{}
	p.edges = []*Edge{}
	p.faces = []*Face{}
	p.sinkid = -1
	p.sourceid = -1
	p.preflowed = 0
	return p
}

func (planar *PlanarGraph) DebugGraph() {
	for i, v := range planar.verts {
		for j := 0; j < v.GetEdgeNum(); j++ {
			fmt.Fprintf(os.Stdout, "[%d].edge[%d] %d\n", i, j, v.GetEdge(j).GetIdx())
		}
	}
	fmt.Fprintf(os.Stdout, "sourceid %d sinkid %d\n", planar.sourceid, planar.sinkid)
	for _, e := range planar.edges {
		scap := fmt.Sprintf("%f", e.GetCap())
		srcap := fmt.Sprintf("%f", e.GetRevCap())
		if e.GetCap() == CAP_INF {
			scap = "1.#INF00"
		}

		if e.GetRevCap() == CAP_INF {
			srcap = "1.#INF00"
		}

		fmt.Fprintf(os.Stdout, "[%d] flags(0x%08x) .cap %s .rcap %s head %d tail %d headdual %d taildual %d\n",
			e.GetIdx(), e.GetFlags(), scap, srcap, e.GetHead().GetIdx(),
			e.GetTail().GetIdx(), e.GetHeadDual().GetIdx(),
			e.GetTailDual().GetIdx())
	}
}

func (planar *PlanarGraph) SetCapNormalize() {
	var capInf, capmin, capeps, curcap float64
	var capzeronum int

	capInf = CAP_ZERO
	capmin = CAP_INF

	for _, e := range planar.edges {
		curcap = e.GetCap()
		if curcap != CAP_INF {
			capInf += curcap
		}

		curcap = e.GetRevCap()
		if curcap != CAP_INF {
			capInf += curcap
		}
	}

	capInf += float64(1.0)

	for _, e := range planar.edges {
		curcap = e.GetCap()
		if curcap == CAP_INF {
			e.SetCap(capInf)
		}

		curcap = e.GetRevCap()
		if curcap == CAP_INF {
			e.SetRevCap(capInf)
		}
	}
	capzeronum = 0
	for _, e := range planar.edges {
		curcap = e.GetCap()
		if curcap == CAP_ZERO {
			capzeronum++
		} else if curcap < capmin {
			capmin = curcap
		}

		curcap = e.GetRevCap()
		if curcap == CAP_ZERO {
			capzeronum++
		} else if curcap < capmin {
			capmin = curcap
		}
	}

	if capzeronum == 0 {
		capeps = CAP_INF
	} else {
		capeps = capmin / float64(capzeronum*2)
	}

	if capeps == CAP_ZERO {
		capeps = CAP_EPSILON
	}

	for _, e := range planar.edges {
		curcap = e.GetCap()
		if curcap == CAP_ZERO {
			e.SetCap(capeps)
			e.SetFlags(e.GetFlags() | EDGE_CAP_EPSILON)
		}

		curcap = e.GetRevCap()
		if curcap == CAP_ZERO {
			e.SetRevCap(capeps)
			e.SetFlags(e.GetFlags() | EDGE_RCAP_EPSILON)
		}
	}
}

func (planar *PlanarGraph) SetCounterClockWise() {
	for _, v := range planar.verts {
		v.CounterClockWise()
	}
}

func MakePlanarGraph(infile string) (planar *PlanarGraph, err error) {
	var sarr []string
	var caps float64
	var linenum int
	var facearr *FaceArr
	var sinkid, sourceid int

	/*open file*/
	planar = nil
	err = nil
	w := -1
	h := -1
	sinkid = -1
	sourceid = -1
	file, e := os.Open(infile)
	if e != nil {
		err = e
		return
	}
	defer file.Close()
	vertshash := NewVertsHash()
	edgehash := NewEdgeHash()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	linenum = 0

	for scanner.Scan() {
		l := scanner.Text()
		l = strings.Trim(l, "\r\n")
		linenum++
		if strings.HasPrefix(l, "#") {
			continue
		}

		if strings.HasPrefix(l, "height=") {
			sarr = strings.Split(l, "=")
			if len(sarr) < 2 {
				continue
			}
			h, err = strconv.Atoi(sarr[1])
			if err != nil {
				err = fmt.Errorf("can not parse (%d) (%s)", linenum, l)
				log.Print(err.Error())
				return
			}
			if w > 0 && facearr == nil {
				facearr = NewFaceArr(w, h)
			}
			continue
		}

		if strings.HasPrefix(l, "width=") {
			sarr = strings.Split(l, "=")
			if len(sarr) < 2 {
				continue
			}
			w, err = strconv.Atoi(sarr[1])
			if err != nil {
				err = fmt.Errorf("can not parse (%d) (%s)", linenum, l)
				log.Print(err.Error())
				return
			}
			if h > 0 && facearr == nil {
				facearr = NewFaceArr(w, h)
			}
			continue
		}

		if strings.HasPrefix(l, "source=") {
			sarr = strings.Split(l, "=")
			if len(sarr) < 2 {
				continue
			}
			sourceid, err = strconv.Atoi(sarr[1])
			if err != nil {
				err = fmt.Errorf("can not parse (%d) (%s)", linenum, l)
				log.Print(err.Error())
				return
			}
			continue
		}

		if strings.HasPrefix(l, "sink=") {
			sarr = strings.Split(l, "=")
			if len(sarr) < 2 {
				continue
			}
			sinkid, err = strconv.Atoi(sarr[1])
			if err != nil {
				err = fmt.Errorf("can not parse (%d) (%s)", linenum, l)
				log.Print(err.Error())
				return
			}
			continue
		}

		sarr = strings.Split(l, ",")
		if len(sarr) < 3 {
			continue
		}

		if w < 0 || h < 0 {
			err = fmt.Errorf("not specified width or height")
			log.Print(err.Error())
			return
		}

		if sarr[2] == "1.#INF00" {
			caps = CAP_INF
		} else {
			caps, e = strconv.ParseFloat(sarr[2], 64)
			if e != nil {
				err = fmt.Errorf("can not parse %d error %s", linenum, e.Error())
				log.Print(err.Error())
				return
			}
		}

		fromvert, verr := vertshash.GetVert(sarr[0], w, h, facearr)
		if verr != nil {
			err = verr
			return
		}
		tovert, verr := vertshash.GetVert(sarr[1], w, h, facearr)
		if verr != nil {
			err = verr
			return
		}
		verr = edgehash.AddEdge(fromvert, tovert, caps)
		if verr != nil {
			err = fmt.Errorf("line(%d) %s", linenum, verr.Error())
			log.Print(err.Error())
			return
		}
	}

	if sourceid < 0 || sinkid < 0 || sinkid == sourceid {
		err = fmt.Errorf("can not find source id and sink id")
		log.Print(err.Error())
		return
	}
	planar = NewPlanarGraph()
	planar.edges = edgehash.edgearr
	planar.verts = vertshash.vertarr
	planar.faces = facearr.faces
	planar.sourceid = sourceid
	planar.sinkid = sinkid
	planar.SetCounterClockWise()
	planar.SetCapNormalize()
	err = nil
	return
}

func (planar *PlanarGraph) preflow() {
	var infedge *Edge
	var infvert *Vertice
	if planar.preflowed > 0 {
		return
	}

	/*we now preflowed the */
	infedge = planar.edges[planar.sinkid]

}
