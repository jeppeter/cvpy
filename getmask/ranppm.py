#!/usr/bin/python

import cv2
import random
import sys
import numpy as np
import logging
import time


def MakeRandomPPM(fname,height,width):
	dimg = np.zeros((height,width,3),np.uint8)
	mask = np.zeros((height,width,3),np.uint8)
	mask.fill(255)

	for i in xrange(width):
		for j in xrange(height):
			dimg[j][i]=( random.randint(0,255),random.randint(0,255),random.randint(0,255))
	# now to set for the random source and random int
	sourceid = 0
	sinkid = 0
	while sourceid == sinkid:
		sourceid = random.randint(0,height*width-1)
		sinkid = random.randint(0,height*width-1)

	# now set the source
	j = sourceid / width
	i = sourceid % width
	mask[j][i] = (1,1,250)

	j = sinkid / width
	i = sinkid % width
	mask[j][i] = (250,1,1)
	dname = '%s.ppm'%(fname)
	mname = '%sseg.ppm'%(fname)
	cv2.imwrite(dname,dimg)
	cv2.imwrite(mname,mask)
	logging.info('%s for ppm and %s for mask'%(dname,mname))
	return

def main():
	if len(sys.argv) < 4:
		sys.stderr.write('%s fnametemplate height widht\n'%(sys.argv[0]))
		sys.exit(4)
	random.seed(time.time())
	MakeRandomPPM(sys.argv[1],int(sys.argv[2]),int(sys.argv[3]))

if __name__ == '__main__':
	logging.basicConfig(level=logging.DEBUG,format='%(filename)s:%(funcName)s:%(lineno)d %(message)s')
	main()