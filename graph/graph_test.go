package graph

import (
	"math"
	"testing"
)

func TestNewGraph(t *testing.T) {
	g := New(true)
	if g.NodeCount() != 0 {
		t.Errorf("expected 0 nodes, got %d", g.NodeCount())
	}
	if g.EdgeCount() != 0 {
		t.Errorf("expected 0 edges, got %d", g.EdgeCount())
	}
	if !g.Directed {
		t.Error("expected directed graph")
	}
}

func TestAddNode(t *testing.T) {
	g := New(true)
	g.AddNode("a", nil)
	g.AddNode("b", nil)

	if g.NodeCount() != 2 {
		t.Errorf("expected 2 nodes, got %d", g.NodeCount())
	}
	if !g.HasNode("a") {
		t.Error("expected node 'a' to exist")
	}
	if !g.HasNode("b") {
		t.Error("expected node 'b' to exist")
	}
}

func TestAddEdge(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("b", "c", 2.0, nil)

	if g.NodeCount() != 3 {
		t.Errorf("expected 3 nodes, got %d", g.NodeCount())
	}
	if g.EdgeCount() != 2 {
		t.Errorf("expected 2 edges, got %d", g.EdgeCount())
	}
	if !g.HasEdge("a", "b") {
		t.Error("expected edge a->b")
	}
	if !g.HasEdge("b", "c") {
		t.Error("expected edge b->c")
	}
	if g.HasEdge("b", "a") {
		t.Error("expected no edge b->a in directed graph")
	}
}

func TestUndirectedGraph(t *testing.T) {
	g := New(false)
	g.AddEdge("a", "b", 1.0, nil)

	if g.EdgeCount() != 1 {
		t.Errorf("expected 1 edge, got %d", g.EdgeCount())
	}
	if !g.HasEdge("a", "b") {
		t.Error("expected edge a->b")
	}
	if !g.HasEdge("b", "a") {
		t.Error("expected edge b->a in undirected graph")
	}
}

func TestNeighbors(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("a", "c", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)

	neighbors := g.Neighbors("a")
	if len(neighbors) != 2 {
		t.Errorf("expected 2 neighbors, got %d", len(neighbors))
	}
}

func TestDegree(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("a", "c", 1.0, nil)
	g.AddEdge("b", "a", 1.0, nil)

	if g.InDegree("a") != 1 {
		t.Errorf("expected in-degree 1, got %d", g.InDegree("a"))
	}
	if g.OutDegree("a") != 2 {
		t.Errorf("expected out-degree 2, got %d", g.OutDegree("a"))
	}
}

func TestBFS(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("a", "c", 1.0, nil)
	g.AddEdge("b", "d", 1.0, nil)
	g.AddEdge("c", "d", 1.0, nil)

	result, err := g.BFS("a")
	if err != nil {
		t.Fatalf("BFS failed: %v", err)
	}

	if len(result.Order) != 4 {
		t.Errorf("expected 4 nodes in order, got %d", len(result.Order))
	}
	if result.Dist["d"] != 2 {
		t.Errorf("expected distance 2 to d, got %d", result.Dist["d"])
	}
}

func TestBFSNotFound(t *testing.T) {
	g := New(true)
	g.AddNode("a", nil)

	_, err := g.BFS("x")
	if err != ErrNodeNotFound {
		t.Errorf("expected ErrNodeNotFound, got %v", err)
	}
}

func TestDFS(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("a", "c", 1.0, nil)
	g.AddEdge("b", "d", 1.0, nil)

	result, err := g.DFS("a")
	if err != nil {
		t.Fatalf("DFS failed: %v", err)
	}

	if len(result.Order) != 4 {
		t.Errorf("expected 4 nodes in order, got %d", len(result.Order))
	}
}

func TestDijkstra(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("a", "c", 4.0, nil)
	g.AddEdge("b", "c", 2.0, nil)
	g.AddEdge("b", "d", 6.0, nil)
	g.AddEdge("c", "d", 3.0, nil)

	result, err := g.Dijkstra("a")
	if err != nil {
		t.Fatalf("Dijkstra failed: %v", err)
	}

	if result.Dist["d"] != 6.0 {
		t.Errorf("expected distance 6 to d, got %f", result.Dist["d"])
	}
	if result.Dist["c"] != 3.0 {
		t.Errorf("expected distance 3 to c, got %f", result.Dist["c"])
	}
}

func TestDijkstraNotFound(t *testing.T) {
	g := New(true)
	g.AddNode("a", nil)

	_, err := g.Dijkstra("x")
	if err != ErrNodeNotFound {
		t.Errorf("expected ErrNodeNotFound, got %v", err)
	}
}

func TestTopologicalSort(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("a", "c", 1.0, nil)
	g.AddEdge("b", "d", 1.0, nil)
	g.AddEdge("c", "d", 1.0, nil)

	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if len(order) != 4 {
		t.Errorf("expected 4 nodes, got %d", len(order))
	}

	// Check ordering constraints
	indexOf := func(s string) int {
		for i, v := range order {
			if v == s {
				return i
			}
		}
		return -1
	}

	if indexOf("a") > indexOf("b") {
		t.Error("a should come before b")
	}
	if indexOf("a") > indexOf("c") {
		t.Error("a should come before c")
	}
	if indexOf("b") > indexOf("d") {
		t.Error("b should come before d")
	}
}

