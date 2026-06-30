// Package io provides graph loading and saving in various formats.
package io

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/EdgarOrtegaRamirez/graphpulse/graph"
)

// LoadEdgeList loads a graph from an edge list file.
// Format: "from,to" or "from,to,weight" per line.
func LoadEdgeList(filename string, directed bool) (*graph.Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	g := graph.New(directed)
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: expected at least 2 fields, got %d", lineNum, len(parts))
		}

		from := strings.TrimSpace(parts[0])
		to := strings.TrimSpace(parts[1])
		weight := 1.0

		if len(parts) >= 3 {
			w, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid weight: %w", lineNum, err)
			}
			weight = w
		}

		g.AddEdge(from, to, weight, nil)
	}

	return g, scanner.Err()
}

// LoadCSV loads a graph from a CSV file with headers.
func LoadCSV(filename string, directed bool) (*graph.Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return graph.New(directed), nil
	}

	g := graph.New(directed)
	headers := records[0]

	// Determine column indices
	fromIdx, toIdx, weightIdx := -1, -1, -1
	for i, h := range headers {
		h = strings.ToLower(strings.TrimSpace(h))
		switch h {
		case "from", "source", "src", "node1":
			fromIdx = i
		case "to", "target", "dst", "node2":
			toIdx = i
		case "weight", "w", "cost":
			weightIdx = i
		}
	}

	if fromIdx == -1 || toIdx == -1 {
		return nil, fmt.Errorf("CSV must have 'from'/'source' and 'to'/'target' columns")
	}

	for rowIdx, row := range records[1:] {
		if fromIdx >= len(row) || toIdx >= len(row) {
			continue
		}

		from := strings.TrimSpace(row[fromIdx])
		to := strings.TrimSpace(row[toIdx])
		weight := 1.0

		if weightIdx >= 0 && weightIdx < len(row) {
			w, err := strconv.ParseFloat(strings.TrimSpace(row[weightIdx]), 64)
			if err == nil {
				weight = w
			}
		}

		_ = rowIdx
		g.AddEdge(from, to, weight, nil)
	}

	return g, nil
}

// GraphJSON is the JSON representation of a graph.
type GraphJSON struct {
	Nodes     []NodeJSON     `json:"nodes"`
	Edges     []EdgeJSON     `json:"edges"`
	Directed  bool           `json:"directed"`
}

// NodeJSON is the JSON representation of a node.
type NodeJSON struct {
	ID    string                 `json:"id"`
	Attrs map[string]interface{} `json:"attrs,omitempty"`
}

// EdgeJSON is the JSON representation of an edge.
type EdgeJSON struct {
	From   string                 `json:"from"`
	To     string                 `json:"to"`
	Weight float64                `json:"weight,omitempty"`
	Attrs  map[string]interface{} `json:"attrs,omitempty"`
}

// LoadJSON loads a graph from a JSON file.
func LoadJSON(filename string) (*graph.Graph, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var gj GraphJSON
	if err := json.Unmarshal(data, &gj); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	g := graph.New(gj.Directed)

	for _, node := range gj.Nodes {
		g.AddNode(node.ID, node.Attrs)
	}

	for _, edge := range gj.Edges {
		g.AddEdge(edge.From, edge.To, edge.Weight, edge.Attrs)
	}

	return g, nil
}

// SaveJSON saves a graph to a JSON file.
func SaveJSON(g *graph.Graph, filename string) error {
	gj := GraphJSON{
		Nodes:    []NodeJSON{},
		Edges:    []EdgeJSON{},
		Directed: g.Directed,
	}

	for _, id := range g.NodeIDs() {
		node := g.Nodes[id]
		gj.Nodes = append(gj.Nodes, NodeJSON{
			ID:    id,
			Attrs: node.Attrs,
		})
	}

	seen := make(map[string]bool)
	for from, edges := range g.Adj {
		for _, e := range edges {
			key := from + "->" + e.To
			if !g.Directed {
				if from > e.To {
					continue
				}
				key = e.From + "->" + e.To
			}
			if !seen[key] {
				seen[key] = true
				gj.Edges = append(gj.Edges, EdgeJSON{
					From:   e.From,
					To:     e.To,
					Weight: e.Weight,
					Attrs:  e.Attrs,
				})
			}
		}
	}

	data, err := json.MarshalIndent(gj, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// LoadGraph loads a graph from a file, auto-detecting format.
func LoadGraph(filename string, directed bool) (*graph.Graph, error) {
	ext := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(ext, ".json"):
		return LoadJSON(filename)
	case strings.HasSuffix(ext, ".csv"):
		return LoadCSV(filename, directed)
	default:
		return LoadEdgeList(filename, directed)
	}
}
