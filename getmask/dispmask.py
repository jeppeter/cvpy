#!/usr/bin/python

import cv2
import sys


def DisplayMatrix(infile):
	try:
		simg = cv2.imread(infile)
		assert(len(simg.shape)>=2)
	except:
		sys.stderr.write('can not open(%s) for image\n'%(infile))
		return

	
	h = simg.shape[0]
	w = simg.shape[1]

	for i in range(h):
		s = '['
		cnt = 0
		for j in range(w):
			if cnt != 0 :
				s += ','
			s += '%s'%(simg[i][j][1])
			cnt += 1
		s += ']\n'
		sys.stdout.write(s)

	return

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s picture\n'%(sys.argv[0]))
		return
	DisplayMatrix(sys.argv[1])

if __name__ == '__main__':
	main()