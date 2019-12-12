package step

import (
	"fmt"

	"github.com/ChristianHuff-DEV/reapy/model"
)

// KindService defines the name for a step that can start/stop a service
const KindService = "Service"

// Service represents a step used to start/stop a service
type Service struct {
	model.RunnableStep
	ServiceName string
	Action      string
}

// GetKind returns the type this step is of
func (service Service) GetKind() string {
	return service.Kind
}

// GetDescription returns a text summarizing what this step does
func (service Service) GetDescription() string {
	return service.Description
}

// FromConfig creates the struct representation of a service step
func (service *Service) FromConfig(configYaml map[string]interface{}) error {
	return fmt.Errorf("the step type \"%s\" is not supported on this operating system", KindService)
}

// Execute starts/stops a service
func (service *Service) Execute() model.Result {
	return model.Result{WasSuccessful: false, Message: "the step type is not supported on this operating system"}
}
