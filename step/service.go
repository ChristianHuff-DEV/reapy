package step

import (
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
