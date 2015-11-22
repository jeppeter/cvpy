#!python
import sys
import logging

class Edge(object):
    def __init__(self, u, v, w):
        self.source = u
        self.sink = v  
        self.capacity = w
    def __repr__(self):
        return "%s->%s:%s" % (self.source, self.sink, self.capacity)

class FlowNetwork(object):
    def __init__(self):
        self.adj = {}
        self.flow = {}
 
    def add_vertex(self, vertex):
        if vertex in self.adj.keys():
            return
        self.adj[vertex] = []
 
    def get_edges(self, v):
        return self.adj[v]
 
    def add_edge(self, u, v, w=0):
        if u == v:
            raise ValueError("u == v")
        self.add_vertex(u)
        self.add_vertex(v)
        edge = Edge(u,v,w)
        redge = Edge(v,u,0)
        edge.redge = redge
        redge.redge = edge
        self.adj[u].append(edge)
        self.adj[v].append(redge)
        self.flow[edge] = 0
        self.flow[redge] = 0
 
    def find_path(self, source, sink, path):
        if source == sink:
            logging.info('source == sink %s'%(source))
            return path
        for edge in self.get_edges(source):
            logging.info('get edge %s cur source %s sink %s path %s'%(edge,source,sink,path))
            residual = edge.capacity - self.flow[edge]
            if residual > 0 and edge not in path:
                result = self.find_path( edge.sink, sink, path + [edge]) 
                if result != None:
                    logging.info('result %s'%(result))
                    return result
 
    def max_flow(self, source, sink):
        path = self.find_path(source, sink, [])
        while path != None:
            residuals = [edge.capacity - self.flow[edge] for edge in path]
            flow = min(residuals)
            logging.info('path %s find residual %s flow %d'%(path,residuals,flow))
            for edge in path:
                self.flow[edge] += flow
                self.flow[edge.redge] -= flow
            path = self.find_path(source, sink, [])
        return sum(self.flow[edge] for edge in self.get_edges(source))

    def get_cap_neighbour(self):
        sortkeys = self.adj.keys()
        cap = {}
        neighbor = {}
        for k in sortkeys:
            for edge in self.adj[k]:
                if k not in cap.keys():
                    cap[k] = {}
                if k not in neighbor.keys():
                    neighbor[k] = []
                if edge.sink not in neighbor.keys():
                    neighbor[edge.sink] = []
                cap[k][edge.sink] = edge.capacity
                neighbor[k].append(edge.sink)
                neighbor[edge.sink].append(k)
        for k in neighbor.keys():
            curnei=[]
            for i in neighbor[k]:
                if i not in curnei:
                    curnei.append(i)
            neighbor[k] = curnei
        return cap ,neighbor


'''
  we use this file to make the max flow min cut algorithm of goldberg tarjan
'''

def UniqueSortArray(neighbours):
	sortarray = []
	for k in neighbours.keys():
		if k not in sortarray:
			sortarray.append(k)
		v = neighbours[k]
		for k2 in v:
			if k2 not in sortarray:
				sortarray.append(k2)
	# now sort the array
	i = 0
	while i < len(sortarray) :
		j = i + 1
		while j < len(sortarray):
			if sortarray[i] > sortarray[j]:
				tmp =sortarray[i]
				sortarray[i] = sortarray[j]
				sortarray[j] = tmp
			j = j + 1
		i = i + 1
	return sortarray

def SetNotUsedValue(capcity,sortarray):
	for k in sortarray:
		if k not in capcity.keys():
			capcity[k] = {}
		for k2 in sortarray:
			if k2 not in capcity[k].keys():
				capcity[k][k2] = 0
	return capcity

def SetNotUsedArr(overflow,sortarray):
	for k in sortarray:
		if k not in overflow.keys():
			overflow[k] = 0
	return overflow


def CanPush(n,neigbours ,nextnodes,capcity,flows):
	for k in neigbours[n]:
		if ((nextnodes[k] + 1 ) == nextnodes[n] ) and \
			(capcity[n][k] - flows[n][k]) > 0 :
			return True
	return False

