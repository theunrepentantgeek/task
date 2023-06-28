package main

import (
	"regexp"
	"strings"
)

type Config struct {
	Task       Attributes `yaml:"task"`
	Dependency Attributes `yaml:"dependency"`
	Call       Attributes `yaml:"invocation"`
	NodeStyles []Style    `yaml:"nodeStyles"`
}

type Style struct {
	matcher    *regexp.Regexp
	Glob       string
	Attributes Attributes
}

func defaultConfig() *Config {
	return &Config{
		Task: Attributes{
			"shape":    "box",
			"style":    "rounded, filled",
			"fontname": "Segoe UI",
			"penwidth": "2",
		},
		Dependency: Attributes{
			"style":     "dashed",
			"arrowhead": "none",
			"arrowtail": "none",
		},
		Call: Attributes{
			"style":     "solid",
			"arrowhead": "normal",
			"arrowtail": "none",
		},
	}
}

func (c *Config) ApplyNodeStyles(node *node) {
	for _, style := range c.NodeStyles {
		if style.Matches(node.label) {
			node.applyStyle(style.Attributes)
		}
	}
}

func (s *Style) Matches(name string) bool {
	if s.matcher == nil {
		g := regexp.QuoteMeta(s.Glob)
		g = strings.ReplaceAll(g, "\\*", ".*")
		g = strings.ReplaceAll(g, "\\?", ".")
		g = "(?i)(^" + g + "$)"

		s.matcher = regexp.MustCompile(g)
	}

	return s.matcher.MatchString(name)
}
