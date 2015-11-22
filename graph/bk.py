import sys
import logging


#########################################
# this code is from the code in the website
# https://en.wikibooks.org/wiki/Algorithm_Implementation/Graphs/Maximum_flow/Boykov_%26_Kolmogorov
# and it will be for the boykov kolmogorov 
#########################################
eps = 0.001

#######################################
# node member 
#           index
#           prevedgeindex
#           mapcaptohere
#           dist
#######################################
class Node:
	def __init__(self):
		self.index = 0
		self.prevedgeindex = -1
		return

	def SetIdx(self,ind):
		self.index = ind


#######################################
# edge member
#           initial_vertex
#           terminal_vertex
#           capacity
#           flow
#           invedgeindex
#######################################
class Edge:
	def __init__(self):
		self.initial_vertex = 0
		self.terminal_vertex = 0
		self.capacity = 0
		self.flow = 0
		self.invedgeindex = 0
		return


class EdgeMat:
	def __init__(self):
		pass

class LinkedListInt:
	def __init__(self):
		self.__listarr = []
		self.__listnum = 0
		return

	def add(self,i):
		self.__listarr.append(i)
		self.__listnum += 1
		return

	def is_empty(self):
		if self.__listnum == 0 :
			return True
		return False
	def peek(self):
		if self.__listnum == 0:
			return None
		return self.__listarr[0]

	def poll(self):
		if self.__listnum == 0:
			return None
		iret = self.__listarr[0]
		self.__listarr = self.__listarr[1:]
		self.__listnum -= 1
		return iret

	def add_first(self,idx):
		tmparr = [idx]
		tmparr.extend(self.__listarr)
		self.__listarr = tmparr
		self.__listnum += 1
		return

	def add_last(self,idx):
		self.__listarr.append(idx)
		self.__listnum += 1
		return

class BoolArray:
	def __init__(self,num):
		self.__arr = num *[0]
		self.__cnt = num
		return 

	def SetTrue(self,idx):
		assert(len(self.__arr) > idx)
		assert(self.__cnt > idx)
		self.__arr[idx] = 1
		return

	def SetFalse(self,idx):
		assert(len(self.__arr) > idx)
		assert(self.__cnt > idx)
		self.__arr[idx] = 0
		return


	def Get(self,idx):
		assert(len(self.__arr) > idx)
		assert(self.__cnt > idx)		
		return self.__arr[idx]



class StartingEdge:
	def __init__(self):
		self.__startingedge=[]
		self.__startingcnt=0
		self.__startcap = 0
		return

	def SetIntArray(self,idx,cnt):
		if self.__startcap <= idx :
			for x in xrange((idx-self.__startcap+1)):
				self.__startingedge.append([])
				self.__startcap += 1
		tmparr = []
		for i in xrange(cnt):
			tmparr.append(0)
		self.__startingedge[idx] = tmparr
		return
	def SetArrayNumber(self,k1,k2,num):
		assert(len(self.__startingedge) == self.__startcap)
		assert(len(self.__startingedge) > k1)
		assert(len(self.__startingedge[k1]) > k2)
		self.__startingedge[k1][k2] = num
		return

	def GetLengthFromIdx(self,idx):
		assert(len(self.__startingedge) > idx)
		assert(len(self.__startingedge) == self.__startcap)
		return len(self.__startingedge[idx])

	def GetArrayNumber(self,k1,k2):
		assert(len(self.__startingedge) > k1)
		assert(len(self.__startingedge[k1]) > k2)
		return self.__startingedge[k1][k2]


