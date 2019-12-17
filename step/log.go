package step

import (
	"fmt"
	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/hpcloud/tail"
	log "github.com/sirupsen/logrus"
)

// KindWatch defines the name for a log step in the config file
const KindWatch = "Watch"

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

	if preferencesYaml, ok := configYaml["Preferences"].(map[string]interface{}); ok {
		if path, ok := preferencesYaml["Path"].(string); ok {
			watch.Path = path
		} else {
			return fmt.Errorf("preference \"Path\" (string) must be set for %s step", watch.GetKind())
		}
		if message, ok := preferencesYaml["Message"].(string); ok {
			watch.Message = message
		} else {
			return fmt.Errorf("preference \"Message\" (string) must be set for %s step", watch.GetKind())
		}
	} else {
		return fmt.Errorf("preferences must be set for %s step", watch.GetKind())
	}

	return nil
}

// Execute starts the file watcher and terminates it if the defined string appeared or a timeout is reached.
func (watch Watch) Execute() (result model.Result) {
	fmt.Printf("Start watching \"%s\" for \"%s\"\n", watch.Path, watch.Message)

	// Setup the channel we listen for the result
	done := make(chan error)

	// Tail the file
	go tailFile(watch.Path, tail.Config{Follow: true}, done)

	// Get the result of the channel
	err := <-done

	fmt.Printf("result: %s\n", err)

	if err != nil {
		result.WasSuccessful = false
		result.Message = err.Error()
		return result
	}

	result.WasSuccessful = true
	result.Message = "Message found"
	return result
}

func tailFile(filename string, config tail.Config, done chan error) {
	t, err := tail.TailFile(filename, config)
	if err != nil {
		log.Println(err)
		done <- err
		return
	}
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
	err = t.Wait()
	if err != nil {
		log.Println(err)
		done <- err
	}
}
