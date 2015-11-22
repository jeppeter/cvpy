import sys
import cv2
import numpy as np
import matplotlib.pyplot as plt
import logging


############

def ShowFFT(infile):
	try:
		simg = cv2.imread(infile,0)
		assert(len(simg) >= 2)
	except:
		sys.stderr.write('can not open (%s) for picture\n'%(infile))
		return
	f = np.fft.fft2(simg)
	fshift = np.fft.fftshift(f)
	fftsimg = np.log(np.abs(fshift))
	# to give inverse transform
	ifshift = np.fft.ifftshift(fshift)
	dimg = np.fft.ifft2(ifshift)
	dimg = np.abs(dimg)
	cv2.imshow('origin',simg)
	cv2.imshow('fft',fftsimg)
	cv2.imshow('inverse back',dimg)
	cv2.waitKey(0)
	return


def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	logging.basicConfig(level=logging.DEBUG,format="%(levelname)-8s [%(filename)-10s:%(funcName)-20s:%(lineno)-5s] %(message)s")
	ShowFFT(sys.argv[1])

if __name__ == '__main__':
	main()
