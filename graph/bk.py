import sys
import logging


#########################################
# this code is from the code in the website
# https://en.wikibooks.org/wiki/Algorithm_Implementation/Graphs/Maximum_flow/Boykov_%26_Kolmogorov
# and it will be for the boykov kolmogorov 
#########################################

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


class StartingEdge:
	def __init__(self):
		self.__startingedge=[]
		self.__startingcnt=0
		self.__startcap = 0
		return

	def SetIntArray(self,idx,cnt):
		if self.__startcap <= idx:
			for x in xrange((idx-self.__startcap-1)):
				self.__startingedge.append([])
				self.__startcap += 1
		self.__startingedge[idx] = cnt*[0]
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
		assert(len(self.__startingedge) > idx)
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
	eps = 0.001
	def __indice_part(self,x,y):
		return x*self.h + y

	def __init__(self,width,height):
		self.w = width
		self.h = height
		voisinsEdgeACreer=[[1,0],[1,-1],[0,-1],[-1,-1]]
		self.nbNode = self.w * self.h + 2
		self.nbEdges = (self.w * self.h * 4) + (self.h - 2)* (self.w - 2)*8 +\
			2*5*(self.h+self.w - 4)+4*3
		self.node = self.nbNode * [Node()]
		self.edge = self.nbEdges * [Edge()]
		self.curNbVoisins = self.nbNode * [0]
		self.node[0].SetIdx(0)
		self.node[1].SetIdx(1)
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

	



