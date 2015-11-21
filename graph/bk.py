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
	pass

class EdgeMat:
	def __init__(self):
		

class StartingEdge:
	def __init__(self):
		self.__startingedge=[]
		self.__startingcnt=0
		self.__startcap = 0
		return

	def SetEdgeCnt(self,idx,cnt):
		if self.__startcap <= idx:
			for x in xrange((idx-self.__startcap-1)):
				self.__startingedge.append([])
				self.__startcap += 1
		self.__startingedge[idx] = cnt*[0]
		return

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
		voisinsEdgeACreer=[[1,0],[1,-1],[0,-1],[-1,-1]]
		self.nbNode = self.w * self.h + 2
		self.nbEdges = (self.w * self.h * 4) + (self.h - 2)* (self.w - 2)*8 +\
			2*5*(self.h+self.w - 4)+4*3
		self.node = self.nbNode * [Node()]
		self.edge = self.nbEdges * [Edge()]
		self.node[0].SetIdx(0)
		self.node[1].SetIdx(1)
		for x in xrange(self.w):
			for y in xrange(self.h):
				i1 = x * self.h + y + 2
				self.node[i1].SetIdx(i1)
		self.curedge = 0
		self.startingedge=StartingEdge()
		self.startingedge.SetEdgeCnt(0,(self.nbNode-2))
		self.startingedge.SetEdgeCnt(1,(self.nbNode-2))

		if self.w == 1 or self.h == 1:
			for x in xrange(self.w):
				for y in xrange(self.h):
					if (x * (x + 1 - self.w)) == 0 and (y * (y+1-self.h)) == 0:
						# for coins
						self.startingedge.SetEdgeCnt((x*self.h+y+2),1+1)
					elif (x*(x+1-self.w)) == 0 or (y * (y + 1-self.h)) == 0 :
						# edges
						self.startingedge.SetEdgeCnt((x*self.h+y+2),2+1)
		else:
			for x in xrange(self.w):
				for y in xrange(self.h):
					if (x*(x+1-self.w)) == 0 and (y*(y+1-self.h)) == 0:
						# coins
						self.startingedge.SetEdgeCnt((x*self.h+y+2),3+2)
					elif (x*(x+1-self.w)) == 0 or (y*(y+1-self.h))==0:
						# edges
						self.startingedge.SetEdgeCnt((x*self.h+y+2),5+2)
					else:
						# in middle
						self.startingedge.SetEdgeCnt((x*self.h+y+2),8+2)
		for x in xrange(self.w):
			for y in xrange(self.h):
				i1 = x*self.h + y + 2
				i2 = 1



