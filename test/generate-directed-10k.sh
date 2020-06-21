# rm generate-directed-10k/* 2>/dev/null
mkdir generate-directed-10k 2>/dev/null
for i in 1 2 4 6 8 10 12 14 16 18 20 22 24 26 28 30 32
do
    echo "Running tests for $i threads"
    for k in $(seq 1 3)
    do
        go run ../bfs/bfs.go -g -v 10000 -directed -t $i | grep -P 'Memory| using .*? took|Serializing graph to disk' >> generate-directed-10k/${i}.res
    done
done
