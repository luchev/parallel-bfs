#include <cstdlib>
#include <iostream>
#include <thread>
#include <random>
#include <stdio.h>
#include <stdlib.h>
#include <chrono>
using namespace std;
using namespace std::chrono;

#define V 50000

bool edges[V][V];

struct thread_data {
    int thread_id;
    int threadCount;
};

bool Mersenne() {
    static std::random_device dev;
    std::mt19937 rng(dev());
    std::uniform_int_distribution<std::mt19937::result_type> dist6(0, 1);

    return dist6(rng);
}

bool MersenneThreaded() {
    static thread_local std::mt19937 generator;
    std::uniform_int_distribution<int> distribution(0, 1);
    return distribution(generator);
}

static unsigned long x = 123456789, y = 362436069, z = 521288629;
bool Marsaglia(void) {
    unsigned long t;
    x ^= x << 16;
    x ^= x >> 5;
    x ^= x << 1;
    t = x;
    x = y;
    y = z;
    z = t ^ x ^ y;

    return z & 1;
}

bool MarsagliaThreaded(void) {
    static thread_local unsigned long x = 123456789, y = 362436069, z = 521288629;
    unsigned long t;
    x ^= x << 16;
    x ^= x >> 5;
    x ^= x << 1;
    t = x;
    x = y;
    y = z;
    z = t ^ x ^ y;

    return z & 1;
}

void work(int id, int threadCount) {
    for (int row = id; row < V; row += threadCount) {
        if (row >= V) {
            break;
        }
        for (int i = 0; i < V; i++) {
            // edges[row][i] = Mersenne();
            // edges[row][i] = MersenneThreaded();
            // edges[row][i] = Marsaglia();
            edges[row][i] = MarsagliaThreaded();
        }
    }
}

int main(int argc, char * argv[]) {
    if (argc < 2) {
        cout << "Provide number of threads as argument\n";
        return(1);
    }

    int threadCount = atoi(argv[1]);
    cout << "Generating " << V << "x" << V << " graph, using " << threadCount << " threads.\n";
    auto start = high_resolution_clock::now(); 

    vector<thread> threads(threadCount);

    for (int i = 0; i < threadCount; i++) {
        threads[i] = thread(work, i, threadCount);
    }

    for (int i = 0; i < threadCount; i++) {
        threads[i].join();
    }
    
    auto stop = high_resolution_clock::now(); 
    auto duration = duration_cast<microseconds>(stop - start); 
    cout << duration.count() / 1000000.0 << endl;
}
