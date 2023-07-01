package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Task       Attributes `yaml:"task"`
	Dependency Attributes `yaml:"dependency"`
	Call       Attributes `yaml:"invocation"`
	NodeStyles []Style    `yaml:"node-styles"`
}

type Style struct {
	matcher    *regexp.Regexp
	Match      string
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

func (c *Config) Load(filepath string) error {
	// Create a reader to load the config file
	reader, err := os.Open(filepath)
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}

	// Load the config and flag any errors
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	return decoder.Decode(c)
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
		g := regexp.QuoteMeta(s.Match)
		g = strings.ReplaceAll(g, "\\*", ".*")
		g = strings.ReplaceAll(g, "\\?", ".")
		g = "(?i)(^" + g + "$)"

		s.matcher = regexp.MustCompile(g)
	}

	return s.matcher.MatchString(name)
}
