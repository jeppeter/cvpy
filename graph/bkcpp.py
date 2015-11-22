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
