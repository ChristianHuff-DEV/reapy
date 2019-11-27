//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func main() {
	log.Print("Starting reapy")
}

func init() {
	log.Print("Read plans definition file")
	readPlansDefinitionFile()
}

func readPlansDefinitionFile() {
	plansDefinitionFilePath := "test.yaml"
	planDefinitionFile, err := ioutil.ReadFile(plansDefinitionFilePath)

	if err != nil {
		log.Panicf("No yaml config found at %s", plansDefinitionFilePath)
	}
	extractPlans(planDefinitionFile)
}

// extractPlans expects a YAML formatted byte array from which it will extract the plans definition
func extractPlans(plansDefinition []byte) {
	plan := Plan{}
	yaml.Unmarshal(plansDefinition, &plan)
}
