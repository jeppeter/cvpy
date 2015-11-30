import sys
import logging

NULL_PTR=-1
MAXFLOW_TERMINAL=-2
MAXFLOW_ORPHAN=-3
MAXFLOW_INFINITE_D=(0xffffffff >> 1)
CPP_OUT=1


def GetIdx(idx):
	if idx == NULL_PTR:
		return 'NULL'
	elif idx == MAXFLOW_ORPHAN:
		return 'MAXFLOW_ORPHAN'
	elif idx == MAXFLOW_TERMINAL:
		return 'MAXFLOW_TERMINAL'
	else:
		return '%d'%(idx)

class Arc:
	def __init__(self):
		self.node_head = NULL_PTR
		self.arc_next = NULL_PTR
		self.arc_sister = NULL_PTR
		self.r_cap = 0
		self.name = ''
		return

####################################################
#  node has 
#
####################################################
class Node:
	def __init__(self):
		self.arc_first = NULL_PTR
		self.arc_parent = NULL_PTR
		self.node_next = NULL_PTR
		self.TS = 0
		self.DIST = 0
		self.is_sink = False
		self.tr_cap = 0
		self.name =''
		return

class NodeBlockPtr:
	def __init__(self):
		self.array_node = NULL_PTR
		return




class BKGraph:
	def __init__(self,nodemax,edgemax):
		self.nodes = []
		self.arcs = []
		self.flow = 0
		self.maxflow_iteration = 0
		self.orphan_list = []
		self.queue_first = NULL_PTR
		self.queue_last = NULL_PTR
		return

	def get_arc_name(self,aidx):
		if aidx == MAXFLOW_TERMINAL :
			return "MAXFLOW_TERMINAL"
		elif aidx == MAXFLOW_ORPHAN :
			return "MAXFLOW_ORPHAN"
		elif aidx == NULL_PTR:
			return "NULL"
		return self.arcs[aidx].name

	def get_node_name(self,nodei):
		if nodei == NULL_PTR:
			return "NULL"
		return self.nodes[nodei].name
	def GetOrphanList(self,orphlist):
		i = 0
		s = 'cnt(%d)['%(len(orphlist))
		for a in orphlist:
			if i != 0:
				s += ','
			s += '%s'%(self.get_node_name(a.array_node))
			i += 1
		s += ']'
		return s

	def add_node(self,num=1):
		for i in xrange(num):
			self.nodes.append(Node())
		return

	def sort_node_arcs(self,nodei):
		arcidxs = []
		arcnames = []
		aidx = self.nodes[nodei].arc_first
		logging.info('aidx %s'%(self.get_arc_name(aidx)))
		while aidx != NULL_PTR:
			arcidxs.append(aidx)
			arcnames.append(self.arcs[aidx].name)
			aidx = self.arcs[aidx].arc_next
		logging.debug('arcidxs %s'%(arcidxs))
		if len(arcidxs) <= 1 :
			return
		for i in xrange(len(arcidxs)):
			for j in range((i+1),len(arcidxs)):
				if arcnames[i] > arcnames[j]:
					tmp = arcnames[i]
					arcnames[i] = arcnames[j]
					arcnames[j] = tmp
					tmp = arcidxs[i]
					arcidxs[i] = arcidxs[j]
					arcidxs[j] = tmp
		self.nodes[nodei].arc_first = arcidxs[0]
		for i in range(0,len(arcidxs)-1):
			self.arcs[arcidxs[i]].arc_next = arcidxs[(i+1)]
		self.arcs[arcidxs[-1]].arc_next = NULL_PTR
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
		logging.info('node[%s].tr_cap (%d -> %d) flow(%d)'%(self.get_node_name(nodeid),self.nodes[nodeid].tr_cap,cap_source - cap_sink,self.flow))
		self.nodes[nodeid].tr_cap = cap_source - cap_sink
		return

	def add_edge(self,nodeidi,nodeidj,cap ,rev_cap):
		assert(len(self.nodes) > nodeidi)
		assert(len(self.nodes) > nodeidj)
		assert(nodeidi != nodeidj)
		assert(cap >= 0)
		assert(rev_cap >= 0)
		aidx = len(self.arcs)
		arevidx = aidx + 1
		# now to add the idx 
		aarc = Arc()
		arevarc = Arc()
		nodei = self.nodes[nodeidi]
		nodei.name = '%d'%(nodeidi)
		nodej = self.nodes[nodeidj]
		nodej.name = '%d'%(nodeidj)

		# now set for the idx
		aarc.arc_sister = arevidx
		arevarc.arc_sister = aidx
		aarc.arc_next = nodei.arc_first
		nodei.arc_first = aidx
		arevarc.arc_next = nodej.arc_first
		nodej.arc_first = arevidx
		aarc.node_head = nodeidj
		aarc.name = '%s->%s'%(nodeidi,nodeidj)
		arevarc.node_head = nodeidi
		arevarc.name = '%s->%s'%(nodeidj,nodeidi)
		aarc.r_cap = cap
		arevarc.r_cap = rev_cap

		# now to set for the arc
		self.arcs.append(aarc)
		self.arcs.append(arevarc)
		self.nodes[nodeidi] = nodei
		self.nodes[nodeidj] = nodej
		return

	def max_flow(self):
		self.maxflow_init()
		curnodeid = NULL_PTR
		nodei = None
		while True:
			nodei = curnodeid
			if nodei != NULL_PTR:
				logging.info('node[%s].next (%s -> %s)'%(self.get_node_name(nodei),self.get_node_name(self.nodes[nodei].node_next),self.get_node_name(NULL_PTR)))
				self.nodes[nodei].node_next = NULL_PTR
				if self.nodes[nodei].arc_parent == NULL_PTR:
					logging.info('node[%s].parent (NULL) nodei -> NULL'%(self.get_node_name(nodei)))
					nodei = NULL_PTR
			if nodei == NULL_PTR:
				nodei = self.next_active()
				if nodei == NULL_PTR:
					break
			logging.info('nodei node[%s].is_sink %s'%(self.get_node_name(nodei),self.nodes[nodei].is_sink))

			if not self.nodes[nodei].is_sink:
				aidx = self.nodes[nodei].arc_first
				logging.info('aidx node[%s].first (%s)'%(self.get_node_name(nodei),self.get_arc_name(aidx)))
				while aidx != NULL_PTR:
					logging.info('arc[%s].r_cap (%d)'%(self.get_arc_name(aidx),self.arcs[aidx].r_cap))
					if self.arcs[aidx].r_cap:
						nodej = self.arcs[aidx].node_head
						logging.info('arc[%s].head -> node[%s].parent (%s)'%(self.get_arc_name(aidx),self.get_node_name(nodej),self.get_arc_name(self.nodes[nodej].arc_parent)))
						if self.nodes[nodej].arc_parent == NULL_PTR:
							logging.info('node[%s].is_sink (%s -> False)'%(self.get_node_name(nodej),self.nodes[nodej].is_sink))
							self.nodes[nodej].is_sink = False
							logging.info('node[%s].parent (%s -> %s)'%(self.get_node_name(nodej),self.get_arc_name(self.nodes[nodej].arc_parent),self.get_arc_name(self.arcs[aidx].arc_sister)))			
							self.nodes[nodej].arc_parent = self.arcs[aidx].arc_sister
							logging.info('node[%s].TS (%d -> node[%s].TS %d)'%(self.get_node_name(nodej),self.nodes[nodej].TS,self.get_node_name(nodei),self.nodes[nodei].TS))
							self.nodes[nodej].TS = self.nodes[nodei].TS
							logging.info('node[%s].DIST (%d -> node[%s].DIST+1 %d)'%(self.get_node_name(nodej),self.nodes[nodej].DIST,self.get_node_name(nodei),self.nodes[nodei].DIST + 1))
							self.nodes[nodej].DIST = self.nodes[nodei].DIST + 1
							self.set_active(nodej)
							self.add_to_change_list(nodej)
						elif self.nodes[nodej].is_sink:
							logging.info('node[%s].is_sink (%s)'%(self.get_node_name(nodej),self.nodes[nodej].is_sink))
							break
						elif self.nodes[nodej].TS <= self.nodes[nodei].TS and \
							self.nodes[nodej].DIST > self.nodes[nodei].DIST :
							logging.info('node[%s].TS (%d) <= node[%s].TS (%d)'%(self.get_node_name(nodej),self.nodes[nodej].TS,self.get_node_name(nodei),self.nodes[nodei].TS))
							logging.info('node[%s].DIST (%d) > node[%s].DIST (%d)'%(self.get_node_name(nodej),self.nodes[nodej].DIST,self.get_node_name(nodei),self.nodes[nodei].DIST))
							logging.info('node[%s].parent (%s -> %s)'%(self.get_node_name(nodej),self.get_arc_name(self.nodes[nodej].arc_parent),self.get_arc_name(self.arcs[aidx].arc_sister)))
							self.nodes[nodej].arc_parent = self.arcs[aidx].arc_sister
							logging.info('node[%s].TS (%d -> node[%s].TS %d)'%(self.get_node_name(nodej),self.nodes[nodej].TS,self.get_node_name(nodei),self.nodes[nodei].TS))
							self.nodes[nodej].TS =self.nodes[nodei].TS
							logging.info('node[%s].DIST (%d -> node[%s].DIST+1 %d)'%(self.get_node_name(nodej),self.nodes[nodej].DIST,self.get_node_name(nodei),self.nodes[nodei].DIST + 1))
							self.nodes[nodej].DIST = self.nodes[nodei].DIST + 1
					logging.info('aidx(%s -> arc[%s].next %s)'%(self.get_arc_name(aidx),self.get_arc_name(aidx),self.get_arc_name(self.arcs[aidx].arc_next)))
					aidx = self.arcs[aidx].arc_next

			else:
				aidx = self.nodes[nodei].arc_first
				logging.info('aidx (node[%s].first (%s))'%(self.get_node_name(nodei),self.get_arc_name(aidx)))
				while aidx != NULL_PTR:					
					sisidx = self.arcs[aidx].arc_sister
					logging.info('arc[%s].sister -> arc[%s].r_cap (%d)'%(self.get_arc_name(aidx),self.get_arc_name(sisidx),self.arcs[sisidx].r_cap))
					if self.arcs[sisidx].r_cap:
						nodej = self.arcs[aidx].node_head
						logging.info('nodej (arc[%s].head (%s))'%(self.get_arc_name(aidx),self.get_node_name(nodej)))
						logging.info('node[%s].parent (%s)'%(self.get_node_name(nodej),self.get_arc_name(self.nodes[nodej].arc_parent)))
						logging.info('node[%s].TS (%d) ?<= node[%s].TS (%d)'%(self.get_node_name(nodej),self.nodes[nodej].TS,self.get_node_name(nodei),self.nodes[nodei].TS))
						logging.info('node[%s].DIST (%d) ?> node[%s].DIST (%d)'%(self.get_node_name(nodej),self.nodes[nodej].DIST,self.get_node_name(nodei),self.nodes[nodei].DIST))
						if self.nodes[nodej].arc_parent == NULL_PTR:
							logging.info('node[%s].is_sink (%s -> True)'%(self.get_node_name(nodej),self.nodes[nodej].is_sink))
							self.nodes[nodej].is_sink = True
							logging.info('node[%s].parent (%s -> arc[%s].sister %s)'%(self.get_node_name(nodej),self.get_arc_name(self.nodes[nodej].arc_parent),self.get_arc_name(aidx),self.get_arc_name(self.arcs[aidx].arc_sister)))
							self.nodes[nodej].arc_parent = self.arcs[aidx].arc_sister
							logging.info('node[%s].TS (%d -> node[%s].TS %d)'%(self.get_node_name(nodej),self.nodes[nodej].TS,self.get_node_name(nodei),self.nodes[nodei].TS))
							self.nodes[nodej].TS = self.nodes[nodei].TS
							logging.info('node[%s].DIST (%d -> node[%s].DIST+1 %d)'%(self.get_node_name(nodej),self.nodes[nodej].DIST,self.get_node_name(nodei),self.nodes[nodei].DIST + 1))
							self.nodes[nodej].DIST = self.nodes[nodei].DIST + 1
							self.set_active(nodej)
							self.add_to_change_list(nodej)
						elif not self.nodes[nodej].is_sink :
							logging.info('aidx (%s -> arc[%s].sister %s)'%(self.get_arc_name(aidx),self.get_arc_name(aidx),self.get_arc_name(self.arcs[aidx].arc_sister)))
							aidx = self.arcs[aidx].arc_sister
							break
						elif self.nodes[nodej].TS <= self.nodes[nodei].TS and \
							self.nodes[nodej].DIST > self.nodes[nodei].DIST:
							logging.info('node[%s].parent (%s -> %s)'%(self.get_node_name(nodej),self.get_arc_name(self.nodes[nodej].arc_parent),self.get_arc_name(self.arcs[aidx].arc_sister)))
							self.nodes[nodej].arc_parent = self.arcs[aidx].arc_sister
							logging.info('node[%s].TS (%d -> node[%s].TS %d)'%(self.get_node_name(nodej),self.nodes[nodej].TS,self.get_node_name(nodei),self.nodes[nodei].TS))
							self.nodes[nodej].TS = self.nodes[nodei].TS
							logging.info('node[%s].DIST (%d -> node[%s].DIST+1 %d)'%(self.get_node_name(nodej),self.nodes[nodej].DIST,self.get_node_name(nodei),self.nodes[nodei].DIST + 1))
							self.nodes[nodej].DIST = self.nodes[nodei].DIST + 1
					logging.info('aidx (%s -> arc[%s].next %s)'%(self.get_arc_name(aidx),self.get_arc_name(aidx),self.get_arc_name(self.arcs[aidx].arc_next)))
					aidx = self.arcs[aidx].arc_next

			self.TIME += 1
			logging.info('TIME %d arc[%s]'%(self.TIME,self.get_arc_name(aidx)))
			self.debug_state('after arcs handle(%d)'%(self.TIME))

			if aidx != NULL_PTR:
				logging.info('node[%s].next (%s -> %s)'%(self.get_node_name(nodei),self.get_node_name(self.nodes[nodei].node_next),self.get_node_name(nodei)))
				self.nodes[nodei].node_next = nodei
				logging.info('curnodeid (%s -> %s)'%(self.get_node_name(curnodeid),self.get_node_name(nodei)))
				curnodeid = nodei

				self.augment(aidx)
				self.debug_state('after augment(%d)'%(self.TIME))

				while len(self.orphan_list) > 0:
					curorphan = self.orphan_list[0]
					self.orphan_list = self.orphan_list[1:]
					curorphnodei = curorphan.array_node
					if self.nodes[curorphnodei].is_sink:
						logging.info('sink orphan %d'%(curorphnodei))
						self.process_sink_orphan(curorphnodei)
						logging.info('sink orphan over %s'%(self.get_node_name(curorphnodei)))
					else:
						logging.info('source orphan %d'%(curorphnodei))
						self.process_source_orphan(curorphnodei)
						logging.info('source orphan over %d'%(curorphnodei))
				self.debug_state('after orphan handle (%d)'%(self.TIME))
				logging.info('curnodeid %s'%(self.get_node_name(curnodeid)))
			else:
				curnodeid = NULL_PTR
		self.maxflow_iteration += 1
		return self.flow

	def get_first_link(self,nodei):
		s = '['
		i = 0
		aidx = self.nodes[nodei].arc_first
		while aidx != NULL_PTR:
			if i != 0:
				s += ','
			s += '%s'%(self.get_arc_name(aidx))
			aidx = self.arcs[aidx].arc_next
			i += 1
		s += ']cnt(%d)'%(i)
		return s

	def get_arc_parent(self,nodei):
		s = '['
		i = 0
		aidx = self.nodes[nodei].arc_parent
		while aidx != NULL_PTR:
			if i != 0:
				s += ','
			if aidx == MAXFLOW_TERMINAL:
				s += 'MAXFLOW_TERMINAL'
				break
			if aidx == MAXFLOW_ORPHAN:
				s += 'MAXFLOW_ORPHAN'
				break
			nodej = self.arcs[aidx].node_head
			s += '%s(%s)'%(self.get_node_name(nodej),self.get_arc_name(aidx))
			i += 1
			if nodej == NULL_PTR:
				break
			aidx = self.nodes[nodej].arc_parent
		s += ']cnt(%d)'%(i)
		return s

	def get_node_next(self,nodei):
		s = '['
		i = 0
		nodej = self.nodes[nodei].node_next
		while nodej != NULL_PTR:
			if i != 0:
				s += ','
			i += 1
			s += '%s'%(self.get_node_name(nodej))
			if nodej == self.nodes[nodej].node_next:
				break
			nodej = self.nodes[nodej].node_next
		s += ']cnt(%d)'%(i)
		return s

	def get_arc_next(self,aidx):
		s = '['
		i = 0
		anext = self.arcs[aidx].arc_next
		while anext != NULL_PTR:
			if i != 0:
				s += ','
			i += 1
			s += '%s'%(self.get_arc_name(anext))
			anext = self.arcs[anext].arc_next
		s += ']cnt(%d)'%(i)
		return s

	def debug_node(self,nodei):
		logging.debug('==============================')
		logging.debug('node[%s].is_sink (%s)'%(self.get_node_name(nodei),self.nodes[nodei].is_sink))
		logging.debug('node[%s].arc_first list(%s)'%(self.get_node_name(nodei),self.get_first_link(nodei)))
		logging.debug('node[%s].arc_parent list(%s)'%(self.get_node_name(nodei),self.get_arc_parent(nodei)))
		logging.debug('node[%s].node_next list(%s)'%(self.get_node_name(nodei),self.get_node_next(nodei)))
		#logging.debug('node[%s].tr_cap (%d)'%(self.get_node_name(nodei),self.nodes[nodei].tr_cap))
		logging.debug('node[%s].TS (%d) node[%s].DIST (%d)'%(self.get_node_name(nodei),self.nodes[nodei].TS,self.get_node_name(nodei),self.nodes[nodei].DIST))
		logging.debug('******************************')
		return

	def debug_arc(self,aidx):
		#logging.debug('+++++++++++++++++++++++++++++++')
		#logging.debug('arc[%s].node_head (%s)'%(self.get_arc_name(aidx),self.get_node_name(self.arcs[aidx].node_head)))
		#logging.debug('arc[%s].arc_next list(%s)'%(self.get_arc_name(aidx),self.get_arc_next(aidx)))
		#logging.debug('arc[%s].arc_sister (%s)'%(self.get_arc_name(aidx),self.get_arc_name(self.arcs[aidx].arc_sister)))
		logging.debug('arc[%s].r_cap (%d)'%(self.get_arc_name(aidx),self.arcs[aidx].r_cap))
		#logging.debug('-------------------------------')
		return

	def debug_queue_state(self,notice,q):
		s =''
		if q == NULL_PTR:
			s += 'NULL'
		else:
			i = 0
			s += '['
			while q != NULL_PTR:
				if i != 0 :
					s += ','
				i += 1
				s += '%s'%(self.get_node_name(q))
				if q == self.nodes[q].node_next:
					break
				q = self.nodes[q].node_next
			s +=']cnt(%d)'%(i)
		logging.debug('%s list(%s)'%(notice,s))
		return
	def sort_nodes_names(self):
		nodeids = []
		nodenames = []
		for i in xrange(len(self.nodes)):
			nodeids.append(i)
			nodenames.append(self.nodes[i].name)

		for i in xrange(len(nodenames)):
			for j in range((i+1),len(nodenames)):
				if nodenames[i] > nodenames[j]:
					tmp = nodenames[i]
					nodenames[i] = nodenames[j]
					nodenames[j] = tmp
					tmp = nodeids[i]
					nodeids[i] = nodeids[j]
					nodeids[j] = tmp
		return nodeids

	def sort_arc_names(self):
		arcids = []
		arcnames = []
		for i in xrange(len(self.arcs)):
			arcids.append(i)
			arcnames.append(self.arcs[i].name)

		for i in xrange(len(arcids)):
			for j in range((i+1),len(arcids)):
				if arcnames[i] > arcnames[j]:
					tmp = arcnames[i]
					arcnames[i] = arcnames[j]
					arcnames[j] = tmp
					tmp = arcids[i]
					arcids[i] = arcids[j]
					arcids[j] = tmp
		return arcids

	def debug_state(self,notice):
		logging.debug('~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~')
		logging.debug('debug state %s'%(notice))
		for nodei in self.sort_nodes_names():
			if self.nodes[nodei].name == '':
				continue
			self.debug_node(nodei)

		for aidx in self.sort_arc_names():
			self.debug_arc(aidx)
		self.debug_queue_state('queue_first',self.queue_first)
		logging.debug('TIME (%d) flow (%d)'%(self.TIME,self.flow))
		logging.debug('orphan_list (%s)'%(self.GetOrphanList(self.orphan_list)))
		logging.debug('~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~')

	def maxflow_init(self):
		self.queue_first = NULL_PTR
		self.queue_last = NULL_PTR
		self.orphan_list = []
		self.TIME = 0

		for nodei in xrange(len(self.nodes)):
			self.nodes[nodei].node_next = NULL_PTR
			self.nodes[nodei].TS = self.TIME
			self.sort_node_arcs(nodei)

			if self.nodes[nodei].tr_cap > 0:
				self.nodes[nodei].is_sink = False
				self.nodes[nodei].arc_parent = MAXFLOW_TERMINAL
				self.set_active(nodei)
				self.nodes[nodei].DIST = 1
			elif self.nodes[nodei].tr_cap < 0 :
				self.nodes[nodei].is_sink = True
				self.nodes[nodei].arc_parent = MAXFLOW_TERMINAL
				self.set_active(nodei)
				self.nodes[nodei].DIST = 1
			else:
				self.nodes[nodei].arc_parent = NULL_PTR

		self.debug_state('after init')

		return

	def next_active(self):
		while True:
			nodei = self.queue_first
			logging.info('nodei (queue_first (%s))'%(self.get_node_name(self.queue_first)))
			if nodei == NULL_PTR:
				return NULL_PTR
			logging.info('node[%s].next (%s)'%(self.get_node_name(nodei),self.get_node_name(self.nodes[nodei].node_next)))
			if self.nodes[nodei].node_next == nodei:
				logging.info('queue_first (%s -> %s)'%(self.get_node_name(self.queue_first),self.get_node_name(NULL_PTR)))
				self.queue_first = NULL_PTR
				logging.info('queue_last (%s -> %s)'%(self.get_node_name(self.queue_last),self.get_node_name(NULL_PTR)))
				self.queue_last = NULL_PTR
			else:
				logging.info('queue_first (%s -> node[%s].next %s)'%(self.get_node_name(self.queue_first),self.get_node_name(nodei),self.get_node_name(self.nodes[nodei].node_next)))
				self.queue_first= self.nodes[nodei].node_next
			logging.info('node[%s].next (%s -> NULL)'%(self.get_node_name(nodei),self.get_node_name(self.nodes[nodei].node_next)))
			self.nodes[nodei].node_next = NULL_PTR

			logging.info('node[%s].parent (%s)'%(self.get_node_name(nodei),self.get_arc_name(self.nodes[nodei].arc_parent)))
			if self.nodes[nodei].arc_parent != NULL_PTR:
				return nodei
		return NULL_PTR

	def set_active(self,nodei):
		logging.info('set node[%s].next (%s) active'%(self.get_node_name(nodei),self.get_node_name(self.nodes[nodei].node_next)))
		logging.info('queue_last (%s) queue_first (%s)'%(self.get_node_name(self.queue_last),self.get_node_name(self.queue_first)))
		if self.nodes[nodei].node_next == NULL_PTR:
			if self.queue_last != NULL_PTR:
				logging.info('set node[%s].next (%s -> %s)'%(self.get_node_name(self.queue_last),self.get_node_name(self.nodes[self.queue_last].node_next),self.get_node_name(nodei)))
				self.nodes[self.queue_last].node_next = nodei
			else:
				logging.info('set queue_first (%s -> %s)'%(self.get_node_name(self.queue_first),self.get_node_name(nodei)))
				self.queue_first = nodei
			logging.info('set queue_last (%s -> %s)'%(self.get_node_name(self.queue_last),self.get_node_name(nodei)))
			self.queue_last = nodei
			logging.info('set node[%s].next (%s -> %s)'%(self.get_node_name(nodei),self.get_node_name(self.nodes[nodei].node_next),self.get_node_name(nodei)))
			self.nodes[nodei].node_next = nodei
		return

	def add_to_change_list(self,nodei):
		# nothing to add change list because we do not set change list before
		return

	def augment(self,aidx):		
		bottlecap = self.arcs[aidx].r_cap
		logging.info('arc[%s].r_cap (%d) bottlecap (%d)'%(self.get_arc_name(aidx),bottlecap,bottlecap))
		# this is source tree
		sisidx = self.arcs[aidx].arc_sister
		nodei = self.arcs[sisidx].node_head
		logging.info('arc[%s].sister (%s) arc[%s].head nodei (%s)'%(self.get_arc_name(aidx),self.get_arc_name(sisidx),self.get_arc_name(sisidx),self.get_node_name(nodei)))
		while True:
			arca = self.nodes[nodei].arc_parent
			logging.info('node[%s].parent (%s)'%(self.get_node_name(nodei),self.get_arc_name(arca)))
			if arca == MAXFLOW_TERMINAL:
				break
			sisidx = self.arcs[arca].arc_sister
			logging.info('arc[%s].sister arc[%s].r_cap (%d) bottlecap(%d)'%(self.get_arc_name(arca),self.get_arc_name(sisidx),self.arcs[sisidx].r_cap,bottlecap))
			if bottlecap > self.arcs[sisidx].r_cap:
				logging.info('bottlecap (%d -> %d)'%(bottlecap,self.arcs[sisidx].r_cap))
				bottlecap = self.arcs[sisidx].r_cap
			logging.info('nodei (%s -> arc[%s].head (%s))'%(self.get_node_name(nodei),self.get_arc_name(arca),self.get_node_name(self.arcs[arca].node_head)))
			nodei = self.arcs[arca].node_head
		logging.info('node[%s].tr_cap (%d) bottlecap(%d)'%(self.get_node_name(nodei),self.nodes[nodei].tr_cap,bottlecap))
		if bottlecap > self.nodes[nodei].tr_cap:
			logging.info('bottlecap (%d -> %d)'%(bottlecap,self.nodes[nodei].tr_cap))
			bottlecap = self.nodes[nodei].tr_cap

		# this is sink tree
		nodei = self.arcs[aidx].node_head
		logging.info('nodei (arc[%s].head (%s))'%(self.get_arc_name(aidx),self.get_node_name(nodei)))
		while True:
			arca = self.nodes[nodei].arc_parent
			logging.info('node[%s].parent (%s)'%(self.get_node_name(nodei),self.get_arc_name(arca)))
			if arca == MAXFLOW_TERMINAL:
				break
			logging.info('arc[%s].r_cap (%d) bottlecap(%d)'%(self.get_arc_name(arca),self.arcs[arca].r_cap,bottlecap))
			if bottlecap > self.arcs[arca].r_cap :
				bottlecap = self.arcs[arca].r_cap
			logging.info('nodei (%s -> arc[%s].head = %s)'%(self.get_node_name(nodei),self.get_arc_name(arca),self.get_node_name(self.arcs[arca].node_head)))
			nodei = self.arcs[arca].node_head
		logging.info('node[%s].tr_cap (%d) bottlecap(%d)'%(self.get_node_name(nodei),self.nodes[nodei].tr_cap,bottlecap))
		if bottlecap > - self.nodes[nodei].tr_cap :
			bottlecap = - self.nodes[nodei].tr_cap

		sisidx = self.arcs[aidx].arc_sister
		logging.info('arc[%s].sister -> arc[%s].r_cap(%d+%d) arc[%s].r_cap(%d-%d)'%(self.get_arc_name(aidx),self.get_arc_name(sisidx),self.arcs[sisidx].r_cap,bottlecap,self.get_arc_name(aidx),self.arcs[aidx].r_cap,bottlecap))
		self.arcs[sisidx].r_cap += bottlecap
		self.arcs[aidx].r_cap -= bottlecap

		nodei = self.arcs[sisidx].node_head
		while True:
			arca = self.nodes[nodei].arc_parent
			logging.info('arca (node[%s].parent (%s))'%(self.get_node_name(nodei),self.get_arc_name(arca)))
			if arca == MAXFLOW_TERMINAL:
				break
			sisidx = self.arcs[arca].arc_sister
			logging.info('arc[%s].r_cap (%d+%d) arc[%s].sister -> arc[%s].r_cap(%d-%d)'%(self.get_arc_name(arca),self.arcs[arca].r_cap,bottlecap,self.get_arc_name(arca),self.get_arc_name(sisidx),self.arcs[sisidx].r_cap,bottlecap))
			self.arcs[arca].r_cap += bottlecap
			self.arcs[sisidx].r_cap -= bottlecap
			if self.arcs[sisidx].r_cap == 0 :
				logging.info('nodei[%s] -> arc[%s] set_orphan_front'%(self.get_node_name(nodei),self.get_arc_name(sisidx)))
				self.set_orphan_front(nodei)
			logging.info('nodei (%s -> arc[%s].head (%s))'%(self.get_node_name(nodei),self.get_arc_name(arca),self.get_node_name(self.arcs[arca].node_head)))
			nodei = self.arcs[arca].node_head

		logging.info('node[%s].tr_cap (%d-%d)'%(self.get_node_name(nodei),self.nodes[nodei].tr_cap,bottlecap))
		self.nodes[nodei].tr_cap -= bottlecap

		if self.nodes[nodei].tr_cap == 0:
			logging.info('node[%s] set_orphan_front'%(self.get_node_name(nodei)))
			self.set_orphan_front(nodei)

		nodei = self.arcs[aidx].node_head
		logging.info('arc[%s].head (%s)'%(self.get_arc_name(aidx),self.get_node_name(nodei)))
		while True:
			arca = self.nodes[nodei].arc_parent
			logging.info('arca (node[%s].parent (%s))'%(self.get_node_name(nodei),self.get_arc_name(arca)))
			if arca == MAXFLOW_TERMINAL:
				break
			sisidx = self.arcs[arca].arc_sister
			logging.info('arc[%s].r_cap (%d+%d) arc[%s].sister -> arc[%s].r_cap(%d-%d)'%(self.get_arc_name(arca),self.arcs[arca].r_cap,bottlecap,self.get_arc_name(arca),self.get_arc_name(sisidx),self.arcs[sisidx].r_cap,bottlecap))
			self.arcs[sisidx].r_cap += bottlecap
			self.arcs[arca].r_cap -= bottlecap
			if self.arcs[arca].r_cap == 0 :
				logging.info('arc[%s] set_orphan_front'%(self.get_arc_name(sisidx)))
				self.set_orphan_front(nodei)
			nodei = self.arcs[arca].node_head
			logging.info('arc[%s].head = %s'%(self.get_arc_name(arca),self.get_node_name(nodei)))
		logging.info('node[%s].tr_cap (%d+%d)'%(self.get_node_name(nodei),self.nodes[nodei].tr_cap,bottlecap))
		self.nodes[nodei].tr_cap += bottlecap
		if self.nodes[nodei].tr_cap == 0:
			logging.info('node[%s] set_orphan_front'%(self.get_node_name(nodei)))
			self.set_orphan_front(nodei)
		logging.info('flow (%d+%d)'%(self.flow,bottlecap))
		self.flow += bottlecap
		return

	def process_sink_orphan(self,nodei):
		d_min = MAXFLOW_INFINITE_D
		arc0 = self.nodes[nodei].arc_first
		arc0_min = NULL_PTR
		logging.info('node[%s].first (%s)'%(self.get_node_name(nodei),self.get_arc_name(arc0)))
		while arc0 != NULL_PTR:
			logging.info('arc[%s].r_cap (%d)'%(self.get_arc_name(arc0),self.arcs[arc0].r_cap))
			if self.arcs[arc0].r_cap != 0 :
				nodej = self.arcs[arc0].node_head
				logging.info('arc[%s].head (%s) is_sink %s'%(self.get_arc_name(arc0),self.get_node_name(nodej),self.nodes[nodej].is_sink))
				if self.nodes[nodej].is_sink:
					arca = self.nodes[nodej].arc_parent
					logging.info('node[%s].parent (%s)'%(self.get_node_name(nodej),self.get_arc_name(arca)))
					if arca != NULL_PTR:
						d = 0
						while True:
							if self.nodes[nodej].TS == self.TIME:
								d += self.nodes[nodej].DIST
								break
							arca = self.nodes[nodej].arc_parent
							d += 1
							if arca == MAXFLOW_TERMINAL:
								self.nodes[nodej].TS = self.TIME
								self.nodes[nodej].DIST = 1
								break
							if arca == MAXFLOW_ORPHAN:
								d = MAXFLOW_INFINITE_D
								break
							nodej = self.arcs[arca].node_head
						if d < MAXFLOW_INFINITE_D:
							if d < d_min:
								arc0_min = arc0
								d_min = d

							nodej = self.arcs[arc0].node_head
							while self.nodes[nodej].TS != self.TIME:
								self.nodes[nodej].TS = self.TIME
								self.nodes[nodej].DIST = d
								d -= 1
								arc_parent = self.nodes[nodej].arc_parent
								nodej = self.arcs[arc_parent].node_head
			logging.info('arc0 (%s -> arc[%s].next = %s)'%(self.get_arc_name(arc0),self.get_arc_name(arc0),self.get_arc_name(self.arcs[arc0].arc_next)))
			arc0 = self.arcs[arc0].arc_next
		logging.info('set node[%s].parent (%s -> %s)'%(self.get_node_name(nodei),self.get_arc_name(self.nodes[nodei].arc_parent),self.get_arc_name(arc0_min)))
		self.nodes[nodei].arc_parent = arc0_min
		if arc0_min != NULL_PTR:
			logging.info('node[%s].TS(%d -> %d) node[%s].DIST (%d -> %d)'%(self.get_node_name(nodei),self.nodes[nodei].TS,self.TIME,self.get_node_name(nodei),self.nodes[nodei].DIST,d_min + 1))
			self.nodes[nodei].TS = self.TIME
			self.nodes[nodei].DIST = d_min + 1
		else:
			self.add_to_change_list(nodei)
			arc0 = self.nodes[nodei].arc_first
			logging.info('arc0 node[%s].first (%s)'%(self.get_node_name(nodei),self.get_arc_name(arc0)))
			while arc0 != NULL_PTR:
				nodej = self.arcs[arc0].node_head
				logging.info('arc[%s].head (%s) is_sink %s'%(self.get_arc_name(arc0),self.get_node_name(nodej),self.nodes[nodej].is_sink))
				if self.nodes[nodej].is_sink:
					arca = self.nodes[nodej].arc_parent
					if arca != NULL_PTR:
						if self.arcs[arc0].r_cap :
							self.set_active(nodej)
						if arca != MAXFLOW_TERMINAL  and arca != MAXFLOW_ORPHAN and \
							self.arcs[arca].node_head == nodei:
							self.set_orphan_rear(nodej)
				logging.info('arc0 ( %s -> arc[%s].next = %s)'%(self.get_arc_name(arc0),self.get_arc_name(arc0),self.get_arc_name(self.arcs[arc0].arc_next)))
				arc0 = self.arcs[arc0].arc_next
		return
	def debug_parent(self,nodei):
		if len(self.nodes) > nodei:
			logging.info('[%d].parent = %d'%(nodei,self.nodes[nodei].arc_parent))
		return

	def get_arc_nodehead(self,aidx):
		if aidx == MAXFLOW_ORPHAN :
			return 'MAXFLOW_ORPHAN'
		elif aidx == MAXFLOW_TERMINAL:
			return 'MAXFLOW_TERMINAL'
		elif aidx == NULL_PTR:
			return 'NULL'
		else:
			return self.get_node_name(self.arcs[aidx].node_head)

	def process_source_orphan(self,nodei):
		arc0_min = NULL_PTR
		d_min = MAXFLOW_INFINITE_D
		arc0 = self.nodes[nodei].arc_first
		logging.info('arc0 node[%s].first (%s)'%(self.get_node_name(nodei),self.get_arc_name(arc0)))
		while arc0 != NULL_PTR:
			sisidx = self.arcs[arc0].arc_sister
			logging.info('arc[%s] sister[%s].r_cap (%d)'%(self.get_arc_name(arc0),self.get_arc_name(sisidx),self.arcs[sisidx].r_cap))
			if self.arcs[sisidx].r_cap:
				nodej = self.arcs[arc0].node_head
				logging.info('nodej arc[%s].head (%s)'%(self.get_arc_name(arc0),self.get_node_name(nodej)))
				logging.info('node[%s].is_sink %s'%(self.get_node_name(nodej),self.nodes[nodej].is_sink))
				if not self.nodes[nodej].is_sink:
					arca = self.nodes[nodej].arc_parent
					logging.info('arca node[%s].parent (%s)'%(self.get_node_name(nodej),self.get_arc_name(arca)))
					if arca != NULL_PTR:
						d = 0 
						while True:
							logging.info('node[%s].TS (%d) ?== TIME(%d)'%(self.get_node_name(nodej),self.nodes[nodej].TS,self.TIME))
							if self.nodes[nodej].TS == self.TIME:
								d += self.nodes[nodej].DIST
								break
							arca = self.nodes[nodej].arc_parent
							logging.info('node[%s].parent (%s) d (%d -> %d)'%(self.get_node_name(nodej),self.get_arc_name(arca),d,d+1))
							d += 1
							if arca == MAXFLOW_TERMINAL:
								logging.info('node[%s].TS (%d -> %d)'%(self.get_node_name(nodej),self.nodes[nodej].TS,self.TIME))
								logging.info('node[%s].DIST (%d -> %d)'%(self.get_node_name(nodej),self.nodes[nodej].DIST,1))
								self.nodes[nodej].TS = self.TIME
								self.nodes[nodej].DIST = 1
								break
							if arca == MAXFLOW_ORPHAN:
								logging.info('d (%d -> %d)'%(d,MAXFLOW_INFINITE_D))
								d = MAXFLOW_INFINITE_D
								logging.info('orphan %s'%(self.get_arc_name(arca)))
								break
							logging.info('nodej (%s -> arc[%s].head %s)'%(self.get_node_name(nodej),self.get_arc_name(arca),self.get_node_name(self.arcs[arca].node_head)))
							nodej = self.arcs[arca].node_head
						logging.info('d (%d)'%(d))
						if d < MAXFLOW_INFINITE_D:
							if d < d_min:
								logging.info('a0_min (%s -> %s)'%(self.get_arc_name(arc0_min),self.get_arc_name(arc0)))
								arc0_min = arc0
								logging.info('d_min (%d -> %d)'%(d_min,d))
								d_min = d
							logging.info('nodej (%s -> arc[%s].head %s)'%(self.get_node_name(nodej),self.get_arc_name(arc0),self.get_node_name(self.arcs[arc0].node_head)))
							nodej = self.arcs[arc0].node_head
							logging.info('node[%s].TS (%d) ? != TIME (%d)'%(self.get_node_name(nodej),self.nodes[nodej].TS,self.TIME))
							while self.nodes[nodej].TS != self.TIME:
								logging.info('node[%s].TS (%d -> %d) node[%s].DIST (%d -> %d)'%(self.get_node_name(nodej),self.nodes[nodej].TS,self.TIME,self.get_node_name(nodej),self.nodes[nodej].DIST,d))
								self.nodes[nodej].TS =self.TIME
								self.nodes[nodej].DIST = d
								d -= 1
								arc_parent = self.nodes[nodej].arc_parent
								logging.info('nodej (%s -> arc[%s].head (%s)'%(self.get_node_name(nodej),self.get_arc_name(arc_parent),self.arcs[arc_parent].node_head))
								nodej = self.arcs[arc_parent].node_head
			logging.info('arc0 (%s -> arc[%s].next (%s))'%(self.get_arc_name(arc0),self.get_arc_name(arc0),self.get_arc_name(self.arcs[arc0].arc_next)))
			arc0 = self.arcs[arc0].arc_next

		logging.info('node[%s].parent (%s -> a0_min (%s))'%(self.get_node_name(nodei),self.get_arc_name(self.nodes[nodei].arc_parent),self.get_arc_name(arc0_min)))
		self.nodes[nodei].arc_parent = arc0_min
		if arc0_min != NULL_PTR:
			logging.info('node[%s].TS (%d -> %d) node[%s].DIST (%d -> %d)'%(self.get_node_name(nodei),self.nodes[nodei].TS,self.TIME,self.get_node_name(nodei),self.nodes[nodei].DIST,d_min+1))
			self.nodes[nodei].TS = self.TIME
			self.nodes[nodei].DIST = d_min + 1
		else:
			self.add_to_change_list(nodei)
			arc0 = self.nodes[nodei].arc_first
			logging.info('arc0 (node[%s].first %s)'%(self.get_node_name(nodei),self.get_arc_name(arc0)))
			while arc0 != NULL_PTR:
				nodej = self.arcs[arc0].node_head
				logging.info('nodej (arc[%s].head (%s))'%(self.get_arc_name(arc0),self.get_node_name(nodej)))
				logging.info('node[%s].is_sink (%s)'%(self.get_node_name(nodej),self.nodes[nodej].is_sink))
				if  not self.nodes[nodej].is_sink:
					arca = self.nodes[nodej].arc_parent
					if arca != NULL_PTR:
						sisidx = self.arcs[arc0].arc_sister
						logging.info('node[%s].first -> arc[%s]sister -> arc[%s].r_cap %d'%(self.get_node_name(nodei),self.get_arc_name(arc0),self.get_arc_name(sisidx),self.arcs[sisidx].r_cap))
						logging.info('arc[%s].head -> node[%s].parent -> arc[%s].head (%s) ?= nodei (%s)'%(self.get_arc_name(arc0),self.get_node_name(nodej),self.get_arc_name(arca),self.get_arc_nodehead(arca),self.get_node_name(nodei)))
						if self.arcs[sisidx].r_cap:
							self.set_active(nodej)
						if arca != MAXFLOW_TERMINAL and arca != MAXFLOW_ORPHAN and self.arcs[arca].node_head == nodei:
							logging.info('add to rear')
							self.set_orphan_rear(nodej)
				logging.info('arc0 (%s -> arc[%s].next %s)'%(self.get_arc_name(arc0),self.get_arc_name(arc0),self.get_arc_name(self.arcs[arc0].arc_next)))
				arc0 = self.arcs[arc0].arc_next
		return

	def set_orphan_front(self,nodei):
		self.nodes[nodei].arc_parent = MAXFLOW_ORPHAN
		blockptr = NodeBlockPtr()
		blockptr.array_node = nodei
		tmparr = [blockptr]
		tmparr.extend(self.orphan_list)
		self.orphan_list=tmparr
		logging.info('set_orphan_front %s orphan_list %s'%(self.get_node_name(nodei),self.GetOrphanList(self.orphan_list)))
		return

	def set_orphan_rear(self,nodei):
		self.nodes[nodei].arc_parent = MAXFLOW_ORPHAN
		blockptr = NodeBlockPtr()
		blockptr.array_node = nodei
		self.orphan_list.append(blockptr)
		logging.info('set_orphan_rear %s orphan_list %s'%(self.get_node_name(nodei),self.GetOrphanList(self.orphan_list)))
		return

