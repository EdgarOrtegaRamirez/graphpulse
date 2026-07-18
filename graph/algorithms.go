package graph

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
)

// BFSResult contains the result of a breadth-first search.
type BFSResult struct {
	Order   []string          // Visit order
	Visited map[string]bool   // Visited nodes
	Parent  map[string]string // Parent map for path reconstruction
	Dist    map[string]int    // Distance from source
}

// BFS performs breadth-first search from a source node.
func (g *Graph) BFS(source string) (*BFSResult, error) {
	if !g.HasNode(source) {
		return nil, ErrNodeNotFound
	}

	visited := make(map[string]bool)
	parent := make(map[string]string)
	dist := make(map[string]int)
	order := []string{}

	queue := []string{source}
	visited[source] = true
	dist[source] = 0

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)

		for _, edge := range g.Adj[node] {
			if !visited[edge.To] {
				visited[edge.To] = true
				parent[edge.To] = node
				dist[edge.To] = dist[node] + 1
				queue = append(queue, edge.To)
			}
		}
	}

	return &BFSResult{
		Order:   order,
		Visited: visited,
		Parent:  parent,
		Dist:    dist,
	}, nil
}

// DFSResult contains the result of a depth-first search.
type DFSResult struct {
	Order     []string          // Visit order (discovery order)
	Finished  []string          // Finish order (reverse topological for DAG)
	Visited   map[string]bool   // Visited nodes
	Parent    map[string]string // Parent map for path reconstruction
	Discovery map[string]int    // Discovery time
	Finish    map[string]int    // Finish time
}

// DFS performs depth-first search.
func (g *Graph) DFS(source string) (*DFSResult, error) {
	if !g.HasNode(source) {
		return nil, ErrNodeNotFound
	}

	visited := make(map[string]bool)
	parent := make(map[string]string)
	discovery := make(map[string]int)
	finish := make(map[string]int)
	order := []string{}
	finished := []string{}
	time := 0

	var dfsVisit func(node string)
	dfsVisit = func(node string) {
		visited[node] = true
		time++
		discovery[node] = time
		order = append(order, node)

		for _, edge := range g.Adj[node] {
			if !visited[edge.To] {
				parent[edge.To] = node
				dfsVisit(edge.To)
			}
		}

		time++
		finish[node] = time
		finished = append(finished, node)
	}

	dfsVisit(source)

	// Visit unconnected nodes
	for _, id := range g.NodeIDs() {
		if !visited[id] {
			dfsVisit(id)
		}
	}

	return &DFSResult{
		Order:     order,
		Finished:  finished,
		Visited:   visited,
		Parent:    parent,
		Discovery: discovery,
		Finish:    finish,
	}, nil
}

// DFSAll performs DFS from all nodes and returns the full traversal.
func (g *Graph) DFSAll() *DFSResult {
	visited := make(map[string]bool)
	parent := make(map[string]string)
	discovery := make(map[string]int)
	finish := make(map[string]int)
	order := []string{}
	finished := []string{}
	time := 0

	var dfsVisit func(node string)
	dfsVisit = func(node string) {
		visited[node] = true
		time++
		discovery[node] = time
		order = append(order, node)

		for _, edge := range g.Adj[node] {
			if !visited[edge.To] {
				parent[edge.To] = node
				dfsVisit(edge.To)
			}
		}

		time++
		finish[node] = time
		finished = append(finished, node)
	}

	for _, id := range g.NodeIDs() {
		if !visited[id] {
			dfsVisit(id)
		}
	}

	return &DFSResult{
		Order:     order,
		Finished:  finished,
		Visited:   visited,
		Parent:    parent,
		Discovery: discovery,
		Finish:    finish,
	}
}

// DijkstraResult contains the result of Dijkstra's algorithm.
type DijkstraResult struct {
	Dist map[string]float64  // Shortest distance from source
	Prev map[string]string   // Previous node in shortest path
	Path map[string][]string // Full path from source to each node
}

// Dijkstra implements Dijkstra's shortest path algorithm.
func (g *Graph) Dijkstra(source string) (*DijkstraResult, error) {
	if !g.HasNode(source) {
		return nil, ErrNodeNotFound
	}

	dist := make(map[string]float64)
	prev := make(map[string]string)
	for id := range g.Nodes {
		dist[id] = math.Inf(1)
	}
	dist[source] = 0

	// Priority queue
	pq := &NodePQ{}
	heap.Init(pq)
	heap.Push(pq, &NodeItem{ID: source, Dist: 0})

	visited := make(map[string]bool)

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*NodeItem)
		u := item.ID

		if visited[u] {
			continue
		}
		visited[u] = true

		for _, edge := range g.Adj[u] {
			v := edge.To
			alt := dist[u] + edge.Weight
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u
				heap.Push(pq, &NodeItem{ID: v, Dist: alt})
			}
		}
	}

	// Build paths
	paths := make(map[string][]string)
	for id := range g.Nodes {
		paths[id] = g.buildPath(prev, source, id)
	}

	return &DijkstraResult{
		Dist: dist,
		Prev: prev,
		Path: paths,
	}, nil
}

