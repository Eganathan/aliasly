// Package config handles all configuration-related functionality for aliasly.
// This includes finding the config file path, loading, and saving configuration.
package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetConfigDir returns the directory where aliasly configuration should be stored.
// It follows the XDG Base Directory Specification on Linux/macOS:
//
//  1. If ALIASLY_CONFIG_DIR environment variable is set, use that
//  2. If XDG_CONFIG_HOME is set, use $XDG_CONFIG_HOME/aliasly
//  3. Otherwise, use $HOME/.config/aliasly
//
// This ensures the config is stored in a standard, predictable location.
func GetConfigDir() string {
	// Check if user has explicitly set a config directory via environment variable
	// This allows power users to customize where their config lives
	if envDir := os.Getenv("ALIASLY_CONFIG_DIR"); envDir != "" {
		return envDir
	}

	// Get the user's home directory
	// os.UserHomeDir() works cross-platform (macOS, Linux, Windows)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If we can't get home dir, fall back to current directory
		// This shouldn't happen in normal circumstances
		return "."
	}

	// Check if XDG_CONFIG_HOME is set (common on Linux)
	// XDG is a standard for where config files should live
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "aliasly")
	}

	// Default: use ~/.config/aliasly
	// filepath.Join handles path separators correctly for each OS
	return filepath.Join(homeDir, ".config", "aliasly")
}

// GetConfigFilePath returns the full path to the config file.
// The config file is always named "config.yaml" inside the config directory.
func GetConfigFilePath() string {
	return filepath.Join(GetConfigDir(), "config.yaml")
}

// EnsureConfigDir creates the config directory if it doesn't exist.
// It uses 0755 permissions (owner can read/write/execute, others can read/execute).
// Returns an error if the directory cannot be created.
func EnsureConfigDir() error {
	configDir := GetConfigDir()

	// os.MkdirAll creates the directory and any necessary parents
	// It's safe to call even if the directory already exists
	// 0755 = rwxr-xr-x (owner full access, others read+execute)
	return os.MkdirAll(configDir, 0755)
}

// GetDefaultShell returns the default shell for the current operating system.
// This is used when executing alias commands.
func GetDefaultShell() string {
	// First, check if user has a preferred shell set via SHELL env var
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}

	// Fall back to OS-specific defaults
	switch runtime.GOOS {
	case "windows":
		// On Windows, use cmd.exe as the default
		return "cmd"
	default:
		// On Unix-like systems (macOS, Linux), use /bin/sh
		// /bin/sh is POSIX-compliant and available on all Unix systems
		return "/bin/sh"
	}
}
