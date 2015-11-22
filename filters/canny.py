import cv2
import numpy as np
import sys


UPKEY=2490368
DOWNKEY=2621440
LEFTKEY=2424832
RIGHTKEY=2555904
AKEY=97
SKEY=115
QKEY=113

class CannyArgs:
	def __init__(self):
		self.min = 50
		self.max = 150
		return

	def name(self):
		return 'min %d max %d'%(self.min,self.max)

	def incmin(self):
		self.min += 1
		if self.min >= self.max:
			self.min = self.max - 1
		return
	def decmin(self):
		self.min -= 1
		if self.min < 0:
			self.min = 0
		return

	def incmax(self):
		self.max += 1
		return

	def decmax(self):
		self.max -= 1
		if self.max <= self.min:
			self.max = self.min + 1
		return

	def setmax(self,maxval):
		self.max = maxval
		if self.max >= 255:
			self.max = 255
		if self.max <= self.min:
			self.max = self.min + 1
		return

	def setmin(self,minval):
		self.min = minval
		if self.min  <= 0:
			self.min = 0
		if self.min >= self.max:
			self.min = self.max - 1
		return

def key_usage():
	sys.stdout.write('UP to increase min\n')
	sys.stdout.write('DOWN to decrease min\n')
	sys.stdout.write('LEFT  RIGHT to max\n')
	return



def CannyChange(img,args):
	return cv2.Canny(img,args.min,args.max)

def nothing(x):
	sys.stdout.write('set value %s\n'%(x))
	pass

def CannyShow(infile):
	try:
		img = cv2.imread(infile,0)
		assert(len(img.shape) >=2)
	except:
		sys.stderr.write('can not open %s as picture\n'%(infile))
		sys.exit(3)
	args = CannyArgs()
	key_usage()
	simg = img

	while True:
		dst = CannyChange(simg,args)
		name = args.name()
		cv2.namedWindow(name)
		cv2.createTrackbar('threshlow',name,0,255,nothing)
		cv2.setTrackbarPos('threshlow',name,args.min)
		cv2.createTrackbar('threshhigh',name,0,255,nothing)	
		cv2.setTrackbarPos('threshhigh',name,args.max)
		cv2.imshow(name,dst)
		k = cv2.waitKey(0)
		minval = cv2.getTrackbarPos('threshlow',name)
		maxval = cv2.getTrackbarPos('threshhigh',name)
		sys.stdout.write('minval %d maxval %d\n'%(minval,maxval))
		cv2.destroyAllWindows()
		args.setmin(minval)
		args.setmax(maxval)
		if k == QKEY or k == -1:
			break
	return

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	CannyShow(sys.argv[1])

if __name__ == '__main__':
	main()


