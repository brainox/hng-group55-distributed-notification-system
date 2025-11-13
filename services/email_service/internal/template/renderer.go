package template

import (
	"fmt"
	"regexp"
	"strings"
)

// RenderTemplate replaces {{variable}} placeholders with actual values
func RenderTemplate(template string, variables map[string]interface{}) (string, error) {
	result := template

	// Find all {{variable}} patterns
	re := regexp.MustCompile(`\{\{(\w+)\}\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	missingVars := []string{}

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		varName := match[1]
		placeholder := match[0]

		value, exists := variables[varName]
		if !exists {
			missingVars = append(missingVars, varName)
			continue
		}

		// Convert value to string
		valueStr := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, valueStr)
	}

	if len(missingVars) > 0 {
		return "", fmt.Errorf("missing required variables: %v", missingVars)
	}

	return result, nil
}
