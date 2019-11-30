package config

import "github.com/ChristianHuff-DEV/reapy/model"

func Execute(steps model.Plan) {
	for _, task := range steps.Tasks {
		for _, step := range task.Steps {
			step.Execute()
		}
	}
}
