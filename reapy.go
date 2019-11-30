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
	log.Print("Welcome to reapy")

	for _, plan := range config.Plans {
		c.Execute(plan)
	}
}

func init() {
	config = readPlanDefinition()
	// Print extracted variables
	for key, value := range config.Variables {
		log.Printf("Variable: %s:%s", key, value)
	}

	// Print extracted plans
	for key, value := range config.Plans {
		log.Printf("Plan: %d:%s", key, value.Name)
		// Print tasks of plan
		for key, value := range value.Tasks {
			log.Printf("Task: %d:%s", key, value.Name)
			//Print steps of a task
			for key, value := range value.Steps {
				log.Printf("Step: %d:%s", key, value.GetKind())
			}
		}
	}
}

func readPlanDefinition() model.Config {
	// Read the plans from the plans.json file
	planDefinition, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		log.Panic(err)
	}
	return c.Extract(planDefinition)
}
