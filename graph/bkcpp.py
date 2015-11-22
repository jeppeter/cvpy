import sys
import logging


class Arc:
	def __init__(self):
		return

class Node:
	def __init__(self):
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
