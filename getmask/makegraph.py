#!/usr/bin/python

import sys
import logging
import cv2
import numpy as np
import math

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
	diff = (ff - tf ) ** 2
	weight = math.exp(- math.sqrt(diff))
	if weight < CONSTANT.EPSILON:
		weight = 0.0
	return weight

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

	h = simg.shape[0]
	w = simg.shape[1]
	refered = np.zeros((h,w,4),np.uint8)

	for i in range(w):
		for j in range(h):
			if (i == 0 and j == 0) :
				# these are the corner
				if refered[i][j][CONSTANT.DIR_EAST] == 0:
					refered[i][j][CONSTANT.DIR_EAST] = 1
					assert(refered[i][j+1][CONSTANT.DIR_WEST] == 0)
					refered[i][j+1][CONSTANT.DIR_WEST] = 1
					val = calc_edge(simg[i][j][2],simg[i][j+1][2])
					fp.write('%d,%d,%f\n'%((i*h+j),(i*h+j+1),val))
				else:
					assert(refered[i][j+1][CONSTANT.DIR_WEST]==1)

				if refered[i][j][CONSTANT.DIR_SOUTH] == 0:
					refered[i][j][CONSTANT.DIR_SOUTH] = 1
					assert(refered[i+1][j][CONSTANT.DIR_NORTH]==0)
					refered[i+1][j][CONSTANT.DIR_NORTH] = 1
					val = calc_edge(simg[i][j][2],simg[i+1][j][2])
					fp.write('%d,%d,%f\n'%((i*h+j),((i+1)*h+j),val))
				else:
					assert(refered[i+1][j][CONSTANT.DIR_NORTH]==1)
			elif (j == 0 and i == (w-1)):


