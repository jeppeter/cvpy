#include <opencv2/opencv.hpp>

using namespace cv;
using namespace std;

int main(int argc,char* argv[])
{
	Mat simg ;

	if (argc < 2){
		cerr << argv[0] << " infile" << endl;
		return -3;
	}

	simg = imread(argv[1],CV_LOAD_IMAGE_COLOR);
	if (simg.data  == NULL){
		cerr << "can not load " << argv[1] << endl;
		return -3;
	}

	namedWindow(argv[1],WINDOW_AUTOSIZE);
	imshow(argv[1],simg);
	waitKey(0);
	return 0;
}