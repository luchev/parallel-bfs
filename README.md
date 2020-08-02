# Parallel Breadth First Traversal using Go

Implementation of a parallel BFS with level barrier using Go. The algorithm outputs the parent function of the resulting spanning tree.

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

## Sources

[A scalable distributed parallel breadth-first search algorithm on BlueGene/L.](http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.1075.3533&rank=1), Yoo Andy, et al. Proceedings of the 2005 ACM/IEEE conference on Supercomputing. IEEE Computer Society, 2005.

[Level-Synchronous Parallel Breadth-First Search Algorithms For Multicore and Multiprocessor Systems](https://www.semanticscholar.org/paper/Level-Synchronous-Parallel-Breadth-First-Search-For-Berrendorf-Makulla/cde0420a117f8643d066cdcd60c95d5ca39a1082), Rudolf, and Mathias Makulla FC 14 (2014)

["Parallel breadth-first search on distributed memory systems."](http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.392.7457&rank=1), Buluç, Aydin, and Kamesh Madduri. Proceedings of 2011 International Conference for High Performance Computing, Networking, Storage and Analysis. ACM, 2011.

[A Tale of BFS: Going Parallel](https://github.com/egonelbre/a-tale-of-bfs), Egon Elbre

[Will Hyper-Threading Improve Processing Performance?](https://medium.com/@ITsolutions/will-hyper-threading-improve-processing-performance-15cba11add74), Bill Jones, Sr. Solution Architect, Dasher Technologies
