package main

import (
	"fmt"
	"io"
)

type Style map[string]string

func (s *Style) writeTo(scope string, buffer io.StringWriter) {
	buffer.WriteString(fmt.Sprintf("    %s [\n", scope))
	for k, v := range *s {
		buffer.WriteString(fmt.Sprintf("        %s=%q\n", k, v))
	}

	buffer.WriteString("    ]\n\n")
}
