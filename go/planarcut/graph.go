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

func (p *VertsHash) AddVerts(name string, w, h int) int {
	_, ok := p.verts[name]
	if ok {
		return 0
	}

	p.verts[name] = NewVertice(name)
	idx, err := strconv.Atoi(name)
	if err != nil {
		log.Fatalf("can not accept name(%s)", name)
	}
	vert := p.verts[name]
	x := idx % w
	y := idx / w
	vert.SetXY(x, y)
	vert.SetIdx(idx)
	p.vertarr = append(p.vertarr, p.verts[name])
	return 1
}

func (p *VertsHash) GetVert(name string, w, h int) *Vertice {
	p.AddVerts(name, w, h)
	return p.verts[name]
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

func (eh *EdgeHash) AddEdge(fromv, tov *Vertice, caps float64) int {
	var ed *Edge

	ed = eh.get_edge(fromv.GetName(), tov.GetName())
	if ed != nil {
		log.Fatalf("set (%s) twice", eh.get_name(fromv.GetName(), tov.GetName()))
		return -1
	}

	ed = eh.get_edge(tov.GetName(), fromv.GetName())
	if ed != nil {
		ed.SetRevCap(caps)
		return 0
	}

	ed = NewEdge()
	ed.SetHead(fromv)
	ed.SetTail(tov)
	ed.SetCap(caps)
	eh.edgemap[eh.get_name(fromv.GetName(), tov.GetName())] = ed
	ed.SetName(eh.get_name(fromv.GetName(), tov.GetName()))
	ed.SetIdx(len(eh.edgearr))
	eh.edgearr = append(eh.edgearr, ed)
	return 1
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

func (fa *FaceArr) GetFaces(x, y int) []*Face {
	var retfaces []*Face
	retfaces = []*Face{}

	if x < 0 || x > fa.ncols {
		log.Fatalf("x(%d) not valid", x)
	}

	if y < 0 || y > fa.nrows {
		log.Fatalf("y(%d) not valid", y)
	}
	if x == 0 || y == 0 || x == fa.ncols || y == fa.nrows {
		/*put the infinite into the faces*/
		retfaces = append(retfaces, fa.faces[fa.totalnum-1])
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
		retfaces = append(retfaces, fa.faces[fa.ncols*(y-1)-1])
		retfaces = append(retfaces, fa.faces[fa.ncols*y-1])
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

	return retfaces

}

/*we get the */
func NewPlanarGraph() *PlanarGraph {
	p := &PlanarGraph{}
	p.verts = []*Vertice{}
	p.edges = []*Edge{}
	p.faces = []*Face{}
	return p
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
				log.Fatalf("can not parse(%d) %s", linenum, l)
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
				log.Fatalf("can not parse(%d) %s", linenum, l)
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
			log.Fatalf("not specified width or height")
		}

		if sarr[2] == "1.#INF00" {
			caps = CAP_INF
		} else {
			caps, e = strconv.ParseFloat(sarr[2], 64)
			if e != nil {
				log.Fatalf("can not parse %d error %s", linenum, e.Error())
				err = e
				return
			}
		}

		fromvert := vertshash.GetVert(sarr[0], w, h)
		tovert := vertshash.GetVert(sarr[1], w, h)
		edgehash.AddEdge(fromvert, tovert, caps)
	}
	planar = NewPlanarGraph()
	planar.edges = edgehash.edgearr
	planar.verts = vertshash.vertarr
	planar.faces = facearr.faces
	err = nil
	return
}
