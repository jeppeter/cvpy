import cv2
import numpy as np
import sys

UPKEY=2490368
DOWNKEY=2621440
LEFTKEY=2424832
RIGHTKEY=2555904


def ShowBorder(infile):
	try:
		img = cv2.imread(infile,1)
		assert(len(img.shape) >=2)
	except:
		sys.stderr.write('can not open %s as picture\n'%(infile))
		sys.exit(3)
	elmsize = 3
	while True:
		elem = cv2.getStructuringElement(cv2.MORPH_RECT,(elmsize,elmsize))
		dilate = cv2.dilate(img,elem)
		erode = cv2.erode(img,elem)

		result = cv2.absdiff(dilate,erode)
		retval,result = cv2.threshold(result,40,255,cv2.THRESH_BINARY)
		result = cv2.bitwise_not(result)
		cv2.imshow('img',result)
		k = cv2.waitKey(0)
		if k not in [UPKEY,RIGHTKEY,LEFTKEY,DOWNKEY]:
			break
		if k == UPKEY :
			elmsize += 1
		elif k == DOWNKEY :
			elmsize -= 1
			if elmsize < 1 :
				elmsize = 1
	return


def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	ShowBorder(sys.argv[1])

if __name__ == '__main__':
	main()

