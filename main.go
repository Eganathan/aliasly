// Package main is the entry point for the aliasly CLI application.
// Aliasly (invoked as 'al') is a command-line tool that simplifies
// running frequently used commands through customizable aliases.
//
// Example usage:
//   al gs          -> runs 'git status'
//   al gc "message" -> runs 'git commit -am "message"'
package main

import (
	// Import our cmd package which contains all CLI commands
	"aliasly/cmd"
)

// main is the entry point of the application.
// It simply calls the Execute function from our cmd package,
// which sets up and runs the Cobra CLI framework.
func main() {
	cmd.Execute()
}
