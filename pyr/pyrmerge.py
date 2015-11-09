import sys
import cv2
import numpy as np
import logging

UPKEY=2490368
DOWNKEY=2621440


def key_usage():
	sys.stdout.write('UP to inc pyrup times\n')
	sys.stdout.write('DOWN to dec pyrup times\n')
	return


def pyr_merge(aimg,bimg,pytimes):
	G = aimg.copy()
	gpA = [G]
	for i in xrange(pytimes):
		G = cv2.pyrDown(G)
		gpA.append(G)
	G = bimg.copy()
	gpB = [G]
	for i in xrange(pytimes):
		G = cv2.pyrDown(G)
		gpB.append(G)

	lpA = [gpA[pytimes-1]]
	for i in xrange(pytimes-1,0,-1):
		GE = cv2.pyrUp(gpA[i])
		L = cv2.subtract(gpA[i-1],GE)
		lpA.append(L)

	lpB = [gpB[pytimes-1]]
	for i in xrange(pytimes - 1,0,-1):
		GE = cv2.pyrUp(gpB[i])
		L = cv2.subtract(gpB[i-1],GE)
		lpB.append(L)

	LS = []
	i = 0
	for la,lb in zip(lpA,lpB):
		rows,cols,dpt = la.shape
		ls = np.hstack((la[:,0:cols/2], lb[:,cols/2:]))
		LS.append(ls)
		i += 1
	ls_ = LS[0]
	for i in xrange(1,pytimes):
		ls_ = cv2.pyrUp(ls_)
		ls_ = cv2.add(ls_,LS[i])
	return ls_

def MergeFilePyramid(afile,bfile):
	try:
		aimg = cv2.imread(afile)
		bimg = cv2.imread(bfile)
		assert(len(aimg.shape)>=3)
		assert(len(bimg.shape)>=3)
	except:
		sys.stderr.write('can not open (%s/%s) for file read\n'%(afile,bfile))
		return -3
	# now for make the whole size
	key_usage()
	arows,acols,achl = aimg.shape
	brows,bcols,bchl = bimg.shape
	mrows = arows
	if mrows < brows:
		mrows = brows
	mcols = acols
	if mcols < bcols:
		mcols = bcols

	aimg = cv2.resize(aimg,(mcols,mrows))
	bimg = cv2.resize(bimg,(mcols,mrows))
	pyrtimes = 3
	maxtimes = 6
	while True:
		mimg = pyr_merge(aimg,bimg,pyrtimes)
		mname = '%d times'%(pyrtimes)
		cv2.imshow(mname,mimg)
		k = cv2.waitKey(0)
		cv2.destroyAllWindows()
		if k not in [UPKEY,DOWNKEY]:
			break
		if k == UPKEY:
			pyrtimes += 1
			if pyrtimes >maxtimes :
				pyrtimes = maxtimes
		else:
			pyrtimes -= 1
			if pyrtimes < 1 :
				pyrtimes = 1
	return 0

def main():
	if len(sys.argv) < 3:
		sys.stderr.write('%s afile bfile\n'%(sys.argv[0]))
		sys.exit(3)
	logging.basicConfig(level=logging.DEBUG,format="%(levelname)-8s [%(filename)-10s:%(funcName)-20s:%(lineno)-5s] %(message)s")
	MergeFilePyramid(sys.argv[1],sys.argv[2])

if __name__ == '__main__':
	main()





