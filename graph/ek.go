package main

func BSF(caps *StringGraph, neighs *Neigbour, flows *StringGraph,
	source string, sink string, maxval int) (max int, parent map[string]string) {
	var queue []string
	parents := make(map[string]string)
	M := make(map[string]int)

	for _, k := range caps.Iter() {
		parents[k] = ""
		M[k] = 0
	}

	M[source] = maxval
	parents[source] = "#"
	queue = append(queue, source)
	for len(queue) > 0 {
		u := queue[len(queue)-1]
		queue = queue[:(len(queue) - 1)]
		for _, v := range neighs.GetValue(u) {
			if (caps.GetValue(u, v)-flows.GetValue(u, v)) > 0 && parents[v] == "" {
				parents[v] = u
				if M[u] < (caps.GetValue(u, v) - flows.GetValue(u, v)) {
					M[v] = M[u]
				} else {
					M[v] = caps.GetValue(u, v) - flows.GetValue(u, v)
				}

				if v != sink {
					queue = append(queue, v)
				} else {
					return M[v], parents
				}

			}
		}

	}
	return 0, parents

}

func EdmondsWarp(caps *StringGraph, neighs *Neigbour, source string, sink string) (flow int, flows *StringGraph) {
	flow = 0
	flows = NewStringGraph()
	maxval := 0
	sortkeys := MakeSortKeys(caps)
	for _, k1 := range sortkeys {
		for _, k2 := range sortkeys {
			maxval += caps.GetValue(k1, k2)
		}
	}

	for {
		max, parents := BSF(caps, neighs, flows, source, sink, maxval)
		if max == 0 {
			break
		}
		flow += max
		v := sink
		for v != source {
			u := parents[v]
			flows.SetValue(u, v, flows.GetValue(u, v)+max)
			flows.SetValue(v, u, flows.GetValue(v, u)-max)
			v = u
		}

	}
	return flow, flows
}
