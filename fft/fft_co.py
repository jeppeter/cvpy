#!python
# coding=utf-8
import cv2
import numpy as np
import matplotlib.pyplot as plt
import sys
import logging
import os

def OutMatrix(name,mat):
	logging.info('out %s [%d:%d]'%(name,mat.shape[0],mat.shape[1]))
	r,c=mat.shape[:2]
	for y in range(0,r):
		s = ''
		for x in range(0,c):
			s += ' %s '%(mat[y][x])
		logging.info('[%d] %s'%(y,s))
	return


def CoBine(afile,bfile):
	afile_base=os.path.basename(afile)
	bfile_base=os.path.basename(bfile)
	img_flower = cv2.imread(afile,0) #直接读为灰度图像
	img_man = cv2.imread(bfile,0) #直接读为灰度图像
	arows,acols = img_flower.shape[:2]
	brows,bcols = img_man.shape[:2]
	mrows = arows
	if mrows < brows :
		mrows = brows
	mcols = acols
	if mcols < bcols :
		mcols = bcols
	img_flower = cv2.resize(img_flower,(mrows,mcols), interpolation = cv2.INTER_CUBIC)
	img_man = cv2.resize(img_man,(mrows,mcols), interpolation = cv2.INTER_CUBIC)
	outname ='%s gray'%(afile)
	cv2.imshow(outname,img_flower)
	outname = '%s gray'%(bfile)
	cv2.imshow(outname,img_man)
	#--------------------------------
	f1 = np.fft.fft2(img_flower)
	f1shift = np.fft.fftshift(f1)
	f1_A = np.abs(f1shift) #取振幅
	f1_P = np.angle(f1shift) #取相位
	#--------------------------------
	f2 = np.fft.fft2(img_man)
	f2shift = np.fft.fftshift(f2)
	f2_A = np.abs(f2shift) #取振幅
	f2_P = np.angle(f2shift) #取相位
	#---图1的振幅--图2的相位--------------------
	img_new1_f = np.zeros(img_flower.shape,dtype=complex) 
	img1_real = f1_A*np.cos(f2_P) #取实部
	img1_imag = f1_A*np.sin(f2_P) #取虚部
	img_new1_f.real = np.array(img1_real) 
	img_new1_f.imag = np.array(img1_imag) 
	f3shift = np.fft.ifftshift(img_new1_f) #对新的进行逆变换
	img_new1 = np.fft.ifft2(f3shift)
	#出来的是复数，无法显示
	img_new1 = np.abs(img_new1)
	img_new1 = np.uint8(img_new1)
	#调整大小范围便于显示
	outname='%s-phase-%s-between'%(afile,bfile)	
	# to just show gray mode
	cv2.imshow(outname,img_new1)
	cv2.imwrite('%s_%s.bmp'%(afile_base,bfile_base),img_new1)
	#---图2的振幅--图1的相位--------------------
	img_new2_f = np.zeros(img_flower.shape,dtype=complex) 
	img2_real = f2_A*np.cos(f1_P) #取实部
	img2_imag = f2_A*np.sin(f1_P) #取虚部
	img_new2_f.real = np.array(img2_real) 
	img_new2_f.imag = np.array(img2_imag) 
	f4shift = np.fft.ifftshift(img_new2_f) #对新的进行逆变换
	img_new2 = np.fft.ifft2(f4shift)
	#出来的是复数，无法显示
	img_new2 = np.abs(img_new2)
	img_new2 = np.uint8(img_new2)
	#调整大小范围便于显示
	outname = '%s-phase-%s-between'%(bfile,afile)
	cv2.imshow(outname,img_new2)
	cv2.imwrite('%s_%s.bmp'%(bfile_base,afile_base),img_new2)
	cv2.waitKey(0)
	return

def main():
	if len(sys.argv) < 3:
		sys.stderr.write('%s afile bfile to fft combine\n'%(sys.argv[0]))
		sys.exit(3)
	logging.basicConfig(level=logging.DEBUG,format="%(levelname)-8s [%(filename)-10s:%(funcName)-20s:%(lineno)-5s] %(message)s")
	CoBine(sys.argv[1],sys.argv[2])

if __name__ == '__main__':
	main()

