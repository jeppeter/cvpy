import sys
import cv2
import numpy as np
import logging

from matplotlib import pyplot as plt

def OutputMat(mat):
	r,c,chl = mat.shape
	logging.info('[%d][%d]'%(r,c))
	for y in range(0,r):
		s = ''
		for x in range(0,c):
			if (x % 16 == 0) and  x != 0 :
				s += '\n'
			s += ' %s '%(mat[y][x])
		logging.info('mat[%d] %s'%(y,s))
	return

def OutputArray(prefix,arr):
	logging.info('%s size(%d)'%(prefix,len(arr)))
	s = ''
	for x in range(0,len(arr)):
		if (x % 16 == 0) and x != 0:
			s += '\n'
		s += ' %s '%(arr[x])
	logging.info('%s'%(s))


def SetBlackForPic(mat):
	r,c,chl = mat.shape
	for y in range(0,r):
		for x in range(0,c):
			cb = x % 256
			cg = x % 256
			cr = x % 256
			mat[y][x] = [cb,cg,cr]
	return mat

def HistoShow(infile):
	try:
		simg = cv2.imread(infile)
		assert(len(simg)>=0)
	except:
		sys.stderr.write('can not open (%s) for picture\n'%(infile))
		return
	color = ('b','g','r')
	simg = SetBlackForPic(simg)
	#OutputMat(simg)
	for i ,col in enumerate(color):
		hist = cv2.calcHist([simg],[i],None,[256],[0,256])		
		plt.plot(hist,color=col)
		plt.xlim([0,256])
	plt.show()

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	logging.basicConfig(level=logging.DEBUG,format="%(levelname)-8s [%(filename)-10s:%(funcName)-20s:%(lineno)-5s] %(message)s")
	HistoShow(sys.argv[1])

if __name__ == '__main__':
	main()



