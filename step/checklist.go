package step

import "github.com/ChristianHuff-DEV/reapy/model"

import "fmt"

// KindChecklist defines the name for a checklist step in the config file
const KindChecklist = "Checklist"

// Checklist is a step used to show the user a list of items to tick
type Checklist struct {
	model.RunnableStep
	// Items represents the individual checks the user has to tick
	Items []string
}

// GetKind returns the type this step is of
func (checklist Checklist) GetKind() string {
	return checklist.Kind
}

// GetDescription returns a summary of what this steps does
func (checklist Checklist) GetDescription() string {
	return checklist.Description
}

// FromConfig creates the struct representation of a step showing a checklist
func (checklist *Checklist) FromConfig(configYaml map[string]interface{}) error {
	checklist.Kind = KindChecklist
	if description, ok := configYaml["Description"]; ok {
		checklist.Description = description.(string)
	}

	preferencesYaml := configYaml["Preferences"].(map[string]interface{})

	if items, ok := preferencesYaml["Items"].([]interface{}); ok {
		for _, item := range items {
			if item, ok := item.(string); ok {
				checklist.Items = append(checklist.Items, item)
			} else {
				return fmt.Errorf("items in checklist must be of type string")
			}
		}
	} else {
		return fmt.Errorf("empty checklist found")
	}
	return nil
}

// Execute shows the checklist to the user
func (checklist Checklist) Execute() (result model.Result) {
	fmt.Println(checklist.GetDescription())

	result.WasSuccessful = true
	result.Message = "all items ticked"
	return result
}