########################################
# GraphCutBoykovKolmogorov member
#           debug = true
#           eps = 0.001
#           nbNode
#           nbEdges
#           w,h
#           node[]
#           edge[]
#           startingedge[][]
#           curedge
#           
########################################
class GraphCutBoykovKolmogorov:
	def __indice_part(self,x,y):
		return x*self.h + y

	def __init__(self,width,height):
		self.w = width
		self.h = height
		self.orphan = LinkedListInt()
		self.active = LinkedListInt()
		self.voisinsEdgeACreer=[[1,0],[1,-1],[0,-1],[-1,-1]]
		self.nbNode = self.w * self.h + 2
		self.nbEdges = (self.w * self.h * 4) + (self.h - 2)* (self.w - 2)*8 +\
			2*5*(self.h+self.w - 4)+4*3
		self.node = self.nbNode * [Node()]
		self.edges = self.nbEdges * [Edge()]
		self.curNbVoisins = self.nbNode * [0]
		self.node[0].SetIdx(0)
		self.node[1].SetIdx(1)
		self.isInS = BoolArray(self.nbNode)
		self.isInA = BoolArray(self.nbNode)
		for x in xrange(self.w):
			for y in xrange(self.h):
				i1 = x * self.h + y + 2
				self.node[i1].SetIdx(i1)
		self.curedge = 0
		self.startingedge=StartingEdge()
		self.startingedge.SetIntArray(0,(self.nbNode-2))
		self.startingedge.SetIntArray(1,(self.nbNode-2))

		if self.w == 1 or self.h == 1:
			for x in xrange(self.w):
				for y in xrange(self.h):
					if (x * (x + 1 - self.w)) == 0 and (y * (y+1-self.h)) == 0:
						# for coins
						self.startingedge.SetIntArray((x*self.h+y+2),1+1)
					elif (x*(x+1-self.w)) == 0 or (y * (y + 1-self.h)) == 0 :
						# edges
						self.startingedge.SetIntArray((x*self.h+y+2),2+1)
		else:
			for x in xrange(self.w):
				for y in xrange(self.h):
					if (x*(x+1-self.w)) == 0 and (y*(y+1-self.h)) == 0:
						# coins
						self.startingedge.SetIntArray((x*self.h+y+2),3+2)
					elif (x*(x+1-self.w)) == 0 or (y*(y+1-self.h))==0:
						# edges
						self.startingedge.SetIntArray((x*self.h+y+2),5+2)
					else:
						# in middle
						self.startingedge.SetIntArray((x*self.h+y+2),8+2)
		# for t-links as sink edges
		for x in xrange(self.w):
			for y in xrange(self.h):
				i1 = x*self.h + y + 2
				i2 = 1
				self.edges[self.curedge].initial_vertex = i1
				self.edges[self.curedge].terminal_vertex = i2
				self.edges[self.curedge].invedgeindex = self.curedge + 1
				self.startingedge.SetArrayNumber(i1,self.curNbVoisins[i1],self.curedge)
				self.curNbVoisins[i1] += 1

				self.curedge += 1

				self.edges[self.curedge].initial_vertex = i2
				self.edges[self.curedge].terminal_vertex = i1
				self.edges[self.curedge].invedgeindex = self.curedge - 1

				self.startingedge.SetArrayNumber(i2,(x*self.h+y),self.curedge)
				self.curedge += 1

		# for n-links
		for x in xrange(self.w):
			for y in xrange(self.h):
				i1 = x*self.h + y +2
				for v in xrange(len(self.voisinsEdgeACreer)):
					vx = x + self.voisinsEdgeACreer[v][0]
					vy = y + self.voisinsEdgeACreer[v][1]
					if vx < 0 or vx >= self.w or vy < 0 or vy >= self.h:
						continue
					i2 = vx*self.h + vy + 2
					self.edges[self.curedge].initial_vertex = i1
					self.edges[self.curedge].terminal_vertex = i2
					self.edges[self.curedge].invedgeindex = self.curedge + 1
					self.startingedge.SetArrayNumber(i1,self.curNbVoisins[i1],self.curedge)
					self.curNbVoisins[i1] += 1

					self.curedge += 1

					self.edges[self.curedge].initial_vertex = i2
					self.edges[self.curedge].terminal_vertex = i1
					self.edges[self.curedge].invedgeindex = self.curedge - 1
					self.startingedge.SetArrayNumber(i2,self.curNbVoisins[i2],self.curedge)
					self.curNbVoisins[i2] += 1
					self.curedge += 1

		# for t-links source edge
		for x in xrange(self.w):
			for y in xrange(self.h):
				i1 = 0 
				i2 = x*self.h + y + 2
				self.edges[self.curedge].initial_vertex = i1
				self.edges[self.curedge].terminal_vertex = i2
				self.edges[self.curedge].invedgeindex = self.curedge + 1
				self.startingedge.SetArrayNumber(i1,(x*self.h+y),self.curedge)
				self.curedge += 1

				self.edges[self.curedge].initial_vertex = i2
				self.edges[self.curedge].terminal_vertex = i1
				self.edges[self.curedge].invedgeindex = self.curedge - 1
				self.startingedge.SetArrayNumber(i2,self.curNbVoisins[i2],self.curedge)
				self.curNbVoisins[i2] += 1
				self.curedge += 1

		assert(self.curedge == self.nbEdges)
		return

	def get_edge(self,x1,y1,x2,y2):
		i1 = x1*self.h + y1 + 2
		i2 = x2*self.h + y2 + 2
		for v in xrange(self.startingedge.GetLengthFromIdx(i1)):
			if self.edges[self.startingedge.GetArrayNumber(i1,v)].terminal_vertex == i2:
				return self.edges[self.startingedge.GetArrayNumber(i1,v)]
		return None

	def set_intern_weight(self,x1,y1,x2,y2,w):
		i1 = x1*self.h + y1 + 2
		i2 = x2*self.h + y2 + 2
		for v in xrange(self.startingedge.GetLengthFromIdx(i1)):
			if self.edges[self.startingedge.GetArrayNumber(i1,v)].terminal_vertex == i2:
				self.edges[self.startingedge.GetArrayNumber(i1,v)].capacity = w
				self.edges[self.edges[self.startingedge.GetArrayNumber(i1,v)].invedgeindex].capacity = w
		return

	def set_source_weight(self,x1,y1,w):
		i1 = 0
		i2 = x1 * self.h + y1
		self.edges[self.startingedge.GetArrayNumber(i1,i2)].capacity = w
		self.edges[self.edges[self.startingedge.GetArrayNumber(i1,i2)].invedgeindex].capacity = w
		return

	def set_sink_weight(self,x1,y1,w):
		i1 = 1
		i2 = x1 *self.h + y1
		self.edges[self.startingedge.GetArrayNumber(i1,i2)].capacity = w
		self.edges[self.edges[self.startingedge.GetArrayNumber(i1,i2)].invedgeindex].capacity = w
		return

	def set_t_weight(self,x1,y1,wsource,wsink):
		i2 = x1*self.h + y1
		self.edges[self.startingedge.GetArrayNumber(0,i2)].capacity = wsource
		self.edges[self.edges[self.startingedge.GetArrayNumber(0,i2)].invedgeindex].capacity = wsource
		self.edges[self.startingedge.GetArrayNumber(1,i2)].capacity = -wsink
		self.edges[self.edges[self.startingedge.GetArrayNumber(1,i2)].invedgeindex].capacity = wsink
		return

	def linked_to(self,x1,y1):
		i2 = x1 * self.h + y1
		if self.edges[self.startingedge.GetArrayNumber(0,i2)].capacity - self.edges[self.startingedge.GetArrayNumber(0,i2)].flow > eps:
			return 0
		elif self.edges[self.edges[self.startingedge.GetArrayNumber(1,i2)].invedgeindex].capacity - self.edges[self.edges[self.startingedge.GetArrayNumber(1,i2)].invedgeindex].flow > esp:
			return 1
		else:
			return 2

	def get_flow(self):
		f= 0
		for x in xrange(self.w):
			for y in xrange(self.h):
				f += self.edges[self.edges[self.startingedge.GetArrayNumber(1,(x*self.h+y))].invedgeindex].flow
		return f

	def reset_flow(self):
		for i in xrange(len(self.edges)):
			self.edges[i].flow = 0
		return

	def do_cut(self):
		self.reset_flow()
		self.isInS = BoolArray(self.nbNode)
		self.isInS.SetTrue(0)
		self.active = LinkedListInt()
		self.active.add(0)
		self.isInA = BoolArray(self.nbNode)
		self.isInA.SetTrue(0)
		self.orphan = LinkedListInt()


		for i in xrange(len(self.node)):
			self.node[i].prevedgeindex = -1

		while True:
			lastedge = self.growth_stage()
			if lastedge == -1:
				return
			self.augmentation_state(lastedge)
			self.adoption_state()
		return
	
	def growth_stage(self):
		if self.isInS.Get(1):
			return self.node[1].prevedgeindex
		while not self.active.is_empty():
			curnodeindex = self.active.peek()
			if self.isInA.Get(curnodeindex) == 0 :
				self.active.poll()
				continue
			for curstartedgeidx in xrange(self.startingedge.GetLengthFromIdx(curnodeindex)):
				curedge = self.edges[self.startingedge.GetArrayNumber(curnodeindex,curstartedgeidx)]
				if (curedge.capacity - curedge.flow) < eps:
					continue
				if self.isInS.Get(curedge.terminal_vertex) == 0 :
					self.active.add(curedge.terminal_vertex)
					self.isInA.SetTrue(curedge.terminal_vertex)
					self.isInS.SetTrue(curedge.terminal_vertex)
					self.node[curedge.terminal_vertex].prevedgeindex = self.startingedge.GetArrayNumber(curnodeindex,curstartedgeidx)
				if curedge.terminal_vertex == 1:
					return self.startingedge.GetArrayNumber(curnodeindex,curstartedgeidx)
			self.active.poll()
			self.isInA.SetFalse(curnodeindex)

		return -1

	def augmentation_stage(self,lastidx):
		bottlecap = self.edges[lastidx].capacity - self.edges[lastidx].flow
		curnodeindex = self.edges[lastidx].initial_vertex
		while curnodeindex != 0 :
			prevedge = self.edges[self.node[curnodeindex].prevedgeindex]
			if bottlecap > (prevedge.capacity - prevedge.flow):
				bottlecap = prevedge.capacity - prevedge.flow
			curnodeindex = prevedge.initial_vertex

		prevedgeindex = -1
		curedgeindex = lastidx
		while curedgeindex != -1:
			self.edges[curedgeindex].flow += bottlecap
			self.edges[self.edges[curnodeindex].invedgeindex].flow -= bottlecap
			prevedgeindex = self.node[self.edges[curedgeindex].initial_vertex].prevedgeindex
			if (self.edges[curedgeindex].capacity - self.edges[curedgeindex].flow) <= esp:
				self.node[self.edges[curedgeindex].terminal_vertex].prevedgeindex = -1
				self.orphan.add_first(self.edges[curedgeindex].terminal_vertex)
		return

	def get_root_of(self,nodeidx):
		currootidx = nodeidx
		while self.node[currootidx].prevedgeindex > 0 :
			currootidx = self.edges[self.node[currootidx].prevedgeindex].initial_vertex
		return currootidx

	def adoption_stage(self):
		while not self.orphan.is_empty():
			curnodeidx = self.orphan.poll()
			hasfindparent = false
			curstartedgeidx = 0
			while  not hasfindparent and curstartedgeidx < self.startingedge.GetLengthFromIdx(curnodeidx):
				curedge = self.edges[self.edges[self.startingedge.GetArrayNumber(curnodeidx,curstartedgeidx)].invedgeindex]
				if (curedge.capacity - curedge.flow) <= eps:
					curstartedgeidx += 1
					continue
				if self.isInS.Get(curedge.initial_vertex) == 0 :
					curstartedgeidx += 1
					continue
				currootnodeidx = curedge.initial_vertex
				while self.node[currootnodeidx].prevedgeindex > 0 :
					currootnodeidx = self.edges[self.node[currootnodeidx].prevedgeindex].initial_vertex
				if currootnodeidx != 0 :
					curstartedgeidx += 1
					continue
				hasfindparent = True
				self.node[curnodeidx].prevedgeindex = self.edges[self.startingedge.GetArrayNumber(curnodeidx,curstartedgeidx)].invedgeindex
				break

			if not hasfindparent :
				self.isInS.SetFalse(curnodeidx)
				self.isInA.SetFalse(curnodeidx)
				for curstartedgeidx in xrange(self.startingedge.GetLengthFromIdx(curnodeidx)):
					curedge = self.edges[self.startingedge.GetArrayNumber(curnodeidx,curstartedgeidx)]
					if self.node[curedge.terminal_vertex].prevedgeindex == self.startingedge.GetArrayNumber(curnodeidx,curstartedgeidx):
						self.node[curedge.terminal_vertex].prevedgeindex = -1
						self.orphan.add_last(curedge.terminal_vertex)
					if self.isInS.Get(curedge.terminal_vertex) == 0 :
						continue
					curedge = self.edges[curedge.invedgeindex]
					if (curedge.capacity - curedge.flow) > eps and self.isInS.Get(curedge.initial_vertex) == 0 :
						self.active.add(curedge.initial_vertex)
						self.isInA.SetTrue(curedge.initial_vertex)




