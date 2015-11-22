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
	GraphType *g = new GraphType(/*estimated # of nodes*/ 4, /*estimated # of edges*/ 10,PrintErrorMsg); 

	g -> add_node(); 
	g -> add_node(); 
	g -> add_node();
	g -> add_node();

	g -> add_tweights( 0,   /* capacities */  1, 2 );
	g -> add_tweights( 1,   /* capacities */  1, 20 );
	g -> add_tweights(2,1,300);
	g -> add_tweights(3,1,400);
	g -> add_edge( 0, 1,    /* capacities */  0, 0 );

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