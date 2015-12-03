package main

import (
	"log"
)

const RB_BLACK int = 1
const RB_RED int = 2

type RBTreeData interface {
	Less(b *RBTreeData) bool
	Stringer() string
}

type RBTreeElem struct {
	parent *RBTreeElem
	left   *RBTreeElem
	right  *RBTreeElem
	color  int
	Data   *RBTreeData
}

type RBTree struct {
	root  *RBTreeElem
	count int
}

func NewRBTreeElem(data *RBTreeData) *RBTreeElem {
	p := &RBTreeElem{}
	p.Data = data
	p.color = RB_BLACK
	p.parent = nil
	p.left = nil
	p.right = nil
	return p
}

func (elem *RBTreeElem) SetParent(parent *RBTreeElem) {
	elem.parent = parent
	return
}

func (elem *RBTreeElem) GetParent() *RBTreeElem {
	return elem.parent
}

func (elem *RBTreeElem) SetLeft(left *RBTreeElem) {
	elem.left = left
	return
}

func (elem *RBTreeElem) GetLeft() *RBTreeElem {
	return elem.left
}

func (elem *RBTreeElem) SetRight(right *RBTreeElem) {
	elem.right = right
	return
}

func (elem *RBTreeElem) GetRight() *RBTreeElem {
	return elem.right
}

func (elem *RBTreeElem) GetColor() int {
	return elem.color
}

func (elem *RBTreeElem) SetColor(color int) {
	elem.color = color
	return
}

func NewRBTree() *RBTree {
	p := &RBTree{}
	p.root = nil
	p.count = 0
	return p
}

func (rb *RBTree) __find_insert_parent(data *RBTreeData, from *RBTreeElem) *RBTreeElem {
	if from.Data.Less(data) {
		if from.GetLeft() != nil {
			return rb.__find_insert_parent(data, from.GetLeft())
		} else {
			return from
		}
	} else {
		if from.GetRight() != nil {
			return rb.__find_insert_parent(data, from.GetRight())
		} else {
			return from
		}
	}
	return nil
}

func (rb *RBTree) find_insert_parent(data *RBTreeData) *RBTreeElem {
	if rb.root.Data.Less(data) {
		if rb.root.GetLeft() != nil {
			return rb.__find_insert_parent(data, rb.root.GetLeft())
		} else {
			return rb.root
		}
	} else {
		if rb.root.GetRight() != nil {
			return rb.__find_insert_parent(data, rb.root.GetRight())
		} else {
			return rb.root
		}
	}
	return nil
}

func (rb *RBTree) Insert(data *RBTreeData) int {
	if rb.count == 0 {
		rb.root = NewRBTreeElem(data)
		rb.count++
	} else {
		parent = rb.find_insert_parent(data)
		if parent != nil {

		} else {
			log.Fatalf("can not find parent to insert (%s)", data.Stringer())
		}

	}
	return rb.count
}
