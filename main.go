package main

import (
	"fmt"

	"github.com/flou/ecs/cmd"
)

var (
	version string = "dev"
	commit  string
	date    string
)

func main() {
	cmd.Execute(
		buildVersion(version, commit, date),
	)
}

func buildVersion(version, commit, date string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	return result
}
