
import sys
import os
import cv2
import numpy as np
import math
import argparse
from matplotlib import pyplot as plt

UPKEY=2490368
DOWNKEY=2621440
LEFTKEY=2424832
RIGHTKEY=2555904


def rotate_about_center(src, angle, scale=1.):
    w = src.shape[1]
    h = src.shape[0]
    rangle = np.deg2rad(angle)  # angle in radians
    # now calculate new image width and height
    nw = (abs(np.sin(rangle)*h) + abs(np.cos(rangle)*w))*scale
    nh = (abs(np.cos(rangle)*h) + abs(np.sin(rangle)*w))*scale
    # ask OpenCV for the rotation matrix
    rot_mat = cv2.getRotationMatrix2D((nw*0.5, nh*0.5), angle, scale)
    # calculate the move from the old center to the new center combined
    # with the rotation
    rot_move = np.dot(rot_mat, np.array([(nw-w)*0.5, (nh-h)*0.5,0]))
    # the move only affects the translation, so update the translation
    # part of the transform
    rot_mat[0,2] += rot_move[0]
    rot_mat[1,2] += rot_move[1]
    return cv2.warpAffine(src, rot_mat, (int(math.ceil(nw)), int(math.ceil(nh))), flags=cv2.INTER_LINEAR)

def add_value(pts,idx):
	c = idx % 2
	r = int((idx - c) /2)
	try:
		pts[r][c] += 1
	except:
		print 'pts[%d][%d] error'%(c,r)
		sys.exit(3)

	return pts

def sub_value(pts,idx):
	c = idx % 2
	r = int((idx - c)/2)
	try:
		pts[r][c] -= 1
	except:
		print 'pts[%d][%d] error'%(c,r)
		sys.exit(3)
	return pts

def add_affine_pts(pts1,pts2,curidx):
	if curidx < 8:
		pts1 = add_value(pts1,curidx)
	else:
		pts2 = add_value(pts2,curidx - 8)
	return pts1,pts2

def sub_affine_pts(pts1,pts2,curidx):
	if curidx < 8:
		pts1 = sub_value(pts1,curidx)
	else:
		pts2 = sub_value(pts2,curidx - 8)
	return pts1,pts2


def TransAffine(infile,args):
	simg = cv2.imread(infile)
	cols = simg.shape[0]
	rows = simg.shape[1]
	if args.rows > 0 :
		rows = args.rows
	if args.cols > 0:
		cols = args.cols
	pts1 = np.float32([[56,65],[368,52],[28,387],[389,390]])
	pts2 = np.float32([[0,0],[300,0],[0,300],[300,300]])
	print 'len(pts1) = %d len(pts2)  = %d'%(len(pts1),len(pts2))
	curidx = 0
	while True:
		M = cv2.getPerspectiveTransform(pts1,pts2)
		dimg = cv2.warpPerspective(simg,M,(rows,cols))
		cv2.imshow('img',dimg)
		k = cv2.waitKey(0)
		cv2.destroyAllWindows()
		if k not in [UPKEY,RIGHTKEY,LEFTKEY,DOWNKEY]:
			break
		if k == UPKEY :
			pts1,pts2 = add_affine_pts(pts1,pts2,curidx)
		elif k == DOWNKEY:
			pts1,pts2 = sub_affine_pts(pts1,pts2,curidx)
		elif k == LEFTKEY:
			curidx += 1
			curidx %= 16			
			print 'curidx %d'%(curidx)
			print 'pts1(%s) pts2(%s)'%(pts1,pts2)
		elif k == RIGHTKEY :
			curidx -= 1
			curidx %= 16
			print 'curidx %d'%(curidx)
			print 'pts1(%s) pts2(%s)'%(pts1,pts2)

	return


def main():
	parser = argparse.ArgumentParser(usage='%s [OPTIONS] infile'%(sys.argv[0]),description='Process picture affine')
	parser.add_argument('-c','--cols',type=int,nargs=1,dest='cols',default=0,help='specify the columns')
	parser.add_argument('-r','--rows',type=int,nargs=1,dest='rows',default=0,help='specify the rows')
	args,files= parser.parse_known_args()
	if len(files) < 1:
		sys.stderr.write('need infile \n')
		parser.print_help(sys.stderr)
		sys.exit(3)
	TransAffine(files[0],args)

if __name__ == '__main__':
	main()

