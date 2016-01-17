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

	dijgraph := NewDijGraph()

	for _,e := range planar.edges {
		/**/
		srcfidx := e.GetHeadDual().GetIdx()
		dstfidx := e.GetTailDual().GetIdx()

		dijgraph.AddEdge(fmt.Sprint("%d",srcfidx),fmt.Sprint("%d",dstfidx),e.GetCap(),
			e.GetRevCap())
	}

	dijgraph.SetSource(fmt.Sprint("%d",planar.sourceid))
	dijgraph.SetSink(fmt.Sprint("%d",planar.sinkid))

	dijgraph.Dijkstra()
	


	/*we now preflowed the */
	infedge = planar.edges[planar.sinkid]

}
