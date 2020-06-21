# rm generate-directed-50k/* 2>/dev/null
mkdir generate-directed-50k 2>/dev/null
for i in 16 18 20 22 24 26 28 30 32
do
    echo "Running tests for $i threads"
    for k in $(seq 1 1)
    do
        go run ../bfs/bfs.go -g -v 50000 -directed -t $i | grep -P 'Memory| using .*? took|Serializing graph to disk' >> generate-directed-50k/${i}.res
    done
done
