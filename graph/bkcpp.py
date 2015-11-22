import sys
import logging


class Arc:
	def __init__(self):
		self.node_head = -1
		self.arc_next = -1
		self.arc_sister = -1
		self.r_cap = 0
		return

####################################################
#  node has 
#
####################################################
class Node:
	def __init__(self):
		self.arc_first = -1
		self.arc_parent = -1
		self.node_next = -1
		self.TS = 0
		self.DIST = 0
		self.is_sink = False
		self.is_marked = False
		self.is_in_changed_list = False
		self.tr_cap = 0
		return

class BKGraph:
	def __init__(self,nodemax,edgemax):
		self.node_num = 0
		self.node_block = []
		if nodemax < 16:
			nodemax = 16
		if edgemax < 16:
			edgemax = 16

		self.nodes = nodemax * [Node()]
		self.arcs = (2*edgemax) * [Arc()]
		self.node_last = self.nodes
		self.nodes_max = nodemax
		self.arcs_last = self.arcs
		self.arcs_max = 2*edgemax
		self.flow = 0
		self.maxflow_iteration = 0
		return

	def add_node(self,num=1):
		for i in xrange(num):
			self.nodes.append(Node())
		nodemax += num
		return

	def add_tweights(self,nodeid,cap_source,cap_sink):
		assert(len(self.nodes) > nodeid)
		delta = self.nodes[nodeid].tr_cap
		if delta > 0 :
			cap_source += delta
		else:
			cap_sink -= delta
		if cap_source < cap_sink:
			self.flow += cap_source
		else:
			self.flow += cap_sink
		self.nodes[nodeid].tr_cap = cap_source - cap_sink
		return

	def add_edge(self,nodeidi,nodeidj,cap ,rev_cap):
		assert(self.nodemax > nodeidi)
		assert(self.nodemax > nodeidj)
		assert(nodeidi != nodeidj)
		assert(cap >= 0)
		assert(rev_cap >= 0)
		aidx = len(self.arcs)
		arevidx = aidx + 1
		# now to add the idx 
		aarc = Arc()
		arevarc = Arc()
		nodei = self.nodes[nodeidi]
		nodej = self.nodes[nodeidj]

		# now set for the idx
		aarc.arc_sister = arevidx
		arevarc.arc_sister = aidx
		aarc.node_next = nodei.arc_first
		nodei.arc_first = aidx
		arevarc.arc_next = nodej.arc_first
		nodej.arc_first = arevidx
		aarc.node_head = nodeidj
		arevarc.node_head = nodeidi
		aarc.r_cap = rev_cap
		arevarc.r_cap = rev_cap

		# now to set for the arc
		self.arcs.append(aarc)
		self.arcs.append(arevarc)
		self.nodes[nodeidi] = nodei
		self.nodes[nodeidj] = nodej
		return

	def max_flow(self):
		self.maxflow_reuse_trees_init()

		curnodeid = -1
		nodei = None
		while True:
			if curnodeid == -1:
				nodei = -1
			else:
				nodei = curnodeid
				self.nodes[nodei].arc_next = -1
				if self.nodes[nodei].arc_parent == -1:
					nodei = -1
			if nodei == -1:
				nodei = self.next_active()
				if nodei == -1:
					break

			if not self.nodes[nodei].is_sink:
				aidx = self.nodes[nodei].arc_first
				while aidx != -1:
					if self.arcs[aidx].r_cap:
						nodej = self.arcs[aidx].node_head
						if self.nodes[nodej].arc_parent == -1:
							self.nodes[nodej].is_sink = 0
							self.nodes[nodej].arc_parent = self.arcs[aidx].arc_sister
							self.nodes[nodej].TS = self.nodes[nodei].TS
							self.nodes[nodej].DIST = self.nodes[nodei].DIST + 1
							self.set_active(nodej)
							self.add_to_change_list(nodej)
						elif self.nodes[nodej].is_sink:
							break
						elif self.nodes[nodej].TS <= self.nodes[nodei].TS and \
							self.nodes[nodej].DIST > self.nodes[nodei].DIST :
							self.nodes[nodej].arc_parent = self.arcs[aidx].arc_sister
							self.nodes[nodej].TS =self.nodes[nodei].TS
							self.nodes[nodej].DIST = self.nodes[nodei].DIST
					aidx = self.arcs[aidx].arc_next
			else:
				aidx = self.nodes[nodei].arc_first
				while aidx != -1:
					sisidx = self.arcs[aidx].arc_sister
					if self.arcs[sisidx].r_cap:
						nodej = self.arcs[aidx].node_head
						if self.nodes[nodej].arc_parent == -1:
							self.nodes[nodej].is_sink = True
							self.nodes[nodej].arc_parent = self.arcs[aidx].arc_sister
							self.nodes[nodej].TS = self.nodes[nodei].TS
							self.nodes[nodej].DIST = self.nodes[nodei].DIST + 1
							self.set_active(nodej)
							self.add_to_change_list(nodej)
						elif not self.nodes[nodej].is_sink :
							aidx = self.arcs[aidx].arc_sister
							break
						elif self.nodes[nodej].TS <= self.nodes[nodei].TS and \
							self.nodes[nodej].DIST > self.nodes[nodei].DIST:
							self.nodes[nodej].arc_parent = self.arcs[aidx].arc_sister
							self.nodes[nodej].TS = self.nodes[nodei].TS
							self.nodes[nodej].DIST = self.nodes[nodei].DIST + 1

					aidx = self.arcs[aidx].arc_next

			self.TIME += 1

			if aidx != -1:
				self.nodes[nodei].node_next = nodei
				curnodeid = nodei

				self.augment(aidx)

				while True:
					np = self.orphan_first
					if np == -1:
						break







