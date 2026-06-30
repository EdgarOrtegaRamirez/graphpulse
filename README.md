# GraphPulse

A lightweight, fast graph analysis CLI tool and Go library. Analyze graph data with comprehensive algorithms, statistics, and visualization — no heavyweight dependencies needed.

## Features

- **Multiple Input Formats**: Edge lists, CSV, JSON
- **Graph Algorithms**: BFS, DFS, Dijkstra's shortest path, topological sort, cycle detection
- **Analysis**: Degree distribution, clustering coefficient, PageRank, connected components, diameter/radius estimation
- **Visualization**: ASCII art rendering, degree histograms, PageRank bar charts
- **Export**: Graphviz DOT format for professional graph visualization
- **Undirected & Directed**: Full support for both graph types

## Quick Start

```bash
# Install
go install github.com/EdgarOrtegaRamirez/graphpulse/cmd/graphpulse@latest

# Or build from source
git clone https://github.com/EdgarOrtegaRamirez/graphpulse.git
cd graphpulse
go build -o graphpulse ./cmd/graphpulse/
```

## Usage

### Create a graph file (edge list format)

```
# edges.txt
a,b
b,c
c,d
d,a
a,c
```

### Compute Statistics

```bash
graphpulse stats edges.txt
# Output:
# Graph Statistics
# ================
# Nodes:          4
# Edges:          5
# Density:        0.4167
# ...
```

### Find Shortest Path

```bash
graphpulse shortest weighted.txt a d --directed
# Shortest path: a → b → c → d
# Distance: 4.00
```

### Detect Cycles

```bash
graphpulse cycle edges.txt
# Cycle detected: a → d
```

### Topological Sort (DAGs)

```bash
graphpulse toposort dag.txt
# Topological order: a → b → c → d
```

### PageRank

```bash
graphpulse pagerank edges.txt
```

### Visualize as ASCII Art

```bash
graphpulse visualize edges.txt
```

### Export to Graphviz DOT

```bash
graphpulse dot edges.txt --directed --title "My Graph" > graph.dot
dot -Tpng graph.dot -o graph.png
```

## Input Formats

### Edge List (default)
```
a,b
b,c,2.5    # with weight
```

### CSV
```csv
source,target,weight
a,b,1.0
b,c,2.0
```

### JSON
```json
{
  "directed": true,
  "nodes": [{"id": "a"}, {"id": "b"}],
  "edges": [{"from": "a", "to": "b", "weight": 1.0}]
}
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `stats` | Compute graph statistics |
| `bfs` | Breadth-first search |
| `dfs` | Depth-first search |
| `shortest` | Dijkstra's shortest path |
| `pagerank` | Compute PageRank scores |
| `toposort` | Topological sort |
| `cycle` | Detect cycles |
| `components` | Find connected/strongly connected components |
| `visualize` | ASCII art visualization |
| `dot` | Export to Graphviz DOT format |
| `histogram` | Degree distribution histogram |
| `info` | Basic graph info |

## Library Usage

```go
package main

import (
    "fmt"
    "github.com/EdgarOrtegaRamirez/graphpulse/graph"
)

func main() {
    g := graph.New(true)
    g.AddEdge("a", "b", 1.0, nil)
    g.AddEdge("b", "c", 2.0, nil)
    g.AddEdge("c", "a", 1.0, nil)

    // BFS
    result, _ := g.BFS("a")
    fmt.Println("BFS order:", result.Order)

    // Dijkstra
    dij, _ := g.Dijkstra("a")
    fmt.Println("Distance a→c:", dij.Dist["c"])

    // PageRank
    ranks := g.PageRank(0.85, 100)
    fmt.Println("PageRank:", ranks)

    // Cycle detection
    hasCycle, cycle := g.DetectCycle()
    fmt.Println("Has cycle:", hasCycle, "Path:", cycle)
}
```

## Algorithms

| Algorithm | Complexity | Description |
|-----------|-----------|-------------|
| BFS | O(V+E) | Breadth-first traversal |
| DFS | O(V+E) | Depth-first traversal |
| Dijkstra | O((V+E) log V) | Shortest path (weighted) |
| Topological Sort | O(V+E) | DAG ordering (Kahn's) |
| Cycle Detection | O(V+E) | DFS-based cycle detection |
| PageRank | O(k·(V+E)) | Link analysis (k iterations) |
| Union-Find | O(α(n)) | Disjoint set (near constant) |
| Tarjan's SCC | O(V+E) | Strongly connected components |
| Clustering Coefficient | O(V·d²) | Triangle counting |

## Architecture

```
graphpulse/
├── graph/          # Core data structures & algorithms
│   ├── graph.go    # Graph type, node/edge management
│   └── algorithms.go  # BFS, DFS, Dijkstra, PageRank, etc.
├── io/             # File I/O (edge list, CSV, JSON)
├── analyze/        # Statistics & analysis
├── visualize/      # ASCII & DOT rendering
└── cmd/graphpulse/ # CLI entry point
```

## License

MIT
