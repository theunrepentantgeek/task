package main

type config struct {
	Task       Style `yaml:"task"`
	Dependency Style `yaml:"dependency"`
	Call       Style `yaml:"invocation"`
}

func defaultConfig() *config {
	return &config{
		Task: Style{
			"shape":    "box",
			"style":    "rounded",
			"fontname": "Segoe UI",
			"penwidth": "2",
		},
		Dependency: Style{
			"style":     "dashed",
			"arrowhead": "none",
			"arrowtail": "none",
		},
		Call: Style{
			"style":     "solid",
			"arrowhead": "normal",
			"arrowtail": "none",
		},
	}
}
