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
	name  string
	tasks []Task
}

type Task struct {
	name  string
	steps []Step
}

type Step struct {
	kind string
}

func main() {
	log.Print("Welcome to reapy")
	config := readPlanDefinition()

	// Print extracted variables
	for key, value := range config.variables {
		log.Printf("Variable: %s:%s", key, value)
	}

	// Print extracted plans
	for key, value := range config.plans {
		log.Printf("Plan: %d:%s", key, value.name)
		// Print tasks of plan
		for key, value := range value.tasks {
			log.Printf("Task: %d:%s", key, value.name)
			//Print steps of a task
			for key, value := range value.steps {
				log.Printf("Step: %d:%s", key, value.kind)
			}
		}
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

		// Parse the tasks belonging to this plan
		plan.tasks = parseTasks(planYaml["Tasks"].([]interface{}))

		plans = append(plans, plan)
	}

	return plans
}

func parseTasks(tasksYaml []interface{}) (tasks []Task) {
	//Iterate tasks
	for _, taskYaml := range tasksYaml {
		var task Task

		taskYaml := taskYaml.(map[string](interface{}))

		if name, ok := taskYaml["Name"].(string); ok {
			task.name = name
		}

		task.steps = parseSteps(taskYaml["Steps"].([]interface{}))

		tasks = append(tasks, task)
	}
	return tasks
}

func parseSteps(stepsYaml []interface{}) (steps []Step) {
	//Iterate tasks
	for _, stepYaml := range stepsYaml {
		var step Step

		stepYaml := stepYaml.(map[string](interface{}))

		if kind, ok := stepYaml["Kind"].(string); ok {
			step.kind = kind
		}

		steps = append(steps, step)
	}
	return steps
}
