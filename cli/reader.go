package cli

import (
	"io/ioutil"
	"log"

	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/ChristianHuff-DEV/reapy/step"
	stepDefinition "github.com/ChristianHuff-DEV/reapy/step"
	"gopkg.in/yaml.v3"
)

// Extract takes the location of the yaml file and delegats it's content to the method reading the content and creating the config.
func Extract(filePath string) (config model.Config) {

	configYaml, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	var configMap map[string]interface{}

	if err := yaml.Unmarshal(configYaml, &configMap); err != nil {
		log.Panicf("Unable to read plan definition: %s", err)
	}

	config = parseConfig(configMap)

	return config
}

// parseConfig takes a map representing the yaml config file content and delegates it to the methods extracting the variables and plans
func parseConfig(configYaml map[string]interface{}) (config model.Config) {
	// Variables
	config.Variables = parseVariables(configYaml["Variables"].(map[string]interface{}))
	// Plans
	config.Plans = parsePlans(configYaml["Plans"].([]interface{}))

	return
}

// parseVariables extracts the first section of the yaml file that defines the variables which might be used in the later definition of tasks/steps
func parseVariables(variablesYaml map[string]interface{}) (variables map[string]string) {
	variables = make(map[string]string)

	// Iterate all variables
	for key, value := range variablesYaml {
		variables[key] = value.(string)
	}

	return variables
}

// parsePlans creates the struct representation of the plans section in the yaml file.
func parsePlans(plansYaml []interface{}) (plans []model.Plan) {
	// Iterate plans
	for _, planYaml := range plansYaml {
		var plan model.Plan

		planYaml := planYaml.(map[string]interface{})

		plan.Name = planYaml["Name"].(string)

		if description, ok := planYaml["Description"].(string); ok {
			plan.Description = description
		}

		// Parse the tasks belonging to this plan
		plan.Tasks = parseTasks(planYaml["Tasks"].([]interface{}))

		plans = append(plans, plan)
	}

	return plans
}

// parseTasks creates the struct representation of the tasks section in the yaml file.
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

// parseSteps creates the struct representation of the steps section in the yaml file.
// It determines which kind of step is defined and create the correct implementation for it.
func parseSteps(stepsYaml []interface{}) (steps []model.Step) {
	//Iterate tasks
	for _, stepYaml := range stepsYaml {

		stepYaml := stepYaml.(map[string]interface{})

		// Create the correct instance based on the type of step
		if kind, ok := stepYaml["Kind"].(string); ok {
			switch kind {
			case step.KindDownload:
				step := stepDefinition.Download{}
				step.FromConfig(stepYaml)
				steps = append(steps, &step)
			case step.KindUnzip:
				step := stepDefinition.Unzip{}
				step.FromConfig(stepYaml)
				steps = append(steps, &step)
			case step.KindDelete:
				step := stepDefinition.Delete{}
				step.FromConfig(stepYaml)
				steps = append(steps, &step)
			case step.KindCreateFolder:
				step := stepDefinition.CreateFolder{}
				step.FromConfig(stepYaml)
				steps = append(steps, &step)
			case step.KindCommand:
				step := stepDefinition.Command{}
				step.FromConfig(stepYaml)
				steps = append(steps, &step)
			}
		}
	}
	return steps
}
