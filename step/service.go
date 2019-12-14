package step

import (
	"github.com/ChristianHuff-DEV/reapy/model"
)

// KindService defines the name for a step that can start/stop a service
const KindService = "Service"

// Service represents a step used to start/stop a service
type Service struct {
	model.RunnableStep
	// Name of the Servie according to it's preferences. (Attention: In Windows there is a "service name" and a "display name". Here the "service name" has to be used.)
	Name string
	// Action whether to "start/stop" the service
	Action string
}

// GetKind returns the type this step is of
func (service Service) GetKind() string {
	return service.Kind
}

// GetDescription returns a text summarizing what this step does
func (service Service) GetDescription() string {
	return service.Description
}
