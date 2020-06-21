import java.util.List;
import java.util.concurrent.ThreadLocalRandom;

public class MatrixGeneratorRunnable implements Runnable {
    Integer id;
    List<Integer> rows;
    MatrixGraph graph;

    @Override
    public void run() {
        for (int row = id; row < graph.vertices; row += graph.params.threads) {
            if (row >= graph.vertices) {
                break;
            }
            for (int i = 0; i < graph.vertices; i++) {
                graph.edges[row][i] = this.randomThreaded(0, 100) <= graph.density;
            }
        }
    }

    MatrixGeneratorRunnable(int id, MatrixGraph graph) {
        this.id = id;
        this.graph = graph;
    }

    int random(int min, int max) {
        return min + (int) (Math.random() * ((max - min) + 1));
    }

    int randomThreaded(int min, int max) {
        return ThreadLocalRandom.current().nextInt(min, max + 1);
    }

    int notRandom(int min, int max) {
        return max;
    }
}
