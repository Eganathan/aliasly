package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"aliasly/internal/alias"
	"aliasly/internal/config"
)

// namePattern validates alias names.
// Alias names can only contain letters, numbers, and hyphens.
var namePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]*$`)

// addCmd represents the add command.
// It interactively guides the user through creating a new alias.
var addCmd = &cobra.Command{
	// Use is the one-line usage
	Use: "add",

	// Aliases for shorter typing
	Aliases: []string{"a", "new"},

	// Short description
	Short: "Add a new alias interactively",

	// Long description
	Long: `Add a new alias through an interactive prompt.

You will be asked to provide:
  - Alias name (short name you'll type, e.g., "gs")
  - Command to run (the full command, e.g., "git status")
  - Description (optional, helps you remember what it does)
  - Parameters (optional, for commands that need input)

For parameterized commands, use {{name}} syntax in your command:
  git commit -am "{{message}}"

Examples:
  al add     # Start interactive alias creation
  al new     # Same as above`,

	// Run function
	Run: runAddCmd,
}

// runAddCmd executes the add command.
func runAddCmd(cmd *cobra.Command, args []string) {
	fmt.Println("Create a new alias")
	fmt.Println("------------------")
	fmt.Println()

	// Step 1: Get alias name
	name, err := promptAliasName()
	if err != nil {
		handlePromptError(err)
		return
	}

	// Step 2: Get command
	command, err := promptCommand()
	if err != nil {
		handlePromptError(err)
		return
	}

	// Step 3: Get description
	description, err := promptDescription()
	if err != nil {
		handlePromptError(err)
		return
	}

	// Step 4: Get parameters (if any {{placeholders}} in command)
	params, err := promptParams(command)
	if err != nil {
		handlePromptError(err)
		return
	}

	// Create the alias
	newAlias := config.Alias{
		Name:        name,
		Command:     command,
		Description: description,
		Params:      params,
	}

	// Save the alias
	if err := alias.Add(newAlias); err != nil {
		printError(fmt.Sprintf("Failed to save alias: %v", err))
		os.Exit(1)
	}

	// Success message
	fmt.Println()
	green := color.New(color.FgGreen, color.Bold)
	green.Printf("Alias '%s' created successfully!\n", name)
	fmt.Println()
	fmt.Printf("Usage: al %s\n", alias.BuildUsageString(newAlias))
}

// promptAliasName asks the user for the alias name.
func promptAliasName() (string, error) {
	// Create a prompt with validation
	prompt := promptui.Prompt{
		Label: "Alias name",
		Validate: func(input string) error {
			// Check if name is valid format
			if !namePattern.MatchString(input) {
				return fmt.Errorf("name must start with a letter and contain only letters, numbers, and hyphens")
			}

			// Check if alias already exists
			if _, exists := alias.Find(input); exists {
				return fmt.Errorf("alias '%s' already exists", input)
			}

			return nil
		},
	}

	return prompt.Run()
}

// promptCommand asks the user for the command to run.
func promptCommand() (string, error) {
	prompt := promptui.Prompt{
		Label: "Command",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("command cannot be empty")
			}
			return nil
		},
	}

	return prompt.Run()
}

// promptDescription asks for an optional description.
func promptDescription() (string, error) {
	prompt := promptui.Prompt{
		Label:   "Description (optional)",
		Default: "",
	}

	return prompt.Run()
}

// promptParams detects {{placeholders}} in the command and asks
// the user to define each parameter.
func promptParams(command string) ([]config.Param, error) {
	// Find all placeholders in the command
	placeholders := alias.ExtractPlaceholders(command)

	// If no placeholders, no parameters needed
	if len(placeholders) == 0 {
		return nil, nil
	}

	fmt.Println()
	fmt.Printf("Found %d parameter(s) in your command:\n", len(placeholders))

	params := make([]config.Param, 0, len(placeholders))

	// For each placeholder, gather parameter details
	for _, name := range placeholders {
		fmt.Printf("\nParameter: {{%s}}\n", name)

		param, err := promptParamDetails(name)
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	return params, nil
}

// promptParamDetails asks for details about a single parameter.
func promptParamDetails(name string) (config.Param, error) {
	// Get description
	descPrompt := promptui.Prompt{
		Label:   "Description",
		Default: "",
	}
	description, err := descPrompt.Run()
	if err != nil {
		return config.Param{}, err
	}

	// Ask if required
	requiredPrompt := promptui.Select{
		Label: "Is this parameter required?",
		Items: []string{"Yes (must be provided)", "No (optional)"},
	}
	requiredIdx, _, err := requiredPrompt.Run()
	if err != nil {
		return config.Param{}, err
	}
	required := requiredIdx == 0

	// If optional, ask for default value
	var defaultVal string
	if !required {
		defaultPrompt := promptui.Prompt{
			Label:   "Default value (leave empty for none)",
			Default: "",
		}
		defaultVal, err = defaultPrompt.Run()
		if err != nil {
			return config.Param{}, err
		}
	}

	return config.Param{
		Name:        name,
		Description: description,
		Required:    required,
		Default:     defaultVal,
	}, nil
}

// handlePromptError handles errors from promptui.
func handlePromptError(err error) {
	// promptui.ErrInterrupt is returned when user presses Ctrl+C
	if err == promptui.ErrInterrupt {
		fmt.Println("\nCancelled.")
		return
	}

	// promptui.ErrEOF is returned when user presses Ctrl+D
	if err == promptui.ErrEOF {
		fmt.Println("\nCancelled.")
		return
	}

	// Other errors
	printError(fmt.Sprintf("Prompt failed: %v", err))
	os.Exit(1)
}