def Calculate(bkgraph):
	bkgraph.do_cut()
	flow = bkgraph.get_flow()
	sys.stdout.write('graph flow is %d\n'%(flow))
	return

def MakeGraph(infile):
	sink = 0
	source = 0
	w = 0
	h = 0
	bkgraph = None
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
			if sink == 0  or w == 0 or h == 0:
				sys.stderr.write('can not define sink or not define width or height\n')
				sys.exit(4)
			if bkgraph is None:
				bkgraph = GraphCutBoykovKolmogorov(w,h)
			curs = int(sarr[0])
			curt = int(sarr[1])
			curw = int(sarr[2])
			y1 = curs / w
			x1 = curs - (y1 * w)

			y2 = curt / w
			x2 = curt - (y2 * w)

			if curs == source :
				bkgraph.set_source_weight(x2,y2,w)
				continue

			if curt == sink :
				bkgraph.set_sink_weight(x1,y1,w)
				continue

			bkgraph.set_intern_weight(x1,y1,x2,y2,w)

	return bkgraph


def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile\n'%(sys.argv[0]))
		sys.exit(4)
	bkgraph = MakeGraph(sys.argv[1])
	assert( not (bkgraph is None))
	Calculate(bkgraph)
	return

if __name__ == '__main__':
	main()

