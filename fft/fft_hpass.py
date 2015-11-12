#!python
# coding=utf-8


import cv2
import numpy as np
import logging
import sys

UPKEY=2490368
DOWNKEY=2621440
LEFTKEY=2424832
RIGHTKEY=2555904

class MaskArgs:
	def __init__(self,row,col):
		self.margin = 30
		self.filter = 0
		self.row = row
		self.col = col

	def name(self):
		return 'margin %d filter %d'%(self.margin,self.filter)
	def upmargin(self):
		self.margin += 1

		if self.margin >= (self.row / 2):
			self.margin = (self.row/2 - 1)

		if self.margin >= (self.col / 2):
			self.margin = (self.col / 2 - 1)
		return

	def downmargin(self):
		self.margin -= 1
		if self.margin <= 0 :
			self.margin = 0
		return

	def upfilter(self):
		self.filter += 1
		if self.filter >= 255:
			self.filter = 255
		return

	def downfilter(self):
		self.filter -= 1
		if self.filter <= 0:
			self.filter = 0
		return



def ChangeMask(maskmat,args):
	rows ,cols=maskmat.shape[:2]
	maskmat[rows/2-args.margin:rows/2+args.margin,cols/2-args.margin:cols/2+args.margin] = args.filter
	return maskmat


def HighPassShow(infile):
	try:
		simg = cv2.imread(infile,0) #直接读为灰度图像
		assert(len(simg.shape)>=2)
	except:
		sys.stderr.write('can not load(%s) for picture\n'%(infile))
		return
	#--------------------------------
	rows,cols = simg.shape
	mask = np.ones(simg.shape,np.uint8)
	args = MaskArgs(rows,cols)
	while True:
		mask = ChangeMask(mask,args)
		#--------------------------------
		f1 = np.fft.fft2(simg)
		f1shift = np.fft.fftshift(f1)
		f1shift = f1shift*mask
		f2shift = np.fft.ifftshift(f1shift) #对新的进行逆变换
		img_new = np.fft.ifft2(f2shift)
		img_new = np.abs(img_new)
		img_new = np.uint8(img_new)
		cv2.imshow(infile,simg)
		cv2.imshow('mask',mask)
		cv2.imshow(args.name(),img_new)
		k = cv2.waitKey(0)
		cv2.destroyAllWindows()
		if k not in [UPKEY,DOWNKEY,RIGHTKEY,LEFTKEY]:
			break
		if k == UPKEY:
			args.upmargin()
		elif k == DOWNKEY :
			args.downmargin()
		elif k == LEFTKEY:
			args.upfilter()
		elif k == RIGHTKEY:
			args.downfilter()
	return

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	logging.basicConfig(level=logging.DEBUG,format="%(levelname)-8s [%(filename)-10s:%(funcName)-20s:%(lineno)-5s] %(message)s")
	HighPassShow(sys.argv[1])

if __name__ == '__main__':
	main()

