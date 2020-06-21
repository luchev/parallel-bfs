import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Semaphore;

public class MatrixGraph implements Graph {
    int vertices;
    int density;
    boolean[][] edges;
    Semaphore activeThreads;
    Parameters params;

    private MatrixGraph() {
    }

    static MatrixGraph fromVertices(Parameters params) {
        final long start = System.nanoTime();

        MatrixGraph graph = new MatrixGraph();
        graph.params = params;
        graph.vertices = params.vertices;
        graph.density = params.density;
        graph.edges = new boolean[params.vertices][params.vertices];

        final Thread[] threads = new Thread[params.threads];
        for (int i = 0; i < params.threads; i++) {
            final MatrixGeneratorRunnable worker = new MatrixGeneratorRunnable(i, graph);
            threads[i] = new Thread(worker);
            threads[i].start();
        }

        for (int i = 0; i < params.threads; i++) {
            try {
                threads[i].join();
            } catch (final InterruptedException e) {
                Logger.err("Failed to join thread " + threads[i].getName());
            }
        }

        final long end = System.nanoTime();
        Logger.info("Threads used in current run " + params.threads);
        Logger.info("Total execution time in seconds " + Time.formatNanoSeconds(end - start));
        return graph;
    }

    public String toString() {
        StringBuilder str = new StringBuilder();
        str.append(vertices);
        str.append('\n');

        for (int i = 0; i < vertices; i++) {
            for (int k = 0; k < vertices; k++) {
                str.append(edges[i][k] ? 1 : 0);
                str.append(' ');
            }
            str.append('\n');
        }

        return str.toString();
    }

    @Override
    public List<Integer> getNeighbours(int vertex) {
        List<Integer> neighbours = new ArrayList<>();
        for (int i = 0; i < this.vertices; i++) {
            if (edges[vertex][i]) {
                neighbours.add(i);
            }
        }
        return neighbours;
    }
}
