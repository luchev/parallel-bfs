package main

import (
	"bufio"
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
	"strconv"
	"strings"
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
		var graph = generateDirectedGraph()
		saveGraph(graph)
	} else if prop.vertices != 0 {
		var graph = generateDirectedGraph()
		saveGraph(graph)
		parentFunction := parallelTraversal(graph)
		saveTraversalParentArray(parentFunction)
	} else if prop.inputFile != "" {
		graph := readGraphFromFile()
		parentFunction := parallelTraversal(graph)
		saveTraversalParentArray(parentFunction)
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

func saveTraversalParentArray(parent []int) {
	file, err := os.OpenFile(prop.outputFile+".parent", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer file.Close()
	if err != nil {
		log.err("Failed to save graph traversal result in ", prop.outputFile)
		return
	}
	for start, end := range parent {
		file.WriteString(fmt.Sprint(start, end, "\n"))
	}
}

func saveGraph(graph matrixGraph) {
	fileName := prop.outputFile + ".graph"
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

func bfsSerial(graph matrixGraph) (parent []int) {
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

func parallelTraversal(graph matrixGraph) (parent []int) {
	startBFS := time.Now()

	visited := make([]bool, prop.vertices)
	parent = make([]int, prop.vertices)
	for i := range parent {
		parent[i] = -1
	}

	var bfsWG sync.WaitGroup
	for threadID := 0; threadID < prop.threads; threadID++ {
		bfsWG.Add(1)
		go parallelTraversalWorker(graph, &bfsWG, threadID, parent, visited)
	}
	bfsWG.Wait()

	log.verbose("Parallel BFS using ", prop.threads, " threads took ", time.Since(startBFS))
	return parent
}

func parallelTraversalWorker(graph matrixGraph, bfsWG *sync.WaitGroup, id int, parent []int, visited []bool) {
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

func (t *logger) fatal(args ...interface{}) {
	if t.outError == nil {
		t.outError = golog.New(os.Stderr, "", 0)
		t.outError.SetPrefix("[ERR] ")
	}
	t.outError.Print(args...)
	os.Exit(1)
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

// Bytes converts a graph to a byte array suitable to be written to a file
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

func generateDirectedGraph() matrixGraph {
	graph := allocateGraphMemory(prop.vertices)

	startGraphGeneration := time.Now()
	var graphGenerateWG sync.WaitGroup
	for threadID := 0; threadID < prop.threads; threadID++ {
		graphGenerateWG.Add(1)
		go generateDirectedGraphWorker(graph, &graphGenerateWG, threadID)
	}
	graphGenerateWG.Wait()
	log.verbose("Graph generation took ", time.Since(startGraphGeneration))
	return graph
}

func generateDirectedGraphWorker(graph matrixGraph, graphGenerateWG *sync.WaitGroup, id int) {
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

func readGraphFromFile() matrixGraph {
	file, err := os.OpenFile(prop.inputFile, os.O_RDONLY, 0644)
	if err != nil {
		log.fatal("Cannot open file ", prop.inputFile)
	}

	reader := bufio.NewReader(file)
	line, _ := reader.ReadString('\n')
	file.Close()
	vertices, err := strconv.Atoi(strings.Trim(line, "\n"))
	if err != nil {
		log.fatal(prop.inputFile, " has incorrect format")
	}
	prop.vertices = vertices

	graph := allocateGraphMemory(prop.vertices)

	startGraphReading := time.Now()
	var graphReadWG sync.WaitGroup
	for threadID := 0; threadID < prop.threads; threadID++ {
		graphReadWG.Add(1)
		go readGraphWorker(graph, &graphReadWG, threadID)
	}
	graphReadWG.Wait()
	log.verbose("Reading graph from file took ", time.Since(startGraphReading))
	return graph
}

func readGraphWorker(graph matrixGraph, readGraphWG *sync.WaitGroup, id int) {
	defer readGraphWG.Done()

	startWorker := time.Now()
	log.verbose("Starting graph reading worker-", id)

	firstLineLength := len(strconv.FormatInt(int64(prop.vertices), 10)) + 1
	file, err := os.OpenFile(prop.inputFile, os.O_RDONLY, 0644)
	if err != nil {
		log.fatal("Failed to open ", prop.inputFile)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineLength := prop.vertices*2 + 1
	seekLength := int64(lineLength * (prop.threads - 1))
	workerStartingSeek := id * lineLength
	// Seek to the starting line for this worker
	file.Seek(int64(firstLineLength), 0)
	file.Seek(int64(workerStartingSeek), 1)

	for row := id; row < prop.vertices; row += prop.threads {
		line, err := reader.ReadBytes('\n')
		if err != nil || len(line) != lineLength {
			log.fatal(prop.inputFile, " has incorrect format")
		}
		// Parse line
		for k := 0; k < prop.vertices; k++ {
			graph[row][k] = int(line[2*k]) == '1'
		}
		// Seek to the next line to be parsed by this worker
		reader.Discard(int(seekLength))
	}

	log.verbose("Graph reading worker-", id, " took ", time.Since(startWorker))
}

func allocateGraphMemory(vertices int) matrixGraph {
	startMemoryAllocation := time.Now()
	graph := make([][]bool, prop.vertices)
	for row := 0; row < cap(graph); row++ {
		graph[row] = make([]bool, prop.vertices)
	}
	log.verbose("Memory alocation took ", time.Since(startMemoryAllocation))
	return graph
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