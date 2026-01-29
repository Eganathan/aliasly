package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"aliasly/internal/config"
)

// uninstallCmd represents the uninstall command.
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Aliasly from your system",
	Long: `Uninstall Aliasly from your system.

This will:
1. Remove shell integration from your shell config
2. Optionally remove the al binary
3. Optionally remove your aliases config file

You will be asked for confirmation before each step.`,

	Run: runUninstallCmd,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstallCmd(cmd *cobra.Command, args []string) {
	red := color.New(color.FgRed, color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	fmt.Println()
	red.Println("Aliasly Uninstaller")
	fmt.Println("===================")
	fmt.Println()

	// Confirm uninstall
	confirmPrompt := promptui.Select{
		Label: "Are you sure you want to uninstall Aliasly?",
		Items: []string{"No, cancel", "Yes, uninstall"},
	}

	idx, _, err := confirmPrompt.Run()
	if err != nil || idx == 0 {
		fmt.Println("Cancelled.")
		return
	}

	fmt.Println()

	// Step 1: Remove shell integration
	shellConfig := getShellConfigFile()
	if shellConfig != "" {
		fmt.Printf("Shell config: %s\n", shellConfig)

		removeShellPrompt := promptui.Select{
			Label: "Remove shell integration from config file?",
			Items: []string{"Yes, remove it", "No, keep it"},
		}

		idx, _, err := removeShellPrompt.Run()
		if err == nil && idx == 0 {
			if err := removeShellIntegration(shellConfig); err != nil {
				yellow.Printf("Warning: Could not remove shell integration: %v\n", err)
			} else {
				green.Println("Shell integration removed.")
			}
		}
		fmt.Println()
	}

	// Step 2: Remove config file
	configPath := config.GetConfigFilePath()
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file: %s\n", configPath)

		removeConfigPrompt := promptui.Select{
			Label: "Remove your aliases config file?",
			Items: []string{"No, keep my aliases", "Yes, delete everything"},
		}

		idx, _, err := removeConfigPrompt.Run()
		if err == nil && idx == 1 {
			configDir := config.GetConfigDir()
			if err := os.RemoveAll(configDir); err != nil {
				yellow.Printf("Warning: Could not remove config: %v\n", err)
			} else {
				green.Println("Config file removed.")
			}
		}
		fmt.Println()
	}

	// Step 3: Remove binary
	binaryPath, _ := os.Executable()
	if binaryPath != "" {
		fmt.Printf("Binary: %s\n", binaryPath)

		removeBinaryPrompt := promptui.Select{
			Label: "Remove the al binary?",
			Items: []string{"Yes, remove it", "No, keep it"},
		}

		idx, _, err := removeBinaryPrompt.Run()
		if err == nil && idx == 0 {
			if err := removeBinary(binaryPath); err != nil {
				yellow.Printf("Warning: Could not remove binary: %v\n", err)
				fmt.Println("You can remove it manually with:")
				fmt.Printf("  sudo rm %s\n", binaryPath)
			} else {
				green.Println("Binary removed.")
			}
		}
		fmt.Println()
	}

	green.Println("Uninstall complete!")
	fmt.Println()
	fmt.Println("Please restart your terminal or run:")
	if shellConfig != "" {
		fmt.Printf("  source %s\n", shellConfig)
	}
}

// getShellConfigFile returns the path to the user's shell config file.
func getShellConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	shell := os.Getenv("SHELL")

	switch {
	case strings.Contains(shell, "zsh"):
		return filepath.Join(home, ".zshrc")
	case strings.Contains(shell, "bash"):
		if runtime.GOOS == "darwin" {
			bashProfile := filepath.Join(home, ".bash_profile")
			if _, err := os.Stat(bashProfile); err == nil {
				return bashProfile
			}
		}
		return filepath.Join(home, ".bashrc")
	case strings.Contains(shell, "fish"):
		return filepath.Join(home, ".config", "fish", "config.fish")
	default:
		return filepath.Join(home, ".bashrc")
	}
}

// removeShellIntegration removes the al init line from shell config.
func removeShellIntegration(configPath string) error {
	// Read the file
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	skipNext := false

	for scanner.Scan() {
		line := scanner.Text()

		// Skip the comment line before al init
		if strings.Contains(line, "Aliasly") && strings.Contains(line, "alias manager") {
			skipNext = true
			continue
		}

		// Skip the al init line
		if strings.Contains(line, "al init") {
			skipNext = false
			continue
		}

		if skipNext {
			skipNext = false
			continue
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Write back
	return os.WriteFile(configPath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

// removeBinary removes the al binary, using sudo if necessary.
func removeBinary(binaryPath string) error {
	// Try to remove directly first
	err := os.Remove(binaryPath)
	if err == nil {
		return nil
	}

	// If permission denied, try with sudo
	if os.IsPermission(err) {
		cmd := exec.Command("sudo", "rm", binaryPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return err
}
