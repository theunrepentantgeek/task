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

var params struct {
	version    bool
	help       bool
	color      bool
	verbose    bool
	configFile string
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

	pflag.BoolVar(&params.version, "version", false, "Show task-graph version.")
	pflag.BoolVarP(&params.help, "help", "h", false, "Shows task-graph usage.")
	pflag.BoolVarP(&params.color, "color", "c", true, "Colored output. Enabled by default. Set flag to false or use NO_COLOR=1 to disable.")
	pflag.BoolVarP(&params.verbose, "verbose", "v", false, "Enables log verbose mode.")
	pflag.StringVarP(&params.configFile, "config", "", "", "Configuration file to use.")

	pflag.Parse()

	if params.version {
		fmt.Printf("task-graph version: %s\n", ver.GetVersion())
		return nil
	}

	files := pflag.Args()
	if params.help || len(files) == 0 {
		pflag.Usage()
		return nil
	}

	if len(files) != 1 {
		return errors.New("task-graph requires exactly one taskfile")
	}

	config := defaultConfig()
	if params.configFile != "" {
		if err := config.Load(params.configFile); err != nil {
			return err
		}
	}

	builder, err := NewGraphBuilder(files[0], config, createLogger())
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
		Verbose: params.verbose,
		Color:   params.color,
	}
}
