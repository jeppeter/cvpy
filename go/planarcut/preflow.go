package main

import (
	"fmt"
)

func (planar *PlanarGraph) preflow() {
	var infedge *Edge
	if planar.preflowed > 0 {
		return
	}

	dijgraph := NewDijGraph()

	for _, e := range planar.edges {
		/**/
		srcfidx := e.GetTailDual().GetIdx()
		dstfidx := e.GetHeadDual().GetIdx()

		dijgraph.AddEdge(fmt.Sprintf("%d", srcfidx), fmt.Sprintf("%d", dstfidx), e.GetCap(),
			e.GetRevCap())
		Debug("%d -> %d .cap %f .rcap %f", srcfidx, dstfidx, e.GetCap(), e.GetRevCap())
	}
	srcidx := -1
	sinkidx := -1

	infedge = planar.verts[planar.sourceid].GetEdge(0)
	if infedge.GetTail().GetIdx() == planar.sourceid {
		srcidx = infedge.GetHeadDual().GetIdx()
	} else {
		srcidx = infedge.GetTailDual().GetIdx()
	}

	infedge = planar.verts[planar.sinkid].GetEdge(0)
	if infedge.GetTail().GetIdx() == planar.sinkid {
		sinkidx = infedge.GetHeadDual().GetIdx()
	} else {
		sinkidx = infedge.GetTailDual().GetIdx()
	}

	Debug("infFaceIdx %d", sinkidx)
	dijgraph.SetSource(fmt.Sprintf("%d", srcidx))
	dijgraph.SetSink(fmt.Sprintf("%d", sinkidx))

	dijgraph.Dijkstra()

	/*we now preflowed the */
	for i, e := range planar.edges {
		tailfidx := e.GetTailDual().GetIdx()
		headfidx := e.GetHeadDual().GetIdx()

		w := e.GetCap()
		rw := e.GetRevCap()
		headw := dijgraph.GetWeigth(fmt.Sprintf("%d", headfidx))
		tailw := dijgraph.GetWeigth(fmt.Sprintf("%d", tailfidx))

		ew := headw - tailw
		Debug("edge[%d].cap (%f) .rcap (%f) dualgraph nodes[%d].dijkWeight (%f) - nodes[%d].dijkWeight (%f) eta (%f)",
			i, w, rw, headfidx, headw, tailfidx,
			tailw, ew)

		w = w - ew
		rw = rw + ew

		if w < CAP_EPSILON {
			w = float64(0.0)
		}

		if rw < CAP_EPSILON {
			rw = float64(0.0)
		}

		e.SetCap(w)
		e.SetRevCap(rw)
	}

	planar.preflowed = 1
	return
}
