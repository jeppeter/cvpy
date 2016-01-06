#!/usr/bin/python

import sys
import logging
import cv2
import numpy as np
import math
import logging

class CONSTANT(object):
	DIR_WEST=0
	DIR_EAST=1
	DIR_NORTH=2
	DIR_SOUTH=3
	EPSILON=1e-6
	def __setattr__(self,*_):
		pass

def calc_edge(fromval ,toval):
	ff = float(fromval)
	tf = float(toval)
	diff = ((ff - tf ) ** 2)
	weight = math.exp(- math.sqrt(diff))
	if weight < CONSTANT.EPSILON:
		weight = CONSTANT.EPSILON
	return weight



class EdgeOut(object):
	def __init__(self,fp,simg,maskimg):
		self.__fp = fp
		self.__w = simg.shape[0]
		self.__h = simg.shape[1]
		self.__simg = simg
		self.__maskimg = maskimg
		self.__fp.write('height=%d\n'%(self.__h))
		self.__fp.write('width=%d\n'%(self.__w))
		self.__refered = np.zeros((self.__w,self.__h,4),np.uint8)
		self.__edgeidx= 0
		return

	def __out_edges(self,fromi,fromj,toi,toj):
		if toi < 0 or toi >= self.__w:
			return 0
		if toj < 0 or toj >= self.__h:
			return 0

		if fromi < 0 or fromi >= self.__w:
			return 0
		if fromj < 0 or fromj >= self.__h:
			return 0


		if fromi > toi and self.__refered[fromi][fromj][CONSTANT.DIR_WEST] != 0:
			assert(self.__refered[toi][toj][CONSTANT.DIR_EAST] != 0)
			return 0

		if fromi < toi and self.__refered[fromi][fromj][CONSTANT.DIR_EAST] != 0 :
			assert(self.__refered[toi][toj][CONSTANT.DIR_WEST] != 0)
			return 0

		if fromj > toj and self.__refered[fromi][fromj][CONSTANT.DIR_NORTH] != 0 :
			assert(self.__refered[toi][toj][CONSTANT.DIR_SOUTH] != 0)
			return 0

		if fromj < toj and self.__refered[fromi][fromj][CONSTANT.DIR_SOUTH] != 0:
			assert(self.__refered[toi][toj][CONSTANT.DIR_NORTH] != 0)
			return 0

		#logging.info('[%d][%d][2] = %d [%d][%d][2] = %d'%(fromi,fromj,self.__simg[fromi][fromj][2],toi,toj,self.__simg[toi][toj][2]))
		val = calc_edge(self.__simg[fromi][fromj][2],self.__simg[toi][toj][2])
		cap = '%f'%(val)
		rcap = '%f'%(val)

		if self.__maskimg[fromi][fromj][1] == 1:
			logging.info('from[%d][%d] mask'%(fromi,fromj))
			rcap = '1.#INF00'

		if self.__maskimg[toi][toj][1] == 1:
			logging.info('to[%d][%d] mask'%(toi,toj))
			cap = '1.#INF00'

		self.__fp.write('# edge[%d] vert[%d][%d] -> vert[%d][%d] .cap(%s) .rcap(%s)\n'%(\
			self.__edgeidx,toi,toj,fromi,fromj,cap,rcap))
		self.__fp.write('%d,%d,%s\n'%((fromi*self.__h + fromj),(toi*self.__h+toj),rcap))
		self.__fp.write('%d,%d,%s\n'%((toi*self.__h+toj),(fromi*self.__h+fromj),cap))
		self.__edgeidx += 1

		if fromi > toi :
			self.__refered[fromi][fromj][CONSTANT.DIR_WEST] = 1
			assert(self.__refered[toi][toj][CONSTANT.DIR_EAST] == 0)
			self.__refered[toi][toj][CONSTANT.DIR_EAST] = 1
		elif fromi < toi :
			self.__refered[fromi][fromj][CONSTANT.DIR_EAST] = 1
			assert(self.__refered[toi][toj][CONSTANT.DIR_WEST] == 0)
			self.__refered[toi][toj][CONSTANT.DIR_WEST] = 1
		elif fromj > toj:
			self.__refered[fromi][fromj][CONSTANT.DIR_NORTH] = 1
			assert(self.__refered[toi][toj][CONSTANT.DIR_SOUTH] == 0)
			self.__refered[toi][toj][CONSTANT.DIR_SOUTH] = 1
		elif fromj < toj:
			self.__refered[fromi][fromj][CONSTANT.DIR_SOUTH] = 1
			assert(self.__refered[toi][toj][CONSTANT.DIR_NORTH] == 0)
			self.__refered[toi][toj][CONSTANT.DIR_NORTH] = 1
		return 1

		
	def out_edges(self,fromi,fromj):
		self.__out_edges(fromi,fromj,fromi-1,fromj)
		self.__out_edges(fromi,fromj,fromi+1,fromj)
		self.__out_edges(fromi,fromj,fromi,fromj-1)
		self.__out_edges(fromi,fromj,fromi,fromj+1)
		return

	def out_all_edges(self):
		for i in range(self.__w):
			for j in range(self.__h):
				self.__out_edges(i,j,i,j+1)
		for i in range(self.__w):
			for j in range(self.__h):
				self.__out_edges(i,j,i+1,j)



def outgraph(infile,maskfile,outfile=None):
	fp = sys.stdout
	if outfile is not None:
		fp = open(outfile,'w')
	try:
		simg = cv2.imread(infile)
		maskimg = cv2.imread(maskfile)
		assert(len(simg.shape) >= 2)
		assert(len(maskimg.shape) >= 2)
		assert(simg.shape[0] == maskimg.shape[0])
		assert(simg.shape[1] == maskimg.shape[1])
	except:
		sys.stderr.write('can not load(%s) or (%s) error\n'%(infile,maskfile))
		return

	w = simg.shape[0]
	h = simg.shape[1]
	eo = EdgeOut(fp,simg,maskimg)
	eo.out_all_edges()
	if fp != sys.stdout:
		fp.close()
	eo = None
	fp = None
	simg = None
	return

def main():
	if len(sys.argv) < 3:
		sys.stderr.write('%s infile maskfile [outfile]\n'%(sys.argv[0]))
		sys.exit(4)
	outfile = None
	maskfile = sys.argv[2]
	infile = sys.argv[1]
	logging.basicConfig(level=logging.DEBUG,format='%(filename)s:%(lineno)d\t%(message)s')
	if len(sys.argv) > 3:
		outfile = sys.argv[3]
	outgraph(infile,maskfile,outfile)
	return

if __name__ == '__main__':
	main()

