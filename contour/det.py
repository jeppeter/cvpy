import sys
import cv2
import numpy as np
import logging

UPKEY=2490368
DOWNKEY=2621440
LEFTKEY=2424832
RIGHTKEY=2555904
AKEY=97
SKEY=115
ZKEY=122
XKEY=120


def key_usage():
	sys.stdout.write('UP DOWN to manage uplow\n')
	sys.stdout.write('LEFT RIGHT to manage uphigh\n')
	sys.stdout.write('AKEY SKEY to manage appmode\n')
	sys.stdout.write('ZKEY XKEY to manage retrmode\n')
	return


class ContourArgs:
	approxymode=[cv2.CHAIN_APPROX_NONE,cv2.CHAIN_APPROX_SIMPLE,cv2.CHAIN_APPROX_TC89_L1,cv2.CHAIN_APPROX_TC89_KCOS]
	appmodename=['CHAIN_APPROX_NONE','CHAIN_APPROX_SIMPLE','CHAIN_APPROX_TC89_L1','CHAIN_APPROX_TC89_KCOS']
	retrievemode=[cv2.RETR_EXTERNAL,cv2.RETR_LIST,cv2.RETR_CCOMP,cv2.RETR_TREE]
	retrievename=['RETR_EXTERNAL','RETR_LIST','RETR_CCOMP','RETR_TREE']
	def __init__(self):
		self.high = 255
		self.low = 127
		self.appidx = 1
		self.appmode = ContourArgs.approxymode[self.appidx]
		self.retridx = 3
		self.retrmode = ContourArgs.retrievemode[self.retridx]
		return
	def uphigh(self):
		self.high += 1
		if self.high >= 255:
			self.high = 255
		return

	def downhigh(self):
		self.high -= 1
		if self.high <= self.low:
			self.high = self.low + 1
		return

	def uplow(self):
		self.low += 1
		if self.low >= self.high:
			self.low = self.high - 1
		return
	def downlow(self):
		self.low -= 1
		if self.low <= 0:
			self.low = 0
		return
	def  __set_appmode(self):
		self.appmode = ContourArgs.approxymode[self.appidx]
		return

	def upappmode(self):
		self.appidx += 1
		self.appidx %= len(ContourArgs.approxymode)
		self.__set_appmode()
		return

	def downappmode(self):
		self.appidx += 1
		self.appidx %= len(ContourArgs.approxymode)
		self.__set_appmode()
		return

	def __set_retrimode(self):
		self.retrmode = ContourArgs.retrievemode[self.retridx]
		return

	def upretrimode(self):
		self.retridx += 1
		self.retridx %= len(ContourArgs.retrievemode)
		self.__set_retrimode()
		return

	def downretrimode(self):
		self.retridx -= 1
		self.retridx %= len(ContourArgs.retrievemode)
		self.__set_retrimode()
		return

	def name(self):
		_name = ''
		_name += ContourArgs.appmodename[self.appidx]

		_name += ' '
		_name += ContourArgs.retrievename[self.retridx]
		_name += '(%d:%d)'%(self.low,self.high)
		return _name



def find_contour(simg,contargs):
	ret,thrsh = cv2.threshold(simg,contargs.low,contargs.high,0)
	dimg,contours,hirearchy = cv2.findContours(thrsh,contargs.retrmode,contargs.appmode)
	return dimg

def ContourShow(infile):
	try:
		simg = cv2.imread(infile,1)
		simg = cv2.cvtColor(simg,cv2.COLOR_BGR2GRAY)
		assert(len(simg.shape)>=2)
	except :
		sys.stderr.write('can not open %s as input'%(infile))
		return
	contargs = ContourArgs()

	while True:
		dimg = find_contour(simg,contargs)
		name = contargs.name()
		cv2.imshow(name,dimg)
		k = cv2.waitKey(0)
		cv2.destroyAllWindows()
		if k not in [UPKEY,DOWNKEY,LEFTKEY,RIGHTKEY,AKEY,SKEY,ZKEY,XKEY]:
			break
		if k == UPKEY:
			contargs.uplow()
		elif k == DOWNKEY:
			contargs.downlow()
		elif k == LEFTKEY:
			contargs.uphigh()
		elif k == RIGHTKEY:
			contargs.downhigh()
		elif k == AKEY:
			contargs.upappmode()
		elif k == SKEY:
			contargs.downappmode()
		elif k == ZKEY:
			contargs.upretrimode()
		elif k == XKEY:
			contargs.downretrimode()

	return


def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	logging.basicConfig(level=logging.DEBUG,format="%(levelname)-8s [%(filename)-10s:%(funcName)-20s:%(lineno)-5s] %(message)s")
	ContourShow(sys.argv[1])

if __name__ == '__main__':
	main()


