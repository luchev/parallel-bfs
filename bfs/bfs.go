package main

import (
	"bufio"
	"bytes"
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

var log logger
var prop properties
var computeTime time.Duration

func main() {
	startProgram := time.Now()
	prop = initProperties()

	// Create graph
	var graph matrixGraph
	if prop.generate || prop.vertices != 0 {
		if prop.directed {
			graph = generateDirectedGraph()
		} else {
			graph = generateUndirectedGraph()
		}
		saveGraphSerial(graph)
		saveGraphParallel(graph)
	} else if prop.inputFile != "" {
		graph = readGraphFromFile()
	} else {
		log.fatal("Either -v or -i should be specified")
	}

	if prop.generate {
		os.Exit(0)
	}

	// Traversal
	parentFunction := bfsSerial(graph)
	parentFunction = bfsLevelBarrier(graph)
	parentFunction = parallelTraversal(graph)

	saveTraversalParentArray(parentFunction)
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
	directed   bool
	inputFile  string
	outputFile string
}

func bfsLevelBarrier(graph matrixGraph) (parent []int) {
	log.verbose("Starting parallel BFS with level barrier using ", prop.threads, " threads")
	startBFS := time.Now()
	threadTimes := make([]time.Duration, prop.threads)

	parent = make([]int, prop.vertices)
	for i := range parent {
		parent[i] = -1
	}
	visited := make([]bool, prop.vertices)
	for index, isVisited := range visited {
		if !isVisited {
			bfsLevelBarrierFromVertex(graph, parent, visited, index, threadTimes)
		}
	}

	for workerID, duration := range threadTimes {
		log.verbose("Parallel BFS worker-", workerID, " took ", duration)
	}

	log.verbose("Parallel BFS with level barrier using ", prop.threads, " threads took ", time.Since(startBFS))
	return parent
}

func bfsLevelBarrierFromVertex(graph matrixGraph, parent []int, visited []bool, start int, threadTimes []time.Duration) {
	// Stores next level from the queue
	futureFrontiers := make([]int, 0)
	futureFrontiers = append(futureFrontiers, start)
	addedNeighbours := make([]bool, prop.vertices)

	for len(futureFrontiers) != 0 {
		// Make channels for communication
		currentFrontiers := make(chan int, prop.vertices)
		nextLevelFrontiers := make(chan int, prop.vertices)

		// Start parallel level traversal workers
		for i := 0; i < prop.threads; i++ {
			go bfsLevelBarrierWorker(graph, parent, visited, currentFrontiers, nextLevelFrontiers, addedNeighbours, i, threadTimes)
		}

		// Initialize queue for vertices from current level
		for _, node := range futureFrontiers {
			currentFrontiers <- node
		}
		close(currentFrontiers)
		futureFrontiers = make([]int, 0)

		// Wait all workers to finish and generate vertices for the next level
		threadsReady := 0
		for vertex := range nextLevelFrontiers {
			if vertex == -1 { // -1 is sent by a thread when it's ready
				threadsReady++
			} else {
				futureFrontiers = append(futureFrontiers, vertex)
			}

			// If all threads have sent the ready signal (-1)
			if threadsReady == prop.threads {
				close(nextLevelFrontiers)
				break
			}
		}
	}
}

func bfsLevelBarrierWorker(graph matrixGraph, parent []int, visited []bool, currentFrontiers <-chan int, nextLevelFrontiers chan<- int, addedNeighbours []bool, workerID int, threadTimes []time.Duration) {
	startWorkerTime := time.Now()
	for vertex := range currentFrontiers {
		if visited[vertex] {
			continue
		}
		visited[vertex] = true
		for neighbour := 0; neighbour < prop.vertices; neighbour++ {
			if !visited[neighbour] && graph[vertex][neighbour] && !addedNeighbours[neighbour] {
				addedNeighbours[neighbour] = true
				parent[neighbour] = vertex
				nextLevelFrontiers <- neighbour
			}
		}
	}
	threadTimes[workerID] += time.Since(startWorkerTime)
	nextLevelFrontiers <- -1
}

func saveTraversalParentArray(parent []int) {
	fileName := prop.outputFile + ".parent"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer file.Close()
	if err != nil {
		log.err("Failed to save graph traversal result in ", fileName)
		return
	}
	for start, end := range parent {
		file.WriteString(fmt.Sprint(start, end, "\n"))
	}
	log.info("Graph traversal result saved as ", fileName)
}

func saveGraphParallel(graph matrixGraph) {
	log.verbose("Starting serializing graph to disk using ", prop.threads, " threads")
	startFileSavingTime := time.Now()

	fileName := prop.outputFile + ".graph"
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		log.fatal("Failed to open ", fileName)
	}
	defer file.Close()

	graphBytes := make([][]byte, prop.vertices)
	rowJobs := make(chan int, prop.vertices)
	readyRowsChannel := make(chan int, prop.vertices)
	for workerID := 0; workerID < prop.threads; workerID++ {
		go graphSerializerWorker(graph, graphBytes, workerID, rowJobs, readyRowsChannel)
	}

	file.WriteString(strconv.Itoa(prop.vertices) + "\n")
	for row := 0; row < prop.vertices; row++ {
		rowJobs <- row
	}
	close(rowJobs)

	readyRows := make([]bool, prop.vertices)
	nextRowToWrite := 0
	for i := 0; i < prop.vertices; i++ {
		readyRows[<-readyRowsChannel] = true
		for ; nextRowToWrite < prop.vertices && readyRows[nextRowToWrite]; nextRowToWrite++ {
			file.Write(graphBytes[nextRowToWrite])
			graphBytes[nextRowToWrite] = nil
		}
	}

	log.info("Graph saved as ", fileName)
	log.verbose("Serializing graph to disk using ", prop.threads, " threads took ", time.Since(startFileSavingTime))
}

