package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

// Config represents the root configuration structure for aliasly.
// It contains application settings and all defined aliases.
type Config struct {
	// Version is the config file format version (for future migrations)
	Version int `mapstructure:"version" yaml:"version" json:"version"`

	// Settings contains global application settings
	Settings Settings `mapstructure:"settings" yaml:"settings" json:"settings"`

	// Aliases is the list of all defined command aliases
	Aliases []Alias `mapstructure:"aliases" yaml:"aliases" json:"aliases"`
}

// Settings contains global configuration options that affect
// how aliasly behaves when running commands.
type Settings struct {
	// Shell is the shell to use for executing commands (e.g., "/bin/bash")
	// If empty, the default shell will be detected automatically
	Shell string `mapstructure:"shell" yaml:"shell" json:"shell"`

	// Verbose, when true, prints the expanded command before running it
	Verbose bool `mapstructure:"verbose" yaml:"verbose" json:"verbose"`
}

// Alias represents a single command alias.
// An alias maps a short name to a longer command, optionally with parameters.
type Alias struct {
	// Name is the short name for the alias (e.g., "gs" for git status)
	Name string `mapstructure:"name" yaml:"name" json:"name"`

	// Command is the actual command to run, may contain {{param}} placeholders
	Command string `mapstructure:"command" yaml:"command" json:"command"`

	// Description is a human-readable explanation of what this alias does
	Description string `mapstructure:"description" yaml:"description" json:"description"`

	// Params defines the parameters that this alias accepts
	Params []Param `mapstructure:"params" yaml:"params,omitempty" json:"params,omitempty"`
}

// Param represents a parameter that can be passed to an alias.
// Parameters are substituted into the command using {{paramName}} syntax.
type Param struct {
	// Name is the parameter name, used in {{name}} placeholders
	Name string `mapstructure:"name" yaml:"name" json:"name"`

	// Description explains what this parameter is for
	Description string `mapstructure:"description" yaml:"description" json:"description"`

	// Required, when true, means this parameter must be provided
	Required bool `mapstructure:"required" yaml:"required" json:"required"`

	// Default is the value to use if the parameter is not provided
	// Only used when Required is false
	Default string `mapstructure:"default" yaml:"default,omitempty" json:"default,omitempty"`
}

// globalConfig holds the currently loaded configuration.
// We use a package-level variable so all parts of the app can access it.
var globalConfig *Config

// configMutex protects concurrent access to globalConfig.
// This is important if multiple goroutines might read/write config.
var configMutex sync.RWMutex

// loaded tracks whether config has been loaded
var loaded bool

// Load reads the configuration from disk and stores it in memory.
// If the config file doesn't exist, it creates a default one.
// Returns an error if the config cannot be read or parsed.
func Load() error {
	configMutex.Lock()
	defer configMutex.Unlock()

	return loadInternal()
}

// loadInternal is the internal load function that assumes the lock is already held.
func loadInternal() error {
	// Ensure the config directory exists before trying to read/write
	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := GetConfigFilePath()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config doesn't exist, create a default one
		globalConfig = createDefaultConfig()
		loaded = true
		return saveInternal()
	}

	// Set up Viper to read our config file
	// Viper is a popular Go library for configuration management
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal (convert) the YAML into our Config struct
	globalConfig = &Config{}
	if err := viper.Unmarshal(globalConfig); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	loaded = true
	return nil
}

// Save writes the current configuration to disk.
// It creates the config file if it doesn't exist.
func Save() error {
	configMutex.Lock()
	defer configMutex.Unlock()

	return saveInternal()
}

