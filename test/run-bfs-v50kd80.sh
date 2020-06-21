# rm run-bfs-v50kd80/* 2>/dev/null
mkdir run-bfs-v50kd80 2>/dev/null
go run ../bfs/bfs.go -g -v 50000 -d 80 -o v50kd80 -q
for i in 1 2 4 6 8 10 12 14 16 18 20 22 24 26 28 30 32
do
    echo "Running tests for $i threads"
    for k in $(seq 1 3)
    do
        go run ../bfs/bfs.go -i v50kd80.graph -t $i | grep -P 'Reading|Serial BFS took|Parallel BFS with level .*? took|Custom parallel traversal using .*? took' >> run-bfs-v50kd80/${i}.res
    done
done
