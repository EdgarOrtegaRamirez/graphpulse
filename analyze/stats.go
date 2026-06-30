// Package analyze provides graph statistics and analysis.
package analyze

import (
	"fmt"
	"math"
	"sort"

	"github.com/EdgarOrtegaRamirez/graphpulse/graph"
)

// Stats holds graph statistics.
type Stats struct {
	NodeCount         int                `json:"node_count"`
	EdgeCount         int                `json:"edge_count"`
	Density           float64            `json:"density"`
	AvgDegree         float64            `json:"avg_degree"`
	MaxDegree         int                `json:"max_degree"`
	MinDegree         int                `json:"min_degree"`
	IsConnected       bool               `json:"is_connected"`
	NumComponents     int                `json:"num_components"`
	HasCycles         bool               `json:"has_cycles"`
	IsDAG             bool               `json:"is_dag"`
	Components        [][]string         `json:"components"`
	DegreeDistribution map[string]int    `json:"degree_distribution"`
	DegreeHistogram   []DegreeBucket     `json:"degree_histogram"`
	DiameterEstimate  int                `json:"diameter_estimate"`
	RadiusEstimate    int                `json:"radius_estimate"`
	ClusteringCoeff   float64            `json:"clustering_coefficient"`
	PageRankTop10     []PageRankEntry    `json:"pagerank_top10"`
}

// DegreeBucket is a bucket in the degree histogram.
type DegreeBucket struct {
	Degree int `json:"degree"`
	Count  int `json:"count"`
}

// PageRankEntry is a node with its PageRank score.
type PageRankEntry struct {
	Node  string  `json:"node"`
	Score float64 `json:"score"`
}

// ComputeStats computes comprehensive statistics for a graph.
func ComputeStats(g *graph.Graph) *Stats {
	stats := &Stats{
		NodeCount:          g.NodeCount(),
		EdgeCount:          g.EdgeCount(),
		DegreeDistribution: make(map[string]int),
	}

	// Density
	n := float64(stats.NodeCount)
	if n > 1 {
		maxEdges := n * (n - 1)
		if !g.Directed {
			maxEdges /= 2
		}
		stats.Density = float64(stats.EdgeCount) / maxEdges
	}

	// Degree statistics
	totalDegree := 0
	stats.MaxDegree = 0
	stats.MinDegree = math.MaxInt32

	degreeCounts := make(map[int]int)
	for _, id := range g.NodeIDs() {
		deg := g.Degree(id)
		totalDegree += deg
		if deg > stats.MaxDegree {
			stats.MaxDegree = deg
		}
		if deg < stats.MinDegree {
			stats.MinDegree = deg
		}
		degreeCounts[deg]++
	}

	if stats.NodeCount > 0 {
		stats.AvgDegree = float64(totalDegree) / float64(stats.NodeCount)
	}

	if stats.NodeCount == 0 {
		stats.MinDegree = 0
	}

	// Degree distribution
	for _, count := range degreeCounts {
		key := fmt.Sprintf("%d", count)
		stats.DegreeDistribution[key] = count
	}

	// Degree histogram
	maxDeg := stats.MaxDegree
	if maxDeg > 20 {
		maxDeg = 20 // Cap histogram
	}
	for d := 0; d <= maxDeg; d++ {
		if count, ok := degreeCounts[d]; ok {
			stats.DegreeHistogram = append(stats.DegreeHistogram, DegreeBucket{
				Degree: d,
				Count:  count,
			})
		}
	}

	// Connected components
	if !g.Directed {
		components := g.ConnectedComponents()
		stats.Components = components
		stats.NumComponents = len(components)
		stats.IsConnected = stats.NumComponents == 1
	} else {
		// For directed graphs, use weakly connected components
		// by creating an undirected copy
		ug := graph.New(false)
		for _, id := range g.NodeIDs() {
			ug.AddNode(id, nil)
		}
		for from, edges := range g.Adj {
			for _, e := range edges {
				ug.AddEdge(from, e.To, 1, nil)
			}
		}
		components := ug.ConnectedComponents()
		stats.Components = components
		stats.NumComponents = len(components)
		stats.IsConnected = stats.NumComponents == 1
	}

	// Cycle detection
	hasCycles, _ := g.DetectCycle()
	stats.HasCycles = hasCycles
	stats.IsDAG = !hasCycles

	// Diameter estimate (BFS from each node for small graphs, sample for large)
	stats.DiameterEstimate = estimateDiameter(g)
	stats.RadiusEstimate = estimateRadius(g)

	// Clustering coefficient
	stats.ClusteringCoeff = computeClusteringCoefficient(g)

	// PageRank (top 10)
	if stats.NodeCount > 0 {
		ranks := g.PageRank(0.85, 100)
		entries := make([]PageRankEntry, 0, len(ranks))
		for node, score := range ranks {
			entries = append(entries, PageRankEntry{Node: node, Score: score})
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Score > entries[j].Score
		})
		if len(entries) > 10 {
			entries = entries[:10]
		}
		stats.PageRankTop10 = entries
	}

	return stats
}

