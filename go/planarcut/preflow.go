package main

import (
	"fmt"
)

func (planar *PlanarGraph) preflow() {
	var infedge *Edge
	var infvert *Vertice
	if planar.preflowed > 0 {
		return
	}

	/*we now preflowed the */
	infedge = planar.edges[planar.sinkid]

}
