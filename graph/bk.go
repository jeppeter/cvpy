package main

import (
	"fmt"
	"log"
)

type Arc struct {
	head   *Node
	next   *Arc
	sister *Arc
	r_cap  int
}

type Node struct {
	first              *Arc
	parent             *Arc
	next               *Node
	TS                 int
	DIST               int
	is_sink            bool
	is_marked          bool
	is_in_changed_list bool
	tr_cap             int
}

type NodePtr struct {
	ptr  *Node
	next *NodePtr
}

type BKGraph struct {
	nodes     *Node
	node_last *Node
	node_max  *Node
	arcs      *Arc
	arc_last  *Arc
	arc_max   *Arc
	node_num  int
}
