package main

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

func (p *Neigbour) GetValue(k1 string) []string {
	val, ok := p.inner[k1]
	if !ok {
		return []string{}
	}
	return val
}

func (p *Neigbour) Iter() []string {
	q := []string{}
	for k, _ := range p.inner {
		q = append(q, k)
	}
	return q
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

type StringNode struct {
	n    string
	next *StringNode
}

type StringStack struct {
	cnt  int
	node *StringNode
}

func NewStringStack() *StringStack {
	p := &StringStack{}
	p.cnt = 0
	p.node = nil
	return p
}

func NewStringNode(n string, nnode *StringNode) *StringNode {
	node := &StringNode{}
	node.n = n
	node.next = nnode
	return node
}

func (p *StringStack) PushValue(n string) {
	p.cnt += 1
	p.node = NewStringNode(n, p.node)
	return
}

func (p *StringStack) PopValue() string {
	if p.cnt == 0 {
		return ""
	}

	n := p.node.n
	p.node = p.node.next
	p.cnt -= 1
	return n
}

func (p *StringStack) Length() int {
	return p.cnt
}
