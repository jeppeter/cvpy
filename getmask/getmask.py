#!/usr/bin/python

import cv2
import sys
import logging
import numpy as np
import os

class MouseRegion:
	def __resetxy (self):
		self.__cursx = -1
		self.__cursy = -1
		self.__curex = -1
		self.__curey = -1
		return
	def __init__(self):
		self.__startx = []
		self.__starty = []
		self.__endx = []
		self.__endy = []
		self.__resetxy()
		return

	def ShiftValue(self):
		if len(self.__startx) == 0 :
			return (None,None,None,None)
		sx = self.__startx[0]
		sy = self.__starty[0]
		ex = self.__endx[0]
		ey = self.__endy[0]
		self.__startx = self.__startx[1:]
		self.__starty = self.__starty[1:]
		self.__endx = self.__endx[1:]
		self.__endy = self.__endy[1:]
		return (sx,sy,ex,ey)

	def Start(self,x,y):
		self.__cursx = x
		self.__cursy = y
		return

	def End(self,x,y):
		self.__curex = x
		self.__curey = y
		if self.__cursx < self.__curex :
			self.__startx.append(self.__cursx)
			self.__endx.append(self.__curex)
		else:
			self.__startx.append(self.__curex)
			self.__endx.append(self.__cursx)

		if self.__cursy < self.__curey:
			self.__starty.append(self.__cursy)
			self.__endy.append(self.__curey)
		else:
			self.__starty.append(self.__curey)
			self.__endy.append(self.__cursy)
		self.__resetxy()
		return


def GetMouseEvent(event,x,y,flags,param):
	if event == cv2.EVENT_LBUTTONDOWN:
		logging.info('start (%d:%d)'%(x,y))
		param.Start(x,y)
	elif event == cv2.EVENT_LBUTTONUP:
		logging.info('end (%d:%d)'%(x,y))
		param.End(x,y)
	return

def GetMask(infile):
	try:
		simg = cv2.imread(infile)
		assert(len(simg.shape) >= 2)
	except:
		sys.stderr.write('can not open(%s)\n'%(infile))
		return
	h = simg.shape[0]
	w = simg.shape[1]
	b = os.path.basename(infile)
	bname,extname = os.path.splitext(b)
	logging.info('w (%d) h (%d)'%(w,h))
	dimg = np.zeros((h,w,3), np.uint8)
	dimg.fill(255)
	selects = MouseRegion()
	cv2.namedWindow(bname)
	cv2.setMouseCallback(bname,GetMouseEvent,selects)
	cv2.imshow(bname,simg)
	cv2.waitKey(0)
	while True:
		sx,sy,ex,ey = selects.ShiftValue()
		if sx is None:
			break
		logging.info('get (%d:%d)->(%d:%d)'%(sx,sy,ex,ey))
		for i in range(sx,ex):
			for j in range(sy,ey):
				dimg[j][i] = (1,1,250)
	sname = '%s.ppm'%(bname)
	mname = '%sseg.ppm'%(bname)
	cv2.imwrite(sname,simg)
	cv2.imwrite(mname,dimg)
	return


def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile\n'%(sys.argv[0]))
		sys.exit(4)
	GetMask(sys.argv[1])

if __name__ == '__main__':
	logging.basicConfig(level=logging.DEBUG,format='%(filename)s:%(funcName)s:%(lineno)d %(message)s')
	main()