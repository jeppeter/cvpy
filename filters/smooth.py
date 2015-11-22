import cv2
import numpy as np
import sys


UPKEY=2490368
DOWNKEY=2621440
LEFTKEY=2424832
RIGHTKEY=2555904
AKEY=97
SKEY=115
GAUSSIANBLUR=1
MEDIABLUR=2
COMMONBLUR=3

SMOOTHTYPE=[GAUSSIANBLUR,MEDIABLUR,COMMONBLUR]
SMOOTHNAME=['GAUSSIANBLUR','MEDIABLUR','COMMONBLUR']

class BlurArgs:
	def __init__(self):
		self.num = 5
		self.flt = 1.5
		self.typesmooth = SMOOTHTYPE[0]
		self.typeint = 0
		return 

	def incnum(self):
		self.num += 1
		return

	def decnum(self):
		self.num -= 1
		if self.num < 0:
			self.num = 0
		return 

	def incflt(self):
		self.flt *= 1.1
		return

	def decflt(self):
		self.flt *= 0.9
		return

	def inctype(self):
		self.typeint += 1
		self.typeint %= len(SMOOTHTYPE)
		self.typesmooth = SMOOTHTYPE[self.typeint]
		return

	def dectype(self):
		self.typeint -= 1
		self.typeint %= len(SMOOTHTYPE)
		self.typesmooth = SMOOTHTYPE[self.typeint]
		return


	def name(self):
		typestr = SMOOTHNAME[self.typeint]
		return 'num %d type %s flt %f'%(self.num,typestr,self.flt)





def key_usage():
	sys.stdout.write('UP to increase \n')
	sys.stdout.write('DOWN to decrease x weight\n')
	sys.stdout.write('LEFT  RIGHT to change y weight\n')
	sys.stdout.write('AKEY SKEY to change blur type\n')
	return


def FilterSmooth(img,args):
	if args.typesmooth == COMMONBLUR:
		dst = cv2.blur(img,(args.num,args.num))
	elif args.typesmooth == MEDIABLUR:
		dst = cv2.medianBlur(img,args.num)
	elif args.typesmooth == GAUSSIANBLUR:
		num = args.num
		if (num % 2) == 0:
			num += 1
		dst = cv2.GaussianBlur(img,(num,num),args.flt)
	return dst

def SaltImage(img,n):
    dst = img
    for k in range(n) :
        i = int(np.random.random() * dst.shape[1]);    
        j = int(np.random.random() * dst.shape[0]);    
        if dst.ndim == 2:     
            dst[j,i] = 255    
        elif dst.ndim == 3:     
            dst[j,i,0]= 255    
            dst[j,i,1]= 255    
            dst[j,i,2]= 255    
    return dst

def key_usage():
	sys.stdout.write('UP to increase num\n')
	sys.stdout.write('DOWN to decrease num\n')
	sys.stdout.write('LEFT to increase flt\n')
	sys.stdout.write('RIGHT to decrease flt\n')
	sys.stdout.write('AKEY SKEY to change blur type\n')
	return

def FilterSmootShow(infile):
	try:
		img = cv2.imread(infile,0)
		assert(len(img.shape) >=2)
	except:
		sys.stderr.write('can not open %s as picture\n'%(infile))
		sys.exit(3)

	args = BlurArgs()
	key_usage()
	saltimg = SaltImage(img,1000)
	while True:
		dst = FilterSmooth(img,args)
		name = args.name()
		cv2.imshow(name,dst)
		cv2.imshow('salt image',saltimg)
		k = cv2.waitKey(0)
		cv2.destroyAllWindows()
		if k not in [UPKEY,RIGHTKEY,LEFTKEY,DOWNKEY,AKEY,SKEY]:
			sys.stdout.write('k %d'%(k))
			break
		if k == UPKEY :
			args.incnum()
		elif k == DOWNKEY :
			args.decnum()
		elif k == LEFTKEY:
			args.incflt()
		elif k == RIGHTKEY:
			args.decflt()
		elif k == AKEY:
			args.inctype()
		elif k == SKEY:
			args.dectype()
	return

def main():
	if len(sys.argv) < 2:
		sys.stderr.write('%s infile to erode'%(sys.argv[0]))
		sys.exit(3)
	FilterSmootShow(sys.argv[1])

if __name__ == '__main__':
	main()


