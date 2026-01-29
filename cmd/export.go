package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"aliasly/internal/config"
)

// exportCmd represents the export command.
// It exports the current configuration to a file or stdout.
var exportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Export aliases to a file",
	Long: `Export your aliases configuration to a YAML file for backup.

If no file is specified, the config is printed to stdout.

Examples:
  al export                    # Print config to terminal
  al export backup.yaml        # Save to backup.yaml
  al export ~/my-aliases.yaml  # Save to home directory`,

	Args: cobra.MaximumNArgs(1),
	Run:  runExportCmd,
}

func init() {
	rootCmd.AddCommand(exportCmd)
}

func runExportCmd(cmd *cobra.Command, args []string) {
	// Get config file path
	configPath := config.GetConfigFilePath()

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		printError(fmt.Sprintf("Failed to read config: %v", err))
		os.Exit(1)
	}

	// If no output file specified, print to stdout
	if len(args) == 0 {
		fmt.Print(string(data))
		return
	}

	// Write to the specified file
	outputPath := args[0]
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		printError(fmt.Sprintf("Failed to write to %s: %v", outputPath, err))
		os.Exit(1)
	}

	fmt.Printf("Config exported to: %s\n", outputPath)
}
