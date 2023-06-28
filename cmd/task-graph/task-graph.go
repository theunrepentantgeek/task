package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"

	"github.com/go-task/task/v3/errors"
	"github.com/go-task/task/v3/internal/logger"
	ver "github.com/go-task/task/v3/internal/version"
)

var flags struct {
	version bool
	help    bool
	color   bool
	verbose bool
}

func main() {
	if err := run(); err != nil {
		// Always log to stderr to allow piping into GraphViz
		log := createLogger()

		log.Errf(logger.Red, "%v\n", err)
		os.Exit(errors.CodeUnknown)
	}

	os.Exit(errors.CodeOk)
}

func run() error {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	pflag.Usage = func() {
		pflag.PrintDefaults()
	}

	pflag.BoolVar(&flags.version, "version", false, "Show task-graph version.")
	pflag.BoolVarP(&flags.help, "help", "h", false, "Shows task-graph usage.")
	pflag.BoolVarP(&flags.color, "color", "c", true, "Colored output. Enabled by default. Set flag to false or use NO_COLOR=1 to disable.")
	pflag.BoolVarP(&flags.verbose, "verbose", "v", false, "Enables log verbose mode.")

	pflag.Parse()

	if flags.version {
		fmt.Printf("task-graph version: %s\n", ver.GetVersion())
		return nil
	}

	if flags.help {
		pflag.Usage()
		return nil
	}

	files := pflag.Args()
	if len(files) != 1 {
		return errors.New("task-graph requires exactly one taskfile")
	}

	builder, err := NewGraphBuilder(files[0], createLogger())
	if err != nil {
		return err
	}

	graph, err := builder.BuildGraph()
	if err != nil {
		return err
	}

	os.Stdout.WriteString(graph)
	return nil
}

func createLogger() *logger.Logger {
	return &logger.Logger{
		Stdout:  os.Stderr,
		Stderr:  os.Stderr,
		Verbose: flags.verbose,
		Color:   flags.color,
	}
}
