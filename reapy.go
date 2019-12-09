//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"github.com/ChristianHuff-DEV/reapy/cli"
	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/c-bata/go-prompt"
)

var config model.Config

func main() {
	// for _, plan := range config.Plans {
	// 	c.Execute(plan)
	// }
	p := prompt.New(cli.Executor, cli.Completer)
	p.Run()
}

// init will read the config yaml before starting the app itself
func init() {
	config = readPlanDefinition()
}

// readPlanDefinition parses a given config yaml file into the config instance
func readPlanDefinition() model.Config {
	return cli.Extract("test.yaml")
}
