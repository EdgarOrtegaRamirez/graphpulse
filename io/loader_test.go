package io

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/EdgarOrtegaRamirez/graphpulse/graph"
)

func TestLoadEdgeList(t *testing.T) {
	content := `# Comment line
a,b
b,c
a,c,3.0
`
	path := filepath.Join(t.TempDir(), "edges.txt")
	os.WriteFile(path, []byte(content), 0644)

	g, err := LoadEdgeList(path, true)
	if err != nil {
		t.Fatalf("LoadEdgeList failed: %v", err)
	}

	if g.NodeCount() != 3 {
		t.Errorf("expected 3 nodes, got %d", g.NodeCount())
	}
	if g.EdgeCount() != 3 {
		t.Errorf("expected 3 edges, got %d", g.EdgeCount())
	}
}

func TestLoadEdgeListUndirected(t *testing.T) {
	content := `a,b
b,c
`
	path := filepath.Join(t.TempDir(), "edges.txt")
	os.WriteFile(path, []byte(content), 0644)

	g, err := LoadEdgeList(path, false)
	if err != nil {
		t.Fatalf("LoadEdgeList failed: %v", err)
	}

	if g.EdgeCount() != 2 {
		t.Errorf("expected 2 edges, got %d", g.EdgeCount())
	}
	if !g.HasEdge("b", "a") {
		t.Error("expected bidirectional edge in undirected graph")
	}
}

func TestLoadCSV(t *testing.T) {
	content := `source,target,weight
a,b,1.0
b,c,2.0
a,c,3.0
`
	path := filepath.Join(t.TempDir(), "graph.csv")
	os.WriteFile(path, []byte(content), 0644)

	g, err := LoadCSV(path, true)
	if err != nil {
		t.Fatalf("LoadCSV failed: %v", err)
	}

	if g.NodeCount() != 3 {
		t.Errorf("expected 3 nodes, got %d", g.NodeCount())
	}
	if g.EdgeCount() != 3 {
		t.Errorf("expected 3 edges, got %d", g.EdgeCount())
	}
}

func TestLoadJSON(t *testing.T) {
	content := `{
  "directed": true,
  "nodes": [
    {"id": "a"},
    {"id": "b"},
    {"id": "c"}
  ],
  "edges": [
    {"from": "a", "to": "b", "weight": 1.0},
    {"from": "b", "to": "c", "weight": 2.0}
  ]
}`
	path := filepath.Join(t.TempDir(), "graph.json")
	os.WriteFile(path, []byte(content), 0644)

	g, err := LoadJSON(path)
	if err != nil {
		t.Fatalf("LoadJSON failed: %v", err)
	}

	if g.NodeCount() != 3 {
		t.Errorf("expected 3 nodes, got %d", g.NodeCount())
	}
	if g.EdgeCount() != 2 {
		t.Errorf("expected 2 edges, got %d", g.EdgeCount())
	}
}

func TestSaveAndLoadJSON(t *testing.T) {
	// Create a graph
	g1 := graph.New(true)
	g1.AddEdge("a", "b", 1.5, nil)
	g1.AddEdge("b", "c", 2.5, nil)

	path := filepath.Join(t.TempDir(), "graph.json")
	if err := SaveJSON(g1, path); err != nil {
		t.Fatalf("SaveJSON failed: %v", err)
	}

	g2, err := LoadJSON(path)
	if err != nil {
		t.Fatalf("LoadJSON failed: %v", err)
	}

	if g2.NodeCount() != g1.NodeCount() {
		t.Errorf("node count mismatch: %d vs %d", g2.NodeCount(), g1.NodeCount())
	}
	if g2.EdgeCount() != g1.EdgeCount() {
		t.Errorf("edge count mismatch: %d vs %d", g2.EdgeCount(), g1.EdgeCount())
	}
}

func TestLoadGraphAutoDetect(t *testing.T) {
	// Test JSON
	jsonContent := `{"directed":true,"nodes":[{"id":"x"}],"edges":[]}`
	jsonPath := filepath.Join(t.TempDir(), "test.json")
	os.WriteFile(jsonPath, []byte(jsonContent), 0644)

	g, err := LoadGraph(jsonPath, true)
	if err != nil {
		t.Fatalf("LoadGraph JSON failed: %v", err)
	}
	if g.NodeCount() != 1 {
		t.Errorf("expected 1 node, got %d", g.NodeCount())
	}

	// Test CSV
	csvContent := "from,to\na,b\n"
	csvPath := filepath.Join(t.TempDir(), "test.csv")
	os.WriteFile(csvPath, []byte(csvContent), 0644)

	g, err = LoadGraph(csvPath, true)
	if err != nil {
		t.Fatalf("LoadGraph CSV failed: %v", err)
	}
	if g.NodeCount() != 2 {
		t.Errorf("expected 2 nodes, got %d", g.NodeCount())
	}
}

func TestLoadEdgeListInvalid(t *testing.T) {
	content := `a`
	path := filepath.Join(t.TempDir(), "bad.txt")
	os.WriteFile(path, []byte(content), 0644)

	_, err := LoadEdgeList(path, true)
	if err == nil {
		t.Error("expected error for invalid edge list")
	}
}

func TestLoadJSONInvalid(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)

	_, err := LoadJSON(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadCSVInvalidHeaders(t *testing.T) {
	content := "foo,bar\n1,2\n"
	path := filepath.Join(t.TempDir(), "bad.csv")
	os.WriteFile(path, []byte(content), 0644)

	_, err := LoadCSV(path, true)
	if err == nil {
		t.Error("expected error for CSV with wrong headers")
	}
}