// estimateDiameter estimates the diameter of the graph.
func estimateDiameter(g *graph.Graph) int {
	if g.NodeCount() == 0 {
		return 0
	}

	maxDist := 0
	nodes := g.NodeIDs()

	// For small graphs, do full BFS from each node
	limit := 100
	if len(nodes) > limit {
		nodes = nodes[:limit]
	}

	for _, source := range nodes {
		result, err := g.BFS(source)
		if err != nil {
			continue
		}
		for _, d := range result.Dist {
			if d > maxDist {
				maxDist = d
			}
		}
	}

	return maxDist
}

// estimateRadius estimates the radius of the graph.
func estimateRadius(g *graph.Graph) int {
	if g.NodeCount() == 0 {
		return 0
	}

	nodes := g.NodeIDs()
	limit := 100
	if len(nodes) > limit {
		nodes = nodes[:limit]
	}

	minEccentricity := math.MaxInt32
	for _, source := range nodes {
		result, err := g.BFS(source)
		if err != nil {
			continue
		}
		maxDist := 0
		for _, d := range result.Dist {
			if d > maxDist {
				maxDist = d
			}
		}
		if maxDist < minEccentricity {
			minEccentricity = maxDist
		}
	}

	if minEccentricity == math.MaxInt32 {
		return 0
	}
	return minEccentricity
}

// computeClusteringCoefficient computes the average clustering coefficient.
func computeClusteringCoefficient(g *graph.Graph) float64 {
	if g.NodeCount() < 3 {
		return 0
	}

	totalCoeff := 0.0
	count := 0

	for _, id := range g.NodeIDs() {
		neighbors := g.Neighbors(id)
		k := len(neighbors)
		if k < 2 {
			continue
		}

		// Count edges between neighbors
		edgesBetween := 0
		neighborSet := make(map[string]bool)
		for _, n := range neighbors {
			neighborSet[n] = true
		}

		for _, n1 := range neighbors {
			for _, n2 := range g.Neighbors(n1) {
				if neighborSet[n2] && n1 != n2 {
					edgesBetween++
				}
			}
		}

		// Each edge counted twice in undirected graph
		if !g.Directed {
			edgesBetween /= 2
		}

		coeff := float64(edgesBetween) / float64(k*(k-1)/2)
		totalCoeff += coeff
		count++
	}

	if count == 0 {
		return 0
	}
	return totalCoeff / float64(count)
}

// FormatStats formats statistics as a human-readable string.
func FormatStats(stats *Stats) string {
	result := fmt.Sprintf("Graph Statistics\n")
	result += fmt.Sprintf("================\n")
	result += fmt.Sprintf("Nodes:          %d\n", stats.NodeCount)
	result += fmt.Sprintf("Edges:          %d\n", stats.EdgeCount)
	result += fmt.Sprintf("Density:        %.4f\n", stats.Density)
	result += fmt.Sprintf("Avg Degree:     %.2f\n", stats.AvgDegree)
	result += fmt.Sprintf("Max Degree:     %d\n", stats.MaxDegree)
	result += fmt.Sprintf("Min Degree:     %d\n", stats.MinDegree)
	result += fmt.Sprintf("Connected:      %v\n", stats.IsConnected)
	result += fmt.Sprintf("Components:     %d\n", stats.NumComponents)
	result += fmt.Sprintf("Has Cycles:     %v\n", stats.HasCycles)
	result += fmt.Sprintf("Is DAG:         %v\n", stats.IsDAG)
	result += fmt.Sprintf("Diameter:       %d\n", stats.DiameterEstimate)
	result += fmt.Sprintf("Radius:         %d\n", stats.RadiusEstimate)
	result += fmt.Sprintf("Clustering Coeff: %.4f\n", stats.ClusteringCoeff)

	if len(stats.PageRankTop10) > 0 {
		result += fmt.Sprintf("\nTop PageRank Nodes:\n")
		for i, entry := range stats.PageRankTop10 {
			result += fmt.Sprintf("  %2d. %-20s %.6f\n", i+1, entry.Node, entry.Score)
		}
	}

	if len(stats.Components) > 0 && stats.NumComponents > 1 {
		result += fmt.Sprintf("\nConnected Components:\n")
		for i, comp := range stats.Components {
			result += fmt.Sprintf("  Component %d: %v\n", i+1, comp)
		}
	}

	return result
}
