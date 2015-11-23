#include <stdio.h>
#include <maxflow/graph.h>

void PrintErrorMsg(const char* errmsg)
{
	fprintf(stderr,errmsg);
	return;
}

int main(int argc,char* argv[])
{
	typedef maxflowLib::Graph<int,int,int> GraphType;
	GraphType *g = new GraphType(/*estimated # of nodes*/ 1, /*estimated # of edges*/ 5,PrintErrorMsg); 

	g -> add_node(20);

	g -> add_edge(1,6,59,0);
	g -> add_edge(1,12,21,0);
	g -> add_edge(1,13,19,0);
	g -> add_edge(1,3,35,0);
	g -> add_edge(1,17,17,0);
	g -> add_edge(17,18,72,0);
	g -> add_edge(17,2,63,0);
	g -> add_edge(17,10,93,0);
	g -> add_edge(17,15,40,0);
	g -> add_edge(17,7,52,0);
	g -> add_edge(17,12,82,0);
	g -> add_edge(5,16,4,0);
	g -> add_edge(8,16,8,0);
	g -> add_edge(8,9,1,0);
	g -> add_edge(8,2,21,0);
	g -> add_edge(4,11,75,0);
	g -> add_edge(4,12,24,0);
	g -> add_edge(4,8,7,0);
	g -> add_edge(2,12,69,0);
	g -> add_edge(2,9,92,0);
	g -> add_edge(2,16,15,0);
	g -> add_edge(12,3,8,0);
	g -> add_edge(12,11,43,0);
	g -> add_edge(12,15,89,0);
	g -> add_edge(14,9,42,0);
	g -> add_edge(14,10,7,0);
	g -> add_edge(14,12,40,0);
	g -> add_edge(14,16,10,0);
	g -> add_edge(14,8,23,0);
	g -> add_edge(10,7,11,0);
	g -> add_edge(10,15,98,0);
	g -> add_edge(10,2,48,0);
	g -> add_edge(11,18,17,0);
	g -> add_edge(11,3,24,0);
	g -> add_edge(9,13,67,0);
	g -> add_edge(9,10,67,0);
	g -> add_edge(13,10,8,0);
	g -> add_edge(15,1,96,0);
	g -> add_edge(15,7,98,0);
	g -> add_edge(15,3,33,0);
	g -> add_edge(15,6,4,0);
	g -> add_edge(15,14,98,0);
	g -> add_edge(15,2,23,0);
	g -> add_edge(15,4,71,0);
	g -> add_edge(15,8,99,0);
	g -> add_edge(16,1,30,0);
	g -> add_edge(16,4,70,0);
	g -> add_edge(16,17,26,0);
	g -> add_edge(16,3,49,0);
	g -> add_edge(16,11,20,0);
	g -> add_edge(6,8,62,0);
	g -> add_edge(6,17,17,0);
	g -> add_edge(3,18,89,0);
	g -> add_edge(3,2,9,0);
	g -> add_edge(3,10,18,0);
	g -> add_edge(3,14,21,0);
	g -> add_tweights(17,0,66);
	g -> add_tweights(18,0,23);
	g -> add_tweights(5,0,33);
	g -> add_tweights(9,28,0);




	int flow = g -> maxflow();

	printf("Flow = %d\n", flow);
	printf("Minimum cut:\n");
	if (g->what_segment(0) == GraphType::SOURCE)
		printf("node0 is in the SOURCE set\n");
	else
		printf("node0 is in the SINK set\n");
	if (g->what_segment(1) == GraphType::SOURCE)
		printf("node1 is in the SOURCE set\n");
	else
		printf("node1 is in the SINK set\n");

	delete g;

	return 0;
}