#!python
# coding=utf-8
import cv2
import numpy as np
import matplotlib.pyplot as plt
import sys
import logging

def ShowFFT2(infile):
	try:
		simg = cv2.imread(infile,0)
		assert(len(simg.shape)>=2)
	except:
		sys.stderr.write('can not open(%s) for picture\n'%(infile))
		return
	f = np.fft.fft2(simg)
	fshift = np.fft.fftshift(f)
	s1 = np.log(np.abs(fshift))
	outfile='%s.gray2.bmp'%(infile)
	cv2.imwrite(outfile,simg)
	cv2.imshow('gray',simg)

	# now get inverse trans
	ifshift = np.fft.ifftshift(np.abs(fshift))
	isimg = np.fft.ifft2(ifshift)
	isimg = np.abs(isimg)
	outfile='%s.inverse2.bmp'%(infile)
	cv2.imwrite(outfile,isimg)
	cv2.imshow('invers',isimg)

	# get phase
	if2shift = np.fft.ifftshift(np.angle(fshift))
	is2img = np.fft.ifft2(if2shift)
	is2img = np.abs(is2img)
	outfile='%s.phase2.bmp'%(infile)
	cv2.imwrite(outfile,is2img)
	cv2.imshow('only phase',is2img)

	# combine inserve and phase
	s2 = np.abs(fshift)
	s2_angle = np.angle(fshift)
	s2_real = s2 * np.cos(s2_angle)
	s2_imag = s2 * np.sin(s2_angle)
	s3 = np.zeros(simg.shape,dtype=complex)
	s3.real = np.array(s2_real)
	s3.imag = np.array(s2_imag)

	if3shift = np.fft.ifftshift(s3)
	is3img = np.fft.ifft2(if3shift)
	is3img = np.abs(is3img)
	outfile='%s.combine2.bmp'%(infile)
	cv2.imwrite(outfile,is3img)
	cv2.imshow('reverse back',is3img)
	cv2.waitKey(0)
	return

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	logging.basicConfig(level=logging.DEBUG,format="%(levelname)-8s [%(filename)-10s:%(funcName)-20s:%(lineno)-5s] %(message)s")
	ShowFFT2(sys.argv[1])

if __name__ == '__main__':
	main()


