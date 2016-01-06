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
	CAP_INF = float64(1.79769313e308)
)

type PlanarGraph struct {
	verts []*Vertice
	edges []*Edge
	faces []*Face
}

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
		return
	}

	if y < 0 || y > fa.nrows {
		err = fmt.Errorf("y(%d) not valid", y)
		return
	}

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
	if fromv.GetY() != tov.GetY() {
		/*if we from the vertical mode return false*/
		return false
	}

	if fromv.GetX() == 0 {
		return true
	}

	return false
}

func (eh *EdgeHash) AddEdge(fromv, tov *Vertice, caps float64) error {
	var ed *Edge

	ed = eh.get_edge(fromv.GetName(), tov.GetName())
	if ed != nil {
		err := fmt.Errorf("set (%s) twice", eh.get_name(fromv.GetName(), tov.GetName()))
		return err
	}

	ed = eh.get_edge(tov.GetName(), fromv.GetName())
	if ed != nil {
		ed.SetCap(caps)
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
		err := fmt.Errorf("%s not valid faces", ed.GetName())
		return err
	}
	if eh.is_upsidedown(fromv, tov) {
		log.Printf("[%d][%d] -> [%d][%d] upsidedown Head %d Tail %d",
			fromv.GetY(), fromv.GetX(), tov.GetY(), tov.GetX(),
			retfaces[1].GetIdx(), retfaces[0].GetIdx())
		ed.SetHeadDual(retfaces[1])
		ed.SetTailDual(retfaces[0])
	} else {
		log.Printf("[%d][%d] -> [%d][%d] not upsidedown Head %d Tail %d",
			fromv.GetY(), fromv.GetX(), tov.GetY(), tov.GetX(),
			retfaces[0].GetIdx(), retfaces[1].GetIdx())
		ed.SetHeadDual(retfaces[0])
		ed.SetTailDual(retfaces[1])
	}
	eh.edgearr = append(eh.edgearr, ed)
	return nil
}

/*we get the */
func NewPlanarGraph() *PlanarGraph {
	p := &PlanarGraph{}
	p.verts = []*Vertice{}
	p.edges = []*Edge{}
	p.faces = []*Face{}
	return p
}

func (planar *PlanarGraph) DebugGraph() {
	for _, e := range planar.edges {
		scap := fmt.Sprintf("%f", e.GetCap())
		srcap := fmt.Sprintf("%f", e.GetRevCap())
		if e.GetCap() == CAP_INF {
			scap = "1.#INF00"
		}

		if e.GetRevCap() == CAP_INF {
			srcap = "1.#INF00"
		}

		fmt.Fprintf(os.Stdout, "[%d] .cap %s .rcap %s head %d tail %d headdual %d taildual %d\n",
			e.GetIdx(), scap, srcap, e.GetHead().GetIdx(),
			e.GetTail().GetIdx(), e.GetHeadDual().GetIdx(),
			e.GetTailDual().GetIdx())
	}
}
func MakePlanarGraph(infile string) (planar *PlanarGraph, err error) {
	var sarr []string
	var caps float64
	var linenum int
	var facearr *FaceArr

	/*open file*/
	planar = nil
	err = nil
	w := -1
	h := -1
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
				return
			}
			if h > 0 && facearr == nil {
				facearr = NewFaceArr(w, h)
			}
			continue
		}

		sarr = strings.Split(l, ",")
		if len(sarr) < 3 {
			continue
		}

		if w < 0 || h < 0 {
			err = fmt.Errorf("not specified width or height")
			return
		}

		if sarr[2] == "1.#INF00" {
			caps = CAP_INF
		} else {
			caps, e = strconv.ParseFloat(sarr[2], 64)
			if e != nil {
				err = fmt.Errorf("can not parse %d error %s", linenum, e.Error())
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
			return
		}
	}
	planar = NewPlanarGraph()
	planar.edges = edgehash.edgearr
	planar.verts = vertshash.vertarr
	planar.faces = facearr.faces

	err = nil
	return
}
