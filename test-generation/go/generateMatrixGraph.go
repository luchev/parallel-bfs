package main

import (
	crand "crypto/rand"
	"fmt"
	rand "math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

const v = 20000

var edges [v][v]bool

func generateMath(wg *sync.WaitGroup, id int, threads int) {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	for row := id; row < v; row += threads {
		if row >= v {
			break
		}
		for k := 0; k < v; k++ {
			edges[row][k] = generator.Intn(3) == 1
		}
	}
	defer wg.Done()
}
func generateMathArray(wg *sync.WaitGroup, id int, threads int) {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	for row := id; row < v; row += threads {
		if row >= v {
			break
		}
		var b []byte = make([]byte, v)
		generator.Read(b)
		for k := 0; k < v; k++ {
			edges[row][k] = b[k]&1 == 1
		}
	}
	defer wg.Done()
}

func generateCrypto(wg *sync.WaitGroup, id int, threads int) {
	for row := id; row < v; row += threads {
		if row >= v {
			break
		}
		var b []byte = make([]byte, v)
		crand.Read(b)
		for k := 0; k < v; k++ {
			edges[row][k] = b[k]&1 == 1
		}
	}
	defer wg.Done()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide number of threads as argument")
		return
	}

	threads, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Provide number of threads as argument")
		return
	}

	fmt.Println("Running with ", threads, " threads")
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		// go generateCrypto(&wg, i, threads)
		// go generateMath(&wg, i, threads)
		go generateMathArray(&wg, i, threads)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Execution time %s\n", elapsed)
}
