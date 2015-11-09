import sys
import cv2
import numpy as np

UPKEY=2490368
DOWNKEY=2621440


def key_usage():
	sys.stdout.write('UP to pyrup\n')
	sys.stdout.write('DOWN to pyrdown\n')
	return


def PyrMid(simg ,updown=1):
	dimg = simg
	if updown > 0 :
		dimg = cv2.pyrUp(simg)
	else:
		dimg = cv2.pyrDown(simg)
	return dimg


def PyrShow(infile):
	try:
		simg = cv2.imread(infile,1)
		assert(len(simg.shape)>=2)
	except :
		sys.stderr.write('can not open %s as input'%(infile))
		return
	key_usage()
	updown = 1
	dimg = simg
	while True:
		dimg = PyrMid(dimg,updown)
		if updown > 0 :
			cv2.imshow('Up',dimg)
		else:
			cv2.imshow('Down',dimg)
		k = cv2.waitKey(0)
		cv2.destroyAllWindows()
		if k not in [UPKEY,DOWNKEY]:
			break
		if k == UPKEY:
			updown = 1
		elif k == DOWNKEY :
			updown = -1
	return

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	PyrShow(sys.argv[1])

if __name__ == '__main__':
	main()



