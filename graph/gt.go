package main

import (
	"os"
)

func FindNextnodesMaxValue(nextnodes map[string]int) int {
	maxval := 0
	for _, kv := range nextnodes {
		if kv > maxval {
			maxval = kv
		}
	}
	return maxval

}

func CanPush(n string, neighbours map[string][]string, nextnodes map[string]int, caps map[string]map[string]int, flows map[string]map[string]int) bool {
	for _, k := range neighbours[n] {
		if ((nextnodes[k] + 1) == nextnodes[n]) && (caps[n][k]-flows[n][k]) > 0 {
			return true
		}
	}
	return false
}

func SetNextNodes(n string, neighbours map[string][]string, nextnodes map[string]int, caps map[string]map[string]int, flows map[string]map[string]int, maxval int) map[string]int {
	minval := maxval
	for _, k := range neighbours[n] {
		if (caps[n][k] - flows[n][k]) > 0 {
			if nextnodes[n] < minval {
				minval = nextnodes[n]
			}
		}
	}
	Debug("set nextnodes[%s] = %d\n", n, (1 + minval))
	nextnodes[n] = 1 + minval
	return nextnodes
}

func FindNextNodes(n string, neighbours map[string][]string, nextnodes map[string]int, caps map[string]map[string]int,
	flows map[string]map[string]int, overflow map[string]int) (rflows map[string]map[string]int, roverflow map[string]int) {

	for _, k := range neighbours[n] {
		if (nextnodes[k] + 1) == nextnodes[n] {
			fval := (caps[n][k] - flows[n][k])
			if fval > overflow[n] {
				fval = overflow[n]
			}
			Debug("Set [%s]->[%s] fval %d\n", n, k, fval)
			overflow[k] += fval
			overflow[n] -= fval
			flows[n][k] += fval
			flows[k][n] -= fval
		}
	}

	return flows, overflow

}

func GoldbergTarjan(caps map[string]map[string]int, neighs map[string][]string, source string, sink string) (flow int, flows map[string]map[string]int) {
	var queue []string
	var n string
	var k string
	flow = 0
	flows = make(map[string]map[string]int)
	sortkeys := MakeSortKeys(caps)
	overflow := make(map[string]int)
	nextnodes := make(map[string]int)
	maxval := 0
	for _, k1 := range sortkeys {
		for _, k2 := range sortkeys {
			flows = SetDictDefValue(flows, k1, k2, 0)
			caps = SetDictDefValue(caps, k1, k2, 0)
			if k1 == k2 && caps[k1][k2] != 0 {
				Debug("can not be set for %s %s value %d\n", k1, k2, caps[k1][k2])
				os.Exit(4)
			}
			maxval += caps[k1][k2]
		}
		overflow[k1] = 0
		nextnodes[k1] = 0
	}
	nextnodes[source] = len(sortkeys)
	for _, n = range neighs[source] {
		flows[source][n] = caps[source][n]
		flows[n][source] = -caps[source][n]
		overflow[n] = caps[source][n]
		queue = append(queue, n)
		Debug("push %s\n", n)
	}

	for len(queue) > 0 {
		maxval = FindNextnodesMaxValue(nextnodes)
		n = queue[len(queue)-1]
		queue = queue[:(len(queue) - 1)]
		Debug("queue %v n %s\n", queue, n)
		Debug("flows %v nextnodes %v overflow %v\n", flows, nextnodes, overflow)
		if !CanPush(n, neighs, nextnodes, caps, flows) {
			Debug("push %s for nextnodes\n", n)
			nextnodes = SetNextNodes(n, neighs, nextnodes, caps, flows, maxval)
		}
		flows, overflow = FindNextNodes(n, neighs, nextnodes, caps, flows, overflow)

		if n != source && n != sink && overflow[n] > 0 {
			queue = append(queue, n)
			Debug("push %s to queue\n", n)
		}

		for _, k := range neighs[n] {
			if k != source && k != sink && overflow[k] > 0 {
				queue = append(queue, k)
				Debug("push %s to queue\n", k)
			}
		}

	}

	flow = 0
	for _, k = range neighs[source] {
		flow += flows[source][k]
		flow -= flows[k][source]
	}

	flow = flow >> 1
	return flow, flows
}
