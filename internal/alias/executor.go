package alias

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"aliasly/internal/config"
)

// ExecuteOptions contains options for command execution.
type ExecuteOptions struct {
	// Shell is the shell to use for executing the command.
	// If empty, the configured shell or system default will be used.
	Shell string

	// Verbose, when true, prints the command before executing it.
	Verbose bool

	// DryRun, when true, prints the command but doesn't execute it.
	// Useful for testing what a command would do.
	DryRun bool
}

// Execute runs a command string in the shell.
// It connects stdin, stdout, and stderr to the terminal so the command
// can interact with the user just like if they ran it directly.
//
// The command is executed using the shell specified in options,
// or the system default shell if not specified.
//
// Returns the exit code of the command, or an error if the command
// couldn't be started.
func Execute(command string, opts ExecuteOptions) (int, error) {
	// Determine which shell to use
	shell := opts.Shell
	if shell == "" {
		// Try to get shell from config
		cfg, err := config.Get()
		if err == nil && cfg.Settings.Shell != "" {
			shell = cfg.Settings.Shell
		} else {
			// Fall back to system default
			shell = config.GetDefaultShell()
		}
	}

	// Check verbose setting from config if not explicitly set
	verbose := opts.Verbose
	if !verbose {
		cfg, err := config.Get()
		if err == nil {
			verbose = cfg.Settings.Verbose
		}
	}

	// If verbose mode is on, print the command we're about to run
	if verbose {
		fmt.Printf("$ %s\n", command)
	}

	// If dry run, just return without executing
	if opts.DryRun {
		fmt.Printf("[dry-run] Would execute: %s\n", command)
		return 0, nil
	}

	// Create the command based on the operating system
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// On Windows, use cmd.exe with /C flag
		// /C means "run this command and then terminate"
		cmd = exec.Command("cmd", "/C", command)
	} else {
		// On Unix-like systems (macOS, Linux), use the shell with -c flag
		// -c means "run the following string as a command"
		cmd = exec.Command(shell, "-c", command)
	}

	// Connect the command's input/output to our terminal
	// This allows the command to:
	// - Read input from the user (stdin)
	// - Print output to the terminal (stdout)
	// - Print errors to the terminal (stderr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Also inherit the environment variables from the current process
	// This ensures commands can access things like PATH, HOME, etc.
	cmd.Env = os.Environ()

	// Run the command and wait for it to complete
	err := cmd.Run()

	// Extract the exit code from the result
	// A nil error means the command succeeded (exit code 0)
	if err == nil {
		return 0, nil
	}

	// If the command failed, try to get the exit code
	// In Go, we need to type-assert to *exec.ExitError to get the exit code
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), nil
	}

	// If we couldn't start the command at all, return the error
	return -1, fmt.Errorf("failed to execute command: %w", err)
}

// Run is a convenience function that parses an alias with arguments
// and executes the resulting command.
// This is the main entry point for running aliases.
func Run(a Alias, args []string) (int, error) {
	// Parse the command by substituting parameters
	command, err := ParseCommand(a, args)
	if err != nil {
		return -1, err
	}

	// Execute the parsed command
	return Execute(command, ExecuteOptions{})
}

// RunWithOptions is like Run but allows specifying execution options.
func RunWithOptions(a Alias, args []string, opts ExecuteOptions) (int, error) {
	// Parse the command by substituting parameters
	command, err := ParseCommand(a, args)
	if err != nil {
		return -1, err
	}

	// Execute the parsed command with the given options
	return Execute(command, opts)
}