func (g *Graph) buildPath(prev map[string]string, source, target string) []string {
	if _, ok := prev[target]; !ok && target != source {
		return nil
	}
	path := []string{}
	for current := target; current != ""; current = prev[current] {
		path = append([]string{current}, path...)
		if current == source {
			break
		}
	}
	return path
}

// NodeItem is used in the priority queue for Dijkstra.
type NodeItem struct {
	ID   string
	Dist float64
}

// NodePQ is a priority queue for nodes.
type NodePQ []*NodeItem

func (pq NodePQ) Len() int            { return len(pq) }
func (pq NodePQ) Less(i, j int) bool  { return pq[i].Dist < pq[j].Dist }
func (pq NodePQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *NodePQ) Push(x interface{}) { *pq = append(*pq, x.(*NodeItem)) }
func (pq *NodePQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

// TopologicalSort performs topological sorting (Kahn's algorithm).
func (g *Graph) TopologicalSort() ([]string, error) {
	if !g.Directed {
		return nil, ErrCycleDetected // Topological sort only for directed graphs
	}

	inDegree := make(map[string]int)
	for id := range g.Nodes {
		inDegree[id] = 0
	}
	for _, edges := range g.Adj {
		for _, edge := range edges {
			inDegree[edge.To]++
		}
	}

	queue := []string{}
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}
	sort.Strings(queue) // Deterministic order

	result := []string{}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		for _, edge := range g.Adj[node] {
			inDegree[edge.To]--
			if inDegree[edge.To] == 0 {
				queue = append(queue, edge.To)
				sort.Strings(queue) // Keep deterministic
			}
		}
	}

	if len(result) != len(g.Nodes) {
		return nil, ErrCycleDetected
	}

	return result, nil
}

// DetectCycle detects if the graph contains a cycle.
func (g *Graph) DetectCycle() (bool, []string) {
	if !g.Directed {
		// For undirected graphs, use DFS to detect back edges
		visited := make(map[string]bool)
		parent := make(map[string]string)

		var dfs func(node string) (bool, []string)
		dfs = func(node string) (bool, []string) {
			visited[node] = true
			for _, edge := range g.Adj[node] {
				if !visited[edge.To] {
					parent[edge.To] = node
					if found, path := dfs(edge.To); found {
						return true, path
					}
				} else if edge.To != parent[node] {
					return true, []string{edge.To, node}
				}
			}
			return false, nil
		}

		for _, id := range g.NodeIDs() {
			if !visited[id] {
				if found, path := dfs(id); found {
					return true, path
				}
			}
		}
		return false, nil
	}

	// For directed graphs, use DFS with coloring
	const (
		White = 0 // Unvisited
		Gray  = 1 // In progress
		Black = 2 // Finished
	)

	color := make(map[string]int)
	parent := make(map[string]string)
	for id := range g.Nodes {
		color[id] = White
	}

	var dfs func(node string) (bool, []string)
	dfs = func(node string) (bool, []string) {
		color[node] = Gray

		for _, edge := range g.Adj[node] {
			if color[edge.To] == Gray {
				// Found cycle - reconstruct
				cycle := []string{edge.To, node}
				current := node
				for current != edge.To {
					current = parent[current]
					if current == "" {
						break
					}
					cycle = append([]string{current}, cycle...)
				}
				return true, cycle
			}
			if color[edge.To] == White {
				parent[edge.To] = node
				if found, path := dfs(edge.To); found {
					return true, path
				}
			}
		}

		color[node] = Black
		return false, nil
	}

	for _, id := range g.NodeIDs() {
		if color[id] == White {
			if found, path := dfs(id); found {
				return true, path
			}
		}
	}

	return false, nil
}

// PageRank computes PageRank scores.
func (g *Graph) PageRank(damping float64, iterations int) map[string]float64 {
	if damping == 0 {
		damping = 0.85
	}
	if iterations == 0 {
		iterations = 100
	}

	n := float64(g.NodeCount())
	rank := make(map[string]float64)
	newRank := make(map[string]float64)

	// Initialize
	for id := range g.Nodes {
		rank[id] = 1.0 / n
	}

	for i := 0; i < iterations; i++ {
		// Reset newRank
		for id := range g.Nodes {
			newRank[id] = (1 - damping) / n
		}

		// Distribute rank
		for id := range g.Nodes {
			outDegree := g.OutDegree(id)
			if outDegree > 0 {
				share := rank[id] / float64(outDegree)
				for _, edge := range g.Adj[id] {
					newRank[edge.To] += damping * share
				}
			} else {
				// Dangling node: distribute evenly
				share := rank[id] / n
				for target := range g.Nodes {
					newRank[target] += damping * share
				}
			}
		}

		rank, newRank = newRank, rank
	}

	return rank
}

