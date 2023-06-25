package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/exp/maps"

	"github.com/go-task/task/v3/internal/logger"
	"github.com/go-task/task/v3/taskfile"
	"github.com/go-task/task/v3/taskfile/read"
)

const indent = "    "

type graphBuilder struct {
	directory string         // directory from which to load a taskfile
	taskfile  string         // name of the taskfile to load (optional)
	log       *logger.Logger // log to use for progress
	config    *config        // configuration for the graph
}

// NewGraphBuilder creates a new graph builder.
// path specifies either a directory or a taskfile.
func NewGraphBuilder(
	path string,
	log *logger.Logger,
) (*graphBuilder, error) {

	result := &graphBuilder{
		log:    log,
		config: defaultConfig(),
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

	nodes := b.createNodes(taskfile)
	dependencies := b.createDependencyEdges(taskfile, nodes)
	calls := b.createCallEdges(taskfile, nodes)

	var buffer strings.Builder
	b.writeFileHeader(&buffer)
	b.writeNodes("Tasks", b.config.Task, nodes, &buffer)
	b.writeEdges("Dependencies", b.config.Dependency, dependencies, &buffer)
	b.writeEdges("Calls", b.config.Call, calls, &buffer)
	b.writeFileTrailer(&buffer)

	return buffer.String(), nil
}

func (b *graphBuilder) createNodes(taskfile *taskfile.Taskfile) map[string]*node {
	result := make(map[string]*node)
	for _, task := range taskfile.Tasks.Values() {
		name := task.Task
		id := b.createId(name)
		node := NewNode(id, name)
		result[id] = node
	}

	return result
}

func (b *graphBuilder) createDependencyEdges(
	taskfile *taskfile.Taskfile,
	nodes map[string]*node,
) []*edge {
	var result []*edge
	for _, task := range taskfile.Tasks.Values() {
		name := task.Task
		fromId := b.createId(name)
		fromNode := nodes[fromId]
		if fromNode == nil {
			b.log.Outf(logger.Yellow, "Didn't find 'from' node for task %s\n", name)
			continue
		}

		for _, dep := range task.Deps {
			toId := b.createId(dep.Task)
			toNode := nodes[toId]
			if toNode == nil {
				b.log.Outf(logger.Yellow, "Didn't find 'to' node for task %s\n", dep.Task)
				continue
			}

			edge := NewEdge(fromNode, toNode)
			result = append(result, edge)
		}
	}

	return result
}

func (b *graphBuilder) createCallEdges(
	taskfile *taskfile.Taskfile,
	nodes map[string]*node,
) []*edge {
	var result []*edge
	for _, task := range taskfile.Tasks.Values() {
		name := task.Task
		fromId := b.createId(name)
		fromNode := nodes[fromId]
		if fromNode == nil {
			b.log.Outf(logger.Yellow, "Didn't find 'from' node for task %s\n", name)
			continue
		}

		for _, cmd := range task.Cmds {
			if cmd.Task == "" {
				// Not a call to a task
				continue
			}

			toId := b.createId(cmd.Task)
			toNode := nodes[toId]
			if toNode == nil {
				b.log.Outf(logger.Yellow, "Didn't find 'to' node for task %s\n", cmd.Task)
				continue
			}

			edge := NewEdge(fromNode, toNode)
			result = append(result, edge)
		}
	}

	return result
}

func (b *graphBuilder) writeFileHeader(buffer io.StringWriter) {
	buffer.WriteString("digraph taskfile {\n")
}

func (b *graphBuilder) writeFileTrailer(buffer io.StringWriter) {
	buffer.WriteString("}\n")
}

func (b *graphBuilder) writeNodes(
	header string,
	style Style,
	nodes map[string]*node,
	buffer io.StringWriter,
) {
	ids := maps.Keys(nodes)
	sort.Strings(ids)

	buffer.WriteString("\n")
	buffer.WriteString(fmt.Sprintf("    // %s\n", header))
	buffer.WriteString("\n")

	if len(style) > 0 {
		style.writeTo("node", buffer)
	}

	for _, id := range ids {
		node := nodes[id]
		node.writeTo(buffer)
	}
}

func (b *graphBuilder) writeEdges(
	header string,
	style Style,
	edges []*edge,
	buffer io.StringWriter,
) {
	buffer.WriteString("\n")
	buffer.WriteString(fmt.Sprintf("    // %s\n", header))
	buffer.WriteString("\n")

	if len(style) > 0 {
		style.writeTo("edge", buffer)
	}

	for _, edge := range edges {
		edge.writeTo(buffer)
	}
}

func (b *graphBuilder) createId(id string) string {
	var builder strings.Builder
	needSeparator := false
	empty := true
	for _, c := range id {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			if needSeparator && !empty {
				builder.WriteRune('_')
			}

			builder.WriteRune(c)
			needSeparator = false
			empty = false
		} else {
			needSeparator = true
		}
	}

	return builder.String()
}
