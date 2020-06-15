package main

import (
	"bytes"
	"container/list"
	"flag"
	"fmt"
	"io/ioutil"
	golog "log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

/* Queue
queue := make([]int, 0)
// Push to the queue
queue = append(queue, 1)
// Top (just get next element, don't remove it)
x = queue[0]
// Discard top element
queue = queue[1:]
// Is empty ?
if len(queue) == 0 {
    fmt.Println("Queue is empty !")
}
*/

var log logger
var prop properties
var computeTime time.Duration

func main() {
	startProgram := time.Now()

	prop = initProperties()

	if prop.generate {
		var graph = generateGraph(prop)
		saveGraph(graph, prop.outputFile)
	} else if prop.vertices != 0 {
		var graph = generateGraph(prop)
		parentSerial := bfsSerial(graph, prop)
		log.info(len(parentSerial))
		parentParallel := bfsParallel(graph, prop)
		log.info(len(parentParallel))
	} else if prop.inputFile != "" {
		// TODO
	} else {
		flag.PrintDefaults()
	}

	log.info("Program execution took ", time.Since(startProgram))
}

type matrixGraph [][]bool
type logger struct {
	outInfo  *golog.Logger
	outError *golog.Logger
}
type properties struct {
	vertices   int
	threads    int
	density    int
	quiet      bool
	generate   bool
	inputFile  string
	outputFile string
}

func saveGraph(graph matrixGraph, fileName string) {
	graphBytes := graph.Bytes()

	startFileSavingTime := time.Now()
	err := ioutil.WriteFile(fileName, graphBytes, 0644)

	if err == nil {
		log.info("Graph saved as ", fileName)
		log.verbose("Writing graph to disk took ", time.Since(startFileSavingTime))
	} else {
		log.err("Cannot save graph as ", fileName)
	}

}

func bfsSerial(graph matrixGraph, prop properties) (parent []int) {
	startBFS := time.Now()
	parent = make([]int, prop.vertices)
	for i := range parent {
		parent[i] = -1
	}
	visited := make([]bool, prop.vertices)
	for index, isVisited := range visited {
		if !isVisited {
			bfsSerialFromNode(graph, parent, visited, index)
		}
	}

	log.verbose("Serial BFS using 1 thread took ", time.Since(startBFS))
	return parent
}

func bfsSerialFromNode(graph matrixGraph, parent []int, visited []bool, start int) {
	queue := list.New()
	queue.PushBack(start)
	vertices := len(graph)
	for queue.Len() != 0 {
		currentVertex := queue.Front().Value.(int)
		queue.Remove(queue.Front())
		for i := 0; i < vertices; i++ {
			if graph[currentVertex][i] {
				if !visited[i] {
					queue.PushBack(i)
					visited[i] = true
					parent[i] = currentVertex
				}
			}
		}
	}
}

func bfsParallel(graph matrixGraph, prop properties) (parent []int) {
	startBFS := time.Now()

	visited := make([]bool, prop.vertices)
	parent = make([]int, prop.vertices)
	for i := range parent {
		parent[i] = -1
	}

	var bfsWG sync.WaitGroup
	for threadID := 0; threadID < prop.threads; threadID++ {
		bfsWG.Add(1)
		go bfsParallelWorker(graph, prop, &bfsWG, threadID, parent, visited)
	}
	bfsWG.Wait()

	log.verbose("Parallel BFS using ", prop.threads, " threads took ", time.Since(startBFS))
	return parent
}

func bfsParallelWorker(graph matrixGraph, prop properties, bfsWG *sync.WaitGroup, id int, parent []int, visited []bool) {
	defer bfsWG.Done()

	startWorker := time.Now()
	log.verbose("Starting graph traversal worker-", id)

	for currentVertex := id; currentVertex < prop.vertices; currentVertex += prop.threads {
		for i := 0; i < prop.vertices; i++ {
			if graph[currentVertex][i] {
				if !visited[i] {
					visited[i] = true
					parent[i] = currentVertex
				}
			}
		}
	}
	log.verbose("Graph traversal worker-", id, " took ", time.Since(startWorker))
}

func (t *logger) info(args ...interface{}) {
	if t.outInfo == nil {
		t.outInfo = golog.New(os.Stdout, "", 0)
		t.outInfo.SetPrefix("[INFO] ")
	}
	t.outInfo.Print(args...)
}

func (t *logger) err(args ...interface{}) {
	if t.outError == nil {
		t.outError = golog.New(os.Stderr, "", 0)
		t.outError.SetPrefix("[ERR] ")
	}
	t.outError.Print(args...)
}

func (t *logger) verbose(args ...interface{}) {
	if prop.quiet {
		return
	}
	if t.outError == nil {
		t.outError = golog.New(os.Stdout, "", 0)
		t.outError.SetPrefix("[INFO] ")
	}
	t.outError.Print(args...)
}

func (graph matrixGraph) String() string {
	return string(graph.Bytes())
}

// Bytes converts a graph to a byte string
func (graph matrixGraph) Bytes() []byte {
	startConvertingGraphToBytes := time.Now()

	buffer := bytes.Buffer{}
	vertices := len(graph)

	buffer.WriteString(fmt.Sprintf("%d\n", vertices))
	for row := 0; row < len(graph); row++ {
		for col := 0; col < len(graph[row]); col++ {
			if graph[row][col] {
				buffer.WriteByte('1')
			} else {
				buffer.WriteByte('0')
			}
			buffer.WriteByte(' ')
		}
		buffer.WriteByte('\n')
	}

	log.verbose("Converting graph to bytes took ", time.Since(startConvertingGraphToBytes))
	return buffer.Bytes()
}

func generateGraph(prop properties) matrixGraph {
	startMemoryAllocation := time.Now()
	graph := make([][]bool, prop.vertices)
	for row := 0; row < cap(graph); row++ {
		graph[row] = make([]bool, prop.vertices)
	}
	log.verbose("Memory alocation took ", time.Since(startMemoryAllocation))

	startGraphGeneration := time.Now()
	var graphGenerateWG sync.WaitGroup
	for threadID := 0; threadID < prop.threads; threadID++ {
		graphGenerateWG.Add(1)
		go generateGraphWorker(graph, prop, &graphGenerateWG, threadID)
	}
	graphGenerateWG.Wait()
	log.verbose("Graph generation took ", time.Since(startGraphGeneration))
	return graph
}

func generateGraphWorker(graph matrixGraph, prop properties, graphGenerateWG *sync.WaitGroup, id int) {
	defer graphGenerateWG.Done()

	startWorker := time.Now()
	log.verbose("Starting graph generating worker-", id)

	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	for row := id; row < prop.vertices; row += prop.threads {
		var b []byte = make([]byte, prop.vertices)
		generator.Read(b)
		for k := 0; k < prop.vertices; k++ {
			graph[row][k] = int(b[k]) < prop.density
		}
	}
	log.verbose("Graph generating worker-", id, " took ", time.Since(startWorker))
}

func initProperties() properties {
	var prop properties
	vertices := flag.Uint("v", 0, "Graph vertices")
	threads := flag.Uint("t", 0, "Threads (0 to use all cpu cores)")
	density := flag.Uint("d", 20, "Graph density in percent (0-100)")
	flag.BoolVar(&prop.quiet, "q", false, "Run quietly")
	flag.BoolVar(&prop.generate, "g", false, "Generate graph only")
	flag.StringVar(&prop.inputFile, "i", "", "Read graph from file")
	flag.StringVar(&prop.outputFile, "o", "graph.out", "Output file")

	flag.Parse()

	prop.vertices = int(*vertices)
	prop.threads = int(*threads)
	prop.density = int(*density)

	// Automatically determine number of threads
	if prop.threads == 0 {
		prop.threads = runtime.NumCPU()
	}
	// Scale the interval of percents 0-100 to the interval of byte 0-255
	// because random generated numbers are in the interval 0-255
	prop.density = int(math.Ceil(float64(prop.density) * 2.56))

	return prop
}
