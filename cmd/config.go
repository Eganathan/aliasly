package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"aliasly/internal/webui"
)

// configCmd represents the config command.
// It starts a local web server for managing aliases through a browser UI.
var configCmd = &cobra.Command{
	// Use is the command name
	Use: "config",

	// Aliases for shorter typing
	Aliases: []string{"cfg", "ui"},

	// Short description
	Short: "Open web UI to manage aliases",

	// Long description
	Long: `Open a web-based configuration interface in your browser.

This starts a local web server that provides a visual interface for:
  - Viewing all your aliases
  - Creating new aliases
  - Editing existing aliases
  - Deleting aliases

The server runs on localhost only and shuts down when you press Ctrl+C.

Examples:
  al config    # Open web configuration UI
  al ui        # Short form`,

	// Run function
	Run: runConfigCmd,
}

// runConfigCmd executes the config command.
func runConfigCmd(cmd *cobra.Command, args []string) {
	// Find an available port by listening on port 0
	// The OS will assign an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		printError(fmt.Sprintf("Failed to find available port: %v", err))
		os.Exit(1)
	}

	// Get the port that was assigned
	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://127.0.0.1:%d", port)

	// Create the HTTP server with our handlers
	server := webui.NewServer()
	httpServer := &http.Server{
		Handler: server.Handler(),
	}

	// Start the server in a goroutine (background thread)
	// This allows us to continue and open the browser
	go func() {
		// Serve accepts connections on the listener
		// It blocks until the server is shut down
		if err := httpServer.Serve(listener); err != http.ErrServerClosed {
			printError(fmt.Sprintf("Server error: %v", err))
		}
	}()

	// Print information
	cyan := color.New(color.FgCyan, color.Bold)
	cyan.Println("Aliasly Configuration UI")
	fmt.Println()
	fmt.Printf("Server running at: %s\n", url)
	fmt.Println()

	// Try to open the browser
	if err := openBrowser(url); err != nil {
		// If browser can't be opened, just show the URL
		fmt.Printf("Could not open browser automatically.\n")
		fmt.Printf("Please open this URL in your browser: %s\n", url)
	} else {
		fmt.Println("Opening in your default browser...")
	}

	fmt.Println()
	fmt.Println("Press Ctrl+C to stop the server")

	// Wait for interrupt signal (Ctrl+C)
	// This keeps the server running until the user decides to stop it
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until we receive a signal

	fmt.Println()
	fmt.Println("Shutting down server...")

	// Gracefully shut down the server
	// context.Background() creates a basic context for the shutdown
	if err := httpServer.Shutdown(context.Background()); err != nil {
		printError(fmt.Sprintf("Server shutdown error: %v", err))
	}

	green := color.New(color.FgGreen)
	green.Println("Server stopped.")
}

// openBrowser opens the specified URL in the default browser.
// It handles different operating systems appropriately.
func openBrowser(url string) error {
	var cmd string
	var args []string

	// Different operating systems have different commands to open URLs
	switch runtime.GOOS {
	case "darwin":
		// macOS uses the "open" command
		cmd = "open"
		args = []string{url}
	case "linux":
		// Linux uses xdg-open (part of xdg-utils package)
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		// Windows uses "start" command through cmd
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Start the command but don't wait for it to finish
	// (the browser will keep running after we return)
	return exec.Command(cmd, args...).Start()
}
