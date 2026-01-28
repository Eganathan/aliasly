// Package alias provides functionality for working with command aliases.
// This includes parsing parameters, validating aliases, and executing commands.
package alias

import (
	"aliasly/internal/config"
)

// Alias is a type alias (pun intended!) for config.Alias.
// This allows us to add methods to the Alias type in this package
// while keeping the actual struct definition in the config package.
type Alias = config.Alias

// Param is a type alias for config.Param.
type Param = config.Param

// Find looks up an alias by name and returns it if found.
// This is a convenience wrapper around config.FindAlias.
func Find(name string) (Alias, bool) {
	return config.FindAlias(name)
}

// GetAll returns all configured aliases.
// This is a convenience wrapper around config.GetAllAliases.
func GetAll() ([]Alias, error) {
	return config.GetAllAliases()
}

// Add creates a new alias.
// Returns an error if the alias name is already taken.
func Add(alias Alias) error {
	return config.AddAlias(alias)
}

// Remove deletes an alias by name.
// Returns an error if the alias doesn't exist.
func Remove(name string) error {
	return config.RemoveAlias(name)
}

// Update modifies an existing alias.
// Returns an error if the alias doesn't exist.
func Update(alias Alias) error {
	return config.UpdateAlias(alias)
}

// GetParamNames returns a list of all parameter names for an alias.
// This is useful for displaying help text or validating input.
func GetParamNames(a Alias) []string {
	names := make([]string, len(a.Params))
	for i, p := range a.Params {
		names[i] = p.Name
	}
	return names
}

// GetRequiredParams returns only the required parameters for an alias.
func GetRequiredParams(a Alias) []Param {
	required := make([]Param, 0)
	for _, p := range a.Params {
		if p.Required {
			required = append(required, p)
		}
	}
	return required
}

// GetOptionalParams returns only the optional parameters for an alias.
func GetOptionalParams(a Alias) []Param {
	optional := make([]Param, 0)
	for _, p := range a.Params {
		if !p.Required {
			optional = append(optional, p)
		}
	}
	return optional
}

// BuildUsageString creates a usage string for an alias.
// Example: "gc <message>" or "gp [branch]"
// Required params are shown in <angle brackets>, optional in [square brackets].
func BuildUsageString(a Alias) string {
	usage := a.Name

	for _, p := range a.Params {
		if p.Required {
			usage += " <" + p.Name + ">"
		} else {
			usage += " [" + p.Name + "]"
		}
	}

	return usage
}
