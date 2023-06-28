package main

import (
	"fmt"
	"io"
)

type Attributes map[string]string

func (s Attributes) writeTo(item string, buffer io.StringWriter) {
	// If we have only one attribute, put it on the same line
	if len(s) == 1 {
		buffer.WriteString(fmt.Sprintf("    %s [", item))
		for k, v := range s {
			buffer.WriteString(fmt.Sprintf("%s=%q", k, v))
		}
		buffer.WriteString("]\n")
		return
	}

	// With multiple attributes, put them on separate lines
	buffer.WriteString(fmt.Sprintf("    %s [\n", item))
	for k, v := range s {
		buffer.WriteString(fmt.Sprintf("        %s=%q\n", k, v))
	}

	buffer.WriteString("    ]\n\n")
}
