package step

import (
	"github.com/ChristianHuff-DEV/reapy/model"
)

// KindWatch defines the name for a log step in the config file
const KindWatch = "Log"

// Watch is a step used to watch a file
type Watch struct {
	model.RunnableStep
	// Path is the location the file to be watched can be found
	Path string
	// Message is the string which the file is watched for
	Message string
}

// GetKind returns the type this step is of
func (watch Watch) GetKind() string {
	return watch.Kind
}

// GetDescription returns a summary of what this step does
func (watch Watch) GetDescription() string {
	return watch.Description
}

// FromConfig creates the struct representation of a step watching a file for a certain string
func (watch *Watch) FromConfig(configYaml map[string]interface{}) error {
	watch.Kind = KindWatch
	if description, ok := configYaml["Description"]; ok {
		watch.Description = description.(string)
	}
	// preferencesYaml := configYaml["Preferences"].(map[string]interface{})
	return nil
}

// Execute starts the file watcher and terminates it if the defined string appeared or a timeout is reached.
func (watch Watch) Execute() (result model.Result) {
	return result
}
