// Package visualize provides graph visualization in ASCII and DOT formats.
package visualize

import (
	"fmt"
	"sort"
	"strings"

	"github.com/EdgarOrtegaRamirez/graphpulse/graph"
)

// ASCIIConfig holds configuration for ASCII visualization.
type ASCIIConfig struct {
	Width    int  // Terminal width
	Height   int  // Terminal height
	Compact  bool // Use compact layout
	ShowWeights bool // Show edge weights
}

// DefaultASCIIConfig returns default ASCII config.
func DefaultASCIIConfig() *ASCIIConfig {
	return &ASCIIConfig{
		Width:    80,
		Height:   24,
		Compact:  false,
		ShowWeights: true,
	}
}

// RenderASCII renders a graph as ASCII art.
func RenderASCII(g *graph.Graph, config *ASCIIConfig) string {
	if config == nil {
		config = DefaultASCIIConfig()
	}

	if g.NodeCount() == 0 {
		return "(empty graph)"
	}

	if g.NodeCount() == 1 {
		for id := range g.Nodes {
			return fmt.Sprintf("[ %s ]", id)
		}
	}

	// Simple layered layout using BFS
	result := []string{}
	result = append(result, fmt.Sprintf("Graph (%d nodes, %d edges)", g.NodeCount(), g.EdgeCount()))
	result = append(result, strings.Repeat("─", 40))

	// Group nodes by BFS layers from first node
	startNode := g.NodeIDs()[0]
	visited := make(map[string]bool)
	layers := [][]string{}

	queue := []string{startNode}
	visited[startNode] = true

	for len(queue) > 0 {
		layer := []string{}
		nextQueue := []string{}

		for _, node := range queue {
			layer = append(layer, node)
			for _, edge := range g.Adj[node] {
				if !visited[edge.To] {
					visited[edge.To] = true
					nextQueue = append(nextQueue, edge.To)
				}
			}
		}

		layers = append(layers, layer)
		queue = nextQueue
	}

	// Add any unvisited nodes
	for _, id := range g.NodeIDs() {
		if !visited[id] {
			layers = append(layers, []string{id})
		}
	}

	// Render layers
	for i, layer := range layers {
		prefix := "  "
		if i == 0 {
			prefix = "▶ "
		} else {
			prefix = "│ "
		}

		nodes := strings.Join(layer, "  ")
		result = append(result, fmt.Sprintf("%s%s", prefix, nodes))

		if i < len(layers)-1 {
			// Draw connections to next layer
			nextLayer := layers[i+1]
			connections := []string{}
			for _, node := range layer {
				for _, edge := range g.Adj[node] {
					for _, next := range nextLayer {
						if edge.To == next {
							conn := fmt.Sprintf("%s→%s", node, edge.To)
							if config.ShowWeights && edge.Weight != 1.0 {
								conn += fmt.Sprintf("(%.1f)", edge.Weight)
							}
							connections = append(connections, conn)
						}
					}
				}
			}
			if len(connections) > 0 {
				result = append(result, fmt.Sprintf("  %s", strings.Join(connections, "  ")))
			}
			result = append(result, "  │")
		}
	}

	// Add edge list
	result = append(result, "")
	result = append(result, "Edges:")
	for _, id := range g.NodeIDs() {
		for _, edge := range g.Adj[id] {
			weight := ""
			if edge.Weight != 1.0 {
				weight = fmt.Sprintf(" [%.1f]", edge.Weight)
			}
			result = append(result, fmt.Sprintf("  %s → %s%s", edge.From, edge.To, weight))
		}
	}

	return strings.Join(result, "\n")
}

// DOTConfig holds configuration for DOT export.
type DOTConfig struct {
	Directed  bool
	RankDir   string // "TB", "LR", "BT", "RL"
	NodeShape string
	NodeColor string
	EdgeColor string
	FontSize  int
	Title     string
}

// DefaultDOTConfig returns default DOT config.
func DefaultDOTConfig() *DOTConfig {
	return &DOTConfig{
		Directed:  true,
		RankDir:   "LR",
		NodeShape: "circle",
		NodeColor: "#4CAF50",
		EdgeColor: "#666666",
		FontSize:  12,
	}
}

