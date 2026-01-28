package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"aliasly/internal/alias"
)

// listCmd represents the list command.
// It displays all configured aliases in a formatted table.
var listCmd = &cobra.Command{
	// Use is the one-line usage for this command
	Use: "list",

	// Aliases are alternative names for this command
	// So "al ls" works the same as "al list"
	Aliases: []string{"ls", "l"},

	// Short is a short description shown in help output
	Short: "List all configured aliases",

	// Long is the detailed description
	Long: `List all configured aliases in a formatted table.

Shows the alias name, the command it runs, and a description.
Parameters are shown in the command with {{name}} syntax.

Examples:
  al list    # Show all aliases
  al ls      # Short form`,

	// Run is the function to execute
	Run: runListCmd,
}

// runListCmd executes the list command.
func runListCmd(cmd *cobra.Command, args []string) {
	// Get all aliases from config
	aliases, err := alias.GetAll()
	if err != nil {
		printError(fmt.Sprintf("Failed to load aliases: %v", err))
		os.Exit(1)
	}

	// Check if there are any aliases
	if len(aliases) == 0 {
		fmt.Println("No aliases configured yet.")
		fmt.Println()
		fmt.Println("Run 'al add' to create your first alias")
		fmt.Println("Or run 'al config' to open the web configuration UI")
		return
	}

	// Print a header
	fmt.Printf("Found %d alias(es):\n\n", len(aliases))

	// Print each alias
	for _, a := range aliases {
		printAlias(a)
	}

	// Print help footer
	fmt.Println()
	fmt.Println("Run 'al <alias>' to execute an alias")
	fmt.Println("Run 'al add' to create a new alias")
	fmt.Println("Run 'al remove <alias>' to delete an alias")
}

// printAlias prints a single alias in a nice format.
func printAlias(a alias.Alias) {
	// Create colored output
	nameColor := color.New(color.FgCyan, color.Bold)
	cmdColor := color.New(color.FgGreen)
	dimColor := color.New(color.Faint)

	// Print alias name (bold cyan)
	nameColor.Printf("  %s", a.Name)

	// Print description if present (dim)
	if a.Description != "" {
		dimColor.Printf(" - %s", a.Description)
	}
	fmt.Println()

	// Print the command (green)
	cmdColor.Printf("    $ %s\n", a.Command)

	// Print parameters if any
	if len(a.Params) > 0 {
		// Build params string
		paramStrs := make([]string, 0, len(a.Params))
		for _, p := range a.Params {
			paramStr := p.Name
			if p.Required {
				paramStr += "*" // Asterisk indicates required
			}
			if p.Default != "" {
				paramStr += fmt.Sprintf("=%s", p.Default)
			}
			paramStrs = append(paramStrs, paramStr)
		}

		dimColor.Printf("    params: %s\n", strings.Join(paramStrs, ", "))
	}

	// Print usage example
	usageStr := alias.BuildUsageString(a)
	dimColor.Printf("    usage:  al %s\n", usageStr)

	fmt.Println() // Empty line between aliases
}
