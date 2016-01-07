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
	EPSILON=float(0.000001)
	STR_EPSILON='0.000001'
	STR_CAP_INF='1.#INF00'
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
		self.__h = simg.shape[0]
		self.__w = simg.shape[1]
		self.__simg = simg
		self.__maskimg = maskimg
		self.__fp.write('width=%d\n'%(self.__w))
		self.__fp.write('height=%d\n'%(self.__h))
		self.__refered = np.zeros((self.__h,self.__w,4),np.uint8)
		self.__edgeidx= 0
		sidx = 0
		brkone = 0
		logging.info('h %d w %d'%(self.__h,self.__w))
		for j in range(self.__h):
			for i in range(self.__w):
				if self.__is_source(j,i):
					brkone = 1
					break
				sidx += 1
			if brkone :
				break
		if brkone == 0 :
			raise Exception('can not find source idx')
		self.__fp.write('source=%d\n'%(sidx))
		eidx=0
		brkone = 0
		for j in range(self.__h):
			for i in range(self.__w):
				if self.__is_sink(j,i):
					brkone = 1
					break
				eidx += 1
			if brkone :
				break
		if brkone == 0 :
			raise Exception('can not find sink idx')
		self.__fp.write('sink=%d\n'%(eidx))
		return

	def __is_source(self,fromi,fromj):
		if self.__maskimg[fromi][fromj][0] <= 10 and self.__maskimg[fromi][fromj][1] <= 10 and self.__maskimg[fromi][fromj][2] >= 240:
			return True
		return False

	def __is_sink(self,fromi,fromj):
		if self.__maskimg[fromi][fromj][0] >= 240 and self.__maskimg[fromi][fromj][1] <= 10 and self.__maskimg[fromi][fromj][2] <= 10:
			return True
		return False

	def __out_edges(self,fromj,fromi,toj,toi):
		if toi < 0 or toi >= self.__w:
			return 0
		if toj < 0 or toj >= self.__h:
			return 0

		if fromi < 0 or fromi >= self.__w:
			return 0
		if fromj < 0 or fromj >= self.__h:
			return 0

		#logging.info('[%d][%d] -> [%d][%d]'%(fromj,fromi,toj,toi))

		if fromi > toi and self.__refered[fromj][fromi][CONSTANT.DIR_WEST] != 0:
			assert(self.__refered[toj][toi][CONSTANT.DIR_EAST] != 0)
			return 0

		if fromi < toi and self.__refered[fromj][fromi][CONSTANT.DIR_EAST] != 0 :
			assert(self.__refered[toj][toi][CONSTANT.DIR_WEST] != 0)
			return 0

		if fromj > toj and self.__refered[fromj][fromi][CONSTANT.DIR_NORTH] != 0 :
			assert(self.__refered[toj][toi][CONSTANT.DIR_SOUTH] != 0)
			return 0

		if fromj < toj and self.__refered[fromj][fromi][CONSTANT.DIR_SOUTH] != 0:
			assert(self.__refered[toj][toi][CONSTANT.DIR_NORTH] != 0)
			return 0

		#logging.info('[%d][%d][2]  [%d][%d][2]'%(fromj,fromi,toj,toi))
		val = calc_edge(self.__simg[fromj][fromi][2],self.__simg[toj][toi][2])
		cap = '%f'%(val)
		rcap = '%f'%(val)

		if self.__is_source(fromj,fromi):
			logging.info('from[%d][%d] mask'%(fromj,fromi))
			rcap = CONSTANT.STR_CAP_INF
			cap = '%f'%(val)
		elif self.__is_source(toj,toi):
			logging.info('to[%d][%d] mask'%(toj,toi))
			cap = CONSTANT.STR_CAP_INF
			rcap = '%f'%(val)

		if self.__is_sink(toj,toi):
			logging.info('to[%d][%d] mask'%(toj,toi))
			rcap = CONSTANT.STR_CAP_INF
			cap = '%f'%(val)
		elif self.__is_sink(fromj,fromi):
			logging.info('from[%d][%d] mask'%(fromj,fromi))
			rcap = '%f'%(val)
			cap = CONSTANT.STR_CAP_INF

		if self.__is_sink(toj,toi) and self.__is_source(fromj,fromi):
			logging.info('from[%d][%d] to[%d][%d] mask'%(fromj,fromi,toj,toi))
			rcap = CONSTANT.STR_CAP_INF
			cap = '%f'%(val)

		if self.__is_source(toj,toi) and self.__is_sink(fromj,fromi):
			rcap = '%f'%(val)
			cap = CONSTANT.STR_CAP_INF

		if self.__is_sink(toj,toi) and self.__is_sink(fromj,fromi):
			cap = CONSTANT.STR_CAP_INF
			rcap = CONSTANT.STR_CAP_INF

		if self.__is_source(toj,toi) and self.__is_source(fromj,fromi):
			cap = CONSTANT.STR_CAP_INF
			rcap = CONSTANT.STR_CAP_INF


		self.__fp.write('# edge[%d] vert[%d][%d] -> vert[%d][%d] .cap(%s) .rcap(%s)\n'%(\
			self.__edgeidx,toj,toi,fromj,fromi,cap,rcap))
		self.__fp.write('%d,%d,%s\n'%((fromj*self.__w + fromi),(toj*self.__w+toi),rcap))
		self.__fp.write('%d,%d,%s\n'%((toj*self.__w+toi),(fromj*self.__w+fromi),cap))
		self.__edgeidx += 1

		if fromi > toi :
			self.__refered[fromj][fromi][CONSTANT.DIR_WEST] = 1
			assert(self.__refered[toj][toi][CONSTANT.DIR_EAST] == 0)
			self.__refered[toj][toi][CONSTANT.DIR_EAST] = 1
		elif fromi < toi :
			self.__refered[fromj][fromi][CONSTANT.DIR_EAST] = 1
			assert(self.__refered[toj][toi][CONSTANT.DIR_WEST] == 0)
			self.__refered[toj][toi][CONSTANT.DIR_WEST] = 1
		elif fromj > toj:
			self.__refered[fromj][fromi][CONSTANT.DIR_NORTH] = 1
			assert(self.__refered[toj][toi][CONSTANT.DIR_SOUTH] == 0)
			self.__refered[toj][toi][CONSTANT.DIR_SOUTH] = 1
		elif fromj < toj:
			self.__refered[fromj][fromi][CONSTANT.DIR_SOUTH] = 1
			assert(self.__refered[toj][toi][CONSTANT.DIR_NORTH] == 0)
			self.__refered[toj][toi][CONSTANT.DIR_NORTH] = 1
		return 1

		
	def out_edges(self,fromj,fromi):
		self.__out_edges(fromj,fromi,fromj-1,fromi)
		self.__out_edges(fromj,fromi,fromj+1,fromi)
		self.__out_edges(fromj,fromi,fromj,fromi-1)
		self.__out_edges(fromj,fromi,fromj,fromi+1)
		return

	def out_all_edges(self):
		for j in range(self.__h):
			for i in range(self.__w):
				self.__out_edges(j,i,j,i+1)
		for j in range(self.__h):
			for i in range(self.__w):
				self.__out_edges(j,i,j+1,i)



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
	logging.basicConfig(level=logging.DEBUG,format='%(filename)s:%(lineno)d %(message)s')
	if len(sys.argv) > 3:
		outfile = sys.argv[3]
	outgraph(infile,maskfile,outfile)
	return

if __name__ == '__main__':
	main()

