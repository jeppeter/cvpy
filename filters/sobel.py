import cv2
import numpy as np
import sys


UPKEY=2490368
DOWNKEY=2621440
LEFTKEY=2424832
RIGHTKEY=2555904
AKEY=97
SKEY=115

def key_usage():
	sys.stdout.write('UP to increase x weight\n')
	sys.stdout.write('DOWN to decrease x weight\n')
	sys.stdout.write('LEFT  RIGHT to change y weight\n')
	return


def SobelImage(img,xw,yw):
	x = cv2.Sobel(img,cv2.CV_16S,1,0)
	y = cv2.Sobel(img,cv2.CV_16S,0,1)
	absx = cv2.convertScaleAbs(x)
	absy = cv2.convertScaleAbs(y)
	dst = cv2.addWeighted(absx,xw,absy,yw,0)
	return dst


def SobelHandle(infile):
	try:
		img = cv2.imread(infile,0)
		assert(len(img.shape) >=2)
	except:
		sys.stderr.write('can not open %s as picture\n'%(infile))
		sys.exit(3)
	xw = 0.5
	yw = 0.5
	key_usage()
	while True:
		dst = SobelImage(img,xw,yw)
		imgname = 'xw %f yw %f'%(xw,yw)
		cv2.imshow(imgname,dst)
		k = cv2.waitKey(0)
		cv2.destroyAllWindows()
		if k not in [UPKEY,RIGHTKEY,LEFTKEY,DOWNKEY]:
			sys.stdout.write('k %d'%(k))
			break
		if k == UPKEY:
			xw *= 0.9
		elif k == DOWNKEY:
			xw *= 1.1
		elif k == LEFTKEY:
			yw *= 0.9
		elif k == RIGHTKEY:
			yw *= 1.1
	return


def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	SobelHandle(sys.argv[1])

if __name__ == '__main__':
	main()
