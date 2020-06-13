import org.apache.commons.cli.*;

public class Main {
    static Options options;

    public static void main(String[] args) {
        initOptions(args);
        CommandLine cmd = initCmd(args);

        if (cmd.hasOption("help")) {
            printHelp(options);
            System.exit(0);
        }

        Parameters parameters = Parameters.fromCmd(cmd);

        if (cmd.hasOption("generate")) {
            Generator generator = new Generator(parameters);
            generator.generateMatrixGraph();
        }
    }

    static void printHelp(Options options) {
        HelpFormatter formatter = new HelpFormatter();
        formatter.printHelp("./run.sh", options);

    }

    static void initOptions(String[] args) {
        options = new Options();
        
        OptionGroup exclusive = new OptionGroup();
        exclusive.addOption(new Option("v", "vertices", true, "Number of vertices (mutually exclusive with -i)"));
        exclusive.addOption(new Option("i", "input-file", true,
                "File containing a graph represented as an adjacency matrix (mutually exclusive with -v)"));
        
        options.addOptionGroup(exclusive);
        options.addOption("d", "density", true, "Density (how many edges to generate when using -v). Possible values 0 to 100 %.");
        options.addOption("t", "threads", true, "Maximum number of threads to use");
        options.addOption("o", "output-file", true, "File to write the graph into");
        options.addOption("q", "quiet", false, "Run in quiet mode");
        options.addOption("h", "help", false, "Show this menu");
        options.addOption("g", "generate", false, "Generate a graph with N vertices. Requires -v N");
    }

    static CommandLine initCmd(String[] args) {
        CommandLineParser argsParser = new DefaultParser();
        try {
            return argsParser.parse(options, args);
        } catch (ParseException e) {
            Logger.err("Error parsing arguments " + e.getMessage());
            printHelp(options);
            System.exit(1);
        }

        return new CommandLine.Builder().build();
    }
}