func graphSerializerWorker(graph matrixGraph, outputBytes [][]byte, id int, rowJobs <-chan int, readyRows chan<- int) {
	log.verbose("Starting graph serializing worker-", id)
	startGraphSerializerWorkerTime := time.Now()

	for row := range rowJobs {
		var buffer bytes.Buffer
		for col := 0; col < len(graph[row]); col++ {
			if graph[row][col] {
				buffer.WriteByte('1')
			} else {
				buffer.WriteByte('0')
			}
			buffer.WriteByte(' ')
		}
		buffer.WriteByte('\n')
		outputBytes[row] = buffer.Bytes()
		readyRows <- row
	}

	log.verbose("Graph serializing worker-", id, " took ", time.Since(startGraphSerializerWorkerTime))
}

func saveGraphSerial(graph matrixGraph) {
	log.verbose("Starting serializing graph to disk")
	startFileSavingTime := time.Now()

	fileName := prop.outputFile + ".graph"
	graphBytes := graph.Bytes()

	err := ioutil.WriteFile(fileName, graphBytes, 0644)
	if err == nil {
		log.info("Graph saved as ", fileName)
		log.verbose("Serializing graph to disk took ", time.Since(startFileSavingTime))
	} else {
		log.err("Cannot save graph as ", fileName)
	}
}

func bfsSerial(graph matrixGraph) (parent []int) {
	log.verbose("Starting serial BFS")
	startBFS := time.Now()

	parent = make([]int, prop.vertices)
	for i := range parent {
		parent[i] = -1
	}

	visited := make([]bool, prop.vertices)
	for index, isVisited := range visited {
		if !isVisited {
			bfsSerialFromVertex(graph, parent, visited, index)
		}
	}

	log.verbose("Serial BFS took ", time.Since(startBFS))
	return parent
}

func bfsSerialFromVertex(graph matrixGraph, parent []int, visited []bool, start int) {
	queue := make([]int, 0)
	queue = append(queue, start)

	vertices := len(graph)
	for len(queue) != 0 {
		currentVertex := queue[0]
		queue = queue[1:]
		for i := 0; i < vertices; i++ {
			if graph[currentVertex][i] && !visited[i] {
				queue = append(queue, i)
				visited[i] = true
				parent[i] = currentVertex
			}
		}
	}
}

func parallelTraversal(graph matrixGraph) (parent []int) {
	log.verbose("Starting custom parallel traversal using ", prop.threads, " threads")
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

	log.verbose("Custom parallel traversal using ", prop.threads, " threads took ", time.Since(startBFS))
	return parent
}

func parallelTraversalWorker(graph matrixGraph, bfsWG *sync.WaitGroup, id int, parent []int, visited []bool) {
	defer bfsWG.Done()
	log.verbose("Starting custom parallel traversal worker-", id)
	startWorker := time.Now()

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

	log.verbose("Custom parallel traversal worker-", id, " took ", time.Since(startWorker))
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
		t.outError.SetPrefix("[FATAL] ")
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
	log.verbose("Starting serializing graph to bytes")
	startConvertingGraphToBytes := time.Now()

	buffer := bytes.Buffer{}
	vertices := len(graph)

	buffer.WriteString(fmt.Sprintf("%d\n", vertices))
	for row := 0; row < vertices; row++ {
		for col := 0; col < vertices; col++ {
			if graph[row][col] {
				buffer.WriteByte('1')
			} else {
				buffer.WriteByte('0')
			}
			buffer.WriteByte(' ')
		}
		buffer.WriteByte('\n')
	}

	log.verbose("Serializing graph to bytes took ", time.Since(startConvertingGraphToBytes))
	return buffer.Bytes()
}

