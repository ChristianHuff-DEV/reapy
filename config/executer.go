package config

import "github.com/ChristianHuff-DEV/reapy/model"

import "log"

// Execute runs all steps in all tasks of the given plan
func Execute(steps model.Plan) {
	for _, task := range steps.Tasks {
		for _, step := range task.Steps {
			result := step.Execute()
			// Print what when wrong if an error occurred
			if !result.WasSuccessful {
				log.Print(result.Message)
			}
		}
	}
}
