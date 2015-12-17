#ifndef __OUT_DEBUG_H__
#define __OUT_DEBUG_H__

#define DEBUG_OUT(...) do{fprintf(stdout,"%s:%s:%d\t",__FILE__,__FUNCTION__,__LINE__);fprintf(stdout,__VA_ARGS__);}while(0)
#define DEBUG_BUFFER_FMT(p,l,...) \
do\
{\
	unsigned char* __pcur = (unsigned char*)p;\
	int __leftlen = (int)l;\
	int __i=0;\
	fprintf(stdout,"%s:%d\tpointer(0x%p)size(%d)\t",__FILE__,__LINE__,__pcur,__leftlen);\
	fprintf(stdout,__VA_ARGS__);\
	for (__i=0;__i < __leftlen;__i ++)\
	{\
		if ((__i % 16) == 0)\
		{\
			fprintf(stdout,"\n[0x%08x]\t",__i);\
		}\
		fprintf(stdout," 0x%02x",*__pcur);\
		__pcur ++;\
	}\
	fprintf(stdout,"\n");\
}while(0)

#endif /*__OUT_DEBUG_H__*/