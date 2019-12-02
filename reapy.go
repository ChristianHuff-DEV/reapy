//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"io/ioutil"
	"log"

	c "github.com/ChristianHuff-DEV/reapy/config"
	"github.com/ChristianHuff-DEV/reapy/model"
)

var config model.Config

func main() {
	for _, plan := range config.Plans {
		c.Execute(plan)
	}
}

func init() {
	config = readPlanDefinition()
}

func readPlanDefinition() model.Config {
	// Read the plans from the plans.json file
	planDefinition, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		log.Panic(err)
	}
	return c.Extract(planDefinition)
}
