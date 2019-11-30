//go:generate goversioninfo -icon=icon.ico -manifest=goversioninfo.exe.manifest
package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"

	"github.com/ChristianHuff-DEV/reapy/model"
)

func main() {
	log.Print("Welcome to reapy")
	config := readPlanDefinition()

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
				log.Printf("Step: %d:%s", key, value.Kind)
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
	return extractConfig(planDefinition)
}

func extractConfig(configYaml []byte) (config model.Config) {
	var configMap map[string]interface{}

	if err := yaml.Unmarshal(configYaml, &configMap); err != nil {
		log.Panicf("Unable to read plan definition: %s", err)
	}
	return parseConfig(configMap)
}

func parseConfig(configYaml map[string]interface{}) (config model.Config) {
	// Variables
	config.Variables = parseVariables(configYaml["Variables"].(map[string]interface{}))
	// Plans
	config.Plans = parsePlans(configYaml["Plans"].([]interface{}))

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

func parsePlans(plansYaml []interface{}) (plans []model.Plan) {
	// Iterate plans
	for _, planYaml := range plansYaml {
		var plan model.Plan

		planYaml := planYaml.(map[string]interface{})

		if name, ok := planYaml["Name"].(string); ok {
			plan.Name = name
		}

		// Parse the tasks belonging to this plan
		plan.Tasks = parseTasks(planYaml["Tasks"].([]interface{}))

		plans = append(plans, plan)
	}

	return plans
}

func parseTasks(tasksYaml []interface{}) (tasks []model.Task) {
	//Iterate tasks
	for _, taskYaml := range tasksYaml {
		var task model.Task

		taskYaml := taskYaml.(map[string](interface{}))

		if name, ok := taskYaml["Name"].(string); ok {
			task.Name = name
		}

		task.Steps = parseSteps(taskYaml["Steps"].([]interface{}))

		tasks = append(tasks, task)
	}
	return tasks
}

func parseSteps(stepsYaml []interface{}) (steps []model.Step) {
	//Iterate tasks
	for _, stepYaml := range stepsYaml {
		var step model.Step

		stepYaml := stepYaml.(map[string](interface{}))

		if kind, ok := stepYaml["Kind"].(string); ok {
			step.Kind = kind
		}

		steps = append(steps, step)
	}
	return steps
}
