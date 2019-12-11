package cli

import (
	"fmt"
	"regexp"
)

// validateVariables ensures that all defined variables are in the correct format
func validateVariables(variablesYaml map[string]interface{}) error {
	for key, value := range variablesYaml {
		// Every varable must have a value
		if value == nil || value == "" {
			return fmt.Errorf("the variable \"%s\" has no value assigned to it", key)
		}
		err := validateVariableName(key)
		if err != nil {
			return err
		}
	}
	return nil
}

// validateVariableName ensures that the given name contains only letter, number, "-" and "_".
func validateVariableName(name string) error {
	regex, _ := regexp.Compile("^[_A-z0-9]*((-|_)*[_A-z0-9])*$")
	isValid := regex.MatchString(name)

	if !isValid {
		return fmt.Errorf("the variable name \"%s\" is invalid. Variable names can only consist out of letters, numbers, \"-\" and \"_\"", name)
	}
	return nil
}
