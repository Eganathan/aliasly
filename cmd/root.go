// Package cmd contains all CLI commands for aliasly.
// It uses the Cobra library to define commands and handle arguments.
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"aliasly/internal/alias"
	"aliasly/internal/config"
)

// Version is the current version of aliasly.
// This can be set at build time using -ldflags.
var Version = "0.1.0"

// rootCmd is the base command when called without any subcommands.
// When the user runs just "al", this command's help is displayed.
// When the user runs "al <something>", we check if <something> is:
//   1. A subcommand (list, add, remove, config)
//   2. An alias name
var rootCmd = &cobra.Command{
	// Use is the one-line usage for this command
	Use: "al [alias] [params...]",

	// Short is a short description shown in the 'help' output
	Short: "Aliasly - A simple command alias manager",

	// Long is the long description shown in the 'help al' output
	Long: `Aliasly (al) is a command-line tool that simplifies running
frequently used commands through customizable aliases.

Instead of typing long commands like:
  git commit -am "fix bug"

You can create an alias and run:
  al gc "fix bug"

Examples:
  al gs              # Run the 'gs' alias (e.g., git status)
  al gc "message"    # Run 'gc' alias with a parameter
  al list            # List all configured aliases
  al add             # Interactively add a new alias
  al config          # Open web UI to manage aliases`,

	// Version will be printed when user runs "al --version"
	Version: Version,

	// Args configures how many arguments this command accepts
	// We use ArbitraryArgs because we accept any number of arguments
	Args: cobra.ArbitraryArgs,

	// SilenceUsage prevents printing usage on errors
	// We handle our own error messages
	SilenceUsage: true,

	// SilenceErrors prevents Cobra from printing errors
	// We handle errors ourselves for better formatting
	SilenceErrors: true,

	// Run is the function to execute when this command is called.
	// This is where we handle running aliases.
	Run: runRootCmd,
}

// runRootCmd is called when the user runs "al <alias> [params...]"
func runRootCmd(cmd *cobra.Command, args []string) {
	// If no arguments provided, show help
	if len(args) == 0 {
		cmd.Help()
		return
	}

	// The first argument should be the alias name
	aliasName := args[0]

	// Get the remaining arguments as parameters for the alias
	// args[1:] gives us everything except the first element
	params := args[1:]

	// Look up the alias
	a, found := alias.Find(aliasName)
	if !found {
		// Alias not found - show a helpful error message
		printError(fmt.Sprintf("Alias '%s' not found", aliasName))
		fmt.Println()
		fmt.Println("Run 'al list' to see available aliases")
		fmt.Println("Run 'al add' to create a new alias")
		os.Exit(1)
	}

	// Run the alias with the provided parameters
	exitCode, err := alias.Run(a, params)
	if err != nil {
		printError(err.Error())

		// If it's a parse error (missing params), show usage help
		if _, ok := err.(*alias.ParseError); ok {
			fmt.Println()
			printAliasUsage(a)
		}

		os.Exit(1)
	}

	// Exit with the same exit code as the executed command
	// This allows aliasly to be used in scripts
	os.Exit(exitCode)
}

// printError prints an error message in red.
func printError(message string) {
	// color.Red is a convenience function from the fatih/color package
	// It prints text in red to make errors stand out
	red := color.New(color.FgRed, color.Bold)
	red.Fprintf(os.Stderr, "Error: %s\n", message)
}

// printAliasUsage prints how to use a specific alias.
func printAliasUsage(a alias.Alias) {
	fmt.Printf("Usage: al %s\n", alias.BuildUsageString(a))

	if a.Description != "" {
		fmt.Printf("       %s\n", a.Description)
	}

	// Show parameters if any
	if len(a.Params) > 0 {
		fmt.Println()
		fmt.Println("Parameters:")
		for _, p := range a.Params {
			requiredStr := ""
			if p.Required {
				requiredStr = " (required)"
			} else if p.Default != "" {
				requiredStr = fmt.Sprintf(" (default: %s)", p.Default)
			}
			fmt.Printf("  %-12s %s%s\n", p.Name, p.Description, requiredStr)
		}
	}
}

// Execute adds all child commands to the root command and runs the application.
// This is called by main.main(). It only needs to happen once.
func Execute() {
	// Load configuration before running any commands
	if err := config.Load(); err != nil {
		// If config can't be loaded, we still want to allow some commands
		// like "al --version" or "al --help"
		// So we just print a warning and continue
		fmt.Fprintf(os.Stderr, "Warning: Could not load config: %v\n", err)
	}

	// Execute the root command (this parses args and runs the appropriate command)
	if err := rootCmd.Execute(); err != nil {
		printError(err.Error())
		os.Exit(1)
	}
}

// init is a special Go function that runs automatically when the package loads.
// We use it to add subcommands to the root command.
func init() {
	// Add subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(configCmd)

	// Add global flags that apply to all commands
	// These can be accessed from any subcommand
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show commands before running them")
}
