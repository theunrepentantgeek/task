package main

import (
	"fmt"
	"io"
)

type node struct {
	id    string
	label string
}

func NewNode(id string, label string) *node {
	return &node{
		id:    id,
		label: label,
	}
}

func (n *node) writeTo(buffer io.StringWriter) {
	buffer.WriteString(fmt.Sprintf("    %s [label=%q];\n", n.id, n.label))
}
