# Parallel BFS

## C++

| Mersenne Twister | Threads | Vertices | Time (seconds) | Speedup |
| ---------------- | ------- | -------- | -------------- | ------- |
| Non-Threaded     | 1       | 10,000   | 32 минути      | 1       |
| Non-Threaded     | 16      | 10,000   | 43 минути      | 0.74    |
| Threaded         | 1       | 10,000   | 6.8707         | 1       |
| Threaded         | 4       | 10,000   | 1.7845         | 3.85    |
| Threaded         | 8       | 10,000   | 0.9108         | 7.54    |
| Threaded         | 16      | 10,000   | 0.5386         | 12.75   |
| Threaded         | 32      | 10,000   | 0.4383         | 15.67   |

| Marsaglia's xorshf | Threads | Vertices | Time (seconds) | Speedup |
| ------------------ | ------- | -------- | -------------- | ------- |
| Non-Threaded       | 1       | 20,000   | 3.6333         | 1       |
| Non-Threaded       | 4       | 20,000   | 10.3713        | 0.35    |
| Threaded           | 1       | 20,000   | 3.4334         | 1       |
| Threaded           | 4       | 20,000   | 0.9047         | 3.7     |
| Threaded           | 8       | 20,000   | 0.4735         | 7.2     |
| Threaded           | 16      | 20,000   | 0.3204         | 10.71   |
| Threaded           | 32      | 20,000   | 0.2857         | 12      |

| Algorithm (threaded) | Threads | Vertices | Time (seconds) | Speedup |
| -------------------- | ------- | -------- | -------------- | ------- |
| Mersenne Twister     | 1       | 50,000   | 172.0770       | 1       |
| Mersenne Twister     | 16      | 50,000   | 12.0559        | 14.27   |
| Mersenne Twister     | 32      | 50,000   | 9.51459        | 18      |
| Marsaglia's xorshf   | 1       | 50,000   | 21.4722        | 1       |
| Marsaglia's xorshf   | 16      | 50,000   | 1.6168         | 13.28   |
| Marsaglia's xorshf   | 32      | 50,000   | 1.6420         | 13      |

## Java

| Math.random | Threads | Vertices | Time (seconds) | Speedup |
| ----------- | ------- | -------- | -------------- | ------- |
|             | 1       | 10,000   | 3.0796         |         |
|             | 4       | 10,000   | 26.1181        |         |

| ThreadLocalRandom | Threads | Vertices | Time (seconds) | Speedup |
| ----------------- | ------- | -------- | -------------- | ------- |
|                   | 1       | 40,000   | 10.3472        | 1       |
|                   | 4       | 40,000   | 3.4738         | 2.97    |
|                   | 8       | 40,000   | 2.3720         | 4.36    |
|                   | 16      | 40,000   | 2.0972         | 4.93    |
|                   | 32      | 40,000   | 1.9776         | 5.23    |

## Go

| Cryptographic pseudo-random number sequence generator | Threads | Vertices | Time (seconds) | Speedup |
| ----------------------------------------------------- | ------- | -------- | -------------- | ------- |
| Threaded                                              | 1       | 40,000   | 12.8282        | 1       |
| Threaded                                              | 4       | 40,000   | 6.3824         | 2       |
| Threaded                                              | 8       | 40,000   | 5.4181         | 2.36    |
| Threaded                                              | 16      | 40,000   | 5.7462         | 2.23    |

| Pseudo-random number generator | Threads | Vertices | Time (seconds) | Speedup |
| ------------------------------ | ------- | -------- | -------------- | ------- |
| Threaded                       | 1       | 40,000   | 33.8937        | 1       |
| Threaded                       | 4       | 40,000   | 8.3982         | 4       |
| Threaded                       | 8       | 40,000   | 4.3595         | 7.77    |
| Threaded                       | 16      | 40,000   | 2.4213         | 14      |

| Pseudo-random number sequence generator | Threads | Vertices | Time (seconds) | Speedup |
| --------------------------------------- | ------- | -------- | -------------- | ------- |
| Threaded                                | 1       | 40,000   | 4.5132         | 1       |
| Threaded                                | 4       | 40,000   | 1.2269         | 3.67    |
| Threaded                                | 8       | 40,000   | 0.6593         | 6.84    |
| Threaded                                | 16      | 40,000   | 0.3996         | 11.29   |

