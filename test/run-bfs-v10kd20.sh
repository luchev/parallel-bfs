# rm run-bfs-v10kd20/* 2>/dev/null
mkdir run-bfs-v10kd20 2>/dev/null
go run ../bfs/bfs.go -g -v 10000 -d 20 -o v10kd20 -q
for i in 1 2 4 6 8 10 12 14 16 18 20 22 24 26 28 30 32
do
    echo "Running tests for $i threads"
    for k in $(seq 1 3)
    do
        go run ../bfs/bfs.go -i v10kd20.graph -t $i | grep -P 'Reading|Serial BFS took|Parallel BFS with level .*? took|Custom parallel traversal using .*? took' >> run-bfs-v10kd20/${i}.res
    done
done
