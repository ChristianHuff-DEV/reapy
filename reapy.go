//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	c "github.com/ChristianHuff-DEV/reapy/config"
	"github.com/ChristianHuff-DEV/reapy/model"
)

var config model.Config

func main() {
	for _, plan := range config.Plans {
		c.Execute(plan)
	}
}

// init will read the config yaml before starting the app itself
func init() {
	config = readPlanDefinition()
}

// readPlanDefinition parses a given config yaml file into the config instance
func readPlanDefinition() model.Config {
	return c.Extract("test.yaml")
}
