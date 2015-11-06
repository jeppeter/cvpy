import numpy as np
import sys
import cv2

UPKEY=2490368
DOWNKEY=2621440
LEFTKEY=2424832
RIGHTKEY=2555904


def DilateFile(infile):
	try:
		simg = cv2.imread(infile,1)
		sys.stdout.write('shape %d'%(len(simg.shape)))
	except :
		sys.stderr.write('can not open %s as input'%(infile))
		return
	kernel = np.ones((5,5),np.uint8)
	iterate = 1
	while True:
		eimg = cv2.dilate(simg,kernel,iterations =iterate)
		cv2.imshow('img',eimg)
		k = cv2.waitKey(0)
		if k not in [UPKEY,RIGHTKEY,LEFTKEY,DOWNKEY]:
			break
		if k == UPKEY :
			iterate += 1
		elif k == DOWNKEY :
			iterate -= 1
			if iterate < 0 :
				iterate = 0
	return

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	DilateFile(sys.argv[1])

if __name__ == '__main__':
	main()