def cpp_command_out(string):
	if CPP_OUT == 1:
		sys.stderr.write(string)
	return

def sortkeys(keys):
	for i in xrange(len(keys)):
		for j in range((i+1),len(keys)):
			if keys[i] > keys[j]:
				tmp = keys[i]
				keys[i] = keys[j]
				keys[j] = tmp
	return keys

def ParseInputFile(infile):
	source=''
	sink=0
	widht=0
	height=0
	bkgraph=None
	sourc_sink_pair={}
	with open(infile) as f:
		for l in f:
			if l.startswith('#'):
				continue
			l = l.rstrip('\r\n')
			if l.startswith('source='):
				sarr = l.split('=')
				if len(sarr) < 2:
					continue
				source = int(sarr[1])
				continue
			elif l.startswith('sink='):
				sarr = l.split('=')
				if len(sarr) < 2:
					continue
				sink = int(sarr[1])
				continue
			elif l.startswith('width='):
				sarr = l.split('=')
				if len(sarr) < 2:
					continue
				w = int(sarr[1])
				continue
			elif l.startswith('height='):
				sarr = l.split('=')
				if len(sarr) < 2:
					continue
				h = int(sarr[1])
				continue
			sarr = l.split(',')
			if len(sarr) < 3:
				continue
			if sink == 0 :
				sys.stderr.write('can not define sink or not define width or height\n')
				sys.exit(4)
			if bkgraph is None:
				bkgraph = BKGraph(sink+1,sink * (sink-1)/2)
				bkgraph.add_node(sink+1)
			curs = int(sarr[0])
			curt = int(sarr[1])
			curw = int(sarr[2])

			if curs == source and curt == sink:
				logging.error('s[%d] == source d[%d] == sink'%(curs,curt))
				sys.exit(4)

			if curs == source and curt != sink :
				if curt not in sourc_sink_pair.keys():
					sourc_sink_pair[curt] = [curw,0]
					#bkgraph.add_tweights(curt,curw,0)
				else:
					# we have put the sink into the graph
					sourc_sink_pair[curt][0]=curw
					logging.info('add t-link[%d] source(%d) sink(%d)'%(curt,sourc_sink_pair[curt][0],sourc_sink_pair[curt][1]))
					cpp_command_out('g -> add_tweights(%d,%d,%d);\n'%(curt,sourc_sink_pair[curt][0],sourc_sink_pair[curt][1]))
					bkgraph.add_tweights(curt,sourc_sink_pair[curt][0],sourc_sink_pair[curt][1])
					del sourc_sink_pair[curt]
				continue

			if curt == sink and curs != source:
				if curs not in sourc_sink_pair.keys():
					sourc_sink_pair[curs]=[0,curw]
					#bkgraph.add_tweights(curs,0,curw)
				else:
					sourc_sink_pair[curs][1]=curw
					logging.info('add t-link[%d] source(%d) sink(%d)'%(curs,sourc_sink_pair[curs][0],sourc_sink_pair[curs][1]))
					cpp_command_out('g -> add_tweights(%d,%d,%d);\n'%(curs,sourc_sink_pair[curs][0],sourc_sink_pair[curs][1]))
					bkgraph.add_tweights(curs,sourc_sink_pair[curs][0],sourc_sink_pair[curs][1])
					del sourc_sink_pair[curs]
				continue

			logging.info('set n-link [%d]->[%d] %d'%(curs,curt,curw))
			cpp_command_out('g -> add_edge(%d,%d,%d,0);\n'%(curs,curt,curw))
			bkgraph.add_edge(curs,curt,curw,0)

	for k in sortkeys(sourc_sink_pair.keys()):
		# now to add t-link weights
		logging.info('add t-link[%d] source(%d) sink(%d)'%(k,sourc_sink_pair[k][0],sourc_sink_pair[k][1]))
		cpp_command_out('g -> add_tweights(%d,%d,%d);\n'%(k,sourc_sink_pair[k][0],sourc_sink_pair[k][1]))
		bkgraph.add_tweights(k,sourc_sink_pair[k][0],sourc_sink_pair[k][1])
	return bkgraph


def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile\n'%(sys.argv[0]))
		sys.exit(4)
	bkgraph = ParseInputFile(sys.argv[1])
	flow = bkgraph.max_flow()
	sys.stdout.write('%d\n'%(flow))
	return

if __name__ == '__main__':
	#logging.basicConfig(level=logging.INFO,format='%(asctime)-15s:%(filename)s:%(lineno)d\t%(message)s')
	#logging.basicConfig(level=logging.INFO,format='%(filename)s:%(funcName)s:%(lineno)d\t%(message)s')
	logging.basicConfig(level=logging.DEBUG,format='%(filename)s:%(funcName)s:%(lineno)d %(message)s')
	#logging.basicConfig(level=logging.DEBUG,format='%(message)s')
	#logging.basicConfig(level=logging.INFO,format='%(message)s')
	#logging.basicConfig(level=logging.ERROR,format='%(asctime)-15s:%(filename)s:%(lineno)d\t%(message)s')
	main()