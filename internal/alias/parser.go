package alias

import (
	"fmt"
	"regexp"
	"strings"
)

// paramPattern is a regular expression that matches {{paramName}} placeholders.
// The \w+ matches one or more word characters (letters, digits, underscore).
// For example, it will match: {{message}}, {{branch}}, {{version_number}}
var paramPattern = regexp.MustCompile(`\{\{(\w+)\}\}`)

// ParseError represents an error that occurred during command parsing.
// It provides detailed information about what went wrong.
type ParseError struct {
	// Message is the human-readable error message
	Message string

	// ParamName is the name of the parameter that caused the error (if applicable)
	ParamName string
}

// Error implements the error interface for ParseError.
// This allows ParseError to be used wherever a regular error is expected.
func (e *ParseError) Error() string {
	return e.Message
}

// ParseCommand takes an alias and a list of arguments, and returns
// the fully expanded command string with all parameters substituted.
//
// For example:
//   Alias command: git commit -am "{{message}}"
//   Args: ["fix bug"]
//   Result: git commit -am "fix bug"
//
// Returns an error if required parameters are missing.
func ParseCommand(a Alias, args []string) (string, error) {
	command := a.Command

	// Build a map of parameter name -> value from the provided arguments.
	// Arguments are positional, so args[0] goes to the first param, etc.
	provided := make(map[string]string)
	for i, param := range a.Params {
		if i < len(args) {
			provided[param.Name] = args[i]
		}
	}

	// Check that all required parameters are provided
	for _, param := range a.Params {
		_, hasValue := provided[param.Name]
		if param.Required && !hasValue {
			return "", &ParseError{
				Message:   fmt.Sprintf("missing required parameter: %s", param.Name),
				ParamName: param.Name,
			}
		}
	}

	// Substitute each parameter placeholder with its value
	for _, param := range a.Params {
		placeholder := fmt.Sprintf("{{%s}}", param.Name)

		// Get the value to substitute
		value, hasValue := provided[param.Name]
		if !hasValue {
			// Use default value for optional parameters
			value = param.Default
		}

		// Replace all occurrences of the placeholder with the value
		command = strings.ReplaceAll(command, placeholder, value)
	}

	return command, nil
}

// ExtractPlaceholders finds all {{paramName}} placeholders in a command string.
// Returns a list of parameter names (without the curly braces).
// This is useful for validating that all placeholders have corresponding params.
func ExtractPlaceholders(command string) []string {
	// FindAllStringSubmatch returns all matches, including capture groups.
	// For "{{foo}} and {{bar}}", it returns:
	// [["{{foo}}", "foo"], ["{{bar}}", "bar"]]
	matches := paramPattern.FindAllStringSubmatch(command, -1)

	// Extract just the parameter names (the captured group)
	names := make([]string, 0, len(matches))
	for _, match := range matches {
		// match[0] is the full match ({{name}}), match[1] is the capture group (name)
		if len(match) >= 2 {
			names = append(names, match[1])
		}
	}

	return names
}

// ValidatePlaceholders checks that all placeholders in a command
// have corresponding parameter definitions.
// Returns a list of undefined placeholders.
func ValidatePlaceholders(a Alias) []string {
	placeholders := ExtractPlaceholders(a.Command)

	// Build a set of defined parameter names for fast lookup
	defined := make(map[string]bool)
	for _, param := range a.Params {
		defined[param.Name] = true
	}

	// Find placeholders that don't have definitions
	undefined := make([]string, 0)
	for _, placeholder := range placeholders {
		if !defined[placeholder] {
			undefined = append(undefined, placeholder)
		}
	}

	return undefined
}

// FormatExample shows what a command would look like with example values.
// This is useful for displaying help text to users.
//
// For example:
//   Command: git commit -am "{{message}}"
//   Params: [message]
//   Result: git commit -am "your message here"
func FormatExample(a Alias) string {
	command := a.Command

	for _, param := range a.Params {
		placeholder := fmt.Sprintf("{{%s}}", param.Name)

		// Use a descriptive example value
		var exampleValue string
		if param.Default != "" {
			exampleValue = param.Default
		} else {
			exampleValue = "<" + param.Name + ">"
		}

		command = strings.ReplaceAll(command, placeholder, exampleValue)
	}

	return command
}
