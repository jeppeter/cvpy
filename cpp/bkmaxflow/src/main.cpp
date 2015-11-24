#include <stdio.h>
#include <maxflow/graph.h>
#include <string>
#include <iostream>
#include <fstream>

typedef maxflowLib::Graph<int, int, int> GraphType;
using namespace std;

void PrintErrorMsg(const char* errmsg)
{
	fprintf(stderr, errmsg);
	return;
}

int GetEdge(char* line,int& s,int& d,int& w)
{
	char* pcur=line;
	s = atoi(pcur);
	while(*pcur != ',' && *pcur != '\r' && *pcur != 0x0 
		&& *pcur != '\n')
	{
		pcur ++;
	}

	if (*pcur != ',')
	{
		return -1;
	}
	pcur ++;
	d = atoi(pcur);
	while(*pcur != ',' && *pcur != '\r' && *pcur != 0x0 
		&& *pcur != '\n')
	{
		pcur ++;
	}

	if (*pcur != ',')
	{
		return -1;
	}
	pcur ++;
	w = atoi(pcur);
	return 3;
}

typedef struct source_pair
{
	int m_setting;
	int m_srcweight;
	int m_dstweight;
} source_pair_t, *psource_pair_t;

int ParseGraph(GraphType* bkgraph, const char* infile)
{
	std::ifstream inf;
	char line[1024];
	int ret;
	int source=-1,sink=-1;
	int s,d,w;
	psource_pair_t psourcepair=NULL;
	int i;

	inf.open(infile,std::ifstream::in);
	if (!inf.is_open())
	{
		ret = -1;
		fprintf(stderr,"can not open (%s)\n",infile);
		goto out;
	}

	while (!inf.eof())
	{
		memset(line,0,1024);
		inf.getline(line,1024);
		if (strncmp(line,"#",1)==0)
		{
			continue;
		}
		else if (strncmp(line,"source=",7) == 0)
		{
			source=atoi(line+7);
			continue;
		}
		else if (strncmp(line,"sink=",5)==0)
		{
			sink=atoi(line+5);
			if (sink < source)
			{
				ret = -1;
				fprintf(stderr,"sink[%d]<source[%d]\n",sink,source);
				goto out;
			}
			bkgraph->add_node(sink);
			psourcepair = (psource_pair_t)malloc(sizeof(*psourcepair)*(sink));
			if (psourcepair == NULL)
			{
				fprintf(stderr,"can not allocate (%d) size\n",sizeof(*psourcepair)*sink);
				ret = -1;
				goto out;
			}
			for (i=0;i<(sink);i++)
			{
				psourcepair[i].m_setting = 0;
				psourcepair[i].m_srcweight = 0;
				psourcepair[i].m_dstweight = 0;
			}

			continue;
		}

		ret = GetEdge(line,s,d,w);
		if (ret < 3)
		{
			continue;
		}

		if (source < 0 || sink <= 0  || source == sink)
		{
			ret = -1;
			fprintf(stderr,"(%s)not define source or sink\n",infile);
			goto out;
		}

		if (s > sink || d > sink)
		{
			fprintf(stderr,"s[%d] or d[%d] > sink[%d]\n",s,d,sink);
			ret = -1;
			goto out;
		}

		if (s == source && d == sink)
		{
			/*now we get the*/
			fprintf(stderr,"s[%d] == source d[%d]== sink\n",s,d);
			ret = -1;
			goto out;
		}

		if (s == source)
		{
			psourcepair[d].m_setting = 1;
			psourcepair[d].m_srcweight = w;
			if (psourcepair[d].m_dstweight != 0)
			{
				fprintf(stderr,"g -> add_tweights(%d,%d,%d);\n",d,psourcepair[d].m_srcweight,psourcepair[d].m_dstweight);
				bkgraph -> add_tweights(d,psourcepair[d].m_srcweight,psourcepair[d].m_dstweight);
				psourcepair[d].m_setting = 0;
				psourcepair[d].m_srcweight = 0;
				psourcepair[d].m_dstweight = 0;
			}
			continue;
		}

		if (d == sink)
		{
			psourcepair[s].m_setting = 1;
			psourcepair[s].m_dstweight = w;
			if (psourcepair[s].m_srcweight != 0)
			{
				fprintf(stderr,"g -> add_tweights(%d,%d,%d);\n",s,psourcepair[s].m_srcweight,psourcepair[s].m_dstweight);
				bkgraph -> add_tweights(s,psourcepair[s].m_srcweight,psourcepair[s].m_dstweight);
				psourcepair[s].m_setting = 0;
				psourcepair[s].m_srcweight = 0;
				psourcepair[s].m_dstweight = 0;
			}
			continue;
		}

		fprintf(stderr,"g -> add_edge(%d,%d,%d,0);\n",s,d,w);
		bkgraph -> add_edge(s,d,w,0);
	}

	for (i = 0;i < sink ;i ++)
	{
		if (psourcepair[i].m_setting)
		{
			bkgraph -> add_tweights(i,psourcepair[i].m_srcweight,psourcepair[i].m_dstweight);
		}
	}

	ret = 0;
out:
	if (psourcepair)
	{
		free(psourcepair);
	}
	psourcepair = NULL;

	if (inf.is_open())
	{
		inf.close();
	}
	return ret;
}

int main(int argc, char* argv[])
{
	int flow ;
	GraphType *g = NULL;
	int ret;
	if (argc < 2)
	{
		fprintf(stderr, "%s infile\n", argv[0]);
		return -1;
	}


	g = new GraphType(/*estimated # of nodes*/ 1, /*estimated # of edges*/ 5, PrintErrorMsg);

	ret = ParseGraph(g, argv[1]);
	if (ret < 0)
	{
		goto out;
	}
	flow = g -> maxflow();
	fprintf(stdout, "%d\n", flow);
	ret = 0;
#if 0
	printf("Minimum cut:\n");
	if (g->what_segment(0) == GraphType::SOURCE)
		printf("node0 is in the SOURCE set\n");
	else
		printf("node0 is in the SINK set\n");
	if (g->what_segment(1) == GraphType::SOURCE)
		printf("node1 is in the SOURCE set\n");
	else
		printf("node1 is in the SINK set\n");
#endif
out:
	delete g;
	return ret;
}