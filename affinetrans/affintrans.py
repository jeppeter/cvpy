
import sys
import os
import cv2
import numpy as np
import math
import argparse

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

def TransAffine(infile,outfile,args):
	simg = cv2.imread(infile,1)	
	if args.rows > 0 :
		rows = args.rows
	if args.cols > 0:
		cols = args.cols
	angle = 15
	scale = 1
	while True:
		dimg = rotate_about_center(simg,angle,scale)	
		cv2.imshow('img',dimg)
		k = cv2.waitKey(0)
		cv2.destroyAllWindows()
		print 'key %s'%(k)
		if k not in [UPKEY,RIGHTKEY,LEFTKEY,DOWNKEY]:
			break
		if k == UPKEY :
			angle += 1
		elif k == DOWNKEY:
			angle -= 1
		elif k == LEFTKEY:
			scale *= 0.9
		elif k == RIGHTKEY :
			scale *= 1.1

	return


def main():
	parser = argparse.ArgumentParser(usage='%s [OPTIONS] infile outfile'%(sys.argv[0]),description='Process picture affine')
	parser.add_argument('-c','--cols',type=int,nargs=1,dest='cols',default=0,help='specify the columns')
	parser.add_argument('-r','--rows',type=int,nargs=1,dest='rows',default=0,help='specify the rows')
	args,files= parser.parse_known_args()
	if len(files) < 2:
		sys.stderr.write('need infile and outfile\n')
		parser.print_help(sys.stderr)
		sys.exit(3)
	TransAffine(files[0],files[1],args)

if __name__ == '__main__':
	main()

