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
	pass


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

########################################
# GraphCutBoykovKolmogorov member
#           debug = true
#           eps = 0.001
#           nbNodes
#           nbEdges
#           w,h
#           node[]
#           edge[]
#           startingedge[][]
#           
########################################
class GraphCutBoykovKolmogorov:
	def __indice_part(self,x,y):
		return x*self.h + y

	def __init__(self,width,height):
		self.w = width
		self.h = height
		voisinsEdgeACreer = 4*[2*[0]]
		voisinsEdgeACreer[0][0]=1
		voisinsEdgeACreer[0][1]=0
		voisinsEdgeACreer[1]