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
	verts map[string]*Vertice
}

func NewVertsHash() *VertsHash {
	p = &VertsHash{}
	p.verts = make(map[string]*Vertice)
	return p
}

func (p *VertsHash) AddVerts(name string) int {
	_, ok := p.verts[name]
	if ok {
		return 0
	}

	p.verts[name] = NewVertice(name)
	return 1
}

func (p *VertsHash) GetVert(name string) *Vertice {
	p.AddVerts(name)
	return p.verts[name]
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
			h = strconv.Atoi(sarr[1])
			continue
		}

		if strings.HasPrefix(l, "width=") {
			sarr = strings.Split(l, "=")
			if len(sarr) < 2 {
				continue
			}
			w = strconv.Atoi(sarr[1])
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

		fromvert := vertshash.GetVert(sarr[0])
		tovert := vertshash.GetVert(sarr[1])

	}
}
