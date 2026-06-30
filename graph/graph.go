// Package graph provides core graph data structures and algorithms.
package graph

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrNodeNotFound      = errors.New("node not found")
	ErrEdgeNotFound      = errors.New("edge not found")
	ErrCycleDetected     = errors.New("cycle detected in graph")
	ErrNegativeWeight    = errors.New("negative edge weight detected")
	ErrSourceNotFound    = errors.New("source node not found")
	ErrTargetNotFound    = errors.New("target node not found")
	ErrEmptyGraph        = errors.New("graph has no nodes")
	ErrDisconnectedGraph = errors.New("graph is not strongly connected")
)

// Node represents a node in the graph.
type Node struct {
	ID    string
	Attrs map[string]interface{}
}

// Edge represents a directed or undirected edge.
type Edge struct {
	From   string
	To     string
	Weight float64
	Attrs  map[string]interface{}
}

// Graph represents a graph with adjacency lists.
type Graph struct {
	Nodes    map[string]*Node
	Adj      map[string][]Edge // adjacency list (outgoing edges)
	Reverse  map[string][]Edge // reverse adjacency list (incoming edges)
	Directed bool
}

// New creates a new empty graph.
func New(directed bool) *Graph {
	return &Graph{
		Nodes:    make(map[string]*Node),
		Adj:      make(map[string][]Edge),
		Reverse:  make(map[string][]Edge),
		Directed: directed,
	}
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(id string, attrs map[string]interface{}) {
	if _, exists := g.Nodes[id]; !exists {
		g.Nodes[id] = &Node{ID: id, Attrs: attrs}
		g.Adj[id] = []Edge{}
		g.Reverse[id] = []Edge{}
	}
}

// AddEdge adds an edge between two nodes.
func (g *Graph) AddEdge(from, to string, weight float64, attrs map[string]interface{}) {
	g.AddNode(from, nil)
	g.AddNode(to, nil)

	edge := Edge{From: from, To: to, Weight: weight, Attrs: attrs}
	g.Adj[from] = append(g.Adj[from], edge)
	g.Reverse[to] = append(g.Reverse[to], edge)

	if !g.Directed {
		revEdge := Edge{From: to, To: from, Weight: weight, Attrs: attrs}
		g.Adj[to] = append(g.Adj[to], revEdge)
		g.Reverse[from] = append(g.Reverse[from], revEdge)
	}
}

// RemoveEdge removes an edge between two nodes.
func (g *Graph) RemoveEdge(from, to string) bool {
	if !g.Directed {
		g.removeEdgeOneWay(from, to)
		g.removeEdgeOneWay(to, from)
		return true
	}
	return g.removeEdgeOneWay(from, to)
}

func (g *Graph) removeEdgeOneWay(from, to string) bool {
	edges := g.Adj[from]
	for i, e := range edges {
		if e.To == to {
			g.Adj[from] = append(edges[:i], edges[i+1:]...)
			// Remove from reverse
			revEdges := g.Reverse[to]
			for j, re := range revEdges {
				if re.From == from {
					g.Reverse[to] = append(revEdges[:j], revEdges[j+1:]...)
					break
				}
			}
			return true
		}
	}
	return false
}

// NodeCount returns the number of nodes.
func (g *Graph) NodeCount() int {
	return len(g.Nodes)
}

// EdgeCount returns the number of edges.
func (g *Graph) EdgeCount() int {
	count := 0
	for _, edges := range g.Adj {
		count += len(edges)
	}
	if !g.Directed {
		count /= 2
	}
	return count
}

// Neighbors returns the neighbors of a node.
func (g *Graph) Neighbors(id string) []string {
	neighbors := []string{}
	for _, e := range g.Adj[id] {
		neighbors = append(neighbors, e.To)
	}
	return neighbors
}

// InDegree returns the in-degree of a node.
func (g *Graph) InDegree(id string) int {
	return len(g.Reverse[id])
}

// OutDegree returns the out-degree of a node.
func (g *Graph) OutDegree(id string) int {
	return len(g.Adj[id])
}

// Degree returns the degree of a node (in-degree for directed, degree for undirected).
func (g *Graph) Degree(id string) int {
	if g.Directed {
		return g.InDegree(id) + g.OutDegree(id)
	}
	return g.OutDegree(id)
}

// HasNode checks if a node exists.
func (g *Graph) HasNode(id string) bool {
	_, exists := g.Nodes[id]
	return exists
}

// HasEdge checks if an edge exists.
func (g *Graph) HasEdge(from, to string) bool {
	for _, e := range g.Adj[from] {
		if e.To == to {
			return true
		}
	}
	return false
}

// NodeIDs returns a sorted list of all node IDs.
func (g *Graph) NodeIDs() []string {
	ids := make([]string, 0, len(g.Nodes))
	for id := range g.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// GetEdge returns an edge between two nodes.
func (g *Graph) GetEdge(from, to string) (*Edge, error) {
	for _, e := range g.Adj[from] {
		if e.To == to {
			return &e, nil
		}
	}
	return nil, ErrEdgeNotFound
}

// Clone creates a deep copy of the graph.
func (g *Graph) Clone() *Graph {
	ng := New(g.Directed)
	for id, node := range g.Nodes {
		attrs := make(map[string]interface{})
		for k, v := range node.Attrs {
			attrs[k] = v
		}
		ng.AddNode(id, attrs)
	}
	for from, edges := range g.Adj {
		for _, e := range edges {
			attrs := make(map[string]interface{})
			for k, v := range e.Attrs {
				attrs[k] = v
			}
			ng.Adj[from] = append(ng.Adj[from], Edge{
				From: e.From, To: e.To, Weight: e.Weight, Attrs: attrs,
			})
		}
	}
	for to, edges := range g.Reverse {
		for _, e := range edges {
			attrs := make(map[string]interface{})
			for k, v := range e.Attrs {
				attrs[k] = v
			}
			ng.Reverse[to] = append(ng.Reverse[to], Edge{
				From: e.From, To: e.To, Weight: e.Weight, Attrs: attrs,
			})
		}
	}
	return ng
}

// String returns a string representation of the graph.
func (g *Graph) String() string {
	directed := "undirected"
	if g.Directed {
		directed = "directed"
	}
	result := fmt.Sprintf("Graph(%s, nodes=%d, edges=%d)\n", directed, g.NodeCount(), g.EdgeCount())
	for _, id := range g.NodeIDs() {
		neighbors := g.Neighbors(id)
		result += fmt.Sprintf("  %s -> %v\n", id, neighbors)
	}
	return result
}
