package config

import "github.com/ChristianHuff-DEV/reapy/model"

// Execute runs all steps in all tasks of the given plan
func Execute(steps model.Plan) {
	for _, task := range steps.Tasks {
		for _, step := range task.Steps {
			step.Execute()
		}
	}
}
