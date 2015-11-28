package main

import "fmt"

type StringGraph struct {
	inner map[string]map[string]int
}

func NewStringGraph() *StringGraph {
	p := &StringGraph{}
	p.inner = make(map[string]map[string]int)
	return p
}

func (p *StringGraph) GetValue(k1, k2 string) int {
	val, ok := p.inner[k1][k2]
	if !ok {
		return 0
	}
	return val
}

func (p *StringGraph) SetValue(k1, k2 string, val int) {
	_, ok := p.inner[k1][k2]
	if !ok {
		_, ok = p.inner[k1]
		if !ok {
			p.inner[k1] = make(map[string]int)
		}
	}
	p.inner[k1][k2] = val
	return
}

func (p *StringGraph) Iter() []string {
	q := []string{}
	for k, _ := range p.inner {
		q = append(q, k)
	}
	return q
}

func (p *StringGraph) IterIdx(k1 string) []string {
	q := []string{}
	if _, ok := p.inner[k1]; ok {
		for k, _ := range p.inner[k1] {
			q = append(q, k)
		}
	}
	return q
}

type Neigbour struct {
	inner map[string][]string
}

func NewNeighbour() *Neigbour {
	p := &Neigbour{}
	p.inner = make(map[string][]string)
	return p
}

func SortArrayString(narr []string) []string {
	var i, j int
	for i = 0; i < len(narr); i++ {
		for j = (i + 1); j < len(narr); j++ {
			if narr[i] > narr[j] {
				tmp := narr[i]
				narr[i] = narr[j]
				narr[j] = tmp
			}
		}
	}
	return narr
}

func (p *Neigbour) GetValue(k1 string) []string {
	val, ok := p.inner[k1]
	if !ok {
		return []string{}
	}
	return SortArrayString(val)
}

func (p *Neigbour) Iter() []string {
	q := []string{}
	for k, _ := range p.inner {
		q = append(q, k)
	}
	return SortArrayString(q)
}

func (p *Neigbour) PushValue(k string, val string) {
	_, ok := p.inner[k]
	if !ok {
		p.inner[k] = []string{}
	}
	p.inner[k] = append(p.inner[k], val)
	return
}

type StringInt struct {
	inner map[string]int
}

func NewStringInt() *StringInt {
	p := &StringInt{}
	p.inner = make(map[string]int)
	return p
}

func (p *StringInt) SetValue(k1 string, val int) {
	p.inner[k1] = val
	return
}

func (p *StringInt) GetValue(k string) int {
	val, ok := p.inner[k]
	if !ok {
		return 0
	}
	return val
}

func (p *StringInt) Iter() []string {
	q := []string{}
	for k, _ := range p.inner {
		q = append(q, k)
	}
	return q
}

type StringStack struct {
	cnt    int
	conmap map[string]int
	conarr []string
}

func NewStringStack() *StringStack {
	p := &StringStack{}
	p.cnt = 0
	p.conmap = make(map[string]int)
	p.conarr = []string{}
	return p
}

func (p *StringStack) PushValue(n string) {
	if _, ok := p.conmap[n]; ok {
		return
	}
	p.conmap[n] = 1
	p.conarr = append(p.conarr, n)
	p.cnt += 1
	return
}

func (p *StringStack) PopValue() string {
	if p.cnt == 0 {
		return ""
	}
	p.cnt -= 1
	n := p.conarr[p.cnt]
	delete(p.conmap, n)
	if p.cnt > 0 {
		p.conarr = p.conarr[:p.cnt]
	} else {
		p.conarr = []string{}
	}
	return n
}

func (p *StringStack) ShiftValue() string {
	if p.cnt == 0 {
		return ""
	}
	p.cnt -= 1
	n := p.conarr[0]
	delete(p.conmap, n)
	if p.cnt > 0 {
		p.conarr = p.conarr[1:]
	} else {
		p.conarr = []string{}
	}
	return n
}

func (p *StringStack) Length() int {
	return p.cnt
}

func (p *StringStack) String() string {
	s := fmt.Sprintf("cnt %d[", p.cnt)
	for _, k := range p.conarr {
		s += fmt.Sprintf(" %s", k)
	}

	s += "]"
	return s

}
