package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"aliasly/internal/alias"
)

// removeCmd represents the remove command.
// It deletes an existing alias after confirmation.
var removeCmd = &cobra.Command{
	// Use shows the expected arguments
	Use: "remove <alias-name>",

	// Aliases for shorter typing
	Aliases: []string{"rm", "delete", "del"},

	// Short description
	Short: "Remove an existing alias",

	// Long description
	Long: `Remove an existing alias from your configuration.

You will be asked to confirm before the alias is deleted.

Examples:
  al remove gs     # Remove the 'gs' alias
  al rm deploy     # Short form
  al delete old    # Alternative form`,

	// Args validates that exactly one argument is provided
	Args: cobra.ExactArgs(1),

	// Run function
	Run: runRemoveCmd,
}

// runRemoveCmd executes the remove command.
func runRemoveCmd(cmd *cobra.Command, args []string) {
	// Get the alias name from arguments
	aliasName := args[0]

	// Check if alias exists
	a, exists := alias.Find(aliasName)
	if !exists {
		printError(fmt.Sprintf("Alias '%s' not found", aliasName))
		fmt.Println()
		fmt.Println("Run 'al list' to see all available aliases")
		os.Exit(1)
	}

	// Show what we're about to delete
	fmt.Printf("Alias: %s\n", a.Name)
	fmt.Printf("Command: %s\n", a.Command)
	if a.Description != "" {
		fmt.Printf("Description: %s\n", a.Description)
	}
	fmt.Println()

	// Ask for confirmation
	confirmed, err := confirmDelete(aliasName)
	if err != nil {
		handlePromptError(err)
		return
	}

	if !confirmed {
		fmt.Println("Cancelled. Alias was not removed.")
		return
	}

	// Remove the alias
	if err := alias.Remove(aliasName); err != nil {
		printError(fmt.Sprintf("Failed to remove alias: %v", err))
		os.Exit(1)
	}

	// Success message
	green := color.New(color.FgGreen, color.Bold)
	green.Printf("Alias '%s' removed successfully!\n", aliasName)
}

// confirmDelete asks the user to confirm deletion.
func confirmDelete(aliasName string) (bool, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Are you sure you want to remove '%s'?", aliasName),
		Items: []string{"No, keep it", "Yes, remove it"},
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return false, err
	}

	// idx 1 = "Yes, remove it"
	return idx == 1, nil
}
