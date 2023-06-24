package main

import (
	"os"
	"path/filepath"

	"github.com/go-task/task/v3/internal/logger"
	"github.com/go-task/task/v3/taskfile/read"
)

type graphBuilder struct {
	directory string         // directory from which to load a taskfile
	taskfile  string         // name of the taskfile to load (optional)
	log       *logger.Logger // log to use for progress

	tasks        map[string]*node
	dependencies []*edge
}

type node struct {
	id    string
	label string
}

type edge struct {
	from *node
	to   *node
}

// NewGraphBuilder creates a new graph builder.
// path specifies either a directory or a taskfile.
func NewGraphBuilder(
	path string,
	log *logger.Logger,
) (*graphBuilder, error) {

	result := &graphBuilder{
		log: log,
	}

	// If the path doesn't exist, we can't load anything
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// If path specifies a folder, then load the taskfile from that folder.
	if stat.IsDir() {
		result.directory = path
	} else {
		result.directory = filepath.Dir(path)
		result.taskfile = filepath.Base(path)
	}

	return result, nil
}

func (b *graphBuilder) BuildGraph() (string, error) {
	reader := &read.ReaderNode{
		Dir:        b.directory,
		Entrypoint: b.taskfile,
		Parent:     nil,
		Optional:   false,
	}

	taskfile, _, err := read.Taskfile(reader)
	if err != nil {
		return "", err
	}

	taskNames := taskfile.Tasks.Keys()

}
