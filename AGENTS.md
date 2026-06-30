# AGENTS.md

## Project Overview

GraphPulse is a lightweight graph analysis CLI tool and Go library. It provides comprehensive graph algorithms, statistics, and visualization.

## Architecture

- `graph/` - Core data structures (Graph, Node, Edge) and algorithms (BFS, DFS, Dijkstra, PageRank, etc.)
- `io/` - File I/O for edge lists, CSV, and JSON formats
- `analyze/` - Graph statistics and analysis (degree distribution, clustering coefficient, etc.)
- `visualize/` - ASCII art rendering and Graphviz DOT export
- `cmd/graphpulse/` - CLI entry point using Cobra

## Development

```bash
# Run all tests
go test ./...

# Build
go build -o graphpulse ./cmd/graphpulse/

# Run specific tests
go test ./graph/ -v
go test ./io/ -v
go test ./analyze/ -v
```

## Testing

- All packages have test coverage
- Run `go test ./...` before pushing
- Integration test: build CLI and run against sample graphs

## Key Design Decisions

- Adjacency list representation for O(1) edge lookup
- Separate reverse adjacency list for efficient in-degree queries
- Path compression in Union-Find for near-constant amortized time
- PageRank uses power iteration with configurable damping factor
