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

def SetDictDefValue(dictarr,k1,k2,v):
    if k1 not in dictarr.keys():
        dictarr[k1] = {}
    if k2 not in dictarr[k1].keys():
        dictarr[k1][k2] = v
    return dictarr

def SetDicArrValue(dictarr,k1,v):
    if k1 not in dictarr.keys():
        dictarr[k1] = v
    return dictarr[k1]

def EdmondsKarp(capacity, neighbors, start, end):
    flow = 0
    flows = {}
    maxval = 0
    for k in capacity.keys():
        for j in capacity.keys():
            flows = SetDictDefValue(flows,k,j,0)
            capacity = SetDictDefValue(capacity,k,j,0)
            maxval += capacity[k][j]
    while True:
        max, parent = BFS(capacity, neighbors, flows, start, end,maxval)
        if max == 0:
            break
        flow = flow + max
        v = end
        while v != start:
            u = parent[v]
            flows[u][v] = flows[u][v] + max
            flows[v][u] = flows[v][u] - max
            v = u
    return (flow, flows)


def BFS(capacity, neighbors, flows, start, end,maxval):
    length = len(capacity)    
    parents = {}
    curmax = 0
    M = {}
    for k in capacity.keys():
        parents[k] = -1
        M[k] = 0
    parents[start] = -2
    M[start] = maxval
    queue = []
    queue.append(start)
    while queue:
        u = queue.pop(0)
        for v in SetDicArrValue(neighbors,u,[]):
            # if there is available capacity and v is is not seen before in search
            if capacity[u][v] - flows[u][v] > 0 and parents[v] == -1:
                parents[v] = u
                # it will work because at the beginning M[u] is Infinity
                M[v] = min(M[u], capacity[u][v] - flows[u][v]) # try to get smallest
                if v != end:
                    queue.append(v)
                else:
                    return M[end], parents
    return 0, parents



def main():
    if len(sys.argv) < 2:
        sys.stderr.write('%s inputfile\n'%(sys.argv[0]))
        sys.exit(4)
    network = FlowNetwork()
    logging.basicConfig(level=logging.INFO,format='%(asctime)-15s:%(filename)s:%(lineno)d\t%(message)s')
    sink = ''
    source = ''
    with open(sys.argv[1],'r') as f:
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
        sys.stderr.write('please specify sink= or source= in %s file\n'%(sys.argv[1]))
        sys.exit(4)

    #print network.max_flow(source,sink)
    cap,neighbor = network.get_cap_neighbour()
    flow,flows = EdmondsKarp(cap,neighbor,source,sink)
    logging.info('flow %d flows %s'%(flow,flows))

if __name__ == '__main__':
    main()