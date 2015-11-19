package main

func BSF(caps map[string]map[string]int, neighs map[string][]string, flows map[string]map[string]int,
	source string, sink string, maxval int) (max int, parent map[string]string) {
	var queue []string
	parents := make(map[string]string)
	M := make(map[string]int)

	for k, _ := range caps {
		parents[k] = ""
		M[k] = 0
	}

	M[source] = maxval
	parents[source] = "#"
	queue = append(queue, source)
	for len(queue) > 0 {
		u := queue[len(queue)-1]
		queue = queue[:(len(queue) - 1)]
		if k, ok := neighs[u]; ok {
			for _, v := range k {
				if (caps[u][v]-flows[u][v]) > 0 && parents[v] == "" {
					parents[v] = u
					if M[u] < (caps[u][v] - flows[u][v]) {
						M[v] = M[u]
					} else {
						M[v] = caps[u][v] - flows[u][v]
					}

					if v != sink {
						queue = append(queue, v)
					} else {
						return M[v], parents
					}

				}
			}

		}
	}
	return 0, parents

}

func EdmondsWarp(caps map[string]map[string]int, neighs map[string][]string, source string, sink string) (flow int, flows map[string]map[string]int) {
	flow = 0
	flows = make(map[string]map[string]int)
	maxval := 0
	sortkeys := MakeSortKeys(caps)
	for _, k1 := range sortkeys {
		for _, k2 := range sortkeys {
			flows = SetDictDefValue(flows, k1, k2, 0)
			caps = SetDictDefValue(caps, k1, k2, 0)
			maxval += caps[k1][k2]
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
			flows[u][v] += max
			flows[v][u] -= max
			v = u
		}

	}
	return flow, flows
}