func generateDirectedGraph() matrixGraph {
	graph := allocateGraphMemory(prop.vertices)
	log.verbose("Starting generating directed graph using ", prop.threads, " threads")
	startGraphGeneration := time.Now()

	var graphGenerateWG sync.WaitGroup
	for threadID := 0; threadID < prop.threads; threadID++ {
		graphGenerateWG.Add(1)
		go generateDirectedGraphWorker(graph, &graphGenerateWG, threadID)
	}
	graphGenerateWG.Wait()

	log.verbose("Generating directed graph using ", prop.threads, " threads took ", time.Since(startGraphGeneration))
	return graph
}

func generateDirectedGraphWorker(graph matrixGraph, graphGenerateWG *sync.WaitGroup, id int) {
	defer graphGenerateWG.Done()
	log.verbose("Starting directed graph generating worker-", id)
	startWorker := time.Now()

	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	for row := id; row < prop.vertices; row += prop.threads {
		var b []byte = make([]byte, prop.vertices)
		generator.Read(b)
		for k := 0; k < prop.vertices; k++ {
			graph[row][k] = int(b[k]) < prop.density
		}
	}
	log.verbose("Generating directed graph worker-", id, " took ", time.Since(startWorker))
}

func generateUndirectedGraph() (graph matrixGraph) {
	graph = allocateGraphMemory(prop.vertices)
	log.verbose("Starting generating undirected graph using ", prop.threads)
	startGraphGeneration := time.Now()

	tasks := make(chan int, prop.vertices)
	var graphGenerateWG sync.WaitGroup
	for i := 0; i < prop.threads; i++ {
		graphGenerateWG.Add(1)
		go generateUnirectedGraphWorker(graph, &graphGenerateWG, tasks, i)
	}

	for vertex := 0; vertex < prop.vertices; vertex++ {
		tasks <- vertex
	}
	close(tasks)

	graphGenerateWG.Wait()
	log.verbose("Generating undirected graph using ", prop.threads, " threads took ", time.Since(startGraphGeneration))
	return graph
}

func generateUnirectedGraphWorker(graph matrixGraph, graphGenerateWG *sync.WaitGroup, tasks <-chan int, id int) {
	defer graphGenerateWG.Done()
	log.verbose("Starting undirected graph generating worker-", id)
	startWorker := time.Now()

	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	for vertex := range tasks {
		var randomBytes []byte = make([]byte, prop.vertices-vertex-1)
		generator.Read(randomBytes)
		graph[vertex][vertex] = true
		for k := vertex + 1; k < prop.vertices; k++ {
			graph[vertex][k] = int(randomBytes[k-vertex-1]) < prop.density
			graph[k][vertex] = int(randomBytes[k-vertex-1]) < prop.density
		}
	}
	log.verbose("Generating undirected graph worker-", id, " took ", time.Since(startWorker))
}

func readGraphFromFile() matrixGraph {
	log.verbose("Starting reading graph from file using ", prop.threads, " threads")
	startGraphReading := time.Now()

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

	var graphReadWG sync.WaitGroup
	for threadID := 0; threadID < prop.threads; threadID++ {
		graphReadWG.Add(1)
		go readGraphWorker(graph, &graphReadWG, threadID)
	}
	graphReadWG.Wait()

	log.verbose("Reading graph from file using ", prop.threads, " threads took ", time.Since(startGraphReading))
	return graph
}

func readGraphWorker(graph matrixGraph, readGraphWG *sync.WaitGroup, id int) {
	defer readGraphWG.Done()
	log.verbose("Starting graph reading worker-", id)
	startWorker := time.Now()

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
	log.verbose("Starting allocating memory for graph with ", prop.vertices, " vertices")
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
	flag.BoolVar(&prop.directed, "directed", false, "Generate directed graph (default is undirected)")
	flag.BoolVar(&prop.generate, "g", false, "Generate graph only")
	flag.StringVar(&prop.inputFile, "i", "", "Read graph from file")
	flag.StringVar(&prop.outputFile, "o", "out", "Output file")

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
