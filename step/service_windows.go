package step

import (
	"errors"
	"fmt"
	"time"

	"github.com/ChristianHuff-DEV/reapy/model"
	"github.com/briandowns/spinner"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// FromConfig creates the struct representation of a service step
func (service *Service) FromConfig(configYaml map[string]interface{}) error {
	service.Kind = KindService
	if description, ok := configYaml["Description"]; ok {
		service.Description = description.(string)
	}

	if preferencesYaml, ok := configYaml["Preferences"].(map[string]interface{}); ok {
		if name, ok := preferencesYaml["Name"]; ok {
			service.Name = name.(string)
		}
		if action, ok := preferencesYaml["Action"].(string); ok {
			if action != "start" && action != "stop" {
				return fmt.Errorf("action for service task has to be \"start\" or \"stop\"")
			}
			service.Action = action
		}
	} else {
		return fmt.Errorf("the preferences \"Name\" and \"Action\" have to he defined for step type %s", KindService)
	}
	return nil
}

// Execute starts/stops a service
func (service *Service) Execute() (result model.Result) {
	fmt.Println(service.Description)

	switch service.Action {
	case "start":
		if err := startService(service.Name); err != nil {
			result.WasSuccessful = false
			result.Message = err.Error()
		} else {
			result.WasSuccessful = true
			result.Message = "Service \"" + service.Name + "\" started"
		}
		return result
	case "stop":
		if err := stopService(service.Name); err != nil {
			result.WasSuccessful = false
			result.Message = err.Error()
		} else {
			result.WasSuccessful = true
			result.Message = "Service \"" + service.Name + "\" stopped"
		}
		return result
	default:
		result.WasSuccessful = false
		result.Message = fmt.Sprintf("unknown action of type %s", service.Action)
		return result
	}
}

func stopService(serviceName string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()

	status, err := s.Control(svc.Stop)
	if err != nil {
		return err
	}

	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spinner.Start()
	timeout := time.Now().Add(120 * time.Second)
	for status.State != svc.Stopped {
		if timeout.Before(time.Now()) {
			spinner.Stop()
			return err
		}
		time.Sleep(5 * time.Second)
		status, err = s.Query()
		if err != nil {
			spinner.Stop()
			return err
		}
	}
	spinner.Stop()

	return nil
}

func startService(serviceName string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()

	err = s.Start("is", "manual-started")
	if err != nil {
		return err
	}
	// Check that the service actually started and keeps running
	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)

	spinner.Start()
	// To start a servie for now we use a fixed timeout. This is due to the fact that the service
	// might start but than quickly fail. Therefore only after the sleep we check that the service
	// is actually running.
	time.Sleep(2 * time.Minute)
	status, err := s.Query()
	if err != nil {
		return err
	}

	// Check that the service is running
	if status.State != svc.Running {
		fmt.Println("Service could not be started")
		spinner.Stop()
		return errors.New("Service stopped running")
	}
	spinner.Stop()

	return nil
}
