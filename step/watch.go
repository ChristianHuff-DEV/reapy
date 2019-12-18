package step

import (
	"fmt"
	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/hpcloud/tail"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
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
	// Timeout in seconds after which watching the file will be interrupted (default time out is )
	Timeout int
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
		// Read path preference
		if path, ok := preferencesYaml["Path"].(string); ok {
			watch.Path = path
		} else {
			return fmt.Errorf("preference \"Path\" (string) must be set for %s step", watch.GetKind())
		}
		// Read message preference
		if message, ok := preferencesYaml["Message"].(string); ok {
			watch.Message = message
		} else {
			return fmt.Errorf("preference \"Message\" (string) must be set for %s step", watch.GetKind())
		}
		// Read timeout preference
		if timeout, ok := preferencesYaml["Timeout"].(int); ok {
			watch.Timeout = timeout
		} else {
			// Default timeout is 300s
			watch.Timeout = 300
		}
	} else {
		return fmt.Errorf("preferences must be set for %s step", watch.GetKind())
	}

	return nil
}

// Execute starts the file watcher and terminates it if the defined string appeared or a timeout is reached.
func (watch Watch) Execute() (result model.Result) {
	fmt.Printf("Start watching \"%s\" for \"%s\" (Timeout: %ds)\n", watch.Path, watch.Message, watch.Timeout)

	// Setup the channel we listen for the result
	watcher := make(chan error)

	// Tail the file in a go routine
	go tailFile(watch.Path, watch.Message, tail.Config{Follow: true}, watcher)

	// Check for the result or timeout if it takes to long
	select {
	case err := <-watcher:
		if err != nil {
			fmt.Println("Message was not found")
			result.WasSuccessful = false
			result.Message = err.Error()
			return result
		}
	case <-time.After(time.Duration(watch.Timeout) * time.Second):
		result.WasSuccessful = false
		result.Message = "Timeout waiting for message"
		return result
	}

	// Close the channel
	close(watcher)

	fmt.Println("Message found")
	result.WasSuccessful = true
	result.Message = "Message found"
	return result
}

// tailFile checks the given file for the given message. It start at the beginning of the file and will return once the string is found. If the string is not yet in the file it will monitor the file and only return once it find the string.
func tailFile(file, message string, config tail.Config, done chan error) {
	t, err := tail.TailFile(file, config)
	if err != nil {
		log.Println(err)
		done <- err
		return
	}
	for line := range t.Lines {
		// Check if we can find the message we are looking for in the current line
		if strings.Contains(line.Text, message) {
			done <- nil
		}
	}
	err = t.Wait()
	if err != nil {
		log.Println(err)
		done <- err
	}
}
