//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	variables map[string]string
	plans     []Plan
}

type Plan struct {
	name string
}

func main() {
	log.Print("Welcome to reapy")
	config := readPlanDefinition()

	// Print extracted variables
	log.Print("Defined Variables:")
	for key, value := range config.variables {
		log.Printf("%s:%s", key, value)
	}

	// Print extracted plans
	log.Print("Defined plans")
	for key, value := range config.plans {
		log.Printf("%d:%s", key, value.name)
	}

}

func readPlanDefinition() Config {
	// Read the plans from the plans.json file
	planDefinition, err := ioutil.ReadFile("test.yaml")
	if err != nil {
		log.Panic(err)
	}
	return extractConfig(planDefinition)
}

func extractConfig(configYaml []byte) (config Config) {
	var configMap map[string]interface{}

	if err := yaml.Unmarshal(configYaml, &configMap); err != nil {
		log.Panicf("Unable to read plan definition: %s", err)
	}
	return parseConfig(configMap)
}

func parseConfig(configYaml map[string]interface{}) (config Config) {
	// Variables
	config.variables = parseVariables(configYaml["Variables"].(map[string]interface{}))
	// Plans
	config.plans = parsePlans(configYaml["Plans"].([]interface{}))

	return
}

func parseVariables(variablesYaml map[string]interface{}) (variables map[string]string) {
	variables = make(map[string]string)

	// Iterate all variables
	for key, value := range variablesYaml {
		variables[key] = value.(string)
	}

	return variables
}

func parsePlans(plansYaml []interface{}) (plans []Plan) {

	// Iterate plans
	for _, planYaml := range plansYaml {
		var plan Plan

		planYaml := planYaml.(map[string]interface{})

		if name, ok := planYaml["Name"].(string); ok {
			plan.name = name
		}

		plans = append(plans, plan)
	}

	return plans
}
