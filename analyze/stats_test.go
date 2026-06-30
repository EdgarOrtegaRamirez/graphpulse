package analyze

import (
	"testing"

	"github.com/EdgarOrtegaRamirez/graphpulse/graph"
)

func TestComputeStats(t *testing.T) {
	g := graph.New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)
	g.AddEdge("c", "a", 1.0, nil)

	stats := ComputeStats(g)

	if stats.NodeCount != 3 {
		t.Errorf("expected 3 nodes, got %d", stats.NodeCount)
	}
	if stats.EdgeCount != 3 {
		t.Errorf("expected 3 edges, got %d", stats.EdgeCount)
	}
	// For directed graph: max edges = n*(n-1) = 3*2 = 6, density = 3/6 = 0.5
	if stats.Density != 0.5 {
		t.Errorf("expected density 0.5, got %f", stats.Density)
	}
	if !stats.HasCycles {
		t.Error("expected cycle detection to find cycle")
	}
	if stats.IsDAG {
		t.Error("expected IsDAG to be false")
	}
}

func TestComputeStatsEmptyGraph(t *testing.T) {
	g := graph.New(true)
	stats := ComputeStats(g)

	if stats.NodeCount != 0 {
		t.Errorf("expected 0 nodes, got %d", stats.NodeCount)
	}
	if stats.Density != 0 {
		t.Errorf("expected density 0, got %f", stats.Density)
	}
}

func TestComputeStatsDAG(t *testing.T) {
	g := graph.New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)

	stats := ComputeStats(g)

	if stats.HasCycles {
		t.Error("expected no cycles in DAG")
	}
	if !stats.IsDAG {
		t.Error("expected IsDAG to be true")
	}
}

func TestComputeStatsDisconnected(t *testing.T) {
	g := graph.New(false)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("c", "d", 1.0, nil)

	stats := ComputeStats(g)

	if stats.IsConnected {
		t.Error("expected disconnected graph")
	}
	if stats.NumComponents != 2 {
		t.Errorf("expected 2 components, got %d", stats.NumComponents)
	}
}

func TestClusteringCoefficient(t *testing.T) {
	// Triangle: a-b-c-a has clustering coefficient 1.0
	g := graph.New(false)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)
	g.AddEdge("c", "a", 1.0, nil)

	stats := ComputeStats(g)
	if stats.ClusteringCoeff != 1.0 {
		t.Errorf("expected clustering coeff 1.0, got %f", stats.ClusteringCoeff)
	}
}

func TestPageRankTop10(t *testing.T) {
	g := graph.New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("a", "c", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)
	g.AddEdge("c", "a", 1.0, nil)

	stats := ComputeStats(g)
	if len(stats.PageRankTop10) == 0 {
		t.Error("expected PageRank results")
	}
}

func TestFormatStats(t *testing.T) {
	g := graph.New(true)
	g.AddEdge("a", "b", 1.0, nil)

	stats := ComputeStats(g)
	result := FormatStats(stats)

	if result == "" {
		t.Error("expected non-empty formatted stats")
	}
}

func TestDegreeHistogram(t *testing.T) {
	g := graph.New(true)
	g.AddEdge("a", "b", 1.0, nil)
	g.AddEdge("a", "c", 1.0, nil)
	g.AddEdge("b", "c", 1.0, nil)

	stats := ComputeStats(g)
	if len(stats.DegreeHistogram) == 0 {
		t.Error("expected degree histogram")
	}
}
