package renderer

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

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		varName := match[1]
		placeholder := match[0]

		value, exists := variables[varName]
		if !exists {
			return "", fmt.Errorf("missing required variable: %s", varName)
		}

		// Convert value to string
		valueStr := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, valueStr)
	}

	return result, nil
}

// ExtractVariables extracts all {{variable}} names from a template
func ExtractVariables(template string) []string {
	re := regexp.MustCompile(`\{\{(\w+)\}\}`)
	matches := re.FindAllStringSubmatch(template, -1)

	variables := []string{}
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		varName := match[1]
		if !seen[varName] {
			variables = append(variables, varName)
			seen[varName] = true
		}
	}

	return variables
}

// ValidateVariables checks if all required variables are provided
func ValidateVariables(required []string, provided map[string]interface{}) []string {
	missing := []string{}

	for _, req := range required {
		if _, exists := provided[req]; !exists {
			missing = append(missing, req)
		}
	}

	return missing
}
