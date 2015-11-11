import sys
import cv2
import numpy as np
import logging

from matplotlib import pyplot as plt


def Show2DHist(infile):
	try:
		simg = cv2.imread(infile)
		assert(len(simg)>=0)
	except:
		sys.stderr.write('can not open (%s) for picture\n'%(infile))
		return
	color = ('b','g','r')
	hsv = cv2.cvtColor(simg,cv2.COLOR_BGR2HSV)
	hist = cv2.calcHist( [hsv], [0, 1], None, [180, 256], [0, 180, 0, 256] )
	plt.imshow(hist,interpolation = 'nearest')
	plt.show()

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	logging.basicConfig(level=logging.DEBUG,format="%(levelname)-8s [%(filename)-10s:%(funcName)-20s:%(lineno)-5s] %(message)s")
	Show2DHist(sys.argv[1])

if __name__ == '__main__':
	main()
