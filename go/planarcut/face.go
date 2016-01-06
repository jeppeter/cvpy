package main

type Face struct {
	idx int
}

func NewFace() *Face {
	p := &Face{}
	p.idx = 0
	return p
}

func (f *Face) SetIdx(idx int) {
	f.idx = idx
	return
}

func (f *Face) GetIdx() int {
	return f.idx
}
