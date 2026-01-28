# Aliasly

A simple, cross-platform CLI tool for managing command aliases. Instead of typing long commands, create short aliases and run them with `al`.

```bash
# Instead of this:
git commit -am "fix login bug"

# Type this:
al gc "fix login bug"
```

## Features

- **Simple aliasing** - Map short names to long commands
- **Parameterized commands** - Use `{{param}}` placeholders for dynamic values
- **Web-based configuration** - Visual UI for managing aliases (`al config`)
- **Interactive CLI** - Add aliases through guided prompts (`al add`)
- **Cross-platform** - Works on macOS and Linux
- **Single binary** - No dependencies, just one executable
- **YAML config** - Human-readable configuration file

## Installation

### From Source

Requires [Go 1.21+](https://go.dev/dl/)

```bash
# Clone the repository
git clone https://github.com/yourusername/aliasly.git
cd aliasly

# Build
go build -o al .

# Move to PATH (optional)
sudo mv al /usr/local/bin/
```

### Pre-built Binaries

Download from the [Releases](https://github.com/yourusername/aliasly/releases) page.

## Quick Start

```bash
# List all aliases
al list

# Run an alias
al gs                    # runs: git status

# Run with parameters
al gc "fix bug"          # runs: git commit -am "fix bug"

# Add a new alias interactively
al add

# Open web UI to manage aliases
al config

# Remove an alias
al remove myalias
```

## Usage

### Running Aliases

```bash
al <alias-name> [parameters...]
```

**Examples:**
```bash
al gs                     # Simple alias
al gc "commit message"    # With required parameter
al gp feature-branch      # With optional parameter
al gp                     # Uses default value for optional param
```

### Managing Aliases

| Command | Description |
|---------|-------------|
| `al list` | List all configured aliases |
| `al add` | Add a new alias interactively |
| `al remove <name>` | Remove an existing alias |
| `al config` | Open web UI for visual management |

### Command Flags

```bash
al --help       # Show help
al --version    # Show version
al -v <alias>   # Verbose mode (shows command before running)
```

## Configuration

Configuration is stored in `~/.config/aliasly/config.yaml`

### Config File Format

```yaml
version: 1
settings:
  shell: /bin/bash    # Shell to use for commands
  verbose: false      # Print commands before running

aliases:
  # Simple alias (no parameters)
  - name: gs
    command: git status
    description: Show git status

  # Alias with required parameter
  - name: gc
    command: git commit -am "{{message}}"
    description: Git commit with message
    params:
      - name: message
        description: Commit message
        required: true

  # Alias with optional parameter and default value
  - name: gp
    command: git push origin {{branch}}
    description: Push to remote branch
    params:
      - name: branch
        description: Branch name
        required: false
        default: main
```

### Parameter Syntax

Use `{{paramName}}` in your command to define parameters:

```yaml
- name: deploy
  command: ./deploy.sh {{env}} --version={{version}}
  params:
    - name: env
      required: true
    - name: version
      required: false
      default: latest
```

Usage: `al deploy production` or `al deploy staging v1.2.3`

### Config Location

The config file location follows XDG standards:

1. `$ALIASLY_CONFIG_DIR/config.yaml` (if set)
2. `$XDG_CONFIG_HOME/aliasly/config.yaml`
3. `~/.config/aliasly/config.yaml` (default)

## Web Configuration UI

Run `al config` to open a browser-based interface for managing aliases:

- View all aliases at a glance
- Add new aliases with a form
- Edit existing aliases
- Delete aliases with confirmation
- Auto-detects parameters from `{{placeholders}}`

The web server runs locally and shuts down when you press `Ctrl+C`.

## Example Aliases

Here are some useful aliases to get you started:

```yaml
aliases:
  # Git shortcuts
  - name: gs
    command: git status
    description: Git status

  - name: gd
    command: git diff
    description: Git diff

  - name: gc
    command: git commit -am "{{message}}"
    description: Git commit with message
    params:
      - name: message
        required: true

  - name: gp
    command: git push origin {{branch}}
    description: Push to branch
    params:
      - name: branch
        required: false
        default: main

  - name: gco
    command: git checkout {{branch}}
    description: Switch branch
    params:
      - name: branch
        required: true

  # Docker shortcuts
  - name: dps
    command: docker ps
    description: List containers

  - name: dex
    command: docker exec -it {{container}} {{cmd}}
    description: Execute in container
    params:
      - name: container
        required: true
      - name: cmd
        required: false
        default: bash

  # Development
  - name: serve
    command: python -m http.server {{port}}
    description: Start HTTP server
    params:
      - name: port
        required: false
        default: "8000"
```

## Building from Source

```bash
# Clone
git clone https://github.com/yourusername/aliasly.git
cd aliasly

# Build for current platform
go build -o al .

# Build for all platforms
GOOS=darwin GOARCH=amd64 go build -o al-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o al-darwin-arm64 .
GOOS=linux GOARCH=amd64 go build -o al-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o al-linux-arm64 .
```

## Shell Completion

Generate shell completion scripts:

```bash
# Bash
al completion bash > /etc/bash_completion.d/al

# Zsh
al completion zsh > "${fpath[1]}/_al"

# Fish
al completion fish > ~/.config/fish/completions/al.fish
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.
