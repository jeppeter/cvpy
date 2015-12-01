#!/usr/bin/python

import cv2
import sys
import logging
import numpy as np
import os

class Rect:
	def __init__(self,x1,y1,x2,y2):
		if x1 < x2:
			self.__left = x1
			self.__right = x2
		else:
			self.__left = x2
			self.__right = x1
		if y1 < y2:
			self.__up = y1
			self.__down = y2
		else:
			self.__up = y2
			self.__down = y1
		return

	def Left(self):
		return self.__left
	def Right(self):
		return self.__right
	def Up(self):
		return self.__up
	def Down(self):
		return self.__down

class MouseRegion:
	def __resetxy (self):
		self.__cursx = -1
		self.__cursy = -1
		self.__curex = -1
		self.__curey = -1
		return
	def __init__(self):
		self.__sourcereg = []
		self.__sinkreg = []
		self.__resetxy()
		return

	def ShiftSource(self):
		if len(self.__sourcereg) == 0 :
			return None
		reg = self.__sourcereg[0]
		self.__sourcereg = self.__sourcereg[1:]
		return reg

	def ShiftSink(self):
		if len(self.__sinkreg) == 0 :
			return None
		reg = self.__sinkreg[0]
		self.__sinkreg = self.__sinkreg[1:]
		return reg

	def SourceStart(self,x,y):
		self.__cursx = x
		self.__cursy = y
		return

	def SourceEnd(self,x,y):
		self.__curex = x
		self.__curey = y
		reg = Rect(self.__cursx,self.__cursy,self.__curex,self.__curey)
		self.__sourcereg.append(reg)
		self.__resetxy()
		return

	def SinkStart(self,x,y):
		self.__cursx = x
		self.__cursy = y
		return

	def SinkEnd(self,x,y):
		self.__curex = x
		self.__curey = y
		reg = Rect(self.__cursx,self.__cursy,self.__curex,self.__curey)
		self.__sinkreg.append(reg)
		self.__resetxy()
		return


def GetMouseEvent(event,x,y,flags,param):
	if event == cv2.EVENT_LBUTTONDOWN and (flags & 0x10) == 0x10:
		logging.info('sourcestart (%d:%d) flags (%d)'%(x,y,flags))
		param.SourceStart(x,y)
	elif event == cv2.EVENT_LBUTTONUP and (flags & 0x10) == 0x10:
		logging.info('sourceend (%d:%d) flags (%d)'%(x,y,flags))
		param.SourceEnd(x,y)
	elif event == cv2.EVENT_LBUTTONDOWN:
		logging.info('sinkstart (%d:%d) flags (%d)'%(x,y,flags))
		param.SinkStart(x,y)
	elif event == cv2.EVENT_LBUTTONUP:
		logging.info('sinkstart (%d:%d) flags (%d)'%(x,y,flags))
		param.SinkEnd(x,y)
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
	i = 0
	while True:
		r = selects.ShiftSource()
		if r is None:
			break
		i += 1
		logging.info('source get (%d:%d)->(%d:%d)'%(r.Left(),r.Up(),r.Right(),r.Down()))
		for i in range(r.Left(),r.Right()):
			for j in range(r.Up(),r.Down()):
				dimg[j][i] = (1,1,250)
	if i == 0 :
		raise Exception('please select a region for source by mouse click with SHIFT')
	i = 0
	while True:
		r = selects.ShiftSink()
		if r is None:
			break
		i += 1
		logging.info('sink get (%d:%d)->(%d:%d)'%(r.Left(),r.Up(),r.Right(),r.Down()))
		for i in range(r.Left(),r.Right()):
			for j in range(r.Up(),r.Down()):
				dimg[j][i] = (250,1,1)

	if i == 0 :
		raise Exception('please select a region for sink by mouse click')
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