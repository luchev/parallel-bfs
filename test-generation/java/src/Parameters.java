import org.apache.commons.cli.CommandLine;

public class Parameters {
    int threads = DEFAULT_THREADS;
    int vertices = DEFAULT_VERTICES;
    int density = DEFAULT_DENSITY;
    String inputFile = DEFAULT_IN_FILE;
    String outputFile = DEFAULT_OUT_FILE;
    Verbosity verbosity = Verbosity.Normal;

    private static final int DEFAULT_THREADS = 1;
    private static final int DEFAULT_VERTICES = 0;
    private static final int DEFAULT_DENSITY = 20;
    private static final String DEFAULT_IN_FILE = "graph.in";
    private static final String DEFAULT_OUT_FILE = "graph.out";
    private static final Verbosity DEFAULT_VERBOSITY = Verbosity.Normal;

    static Parameters fromCmd(CommandLine args) {
        Parameters parameters = new Parameters();
        parameters.threads = parseThreads(args);
        if (args.hasOption("input-file")) {
            parameters.inputFile = parseInputFile(args);
        } else {
            parameters.vertices = parseVertices(args);
        }
        parameters.density = parseDensity(args);
        parameters.outputFile = parseOutputFile(args);
        parameters.verbosity = parseVerbosity(args);
        return parameters;
    }

    private Parameters() { }

    private static int parseThreads(CommandLine args) {
        String arg = args.getOptionValue("threads");
        if (arg != null) {
            if (arg.equals("auto")) {
                int cpuThreads = Runtime.getRuntime().availableProcessors();
                Logger.info("System has " + Integer.toString(cpuThreads) + " logical threads");
                Logger.info("Initializing number of threads " + Integer.toString(cpuThreads));
                return cpuThreads;
            }
            try {
                int threads = Integer.parseInt(arg);
                if (threads < 1) {
                    Logger.warn("Number of threads < 1. Using default number of threads " + Integer.toString(DEFAULT_THREADS));
                    return DEFAULT_THREADS;
                } else {
                    Logger.info("Initializing number of threads " + Integer.toString(threads));
                    return threads;
                }

            } catch (NumberFormatException e) {
                Logger.warn("Failed to parse <threads>. Using default number of threads " + Integer.toString(DEFAULT_THREADS));
                return DEFAULT_THREADS;
            }
        }

        Logger.info("No thread count provided. Using default number of threads " + Integer.toString(DEFAULT_THREADS));
        return DEFAULT_THREADS;
    }

    private static int parseVertices(CommandLine args) {
        String arg = args.getOptionValue("vertices");
        if (arg != null) {
            try {
                int vertices = Integer.parseInt(arg);
                if (vertices < 0) {
                    Logger.warn("Vertices < 0. Using default vertices " + Integer.toString(DEFAULT_VERTICES));
                    return DEFAULT_VERTICES;
                } else {
                    Logger.info("Initializing vertices to " + Integer.toString(vertices));
                    return vertices;
                }

            } catch (NumberFormatException e) {
                Logger.warn("Failed to parse <vertices>. Using default vertices " + DEFAULT_VERTICES);
                return DEFAULT_VERTICES;
            }
        }

        Logger.info("No vertices provided. Using default vertices " + Integer.toString(DEFAULT_VERTICES));
        return DEFAULT_VERTICES;
    }

    private static int parseDensity(CommandLine args) {
        String arg = args.getOptionValue("density");
        if (arg != null) {
            try {
                int density = Integer.parseInt(arg);
                if (density < 0) {
                    Logger.warn("Density < 0. Using default density " + Integer.toString(DEFAULT_DENSITY));
                    return DEFAULT_DENSITY;
                } else {
                    Logger.info("Initializing density to " + Integer.toString(density));
                    return density;
                }

            } catch (NumberFormatException e) {
                Logger.warn("Failed to parse <density>. Using default density " + DEFAULT_DENSITY);
                return DEFAULT_DENSITY;
            }
        }

        Logger.info("No density provided. Using default density " + Integer.toString(DEFAULT_DENSITY));
        return DEFAULT_DENSITY;
    }

    private static String parseOutputFile(CommandLine args) {
        String arg = args.getOptionValue("output-file");
        if (arg != null) {
            Logger.info("Initializing output filename " + arg);
            return arg;
        }

        Logger.info("No output filename provided. Using default output filename " + DEFAULT_OUT_FILE);
        return DEFAULT_OUT_FILE;
    }

    private static String parseInputFile(CommandLine args) {
        String arg = args.getOptionValue("input-file");
        if (arg != null) {
            Logger.info("Initializing input filename " + arg);
            return arg;
        }

        Logger.info("No input filename provided. Using default input filename " + DEFAULT_IN_FILE);
        return DEFAULT_IN_FILE;
    }

    private static Verbosity parseVerbosity(CommandLine args) {
        if (args.hasOption("quiet")) {
            Logger.info("Initializing verbosity level to Quiet");
            return Verbosity.Quiet;
        } else if (args.hasOption("verbose")) {
            Logger.info("Initializing verbosity level to Verbose");
            return Verbosity.Verbose;
        } else {
            Logger.info("Initializing verbosity level to " + DEFAULT_VERBOSITY.toString());
            return DEFAULT_VERBOSITY;
        }
    }
}

enum Verbosity {
    Quiet,
    Normal,
    Verbose
}
