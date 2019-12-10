package cli

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

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
	config.Plans = parsePlans(configYaml["Plans"].([]interface{}), config.Variables)

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
func parsePlans(plansYaml []interface{}, variables map[string]string) (plans map[string]model.Plan) {
	plans = make(map[string]model.Plan)

	// Iterate plans
	for _, planYaml := range plansYaml {
		var plan model.Plan

		planYaml := planYaml.(map[string]interface{})

		plan.Name = planYaml["Name"].(string)

		if description, ok := planYaml["Description"].(string); ok {
			plan.Description = description
		}

		// Parse the tasks belonging to this plan
		plan.Tasks = parseTasks(planYaml["Tasks"].([]interface{}), variables)

		plans[plan.Name] = plan
	}

	return plans
}

// parseTasks creates the struct representation of the tasks section in the yaml file.
func parseTasks(tasksYaml []interface{}, variables map[string]string) (tasks []model.Task) {
	//Iterate tasks
	for _, taskYaml := range tasksYaml {
		var task model.Task

		taskYaml := taskYaml.(map[string](interface{}))

		if name, ok := taskYaml["Name"].(string); ok {
			task.Name = name
		}

		task.Steps = parseSteps(taskYaml["Steps"].([]interface{}), variables)

		tasks = append(tasks, task)
	}
	return tasks
}

// parseSteps creates the struct representation of the steps section in the yaml file.
// It determines which kind of step is defined and create the correct implementation for it.
func parseSteps(stepsYaml []interface{}, variables map[string]string) (steps []model.Step) {
	//Iterate tasks
	for _, stepYaml := range stepsYaml {

		stepYaml := stepYaml.(map[string]interface{})

		// Create the correct instance based on the type of step
		if kind, ok := stepYaml["Kind"].(string); ok {
			// Check the preferences if any variables where used that need to be expanded
			if _, ok := stepYaml["Preferences"]; ok {
				stepYaml["Preferences"] = expandPreferences(stepYaml["Preferences"].(map[string]interface{}), variables)
			}
			// Create the correct type of step
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

// expandPreferences checks each preference if it contains a variable "${...}" if it finds one
// it will check that this variable has been defined and fills it accordingly. If a string contains
// a variable that is not defined an error is returned.
func expandPreferences(preferences map[string]interface{}, variables map[string]string) (expandedPreferences map[string]interface{}) {
	expandedPreferences = make(map[string]interface{})
	for key, value := range preferences {
		// A slice has to be handled differently than a map
		switch value.(type) {
		default:
			log.Println("Unknown preference type in config")
		case []interface{}:
			// For an array
			for k, v := range value.([]interface{}) {
				value.([]interface{})[k] = expandVariable(v.(string), variables)
			}
			expandedPreferences[key] = value
		case interface{}:
			// For a single field
			expandedPreferences[key] = expandVariable(value.(string), variables)
		}
	}
	return expandedPreferences
}

func expandVariable(preference string, variables map[string]string) (expandedPreference string) {
	expandedPreference = preference
	// See if there are variables in the preference
	r, _ := regexp.Compile(`\${(.*?)\}`)
	hits := r.FindAllStringIndex(preference, -1)
	// Iterate over all hits
	for _, hit := range hits {
		h := preference[hit[0]:hit[1]]
		// Strip the beginning "${" and end "}" of the variable to get it's name
		variableName := h[2 : len(h)-1]
		// Do we have a variables with that name
		variableValue := variables[variableName]
		expandedPreference = strings.Replace(preference, h, variableValue, -1)
	}
	return expandedPreference
}
