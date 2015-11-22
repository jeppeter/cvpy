import cv2
import sys
import numpy as np


def calcDrawHisto(img):
	h = np.zeros((256,256,3))
	bins = np.arange(256).reshape(256,1)
	color = [(255,0,0),(0,255,0),(0,0,255)]

	for ch,clr in enumerate(color):
		origHist = cv2.calcHist([img],[ch],None,[256],[0,256])
		cv2.normalize(origHist,origHist,0,255*0.9,cv2.NORM_MINMAX)
		hist = np.int32(np.around(origHist))
		pts = np.column_stack((bins,hist))
		cv2.polylines(h,[pts],False,clr)
	return np.flipud(h)

def ShowHist(infile):
	try:
		simg = cv2.imread(infile,1)
		assert(len(simg.shape)>=2)
	except:
		sys.stderr.write('can not read %s'%(infile))
		return -1

	himg = calcDrawHisto(simg)
	cv2.imshow('histo',himg)
	cv2.waitKey(0)
	return 0

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile\n'%(sys.argv[0]))
		sys.exit(3)
	ShowHist(sys.argv[1])

if __name__ == '__main__':
	main()
