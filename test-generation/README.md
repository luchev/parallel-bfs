# Parallel BFS

## C++

| Mersenne Twister | Threads | Vertices | CPU load % | Thread load % | Time (seconds) | Speedup |
| ---------------- | ------- | -------- | ---------- | ------------- | -------------- | ------- |
| Non-Threaded     | 1       | 10,000   | 100        | 100           | 32 минути      | 1       |
| Non-Threaded     | 16      | 10,000   | 1200       | 80            | 43 минути      | 0.74    |
| Threaded         | 1       | 10,000   | 100        | 100           | 6.87077        | 1       |
| Threaded         | 4       | 10,000   | 400        | 100           | 1.78457        | 3.85    |
| Threaded         | 8       | 10,000   | 800        | 100           | 0.910811       | 7.54    |
| Threaded         | 16      | 10,000   | 1600       | 100           | 0.538676       | 12.75   |

| Marsaglia's xorshf | Threads | Vertices | CPU load % | Thread load % | Time (seconds) | Speedup |
| ------------------ | ------- | -------- | ---------- | ------------- | -------------- | ------- |
| Non-Threaded       | 1       | 20,000   | 100        | 100           | 3.63337        | 1       |
| Non-Threaded       | 4       | 20,000   | 400        | 100           | 10.3713        | 0.35    |
| Threaded           | 1       | 20,000   | 100        | 100           | 3.43346        | 1       |
| Threaded           | 4       | 20,000   | 400        | 100           | 0.904795       | 3.7     |
| Threaded           | 8       | 20,000   | 800        | 100           | 0.473585       | 7.2     |
| Threaded           | 16      | 20,000   | 1600       | 100           | 0.323407       | 10.6    |

| Algorithm (threaded) | Threads | Vertices | CPU load % | Thread load % | Time (seconds) | Speedup |
| -------------------- | ------- | -------- | ---------- | ------------- | -------------- | ------- |
| Mersenne Twister     | 1       | 50,000   | 100        | 100           | 172.077        | 1       |
| Mersenne Twister     | 16      | 50,000   | 1600       | 100           | 13.0528        | 13.1    |
|                      |         |          |            |               |                |         |
| Marsaglia's xorshf   | 1       | 50,000   | 100        | 100           | 21.4722        | 1       |
| Marsaglia's xorshf   | 16      | 50,000   | 1600       | 100           | 1.75111        | 12.2    |

## Go

| Cryptographic pseudo-random number generator (crypto/rand) | Threads | Vertices | CPU load % | Thread load % | Time (seconds) | Speedup |
| ---------------------------------------------------------- | ------- | -------- | ---------- | ------------- | -------------- | ------- |
| Threaded                                                   | 1       | 50,000   |            |               |                |         |
| Threaded                                                   | 4       |          |            |               |                |         |
| Threaded                                                   | 8       |          |            |               |                |         |
| Threaded                                                   | 16      |          |            |               |                |         |

| Pseudo-random number generator (math/rand) | Threads | Vertices | CPU load % | Thread load % | Time (seconds) | Speedup |
| ------------------------------------------ | ------- | -------- | ---------- | ------------- | -------------- | ------- |
| Threaded                                   | 1       | 50,000   |            |               |                |         |
| Threaded                                   | 4       |          |            |               |                |         |
| Threaded                                   | 8       |          |            |               |                |         |
| Threaded                                   | 16      |          |            |               |                |         |

| Pseudo-random number sequence generator (math/rand) | Threads | Vertices | CPU load % | Thread load % | Time (seconds) | Speedup |
| --------------------------------------------------- | ------- | -------- | ---------- | ------------- | -------------- | ------- |
| Threaded                                            | 1       | 50,000   |            |               |                |         |
| Threaded                                            | 4       |          |            |               |                |         |
| Threaded                                            | 8       |          |            |               |                |         |
| Threaded                                            | 16      |          |            |               |                |         |

