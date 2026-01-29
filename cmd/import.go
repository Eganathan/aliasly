package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"

	"aliasly/internal/config"
)

// importCmd represents the import command.
// It imports configuration from a file.
var importCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import aliases from a file",
	Long: `Import aliases from a YAML configuration file.

By default, this merges new aliases with your existing ones.
Existing aliases with the same name will be skipped.

Use --replace to completely replace your config instead.

Examples:
  al import backup.yaml           # Merge aliases from backup.yaml
  al import ~/my-aliases.yaml     # Merge from home directory
  al import backup.yaml --replace # Replace entire config`,

	Args: cobra.ExactArgs(1),
	Run:  runImportCmd,
}

// replaceFlag determines whether to replace instead of merge
var replaceFlag bool

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().BoolVarP(&replaceFlag, "replace", "r", false, "Replace entire config instead of merging")
}

func runImportCmd(cmd *cobra.Command, args []string) {
	inputPath := args[0]

	// Check if input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		printError(fmt.Sprintf("File not found: %s", inputPath))
		os.Exit(1)
	}

	// Read the input file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		printError(fmt.Sprintf("Failed to read file: %v", err))
		os.Exit(1)
	}

	// Validate YAML structure
	var newConfig config.Config
	if err := yaml.Unmarshal(data, &newConfig); err != nil {
		printError(fmt.Sprintf("Invalid YAML format: %v", err))
		os.Exit(1)
	}

	// Show what will be imported
	fmt.Printf("Found %d alias(es) in %s\n", len(newConfig.Aliases), inputPath)
	fmt.Println()

	if replaceFlag {
		// Replace mode - ask for confirmation
		if err := replaceConfig(inputPath, data); err != nil {
			printError(err.Error())
			os.Exit(1)
		}
	} else {
		// Merge mode (default)
		if err := mergeConfig(&newConfig); err != nil {
			printError(err.Error())
			os.Exit(1)
		}
	}
}

func replaceConfig(inputPath string, data []byte) error {
	// Ask if user wants to backup current config
	backupPrompt := promptui.Select{
		Label: "Do you want to backup your current config first?",
		Items: []string{"Yes, create backup", "No, just replace"},
	}

	backupIdx, _, err := backupPrompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("Cancelled.")
			return nil
		}
		return err
	}

	if backupIdx == 0 {
		// Create backup
		configPath := config.GetConfigFilePath()
		backupPath := configPath + ".backup"

		currentData, err := os.ReadFile(configPath)
		if err == nil {
			if err := os.WriteFile(backupPath, currentData, 0644); err != nil {
				return fmt.Errorf("failed to create backup: %w", err)
			}
			fmt.Printf("Backup saved to: %s\n", backupPath)
		}
	}

	// Confirm import
	confirmPrompt := promptui.Select{
		Label: "Replace current config with imported file?",
		Items: []string{"No, cancel", "Yes, replace"},
	}

	confirmIdx, _, err := confirmPrompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("Cancelled.")
			return nil
		}
		return err
	}

	if confirmIdx == 0 {
		fmt.Println("Cancelled.")
		return nil
	}

	// Write the new config
	configPath := config.GetConfigFilePath()
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Reload config
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	green := color.New(color.FgGreen, color.Bold)
	green.Println("Config replaced successfully!")

	return nil
}

func mergeConfig(newConfig *config.Config) error {
	// Get current aliases
	currentAliases, err := config.GetAllAliases()
	if err != nil {
		return fmt.Errorf("failed to load current config: %w", err)
	}

	// Build map of existing aliases
	existing := make(map[string]bool)
	for _, a := range currentAliases {
		existing[a.Name] = true
	}

	// Count new and duplicate aliases
	newCount := 0
	duplicates := []string{}

	for _, a := range newConfig.Aliases {
		if existing[a.Name] {
			duplicates = append(duplicates, a.Name)
		} else {
			newCount++
		}
	}

	fmt.Printf("New aliases to add: %d\n", newCount)
	if len(duplicates) > 0 {
		fmt.Printf("Already exist (will skip): %v\n", duplicates)
	}
	fmt.Println()

	if newCount == 0 {
		fmt.Println("No new aliases to import. All aliases already exist.")
		return nil
	}

	// Confirm
	confirmPrompt := promptui.Select{
		Label: fmt.Sprintf("Add %d new alias(es)?", newCount),
		Items: []string{"No, cancel", "Yes, add them"},
	}

	confirmIdx, _, err := confirmPrompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("Cancelled.")
			return nil
		}
		return err
	}

	if confirmIdx == 0 {
		fmt.Println("Cancelled.")
		return nil
	}

	// Add new aliases
	added := 0
	for _, a := range newConfig.Aliases {
		if !existing[a.Name] {
			if err := config.AddAlias(a); err != nil {
				fmt.Printf("Warning: Failed to add '%s': %v\n", a.Name, err)
			} else {
				added++
			}
		}
	}

	green := color.New(color.FgGreen, color.Bold)
	green.Printf("Added %d new alias(es)!\n", added)

	return nil
}
