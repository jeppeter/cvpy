package main

import (
	"time"
)

/*********************************************************
      this file is for the Goldberg-Tarjan algorithm for
      maxflow -mincut graph search
*********************************************************/

import (
	"os"
)

func FindNextnodesMaxValue(nextnodes *StringInt) int {
	maxval := 0
	for _, k := range nextnodes.Iter() {
		kv := nextnodes.GetValue(k)
		if kv > maxval {
			maxval = kv
		}
	}
	return maxval

}

func CanPush(n string, neighbours *Neigbour, nextnodes *StringInt, caps *StringGraph, flows *StringGraph) bool {
	for _, k := range neighbours.GetValue(n) {
		if ((nextnodes.GetValue(k) + 1) == nextnodes.GetValue(n)) && (caps.GetValue(n, k)-flows.GetValue(n, k)) > 0 {
			return true
		}
	}
	return false
}

func SetNextNodes(n string, neighbours *Neigbour, nextnodes *StringInt, caps *StringGraph, flows *StringGraph, maxval int) {
	minval := maxval
	for _, k := range neighbours.GetValue(n) {
		if (caps.GetValue(n, k) - flows.GetValue(n, k)) > 0 {
			if nextnodes.GetValue(n) < minval {
				minval = nextnodes.GetValue(n)
			}
		}
	}
	//Debug("set nextnodes[%s] = %d\n", n, (1 + minval))
	nextnodes.SetValue(n, 1+minval)
	return
}

func FindNextNodes(n string, neighbours *Neigbour, nextnodes *StringInt, caps *StringGraph,
	flows *StringGraph, overflow *StringInt) {

	for _, k := range neighbours.GetValue(n) {
		if (nextnodes.GetValue(k) + 1) == nextnodes.GetValue(n) {
			fval := (caps.GetValue(n, k) - flows.GetValue(n, k))
			if fval > overflow.GetValue(n) {
				fval = overflow.GetValue(n)
			}
			//Debug("Set [%s]->[%s] fval %d\n", n, k, fval)
			overflow.SetValue(k, overflow.GetValue(k)+fval)
			overflow.SetValue(n, overflow.GetValue(n)-fval)
			flows.SetValue(n, k, flows.GetValue(n, k)+fval)
			flows.SetValue(k, n, flows.GetValue(k, n)-fval)
		}
	}
	return

}

func GoldbergTarjan(caps *StringGraph, neighs *Neigbour, source string, sink string) (flow int, flows *StringGraph) {

	var n string
	var k string

	flow = 0
	flows = NewStringGraph()
	sortkeys := MakeSortKeys(caps)
	overflow := NewStringInt()
	nextnodes := NewStringInt()
	maxval := 0
	for _, k1 := range sortkeys {
		for _, k2 := range sortkeys {
			if k1 == k2 && caps.GetValue(k1, k2) != 0 {
				Debug("can not be set for %s %s value %d\n", k1, k2, caps.GetValue(k1, k2))
				os.Exit(4)
			}
			maxval += caps.GetValue(k1, k2)
		}
	}
	nextnodes.SetValue(source, len(sortkeys))
	queue := NewStringStack()
	for _, n = range neighs.GetValue(source) {
		flows.SetValue(source, n, caps.GetValue(source, n))
		flows.SetValue(n, source, -caps.GetValue(source, n))
		overflow.SetValue(n, caps.GetValue(source, n))
		queue.PushValue(n)
		//Debug("push %s\n", n)
	}

	stime := time.Now()

	for queue.Length() > 0 {
		maxval = FindNextnodesMaxValue(nextnodes)
		n = queue.PopValue()
		etime := time.Now()
		if etime.Sub(stime) > time.Second*2 {
			stime = etime
			Debug("time %s queue len(%d) n %s\n", etime, queue.Length(), n)
		}
		//Debug("queue %v n %s\n", queue, n)
		//Debug("flows %v nextnodes %v overflow %v\n", flows, nextnodes, overflow)
		if !CanPush(n, neighs, nextnodes, caps, flows) {
			//Debug("push %s for nextnodes\n", n)
			SetNextNodes(n, neighs, nextnodes, caps, flows, maxval)
		}
		FindNextNodes(n, neighs, nextnodes, caps, flows, overflow)

		if n != source && n != sink && overflow.GetValue(n) > 0 {
			queue.PushValue(n)
			//Debug("push %s to queue\n", n)
		}

		for _, k := range neighs.GetValue(n) {
			if k != source && k != sink && overflow.GetValue(k) > 0 {
				queue.PushValue(k)
				//Debug("push %s to queue\n", k)
			}
		}

	}

	flow = 0
	for _, k = range neighs.GetValue(source) {
		flow += flows.GetValue(source, k)
		flow -= flows.GetValue(k, source)
	}

	flow = flow >> 1
	return flow, flows
}
