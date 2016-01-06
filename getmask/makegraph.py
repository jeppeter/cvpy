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
	def __init__(self,fp,simg):
		self.__fp = fp
		self.__w = simg.shape[0]
		self.__h = simg.shape[1]
		self.__simg = simg
		self.__fp.write('height=%d\n'%(self.__h))
		self.__fp.write('width=%d\n'%(self.__w))
		self.__refered = np.zeros((self.__w,self.__h,4),np.uint8)
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
		self.__fp.write('# [%d][%d][2] %d [%d][%d][2] %d\n'%(fromi,fromj,self.__simg[fromi][fromj][2],toi,toj,self.__simg[toi][toj][2]))
		self.__fp.write('%d,%d,%f\n'%((fromi*self.__h + fromj),(toi*self.__h+toj),val))

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



def outgraph(infile,outfile=None):
	fp = sys.stdout
	if outfile is not None:
		fp = open(outfile,'w')
	try:
		simg = cv2.imread(infile)
		assert(len(simg.shape) >= 2)
	except:
		sys.stderr.write('can not load(%s) error\n'%(infile))
		return

	w = simg.shape[0]
	h = simg.shape[1]
	eo = EdgeOut(fp,simg)
	for i in range(w):
		for j in range(h):
			eo.out_edges(i,j)
	if fp != sys.stdout:
		fp.close()
	eo = None
	fp = None
	simg = None
	return

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile [outfile]\n'%(sys.argv[0]))
		sys.exit(4)
	outfile = None
	infile = sys.argv[1]
	logging.basicConfig(level=logging.DEBUG,format='%(filename)s:%(lineno)d\t%(message)s')
	if len(sys.argv) > 2:
		outfile = sys.argv[2]
	outgraph(infile,outfile)
	return

if __name__ == '__main__':
	main()

