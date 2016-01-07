#!python


import sys
import logging
import cv2

def DisplayImage(infile):
	try:
		simg = cv2.imread(infile)
		assert(len(simg.shape) >= 2)
	except:
		sys.stderr.write('can not load(%s) image\n'%(infile))
		return
	h = simg.shape[0]
	w = simg.shape[1]
	logging.info('w %d h %d'%(w,h))
	cv2.imshow(infile,simg)
	cv2.waitKey(0)
	return

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile\n'%(sys.argv[0]))
		sys.exit(4)
	logging.basicConfig(level=logging.DEBUG,format='%(filename)s:%(lineno)d %(message)s')
	DisplayImage(sys.argv[1])
	return

if __name__ == '__main__':
	main()