package main

type Face struct {
	idx int64
}

func NewFace() *Face {
	p := &Face{}
	p.idx = 0
	return p
}