def SetNextNodes(n,neighbours,nextnodes,capcity,flows,maxval):
	minval = maxval
	for k in neighbours[n]:
		if capcity[n][k] - flows[n][k] > 0:
			minval = min(minval,nextnodes[k])
	nextnodes[n] = 1 + minval
	return nextnodes

def FindMaxValueInNextNodes(nextnodes):
	maxval = 0
	for k in nextnodes.keys():
		if nextnodes[k] > maxval :
			maxval = nextnodes[k]
	return maxval

def FindNextNode(n,neighbours,nextnodes,capacity,flows,overflow):
	for k in neighbours[n]:
		if nextnodes[k] + 1 == nextnodes[n]:
			fval = (capacity[n][k] - flows[n][k])
			if fval > overflow[n]:
				fval = overflow[n]
			overflow[k] += fval
			overflow[n] -= fval
			flows[n][k] += fval
			flows[k][n] -= fval

	return flows,overflow


def GoldbergTarjan(capcity,neighbours,source,sink):
	# first to set for not used value
	sortarray = UniqueSortArray(neighbours)
	capcity = SetNotUsedValue(capcity,sortarray)
	flows = {}
	flows = SetNotUsedValue(flows,sortarray)
	overflow = {}
	overflow = SetNotUsedArr(overflow,sortarray)
	nextnodes = {}
	nextnodes = SetNotUsedArr(nextnodes,sortarray)
	active_nodes = set([])
	nextnodes[source] = len(sortarray)
	for n in neighbours[source]:
		flows[source][n] = capcity[source][n]
		flows[n][source] = - capcity[source][n]
		overflow[n] = capcity[source][n]
		active_nodes.add(n)

	while len(active_nodes) > 0:
		maxval = FindMaxValueInNextNodes(nextnodes)
		n = active_nodes.pop()
		if not CanPush(n,neighbours,nextnodes,capcity,flows):
			nextnodes = SetNextNodes(n,neighbours,nextnodes,capcity,flows,maxval)
		flows,overflow = FindNextNode(n,neighbours,nextnodes,capcity,flows,overflow)

		if n != source and n != sink and overflow[n] > 0:
			active_nodes.add(n)
		for k in neighbours[n]:
			if k != source and k != sink and overflow[k] > 0 :
				active_nodes.add(k)

	sumval = 0
	for k in neighbours[source]:
		sumval += flows[source][k]
		sumval -= flows[k][source]
	sumval = sumval / 2
	return sumval , flows

def ParseAndGetValue(infile):
    sink = ''
    source = ''
    network = FlowNetwork()
    with open(infile,'r') as f:
        for l in f:
            if l.startswith('#'):
                continue
            l = l.rstrip('\r\n')
            if l.startswith('source='):
                sarr = l.split('=')
                if len(sarr) < 2:
                    continue
                source = sarr[1]
                continue
            if l.startswith('sink='):
                sarr = l.split('=')
                if len(sarr) < 2:
                    continue
                sink = sarr[1]
                continue
            sarr = l.split(',')
            if len(sarr) < 3:
                continue
            network.add_edge(sarr[0],sarr[1],int(sarr[2]))
    if sink == '' or source == '':
        sys.stderr.write('please specify sink= or source= in %s file\n'%(infile))
        sys.exit(4)

    #print network.max_flow(source,sink)
    cap,neighbor = network.get_cap_neighbour()
    return cap,neighbor,source,sink


def main():
    if len(sys.argv) < 2:
        sys.stderr.write('%s inputfile\n'%(sys.argv[0]))
        sys.exit(4)
    logging.basicConfig(level=logging.INFO,format='%(asctime)-15s:%(filename)s:%(lineno)d\t%(message)s')
    cap,neighbor,source,sink = ParseAndGetValue(sys.argv[1])
    flow,flows = GoldbergTarjan(cap,neighbor,source,sink)
    logging.info('flow %d flows %s'%(flow,flows))

if __name__ == '__main__':
    main()





