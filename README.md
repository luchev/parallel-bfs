# Parallel Breadth First Traversal using Go

Implementation of a parallel BFS with level barrier using Go.

## Usage

If you don't know where to start, run `go run bfs/bfs.go --help`. This will give you a list of command line arguments which you can use.

## Examples

Generate a `directed` graph with `100 vertices` and `30% density` (30% chance to generate an edge between 2 vertices). Output the resulting graph in a file `myGraph.graph`.

```go
go run bfs.go -g -v 100 -o 'myGraph' -d 30 -q -directed
```

Run a Breadth first traversal on graph read from file `myGraphFile.graph` using 64 threads.

```go
go run bfs.go -i 'myGraphFile.graph' -t 64 
```

## Performance analysis

[Test results](https://github.com/luchev/parallel-bfs/tree/master/Test%20results) contains the results from running tests on a 16-core machine.

[test/](https://github.com/luchev/parallel-bfs/tree/master/test) contains bash scripts to assist running multiple tests with different number of threads.

[Analysis.pdf](https://github.com/luchev/parallel-bfs/blob/master/Analysis.pdf) presents the performance results of this parallel BFS algorithm implementation.