// ExportDOT exports a graph in Graphviz DOT format.
func ExportDOT(g *graph.Graph, config *DOTConfig) string {
	if config == nil {
		config = DefaultDOTConfig()
	}

	directed := "digraph"
	if !g.Directed {
		directed = "graph"
	}

	lines := []string{}
	lines = append(lines, fmt.Sprintf("%s G {", directed))
	lines = append(lines, fmt.Sprintf("  rankdir=%s;", config.RankDir))
	lines = append(lines, fmt.Sprintf("  node [shape=%s, style=filled, fillcolor=%s, fontsize=%d];",
		config.NodeShape, config.NodeColor, config.FontSize))
	lines = append(lines, fmt.Sprintf("  edge [color=%s];", config.EdgeColor))

	if config.Title != "" {
		lines = append(lines, fmt.Sprintf("  labelloc=t; label=\"%s\";", config.Title))
	}

	// Nodes
	lines = append(lines, "")
	lines = append(lines, "  // Nodes")
	for _, id := range g.NodeIDs() {
		label := id
		lines = append(lines, fmt.Sprintf("  \"%s\" [label=\"%s\"];", id, label))
	}

	// Edges
	lines = append(lines, "")
	lines = append(lines, "  // Edges")
	seen := make(map[string]bool)
	for _, id := range g.NodeIDs() {
		for _, edge := range g.Adj[id] {
			key := edge.From + "->" + edge.To
			if !g.Directed && edge.From > edge.To {
				key = edge.To + "->" + edge.From
			}
			if seen[key] {
				continue
			}
			seen[key] = true

			arrow := "->"
			if !g.Directed {
				arrow = "--"
			}

			edgeAttr := ""
			if edge.Weight != 1.0 {
				edgeAttr = fmt.Sprintf(" [label=\"%.1f\"]", edge.Weight)
			}

			lines = append(lines, fmt.Sprintf("  \"%s\" %s \"%s\"%s;", edge.From, arrow, edge.To, edgeAttr))
		}
	}

	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

// RenderDegreeHistogram renders an ASCII histogram of degree distribution.
func RenderDegreeHistogram(g *graph.Graph) string {
	degreeCounts := make(map[int]int)
	maxDegree := 0

	for _, id := range g.NodeIDs() {
		deg := g.Degree(id)
		degreeCounts[deg]++
		if deg > maxDegree {
			maxDegree = deg
		}
	}

	if maxDegree == 0 {
		return "(no edges)"
	}

	// Limit to reasonable range
	if maxDegree > 20 {
		maxDegree = 20
	}

	lines := []string{}
	lines = append(lines, "Degree Distribution")
	lines = append(lines, strings.Repeat("─", 40))

	maxCount := 0
	for _, count := range degreeCounts {
		if count > maxCount {
			maxCount = count
		}
	}

	barWidth := 30
	for d := 0; d <= maxDegree; d++ {
		count := degreeCounts[d]
		barLen := 0
		if maxCount > 0 {
			barLen = (count * barWidth) / maxCount
		}
		bar := strings.Repeat("█", barLen)
		lines = append(lines, fmt.Sprintf("  %2d │ %-30s %d", d, bar, count))
	}

	return strings.Join(lines, "\n")
}

// RenderPageRankBar renders PageRank scores as a bar chart.
func RenderPageRankBar(g *graph.Graph, top int) string {
	if g.NodeCount() == 0 {
		return "(empty graph)"
	}

	ranks := g.PageRank(0.85, 100)

	type entry struct {
		node  string
		score float64
	}

	entries := make([]entry, 0, len(ranks))
	for node, score := range ranks {
		entries = append(entries, entry{node, score})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score > entries[j].score
	})

	if top > 0 && len(entries) > top {
		entries = entries[:top]
	}

	maxScore := 0.0
	for _, e := range entries {
		if e.score > maxScore {
			maxScore = e.score
		}
	}

	lines := []string{}
	lines = append(lines, fmt.Sprintf("PageRank (top %d)", len(entries)))
	lines = append(lines, strings.Repeat("─", 50))

	barWidth := 25
	for i, e := range entries {
		barLen := 0
		if maxScore > 0 {
			barLen = int((e.score / maxScore) * float64(barWidth))
		}
		bar := strings.Repeat("█", barLen)
		lines = append(lines, fmt.Sprintf("  %2d. %-16s %-25s %.6f", i+1, e.node, bar, e.score))
	}

	return strings.Join(lines, "\n")
}

// ShortestPathDOT highlights the shortest path in a DOT graph.
func ShortestPathDOT(g *graph.Graph, path []string, config *DOTConfig) string {
	if config == nil {
		config = DefaultDOTConfig()
	}

	directed := "digraph"
	if !g.Directed {
		directed = "graph"
	}

	lines := []string{}
	lines = append(lines, fmt.Sprintf("%s G {", directed))
	lines = append(lines, fmt.Sprintf("  rankdir=%s;", config.RankDir))
	lines = append(lines, "  node [shape=circle, style=filled, fontsize=12];")
	lines = append(lines, "  edge [color=#666666];")

	// Highlight path nodes
	pathSet := make(map[string]bool)
	for _, node := range path {
		pathSet[node] = true
	}

	for _, id := range g.NodeIDs() {
		color := "#E0E0E0"
		if pathSet[id] {
			color = "#FF9800"
		}
		lines = append(lines, fmt.Sprintf("  \"%s\" [label=\"%s\", fillcolor=\"%s\"];", id, id, color))
	}

	// Edges
	edgePairs := make(map[string]bool)
	for i := 0; i < len(path)-1; i++ {
		edgePairs[path[i]+"->"+path[i+1]] = true
	}

	seen := make(map[string]bool)
	for _, id := range g.NodeIDs() {
		for _, edge := range g.Adj[id] {
			key := edge.From + "->" + edge.To
			if !g.Directed && edge.From > edge.To {
				key = edge.To + "->" + edge.From
			}
			if seen[key] {
				continue
			}
			seen[key] = true

			arrow := "->"
			if !g.Directed {
				arrow = "--"
			}

			color := "#999999"
			penwidth := "1.0"
			if edgePairs[key] || (g.Directed && edgePairs[edge.To+"->"+edge.From]) {
				color = "#FF5722"
				penwidth = "3.0"
			}

			lines = append(lines, fmt.Sprintf("  \"%s\" %s \"%s\" [color=\"%s\", penwidth=%s];",
				edge.From, arrow, edge.To, color, penwidth))
		}
	}

	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}