func TestTopologicalSortCycle(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)
	g.AddEdge("c", "a", 1.0, nil)

	_, err := g.TopologicalSort()
	if err != ErrCycleDetected {
		t.Errorf("expected ErrCycleDetected, got %v", err)
	}
}

func TestDetectCycleDirected(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)
	g.AddEdge("c", "a", 1.0, nil)

	hasCycle, path := g.DetectCycle()
	if !hasCycle {
		t.Error("expected cycle to be detected")
	}
	if len(path) < 2 {
		t.Error("expected cycle path with at least 2 nodes")
	}
}

func TestDetectCycleNoCycle(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)

	hasCycle, _ := g.DetectCycle()
	if hasCycle {
		t.Error("expected no cycle")
	}
}

func TestPageRank(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("a", "c", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)
	g.AddEdge("c", "a", 1.0, nil)

	ranks := g.PageRank(0.85, 100)

	// Sum should be approximately 1.0
	sum := 0.0
	for _, score := range ranks {
		sum += score
	}
	if math.Abs(sum-1.0) > 0.01 {
		t.Errorf("expected PageRank sum ~1.0, got %f", sum)
	}

	// Node with more incoming links should have higher rank
	if ranks["c"] <= ranks["b"] {
		t.Error("expected c to have higher rank than b")
	}
}

func TestConnectedComponents(t *testing.T) {
	g := New(false)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("c", "d", 1.0, nil)

	components := g.ConnectedComponents()
	if len(components) != 2 {
		t.Errorf("expected 2 components, got %d", len(components))
	}
}

func TestStronglyConnectedComponents(t *testing.T) {
	g := New(true)
	// Strongly connected: a -> b -> c -> a
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)
	g.AddEdge("c", "a", 1.0, nil)
	// Separate cycle: d -> e -> d
	g.AddEdge("d", "e", 1.0, nil)
	g.AddEdge("e", "d", 1.0, nil)

	sccs := g.StronglyConnectedComponents()
	if len(sccs) != 2 {
		t.Errorf("expected 2 SCCs, got %d: %v", len(sccs), sccs)
	}
}

func TestUnionFind(t *testing.T) {
	uf := NewUnionFind()
	uf.MakeSet("a")
	uf.MakeSet("b")
	uf.MakeSet("c")

	uf.Union("a", "b")

	if !uf.SameSet("a", "b") {
		t.Error("expected a and b to be in same set")
	}
	if uf.SameSet("a", "c") {
		t.Error("expected a and c to be in different sets")
	}
	if uf.CountSets() != 2 {
		t.Errorf("expected 2 sets, got %d", uf.CountSets())
	}
}

func TestUnionFindPathCompression(t *testing.T) {
	uf := NewUnionFind()
	uf.MakeSet("a")
	uf.MakeSet("b")
	uf.MakeSet("c")
	uf.MakeSet("d")

	uf.Union("a", "b")
	uf.Union("b", "c")
	uf.Union("c", "d")

	if !uf.SameSet("a", "d") {
		t.Error("expected a and d to be in same set")
	}
}

func TestClone(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("b", "c", 2.0, nil)

	clone := g.Clone()

	if clone.NodeCount() != g.NodeCount() {
		t.Error("clone has different node count")
	}
	if clone.EdgeCount() != g.EdgeCount() {
		t.Error("clone has different edge count")
	}

	// Modify original, clone should be unaffected
	g.AddNode("d", nil)
	if clone.HasNode("d") {
		t.Error("clone should not have node 'd'")
	}
}

func TestShortestPath(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)
	// No direct a->c edge, so shortest path is a->b->c

	path, err := g.ShortestPath("a", "c")
	if err != nil {
		t.Fatalf("ShortestPath failed: %v", err)
	}

	// BFS shortest path (unweighted) should be a->b->c (3 nodes)
	if len(path) != 3 {
		t.Errorf("expected path of length 3, got %d: %v", len(path), path)
	}
}

func TestNodeIDs(t *testing.T) {
	g := New(true)
	g.AddNode("c", nil)
	g.AddNode("a", nil)
	g.AddNode("b", nil)

	ids := g.NodeIDs()
	if len(ids) != 3 {
		t.Errorf("expected 3 IDs, got %d", len(ids))
	}
	if ids[0] != "a" || ids[1] != "b" || ids[2] != "c" {
		t.Errorf("expected sorted IDs, got %v", ids)
	}
}

func TestGetEdge(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 3.5, nil)

	edge, err := g.GetEdge("a", "b")
	if err != nil {
		t.Fatalf("GetEdge failed: %v", err)
	}
	if edge.Weight != 3.5 {
		t.Errorf("expected weight 3.5, got %f", edge.Weight)
	}

	_, err = g.GetEdge("a", "c")
	if err != ErrEdgeNotFound {
		t.Errorf("expected ErrEdgeNotFound, got %v", err)
	}
}

func TestString(t *testing.T) {
	g := New(true)
	g.AddEdge("a", "b", 1.0, nil)

	s := g.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
