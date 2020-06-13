public class Time {
    static String formatNanoSeconds(final long nanoTime) {
        return String.format("%.9f", (double) nanoTime / 1_000_000_000);
    }
}
