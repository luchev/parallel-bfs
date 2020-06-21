public class Logger {
    static void info(String s) {
        System.out.println("[INFO] " + s);
    }

    static void err(String s) {
        System.err.println("[ERR]" + s);
    }

    static void warn(String s) {
        System.err.println("[WARN]" + s);
    }
}