// saveInternal is the internal save function that assumes the lock is already held.
// This prevents deadlocks when called from loadInternal() or other functions.
func saveInternal() error {
	if globalConfig == nil {
		return fmt.Errorf("no configuration loaded")
	}

	// Ensure config directory exists
	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal (convert) our Config struct to YAML format
	// yaml.Marshal converts Go structs to YAML text
	data, err := yaml.Marshal(globalConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write the YAML to the config file
	// 0644 = rw-r--r-- (owner can read/write, others can read)
	configPath := GetConfigFilePath()
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ensureLoaded makes sure the config is loaded before proceeding.
// Must be called while holding the write lock.
func ensureLoaded() error {
	if !loaded {
		return loadInternal()
	}
	return nil
}

// Get returns the current configuration.
// It loads the config from disk if not already loaded.
// Returns an error if the config cannot be loaded.
func Get() (*Config, error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	if err := ensureLoaded(); err != nil {
		return nil, err
	}

	return globalConfig, nil
}

// FindAlias searches for an alias by name.
// Returns the alias and true if found, or an empty alias and false if not found.
func FindAlias(name string) (Alias, bool) {
	configMutex.Lock()
	defer configMutex.Unlock()

	if err := ensureLoaded(); err != nil {
		return Alias{}, false
	}

	// Linear search through aliases
	// For a typical number of aliases (< 100), this is fast enough
	for _, alias := range globalConfig.Aliases {
		if alias.Name == name {
			return alias, true
		}
	}

	return Alias{}, false
}

// AddAlias adds a new alias to the configuration.
// Returns an error if an alias with the same name already exists.
func AddAlias(alias Alias) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	if err := ensureLoaded(); err != nil {
		return err
	}

	// Check if alias already exists
	for _, a := range globalConfig.Aliases {
		if a.Name == alias.Name {
			return fmt.Errorf("alias '%s' already exists", alias.Name)
		}
	}

	globalConfig.Aliases = append(globalConfig.Aliases, alias)

	return saveInternal()
}

// RemoveAlias removes an alias from the configuration by name.
// Returns an error if the alias doesn't exist.
func RemoveAlias(name string) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	if err := ensureLoaded(); err != nil {
		return err
	}

	// Find and remove the alias
	found := false
	newAliases := make([]Alias, 0, len(globalConfig.Aliases))
	for _, alias := range globalConfig.Aliases {
		if alias.Name == name {
			found = true
			continue // Skip this alias (remove it)
		}
		newAliases = append(newAliases, alias)
	}

	if !found {
		return fmt.Errorf("alias '%s' not found", name)
	}

	globalConfig.Aliases = newAliases

	return saveInternal()
}

// UpdateAlias updates an existing alias in the configuration.
// Returns an error if the alias doesn't exist.
func UpdateAlias(alias Alias) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	if err := ensureLoaded(); err != nil {
		return err
	}

	// Find and update the alias
	found := false
	for i, a := range globalConfig.Aliases {
		if a.Name == alias.Name {
			globalConfig.Aliases[i] = alias
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("alias '%s' not found", alias.Name)
	}

	return saveInternal()
}

// GetAllAliases returns a copy of all aliases.
func GetAllAliases() ([]Alias, error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	if err := ensureLoaded(); err != nil {
		return nil, err
	}

	// Return a copy to prevent external modification
	aliases := make([]Alias, len(globalConfig.Aliases))
	copy(aliases, globalConfig.Aliases)

	return aliases, nil
}

// createDefaultConfig creates a new Config with sensible defaults
// and some example aliases to help users get started.
func createDefaultConfig() *Config {
	return &Config{
		Version: 1,
		Settings: Settings{
			Shell:   GetDefaultShell(),
			Verbose: false,
		},
		Aliases: []Alias{
			{
				Name:        "gs",
				Command:     "git status",
				Description: "Show git status",
			},
			{
				Name:        "gc",
				Command:     `git commit -am "{{message}}"`,
				Description: "Git commit with message",
				Params: []Param{
					{
						Name:        "message",
						Description: "Commit message",
						Required:    true,
					},
				},
			},
			{
				Name:        "gp",
				Command:     "git push origin {{branch}}",
				Description: "Push to remote branch",
				Params: []Param{
					{
						Name:        "branch",
						Description: "Branch name",
						Required:    false,
						Default:     "main",
					},
				},
			},
		},
	}
}
