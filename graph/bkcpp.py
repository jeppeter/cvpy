import sys
import logging

NULL_PTR=-1
MAXFLOW_TERMINAL=-2
MAXFLOW_ORPHAN=-3
MAXFLOW_INFINITE_D=(0xffffffff >> 1)
CPP_OUT=0


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
		self.is_marked = False
		self.is_in_changed_list = False
		self.tr_cap = 0
		return

class NodeBlockPtr:
	def __init__(self):
		self.array_node = NULL_PTR
		return


def GetOrphanList(orphlist):
	i = 0
	s = 'cnt(%d)['%(len(orphlist))
	for a in orphlist:
		if i != 0:
			s += ','
		s += '%s'%(GetIdx(a.array_node))
		i += 1
	s += ']'
	return s

class BKGraph:
	def __init__(self,nodemax,edgemax):
		self.nodes = []
		self.arcs = []
		self.flow = 0
		self.maxflow_iteration = 0
		self.orphan_list = []
		self.queue_first = [NULL_PTR,NULL_PTR]
		self.queue_last = [NULL_PTR,NULL_PTR]
		return

	def add_node(self,num=1):
		for i in xrange(num):
			self.nodes.append(Node())
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
		logging.info('node[%s].tr_cap = %d flow %d'%(GetIdx(nodeid),self.nodes[nodeid].tr_cap,self.flow))
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
		nodej = self.nodes[nodeidj]

		# now set for the idx
		aarc.arc_sister = arevidx
		arevarc.arc_sister = aidx
		aarc.arc_next = nodei.arc_first
		nodei.arc_first = aidx
		logging.info('arc[%s].next = %s nodei[%s].first = %s'%(GetIdx(aidx),GetIdx(aarc.arc_next),GetIdx(nodeidi),GetIdx(nodei.arc_first)))
		arevarc.arc_next = nodej.arc_first
		nodej.arc_first = arevidx
		logging.info('arevarc[%s].next = %s nodej[%s].first = %s'%(GetIdx(arevidx),GetIdx(arevarc.arc_next),GetIdx(nodeidj),GetIdx(nodej.arc_first)))
		aarc.node_head = nodeidj
		arevarc.node_head = nodeidi
		aarc.r_cap = cap
		arevarc.r_cap = rev_cap
		logging.info('arc[%s].r_cap = %d arc[%s].r_cap = %d'%(GetIdx(aidx),cap,GetIdx(arevidx),rev_cap))

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
				logging.info('nodei %s'%(GetIdx(nodei)))
				self.nodes[nodei].arc_next = NULL_PTR
				if self.nodes[nodei].arc_parent == NULL_PTR:
					nodei = NULL_PTR
			if nodei == NULL_PTR:
				nodei = self.next_active()
				if nodei == NULL_PTR:
					break
			logging.info('nodei %s'%(GetIdx(nodei)))

			if not self.nodes[nodei].is_sink:
				aidx = self.nodes[nodei].arc_first
				logging.info('[%s].first = %s'%(GetIdx(nodei),GetIdx(aidx)))
				while aidx != NULL_PTR:
					logging.info('[%s].r_cap = %d'%(GetIdx(aidx),self.arcs[aidx].r_cap))
					if self.arcs[aidx].r_cap:
						nodej = self.arcs[aidx].node_head
						logging.info('[%s].head -> [%s].parent %s'%(GetIdx(aidx),GetIdx(nodej),GetIdx(self.nodes[nodej].arc_parent)))
						if self.nodes[nodej].arc_parent == NULL_PTR:
							logging.info('[%d].parent = %d'%(nodej,self.arcs[aidx].arc_sister))
							self.nodes[nodej].is_sink = False
							self.nodes[nodej].arc_parent = self.arcs[aidx].arc_sister
							self.nodes[nodej].TS = self.nodes[nodei].TS
							self.nodes[nodej].DIST = self.nodes[nodei].DIST + 1
							self.set_active(nodej)
							self.add_to_change_list(nodej)
						elif self.nodes[nodej].is_sink:
							logging.info('[%s].is_sink'%(GetIdx(nodej)))
							break
						elif self.nodes[nodej].TS <= self.nodes[nodei].TS and \
							self.nodes[nodej].DIST > self.nodes[nodei].DIST :
							logging.info('[%d].parent = %d'%(nodej,self.arcs[aidx].arc_sister))
							self.nodes[nodej].arc_parent = self.arcs[aidx].arc_sister
							self.nodes[nodej].TS =self.nodes[nodei].TS
							self.nodes[nodej].DIST = self.nodes[nodei].DIST + 1
					logging.info('[%s].next = %s'%(GetIdx(aidx),GetIdx(self.arcs[aidx].arc_next)))
					aidx = self.arcs[aidx].arc_next

			else:
				aidx = self.nodes[nodei].arc_first
				while aidx != NULL_PTR:					
					sisidx = self.arcs[aidx].arc_sister
					logging.info('[%d].r_cap = %d'%(sisidx,self.arcs[sisidx].r_cap))
					if self.arcs[sisidx].r_cap:
						nodej = self.arcs[aidx].node_head
						if self.nodes[nodej].arc_parent == NULL_PTR:
							logging.info('[%d].parent = %d'%(nodej,self.arcs[aidx].arc_sister))
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
							logging.info('[%d].parent = %d'%(nodej,self.arcs[aidx].arc_sister))
							self.nodes[nodej].arc_parent = self.arcs[aidx].arc_sister
							self.nodes[nodej].TS = self.nodes[nodei].TS
							self.nodes[nodej].DIST = self.nodes[nodei].DIST + 1

					aidx = self.arcs[aidx].arc_next

			self.TIME += 1
			logging.info('TIME %d arc[%s]'%(self.TIME,GetIdx(aidx)))

			if aidx != NULL_PTR:
				logging.info('[%s].next %s -> %s'%(GetIdx(nodei),GetIdx(self.nodes[nodei].node_next),GetIdx(nodei)))
				self.nodes[nodei].node_next = nodei
				curnodeid = nodei
				logging.info('curnodeid %s'%(GetIdx(curnodeid)))

				self.augment(aidx)

				while len(self.orphan_list) > 0:
					curorphan = self.orphan_list[0]
					self.orphan_list = self.orphan_list[1:]
					curorphnodei = curorphan.array_node
					if self.nodes[curorphnodei].is_sink:
						logging.info('sink orphan %d'%(curorphnodei))
						self.process_sink_orphan(curorphnodei)
					else:
						logging.info('source orphan %d'%(curorphnodei))
						self.process_source_orphan(curorphnodei)
						logging.info('source orphan over %d'%(curorphnodei))
			else:
				curnodeid = NULL_PTR
		self.maxflow_iteration += 1
		return self.flow

	def maxflow_init(self):
		self.queue_first = [NULL_PTR,NULL_PTR]
		self.queue_last = [NULL_PTR,NULL_PTR]
		self.orphan_list = []
		self.TIME = 0

		for nodei in xrange(len(self.nodes)):
			self.nodes[nodei].node_next = NULL_PTR
			self.nodes[nodei].is_marked = 0
			self.nodes[nodei].is_in_changed_list = 0
			self.nodes[nodei].TS = self.TIME

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
		return

	def next_active(self):
		while True:
			nodei = self.queue_first[0]
			logging.info('queue_first[0] = %s'%(GetIdx(self.queue_first[0])))
			if nodei == NULL_PTR:
				nodei = self.queue_first[1]
				self.queue_first[0] = nodei
				self.queue_last[0] = self.queue_last[1]
				self.queue_first[1] = NULL_PTR
				self.queue_last[1] = NULL_PTR
				if nodei == NULL_PTR:
					return NULL_PTR
			logging.info('[%s].next = %s'%(GetIdx(nodei),GetIdx(self.nodes[nodei].node_next)))
			if self.nodes[nodei].node_next == nodei:
				self.queue_first[0] = NULL_PTR
				self.queue_last[0] = NULL_PTR
			else:
				self.queue_first[0]= self.nodes[nodei].node_next
			self.nodes[nodei].node_next = NULL_PTR

			logging.info('[%s].parent = %s'%(GetIdx(nodei),GetIdx(self.nodes[nodei].arc_parent)))
			if self.nodes[nodei].arc_parent != NULL_PTR:
				return nodei
		return NULL_PTR

	def set_active(self,nodei):
		logging.info('set node[%s].next = %s active'%(GetIdx(nodei),GetIdx(self.nodes[nodei].node_next)))
		logging.info('queue_last[1] = %s queue_first[1] = %s'%(GetIdx(self.queue_last[1]),GetIdx(self.queue_first[1])))
		if self.nodes[nodei].node_next == NULL_PTR:
			if self.queue_last[1] != NULL_PTR:
				self.nodes[self.queue_last[1]].node_next = nodei
				logging.info('set[%s].next = %s'%(GetIdx(self.queue_last[1]),GetIdx(nodei)))
			else:
				self.queue_first[1] = nodei
				logging.info('set queue_first[1] = %s'%(GetIdx(nodei)))
			self.queue_last[1] = nodei
			self.nodes[nodei].node_next = nodei
			logging.info('[%s].next = %s'%(GetIdx(nodei),GetIdx(nodei)))
		return

	def add_to_change_list(self,nodei):
		# nothing to add change list because we do not set change list before
		return

	def augment(self,aidx):		
		bottlecap = self.arcs[aidx].r_cap
		logging.info('[%s].r_cap bottlecap %d'%(GetIdx(aidx),bottlecap))
		# this is source tree
		sisidx = self.arcs[aidx].arc_sister
		nodei = self.arcs[sisidx].node_head
		logging.info('[%s].sister %s nodei %s'%(GetIdx(aidx),GetIdx(sisidx),GetIdx(nodei)))
		while True:
			arca = self.nodes[nodei].arc_parent
			logging.info('[%s].parent %s'%(GetIdx(nodei),GetIdx(arca)))
			if arca == MAXFLOW_TERMINAL:
				break
			sisidx = self.arcs[arca].arc_sister
			logging.info('[%s].sister [%s].r_cap %d bottlecap(%d)'%(GetIdx(arca),GetIdx(sisidx),self.arcs[sisidx].r_cap,bottlecap))
			if bottlecap > self.arcs[sisidx].r_cap:
				bottlecap = self.arcs[sisidx].r_cap
			nodei = self.arcs[arca].node_head
			logging.info('[%s].head = %s'%(GetIdx(arca),GetIdx(nodei)))
		logging.info('[%s].tr_cap = %d bottlecap(%d)'%(GetIdx(nodei),self.nodes[nodei].tr_cap,bottlecap))
		if bottlecap > self.nodes[nodei].tr_cap:
			bottlecap = self.nodes[nodei].tr_cap

		# this is sink tree
		nodei = self.arcs[aidx].node_head
		logging.info('[%s].head = %s'%(GetIdx(aidx),GetIdx(nodei)))
		while True:
			arca = self.nodes[nodei].arc_parent
			logging.info('[%s].parent = %s'%(GetIdx(nodei),GetIdx(arca)))
			if arca == MAXFLOW_TERMINAL:
				break
			logging.info('[%s].r_cap = %d bottlecap(%d)'%(GetIdx(arca),self.arcs[arca].r_cap,bottlecap))
			if bottlecap > self.arcs[arca].r_cap :
				bottlecap = self.arcs[arca].r_cap
			nodei = self.arcs[arca].node_head
			logging.info('[%s].head = %s'%(GetIdx(arca),GetIdx(nodei)))
		logging.info('[%s].tr_cap = %d bottlecap(%d)'%(GetIdx(nodei),self.nodes[nodei].tr_cap,bottlecap))
		if bottlecap > - self.nodes[nodei].tr_cap :
			bottlecap = - self.nodes[nodei].tr_cap

		sisidx = self.arcs[aidx].arc_sister
		logging.info('[%s].sister -> [%s].r_cap(%d+%d) [%s].r_cap(%d-%d)'%(GetIdx(aidx),GetIdx(sisidx),self.arcs[sisidx].r_cap,bottlecap,GetIdx(aidx),self.arcs[aidx].r_cap,bottlecap))
		self.arcs[sisidx].r_cap += bottlecap
		self.arcs[aidx].r_cap -= bottlecap

		nodei = self.arcs[sisidx].node_head
		while True:
			arca = self.nodes[nodei].arc_parent
			if arca == MAXFLOW_TERMINAL:
				break
			sisidx = self.arcs[arca].arc_sister
			logging.info('[%s].r_cap (%d+%d) [%s].sister -> [%s].r_cap(%d-%d)'%(GetIdx(arca),self.arcs[arca].r_cap,bottlecap,GetIdx(arca),GetIdx(sisidx),self.arcs[sisidx].r_cap,bottlecap))
			self.arcs[arca].r_cap += bottlecap
			self.arcs[sisidx].r_cap -= bottlecap
			if self.arcs[sisidx].r_cap == 0 :
				logging.info('[%s] set_orphan_front'%(GetIdx(sisidx)))
				self.set_orphan_front(nodei)			
			nodei = self.arcs[arca].node_head
			logging.info('[%s].head = %s'%(GetIdx(arca),GetIdx(nodei)))

		logging.info('[%s].tr_cap (%d-%d)'%(GetIdx(nodei),self.nodes[nodei].tr_cap,bottlecap))
		self.nodes[nodei].tr_cap -= bottlecap

		if self.nodes[nodei].tr_cap == 0:
			logging.info('[%s] set_orphan_front'%(GetIdx(nodei)))
			self.set_orphan_front(nodei)

		nodei = self.arcs[aidx].node_head
		logging.info('[%s].head = %s'%(GetIdx(aidx),GetIdx(nodei)))
		while True:
			arca = self.nodes[nodei].arc_parent
			logging.info('[%s].parent = %s'%(GetIdx(nodei),GetIdx(arca)))
			if arca == MAXFLOW_TERMINAL:
				break
			sisidx = self.arcs[arca].arc_sister
			logging.info('[%s].r_cap (%d+%d) [%s].sister -> [%s].r_cap(%d-%d)'%(GetIdx(arca),self.arcs[arca].r_cap,bottlecap,GetIdx(arca),GetIdx(sisidx),self.arcs[sisidx].r_cap,bottlecap))
			self.arcs[sisidx].r_cap += bottlecap
			self.arcs[arca].r_cap -= bottlecap
			if self.arcs[arca].r_cap == 0 :
				logging.info('[%s] set_orphan_front'%(GetIdx(sisidx)))
				self.set_orphan_front(nodei)
			nodei = self.arcs[arca].node_head
			logging.info('[%s].head = %s'%(GetIdx(arca),GetIdx(nodei)))
		logging.info('[%s].tr_cap (%d+%d)'%(GetIdx(nodei),self.nodes[nodei].tr_cap,bottlecap))
		self.nodes[nodei].tr_cap += bottlecap
		if self.nodes[nodei].tr_cap == 0:
			logging.info('[%s] set_orphan_front'%(GetIdx(nodei)))
			self.set_orphan_front(nodei)
		logging.info('flow (%d+%d)'%(self.flow,bottlecap))
		self.flow += bottlecap
		return

	def process_sink_orphan(self,nodei):
		d_min = MAXFLOW_INFINITE_D
		arc0 = self.nodes[nodei].arc_first
		arc0_min = NULL_PTR
		while arc0 != NULL_PTR:
			if self.arcs[arc0].r_cap != 0 :
				nodej = self.arcs[arc0].node_head
				if self.nodes[nodej].is_sink:
					arca = self.nodes[nodej].arc_parent
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
			arc0 = self.arcs[arc0].arc_next

		self.nodes[nodei].arc_parent = arc0_min
		if arc0_min != NULL_PTR:
			self.nodes[nodei].TS = self.TIME
			self.nodes[nodei].DIST = d_min + 1
		else:
			self.add_to_change_list(nodei)
			arc0 = self.nodes[nodei].arc_first
			while arc0 != NULL_PTR:
				nodej = self.arcs[arc0].node_head
				if self.nodes[nodej].is_sink:
					arca = self.nodes[nodej].arc_parent
					if arca != NULL_PTR:
						if self.arcs[arc0].r_cap :
							self.set_active(nodej)
						if arca != MAXFLOW_TERMINAL  and arca != MAXFLOW_ORPHAN and \
							self.arcs[arca].node_head == nodei:
							self.set_orphan_rear(nodej)
				arc0 = self.arcs[arc0].arc_next
		return
	def debug_parent(self,nodei):
		if len(self.nodes) > nodei:
			logging.info('[%d].parent = %d'%(nodei,self.nodes[nodei].arc_parent))
		return

	def process_source_orphan(self,nodei):
		arc0_min = NULL_PTR
		d_min = MAXFLOW_INFINITE_D
		arc0 = self.nodes[nodei].arc_first
		while arc0 != NULL_PTR:
			sisidx = self.arcs[arc0].arc_sister
			logging.info('[%d] sister[%d].r_cap %d'%(arc0,sisidx,self.arcs[sisidx].r_cap))
			if self.arcs[sisidx].r_cap:
				nodej = self.arcs[arc0].node_head
				if not self.nodes[nodej].is_sink:
					arca = self.nodes[nodej].arc_parent
					if arca != NULL_PTR:
						d = 0 
						while True:
							if self.nodes[nodej].TS == self.TIME:
								d += self.nodes[nodej].DIST
								break
							arca = self.nodes[nodej].arc_parent
							logging.info('[%d].parent %d'%(nodej,arca))
							d += 1
							if arca == MAXFLOW_TERMINAL:
								self.nodes[nodej].TS = self.TIME
								self.nodes[nodej].DIST = 1
								break
							if arca == MAXFLOW_ORPHAN:								
								d = MAXFLOW_INFINITE_D
								logging.info('orphan %d'%(arca))
								break
							nodej = self.arcs[arca].node_head
						logging.info('d %d'%(d))
						if d < MAXFLOW_INFINITE_D:
							if d < d_min:
								logging.info('a0_min %d'%(arc0))
								arc0_min = arc0
								d_min = d
							nodej = self.arcs[arc0].node_head
							while self.nodes[nodej].TS != self.TIME:
								self.nodes[nodej].TS =self.TIME
								self.nodes[nodej].DIST = d
								d -= 1
								arc_parent = self.nodes[nodej].arc_parent
								nodej = self.arcs[arc_parent].node_head
			arc0 = self.arcs[arc0].arc_next			

		logging.info('a0_min %d'%(arc0_min))
		self.nodes[nodei].arc_parent = arc0_min
		if arc0_min != NULL_PTR:
			self.nodes[nodei].TS = self.TIME
			self.nodes[nodei].DIST = d_min + 1
		else:
			self.add_to_change_list(nodei)
			arc0 = self.nodes[nodei].arc_first
			while arc0 != NULL_PTR:
				nodej = self.arcs[arc0].node_head
				logging.info('nodej %d'%(nodej))
				if  not self.nodes[nodej].is_sink:
					arca = self.nodes[nodej].arc_parent
					if arca != NULL_PTR:
						sisidx = self.arcs[arc0].arc_sister
						logging.info('[%d][%d] sister[%d].r_cap %d'%(nodej,arc0,sisidx,self.arcs[sisidx].r_cap))
						if self.arcs[sisidx].r_cap:
							self.set_active(nodej)
						if arca != MAXFLOW_TERMINAL and arca != MAXFLOW_ORPHAN and self.arcs[arca].node_head == nodei:
							logging.info('add to rear')
							self.set_orphan_rear(nodej)
				arc0 = self.arcs[arc0].arc_next
		return

	def set_orphan_front(self,nodei):
		self.nodes[nodei].arc_parent = MAXFLOW_ORPHAN
		blockptr = NodeBlockPtr()
		blockptr.array_node = nodei
		tmparr = [blockptr]
		tmparr.extend(self.orphan_list)
		self.orphan_list=tmparr
		logging.info('set_orphan_front %s orphan_list %s'%(GetIdx(nodei),GetOrphanList(self.orphan_list)))
		return

	def set_orphan_rear(self,nodei):
		self.nodes[nodei].arc_parent = MAXFLOW_ORPHAN
		blockptr = NodeBlockPtr()
		blockptr.array_node = nodei
		self.orphan_list.append(blockptr)
		logging.info('set_orphan_rear %d'%(nodei))
		return

def cpp_command_out(string):
	if CPP_OUT == 1:
		sys.stderr.write(string)
	return

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
	

	for k in sourc_sink_pair.keys():
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
	logging.basicConfig(level=logging.INFO,format='%(filename)s:%(lineno)d\t%(message)s')
	#logging.basicConfig(level=logging.INFO,format='%(message)s')
	#logging.basicConfig(level=logging.ERROR,format='%(asctime)-15s:%(filename)s:%(lineno)d\t%(message)s')
	main()