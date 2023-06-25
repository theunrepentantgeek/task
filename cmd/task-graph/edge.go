package main

import (
	"fmt"
	"io"
)

type edge struct {
	from *node
	to   *node
}

func NewEdge(from *node, to *node) *edge {
	return &edge{
		from: from,
		to:   to,
	}
}

func (e *edge) writeTo(buffer io.StringWriter) {
	buffer.WriteString(fmt.Sprintf("    %s -> %s;\n", e.from.id, e.to.id))
}