// ConnectedComponents returns the connected components of an undirected graph.
func (g *Graph) ConnectedComponents() [][]string {
	visited := make(map[string]bool)
	components := [][]string{}

	var dfs func(node string, component *[]string)
	dfs = func(node string, component *[]string) {
		visited[node] = true
		*component = append(*component, node)
		for _, edge := range g.Adj[node] {
			if !visited[edge.To] {
				dfs(edge.To, component)
			}
		}
	}

	for _, id := range g.NodeIDs() {
		if !visited[id] {
			component := []string{}
			dfs(id, &component)
			sort.Strings(component)
			components = append(components, component)
		}
	}

	return components
}

// StronglyConnectedComponents returns the SCCs using Tarjan's algorithm.
func (g *Graph) StronglyConnectedComponents() [][]string {
	if !g.Directed {
		return [][]string{g.ConnectedComponents()[0]}
	}

	index := 0
	stack := []string{}
	onStack := make(map[string]bool)
	indices := make(map[string]int)
	lowlinks := make(map[string]int)
	sccs := [][]string{}

	var strongConnect func(v string)
	strongConnect = func(v string) {
		indices[v] = index
		lowlinks[v] = index
		index++
		stack = append(stack, v)
		onStack[v] = true

		for _, edge := range g.Adj[v] {
			w := edge.To
			if _, ok := indices[w]; !ok {
				strongConnect(w)
				if lowlinks[w] < lowlinks[v] {
					lowlinks[v] = lowlinks[w]
				}
			} else if onStack[w] {
				if indices[w] < lowlinks[v] {
					lowlinks[v] = indices[w]
				}
			}
		}

		if lowlinks[v] == indices[v] {
			scc := []string{}
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				scc = append(scc, w)
				if w == v {
					break
				}
			}
			sort.Strings(scc)
			sccs = append(sccs, scc)
		}
	}

	for _, id := range g.NodeIDs() {
		if _, ok := indices[id]; !ok {
			strongConnect(id)
		}
	}

	return sccs
}

// UnionFind is a disjoint set data structure.
type UnionFind struct {
	parent map[string]string
	rank   map[string]int
}

// NewUnionFind creates a new UnionFind structure.
func NewUnionFind() *UnionFind {
	return &UnionFind{
		parent: make(map[string]string),
		rank:   make(map[string]int),
	}
}

// MakeSet creates a new set with a single element.
func (uf *UnionFind) MakeSet(x string) {
	if _, ok := uf.parent[x]; !ok {
		uf.parent[x] = x
		uf.rank[x] = 0
	}
}

// Find finds the representative of the set containing x (with path compression).
func (uf *UnionFind) Find(x string) string {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // Path compression
	}
	return uf.parent[x]
}

// Union merges the sets containing x and y (by rank).
func (uf *UnionFind) Union(x, y string) {
	rx := uf.Find(x)
	ry := uf.Find(y)

	if rx == ry {
		return
	}

	// Union by rank
	if uf.rank[rx] < uf.rank[ry] {
		uf.parent[rx] = ry
	} else if uf.rank[rx] > uf.rank[ry] {
		uf.parent[ry] = rx
	} else {
		uf.parent[ry] = rx
		uf.rank[rx]++
	}
}

// SameSet checks if x and y are in the same set.
func (uf *UnionFind) SameSet(x, y string) bool {
	return uf.Find(x) == uf.Find(y)
}

// CountSets returns the number of disjoint sets.
func (uf *UnionFind) CountSets() int {
	roots := make(map[string]bool)
	for x := range uf.parent {
		roots[uf.Find(x)] = true
	}
	return len(roots)
}

// ShortestPath returns the shortest path between two nodes using BFS (unweighted).
func (g *Graph) ShortestPath(source, target string) ([]string, error) {
	if !g.HasNode(source) {
		return nil, ErrSourceNotFound
	}
	if !g.HasNode(target) {
		return nil, ErrTargetNotFound
	}

	result, err := g.BFS(source)
	if err != nil {
		return nil, err
	}

	path, ok := result.Path()[target]
	if !ok || len(path) == 0 {
		return nil, fmt.Errorf("no path from %s to %s", source, target)
	}

	return path, nil
}

// Path returns the path from BFS result.
func (r *BFSResult) Path() map[string][]string {
	paths := make(map[string][]string)
	for node := range r.Parent {
		path := []string{}
		current := node
		for current != "" {
			path = append([]string{current}, path...)
			if parent, ok := r.Parent[current]; ok {
				current = parent
			} else {
				break
			}
		}
		paths[node] = path
	}
	// Source has a path of just itself
	paths[""] = nil
	return paths
}
