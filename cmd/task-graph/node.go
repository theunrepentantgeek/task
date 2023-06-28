package main

import (
	"io"
)

type node struct {
	id         string
	label      string
	attributes Attributes
}

func NewNode(id string, label string) *node {
	result := &node{
		id:         id,
		label:      label,
		attributes: make(Attributes),
	}

	result.attributes["label"] = label

	return result
}

func (n *node) applyStyle(style Attributes) {
	for k, v := range style {
		n.attributes[k] = v
	}
}

func (n *node) writeTo(buffer io.StringWriter) {
	n.attributes.writeTo(n.id, buffer)
}
