package docker

import "fmt"

type Runtime struct {
	Name     string
	Image    string
	Command  []string
	FileName string
}

// TODO: will update runtime configs
var SupportedRuntimes = map[string]Runtime{
	"python": {
		Name:     "python",
		Image:    "python:3.9-alpine",
		Command:  []string{"python", "/app/main.py"},
		FileName: "main.py",
	},
	"node": {
		Name:     "node",
		Image:    "node:16-alpine",
		Command:  []string{"node", "/app/main.js"},
		FileName: "index.js",
	},
	"go": {
		Name:     "go",
		Image:    "golang:1.24-alpine",
		Command:  []string{"go", "run", "/app/main.go"},
		FileName: "main.go",
	},
}

func GetRuntime(language string) (Runtime, error) {
	runtime, exists := SupportedRuntimes[language]
	if !exists {
		return Runtime{}, fmt.Errorf("unsupported language: %s", language)
	}
	return runtime, nil
}
