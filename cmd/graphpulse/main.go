// Package main provides the GraphPulse CLI entry point.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/graphpulse/analyze"
	gio "github.com/EdgarOrtegaRamirez/graphpulse/io"
	"github.com/EdgarOrtegaRamirez/graphpulse/visualize"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "graphpulse",
		Short: "GraphPulse - Lightweight Graph Analysis CLI",
		Long: `GraphPulse is a lightweight, fast graph analysis tool.
It supports multiple input formats and provides comprehensive
graph algorithms, statistics, and visualization.`,
		Version: version,
	}

	// Stats command
	var statsDirected bool
	var statsJSON bool
	statsCmd := &cobra.Command{
		Use:   "stats [file]",
		Short: "Compute graph statistics",
		Long:  `Compute comprehensive statistics for a graph including degree distribution, connectivity, cycles, PageRank, and more.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], statsDirected)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			stats := analyze.ComputeStats(g)
			if statsJSON {
				fmt.Println(analyze.FormatStatsJSON(stats))
			} else {
				fmt.Print(analyze.FormatStats(stats))
			}
			return nil
		},
	}
	statsCmd.Flags().BoolVar(&statsDirected, "directed", false, "Treat graph as directed")
	statsCmd.Flags().BoolVar(&statsJSON, "json", false, "Output statistics as JSON")
	rootCmd.AddCommand(statsCmd)

	// BFS command
	var bfsDirected bool
	bfsCmd := &cobra.Command{
		Use:   "bfs [file] [source]",
		Short: "Breadth-first search",
		Long:  `Perform BFS from a source node and show visit order and distances.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], bfsDirected)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			result, err := g.BFS(args[1])
			if err != nil {
				return fmt.Errorf("BFS failed: %w", err)
			}

			fmt.Printf("BFS from %s\n", args[1])
			fmt.Printf("Visit order: %s\n", strings.Join(result.Order, " → "))
			fmt.Printf("Distances:\n")
			for _, id := range g.NodeIDs() {
				if d, ok := result.Dist[id]; ok {
					fmt.Printf("  %s: %d\n", id, d)
				}
			}
			return nil
		},
	}
	bfsCmd.Flags().BoolVar(&bfsDirected, "directed", false, "Treat graph as directed")
	rootCmd.AddCommand(bfsCmd)

	// DFS command
	var dfsDirected bool
	dfsCmd := &cobra.Command{
		Use:   "dfs [file] [source]",
		Short: "Depth-first search",
		Long:  `Perform DFS from a source node and show discovery/finish times.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], dfsDirected)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			result, err := g.DFS(args[1])
			if err != nil {
				return fmt.Errorf("DFS failed: %w", err)
			}

			fmt.Printf("DFS from %s\n", args[1])
			fmt.Printf("Discovery order: %s\n", strings.Join(result.Order, " → "))
			fmt.Printf("Finish order:    %s\n", strings.Join(result.Finished, " → "))
			fmt.Printf("Discovery/Finish times:\n")
			for _, id := range g.NodeIDs() {
				fmt.Printf("  %s: %d/%d\n", id, result.Discovery[id], result.Finish[id])
			}
			return nil
		},
	}
	dfsCmd.Flags().BoolVar(&dfsDirected, "directed", false, "Treat graph as directed")
	rootCmd.AddCommand(dfsCmd)

	// Shortest path command
	var spDirected bool
	spCmd := &cobra.Command{
		Use:   "shortest [file] [source] [target]",
		Short: "Find shortest path",
		Long:  `Find the shortest path between two nodes using Dijkstra's algorithm.`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], spDirected)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			result, err := g.Dijkstra(args[1])
			if err != nil {
				return fmt.Errorf("Dijkstra failed: %w", err)
			}

			dist := result.Dist[args[2]]
			path := result.Path[args[2]]

			if path == nil {
				fmt.Printf("No path from %s to %s\n", args[1], args[2])
			} else {
				fmt.Printf("Shortest path: %s\n", strings.Join(path, " → "))
				fmt.Printf("Distance: %.2f\n", dist)
			}
			return nil
		},
	}
	spCmd.Flags().BoolVar(&spDirected, "directed", false, "Treat graph as directed")
	rootCmd.AddCommand(spCmd)

	// PageRank command
	var prDirected bool
	var prIterations int
	var prDamping float64
	prCmd := &cobra.Command{
		Use:   "pagerank [file]",
		Short: "Compute PageRank scores",
		Long:  `Compute PageRank scores for all nodes in the graph.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], prDirected)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			ranks := g.PageRank(prDamping, prIterations)

			fmt.Println(visualize.RenderPageRankBar(g, 20))
			fmt.Printf("\nAll scores:\n")
			for _, id := range g.NodeIDs() {
				fmt.Printf("  %-20s %.6f\n", id, ranks[id])
			}
			return nil
		},
	}
	prCmd.Flags().BoolVar(&prDirected, "directed", false, "Treat graph as directed")
	prCmd.Flags().IntVarP(&prIterations, "iterations", "i", 100, "Number of iterations")
	prCmd.Flags().Float64VarP(&prDamping, "damping", "d", 0.85, "Damping factor")
	rootCmd.AddCommand(prCmd)

	// Toposort command
	tsCmd := &cobra.Command{
		Use:   "toposort [file]",
		Short: "Topological sort",
		Long:  `Perform topological sort on a directed acyclic graph.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], true)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			order, err := g.TopologicalSort()
			if err != nil {
				return fmt.Errorf("topological sort failed: %w", err)
			}

			fmt.Printf("Topological order: %s\n", strings.Join(order, " → "))
			return nil
		},
	}
	rootCmd.AddCommand(tsCmd)

	// Cycle detection command
	var cycleDirected bool
	cycleCmd := &cobra.Command{
		Use:   "cycle [file]",
		Short: "Detect cycles",
		Long:  `Detect if the graph contains a cycle and show the cycle path.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], cycleDirected)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			hasCycle, path := g.DetectCycle()
			if hasCycle {
				fmt.Printf("Cycle detected: %s\n", strings.Join(path, " → "))
			} else {
				fmt.Println("No cycle detected")
			}
			return nil
		},
	}
	cycleCmd.Flags().BoolVar(&cycleDirected, "directed", false, "Treat graph as directed")
	rootCmd.AddCommand(cycleCmd)

	// Components command
	var compDirected bool
	compCmd := &cobra.Command{
		Use:   "components [file]",
		Short: "Find connected components",
		Long:  `Find connected components (undirected) or strongly connected components (directed).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], compDirected)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			var components [][]string
			if g.Directed {
				components = g.StronglyConnectedComponents()
				fmt.Printf("Strongly Connected Components (%d):\n", len(components))
			} else {
				components = g.ConnectedComponents()
				fmt.Printf("Connected Components (%d):\n", len(components))
			}

			for i, comp := range components {
				fmt.Printf("  %d: %v\n", i+1, comp)
			}
			return nil
		},
	}
	compCmd.Flags().BoolVar(&compDirected, "directed", false, "Treat graph as directed")
	rootCmd.AddCommand(compCmd)

	// Visualize command
	var vizDirected bool
	var vizCompact bool
	var vizWeights bool
	vizCmd := &cobra.Command{
		Use:   "visualize [file]",
		Short: "Visualize graph as ASCII art",
		Long:  `Render the graph as ASCII art in the terminal.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], vizDirected)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			config := visualize.DefaultASCIIConfig()
			config.Compact = vizCompact
			config.ShowWeights = vizWeights

			fmt.Println(visualize.RenderASCII(g, config))
			return nil
		},
	}
	vizCmd.Flags().BoolVar(&vizDirected, "directed", false, "Treat graph as directed")
	vizCmd.Flags().BoolVarP(&vizCompact, "compact", "c", false, "Use compact layout")
	vizCmd.Flags().BoolVarP(&vizWeights, "weights", "w", true, "Show edge weights")
	rootCmd.AddCommand(vizCmd)

	// DOT export command
	var dotDirected bool
	var dotTitle string
	dotCmd := &cobra.Command{
		Use:   "dot [file]",
		Short: "Export as Graphviz DOT format",
		Long:  `Export the graph in Graphviz DOT format for visualization with graphviz tools.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], dotDirected)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			config := visualize.DefaultDOTConfig()
			config.Directed = dotDirected
			config.Title = dotTitle

			fmt.Println(visualize.ExportDOT(g, config))
			return nil
		},
	}
	dotCmd.Flags().BoolVar(&dotDirected, "directed", false, "Treat graph as directed")
	dotCmd.Flags().StringVarP(&dotTitle, "title", "t", "", "Graph title")
	rootCmd.AddCommand(dotCmd)

	// Histogram command
	var histDirected bool
	histCmd := &cobra.Command{
		Use:   "histogram [file]",
		Short: "Show degree distribution histogram",
		Long:  `Render an ASCII histogram of the degree distribution.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], histDirected)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			fmt.Println(visualize.RenderDegreeHistogram(g))
			return nil
		},
	}
	histCmd.Flags().BoolVar(&histDirected, "directed", false, "Treat graph as directed")
	rootCmd.AddCommand(histCmd)

	// Info command
	infoCmd := &cobra.Command{
		Use:   "info [file]",
		Short: "Show basic graph info",
		Long:  `Show basic information about a graph.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := gio.LoadGraph(args[0], false)
			if err != nil {
				return fmt.Errorf("failed to load graph: %w", err)
			}

			fmt.Printf("Nodes:    %d\n", g.NodeCount())
			fmt.Printf("Edges:    %d\n", g.EdgeCount())
			fmt.Printf("Directed: %v\n", g.Directed)
			fmt.Printf("\nNodes:\n")
			for _, id := range g.NodeIDs() {
				inDeg := g.InDegree(id)
				outDeg := g.OutDegree(id)
				if g.Directed {
					fmt.Printf("  %s (in=%d, out=%d)\n", id, inDeg, outDeg)
				} else {
					fmt.Printf("  %s (degree=%d)\n", id, outDeg)
				}
			}
			return nil
		},
	}
	rootCmd.AddCommand(infoCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
