import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.concurrent.*;

public class Generator {
    Parameters params;

    Generator(final Parameters parameters) {
        this.params = parameters;
    }

    void generateMatrixGraph() {
        MatrixGraph graph = MatrixGraph.fromVertices(params);
        // try {
        //     Files.write(Paths.get(params.outputFile), graph.toString().getBytes());
        // } catch (IOException e) {
        //     Logger.err("Failed to save file " + params.outputFile);
        // }
    }

    String formatNanoSeconds(final long nanoTime) {
        return String.format("%.9f", (double) nanoTime / 1_000_000_000);
    }

    int random(int min, int max) {
        return min + (int) (Math.random() * ((max - min) + 1));
    }

    int randomThreaded(int min, int max) {
        return ThreadLocalRandom.current().nextInt(min, max + 1);
    }
}
