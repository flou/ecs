package main

import (
	"fmt"
	"os"

	"github.com/flou/ecs/cmd"
)

var (
	version = "dev"
	commit  = ""
)

func main() {
	cmd.Execute(
		buildVersion(version, commit),
		os.Exit,
		os.Args[1:],
	)
}

func buildVersion(version, commit string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	return result
}
