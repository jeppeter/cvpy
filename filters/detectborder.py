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
	sys.stdout.write('UP to increase array size\n')
	sys.stdout.write('DOWN to decrease array size\n')
	sys.stdout.write('LEFT  RIGHT to change type of MORPH from [RECT ,CROSS,ELLIPSE]\n')
	sys.stdout.write('A  to decrease thresh hold\n')
	sys.stdout.write('S to increase thresh hold\n')
	return


def ShowBorder(infile):
	try:
		img = cv2.imread(infile,1)
		assert(len(img.shape) >=2)
	except:
		sys.stderr.write('can not open %s as picture\n'%(infile))
		sys.exit(3)
	key_usage()
	elmsize = 3
	typearr = [cv2.MORPH_RECT,cv2.MORPH_CROSS,cv2.MORPH_ELLIPSE]
	typename = ['RECT','CROSS','ELLIPSE']
	typeint = 0
	thrsh = 40	
	while True:
		elem = cv2.getStructuringElement(typearr[typeint],(elmsize,elmsize))
		dilate = cv2.dilate(img,elem)
		erode = cv2.erode(img,elem)

		result = cv2.absdiff(dilate,erode)
		retval,result = cv2.threshold(result,thrsh,255,cv2.THRESH_BINARY)
		result = cv2.bitwise_not(result)
		imgshow = 'img %s %d thrsh %d'%(typename[typeint],elmsize,thrsh)
		cv2.imshow(imgshow,result)
		k = cv2.waitKey(0)
		cv2.destroyAllWindows()
		if k not in [UPKEY,RIGHTKEY,LEFTKEY,DOWNKEY,AKEY,SKEY]:
			sys.stdout.write('k %d'%(k))
			break
		if k == UPKEY :
			elmsize += 1
		elif k == DOWNKEY :
			elmsize -= 1
			if elmsize < 1 :
				elmsize = 1
		elif k == LEFTKEY:
			typeint -= 1
			typeint %= len(typearr)
		elif k == RIGHTKEY:
			typeint += 1
			typeint %= len(typearr)
		elif k == AKEY:
			thrsh -= 1
			if thrsh < 0 :
				thrsh = 0
		elif k == SKEY:
			thrsh += 1			

	return


def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	ShowBorder(sys.argv[1])

if __name__ == '__main__':
	main()

