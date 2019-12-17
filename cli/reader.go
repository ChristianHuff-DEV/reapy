package cli

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/ChristianHuff-DEV/reapy/step"
	stepDefinition "github.com/ChristianHuff-DEV/reapy/step"
	"gopkg.in/yaml.v3"
)

// Extract takes the location of the yaml file and delegats it's content to the method reading the content and creating the config.
func Extract(filePath string) (config model.Config, err error) {
	log.Printf("read config from %s", filePath)

	configYaml, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	var configMap map[string]interface{}

	if err := yaml.Unmarshal(configYaml, &configMap); err != nil {
		log.Panicf("Unable to read plan definition: %s", err)
	}

	config, err = parseConfig(configMap)
	if err != nil {
		return config, err
	}

	return config, nil
}

// parseConfig takes a map representing the yaml config file content and delegates it to the methods extracting the variables and plans
func parseConfig(configYaml map[string]interface{}) (config model.Config, err error) {

	if variablesYaml, ok := configYaml["Variables"].(map[string]interface{}); ok {
		err = validateVariables(variablesYaml)
		if err != nil {
			return config, err
		}
		config.Variables = parseVariables(variablesYaml)
	}

	if plansYaml, ok := configYaml["Plans"].([]interface{}); ok {
		config.Plans, err = parsePlans(plansYaml, config.Variables)
		if err != nil {
			return config, err
		}
	} else {
		return config, fmt.Errorf("no plans defined")
	}

	return config, nil
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
func parsePlans(plansYaml []interface{}, variables map[string]string) (plans map[string]model.Plan, err error) {
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
		plan.Tasks, err = parseTasks(planYaml["Tasks"].([]interface{}), variables)
		if err != nil {
			return plans, err
		}

		plans[plan.Name] = plan
	}

	return plans, nil
}

// parseTasks creates the struct representation of the tasks section in the yaml file.
func parseTasks(tasksYaml []interface{}, variables map[string]string) (tasks []model.Task, err error) {
	//Iterate tasks
	for _, taskYaml := range tasksYaml {
		var task model.Task

		taskYaml := taskYaml.(map[string](interface{}))

		if name, ok := taskYaml["Name"].(string); ok {
			task.Name = name
		}

		task.Steps, err = parseSteps(taskYaml["Steps"].([]interface{}), variables)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

// parseSteps creates the struct representation of the steps section in the yaml file.
// It determines which kind of step is defined and create the correct implementation for it.
func parseSteps(stepsYaml []interface{}, variables map[string]string) (steps []model.Step, err error) {
	//Iterate tasks
	for _, stepYaml := range stepsYaml {

		stepYaml := stepYaml.(map[string]interface{})

		// Create the correct instance based on the type of step
		if kind, ok := stepYaml["Kind"].(string); ok {
			// Check the preferences if any variables where used that need to be expanded
			if _, ok := stepYaml["Preferences"]; ok {
				stepYaml["Preferences"], err = expandPreferences(stepYaml["Preferences"].(map[string]interface{}), variables)
				if err != nil {
					return steps, err
				}
			}
			// Create the correct type of step
			switch kind {
			case step.KindDownload:
				step := stepDefinition.Download{}
				err = step.FromConfig(stepYaml)
				if err != nil {
					return steps, err
				}
				steps = append(steps, &step)
			case step.KindUnzip:
				step := stepDefinition.Unzip{}
				err = step.FromConfig(stepYaml)
				if err != nil {
					return steps, err
				}
				steps = append(steps, &step)
			case step.KindDelete:
				step := stepDefinition.Delete{}
				err = step.FromConfig(stepYaml)
				if err != nil {
					return steps, err
				}
				steps = append(steps, &step)
			case step.KindCreateFolder:
				step := stepDefinition.CreateFolder{}
				err = step.FromConfig(stepYaml)
				if err != nil {
					return steps, err
				}
				steps = append(steps, &step)
			case step.KindCommand:
				step := stepDefinition.Command{}
				err = step.FromConfig(stepYaml)
				if err != nil {
					return steps, err
				}
				steps = append(steps, &step)
			case step.KindService:
				step := stepDefinition.Service{}
				err = step.FromConfig(stepYaml)
				if err != nil {
					return steps, err
				}
				steps = append(steps, &step)
			case step.KindChecklist:
				step := stepDefinition.Checklist{}
				err = step.FromConfig(stepYaml)
				if err != nil {
					return steps, err
				}
				steps = append(steps, &step)
			case step.KindWatch:
				step := stepDefinition.Watch{}
				err = step.FromConfig(stepYaml)
				if err != nil {
					return steps, err
				}
				steps = append(steps, &step)
			}
		}
	}
	return steps, nil
}

// expandPreferences checks each preference if it contains a variable "${...}" if it finds one
// it will check that this variable has been defined and fills it accordingly. If a string contains
// a variable that is not defined an error is returned.
func expandPreferences(preferences map[string]interface{}, variables map[string]string) (expandedPreferences map[string]interface{}, err error) {
	// FIXME: Add check to handle empty preferences map
	expandedPreferences = make(map[string]interface{})
	for key, value := range preferences {
		// A slice has to be handled differently than a map
		switch value.(type) {
		default:
			log.Println("Unknown preference type in config")
		case []interface{}:
			// For an array
			for k, v := range value.([]interface{}) {
				value.([]interface{})[k], err = expandVariable(v.(string), variables)
				if err != nil {
					return expandedPreferences, err
				}
			}
			expandedPreferences[key] = value
			// For a single field
		case interface{}:
			// Handle string variables (Only string variables can be expanded)
			if _, ok := value.(string); ok {
				expandedPreferences[key], err = expandVariable(value.(string), variables)
				if err != nil {
					return expandedPreferences, err
				}
				// Handle bool variables
			} else if _, ok := value.(bool); ok {
				expandedPreferences[key] = value
			}
		}
	}
	return expandedPreferences, nil
}

func expandVariable(preference string, variables map[string]string) (expandedPreference string, err error) {
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
		if variableValue, ok := variables[variableName]; ok {
			expandedPreference = strings.Replace(expandedPreference, h, variableValue, -1)
		} else {
			return expandedPreference, fmt.Errorf("variable \"%s\" used but not defined", variableName)
		}
	}
	return expandedPreference, nil
}
