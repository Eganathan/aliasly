#!/bin/bash
# Aliasly Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/Eganathan/aliasly/master/scripts/install.sh | bash

set -e

# Configuration
REPO="Eganathan/aliasly"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="al"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo ""
echo -e "${BLUE}"
echo "  ___  _ _           _       "
echo " / _ \| (_)         | |      "
echo "/ /_\ \ |_  __ _ ___| |_   _ "
echo "|  _  | | |/ _\` / __| | | | |"
echo "| | | | | | (_| \__ \ | |_| |"
echo "\_| |_/_|_|\__,_|___/_|\__, |"
echo "                        __/ |"
echo "                       |___/ "
echo -e "${NC}"
echo "Aliasly Installer"
echo ""

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    *)
        echo -e "${RED}Error: Unsupported operating system: $OS${NC}"
        echo "Aliasly supports macOS and Linux only."
        exit 1
        ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    amd64) ARCH="amd64" ;;
    arm64) ARCH="arm64" ;;
    aarch64) ARCH="arm64" ;;
    *)
        echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo "Detected: ${OS}/${ARCH}"

# Get latest release version from GitHub
echo "Fetching latest release..."

# Try to get latest version, with fallback
LATEST_VERSION=""

# Method 1: Try GitHub API
LATEST_VERSION=$(curl -fsSL -H "Accept: application/vnd.github.v3+json" \
    "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | \
    grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")

# Method 2: If API fails, try parsing releases page
if [ -z "$LATEST_VERSION" ]; then
    LATEST_VERSION=$(curl -fsSL "https://github.com/${REPO}/releases" 2>/dev/null | \
        grep -oE 'releases/tag/v[0-9]+\.[0-9]+\.[0-9]+' | head -1 | \
        sed 's/releases\/tag\///' || echo "")
fi

# Fallback to known working version
if [ -z "$LATEST_VERSION" ]; then
    echo -e "${YELLOW}Warning: Could not fetch latest version, using v0.1.5${NC}"
    LATEST_VERSION="v0.1.5"
fi

echo "Latest version: ${LATEST_VERSION}"

# Construct download URL
BINARY="al-${OS}-${ARCH}"
if [ "$OS" = "darwin" ]; then
    ARCHIVE="${BINARY}.zip"
else
    ARCHIVE="${BINARY}.tar.gz"
fi
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${ARCHIVE}"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download
echo "Downloading ${ARCHIVE}..."
if ! curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${ARCHIVE}" 2>/dev/null; then
    echo -e "${RED}Error: Failed to download from ${DOWNLOAD_URL}${NC}"
    echo ""
    echo "The release may not exist yet. You can build from source:"
    echo "  git clone https://github.com/${REPO}.git"
    echo "  cd aliasly && go build -o al ."
    exit 1
fi

# Extract
echo "Extracting..."
cd "$TMP_DIR"
if [ "$OS" = "darwin" ]; then
    unzip -q "$ARCHIVE"
else
    tar -xzf "$ARCHIVE"
fi

# Install binary
echo "Installing to ${INSTALL_DIR}..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo -e "${YELLOW}Need sudo permission to install to ${INSTALL_DIR}${NC}"
    sudo mv "$BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
fi

chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

# Verify binary works
if ! "${INSTALL_DIR}/${BINARY_NAME}" --version &>/dev/null; then
    echo -e "${RED}Error: Installation failed - binary not working${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}Binary installed successfully!${NC}"

# =============================================
# Shell Integration
# =============================================

echo ""
echo "Setting up shell integration..."

# Detect shell config file
SHELL_NAME=$(basename "$SHELL")
case "$SHELL_NAME" in
    zsh)
        SHELL_CONFIG="$HOME/.zshrc"
        ;;
    bash)
        if [ "$(uname)" = "Darwin" ] && [ -f "$HOME/.bash_profile" ]; then
            SHELL_CONFIG="$HOME/.bash_profile"
        else
            SHELL_CONFIG="$HOME/.bashrc"
        fi
        ;;
    fish)
        SHELL_CONFIG="$HOME/.config/fish/config.fish"
        mkdir -p "$(dirname "$SHELL_CONFIG")"
        ;;
    *)
        SHELL_CONFIG="$HOME/.bashrc"
        ;;
esac

# The line we need to add
INIT_LINE='eval "$(al init)"'
FISH_INIT_LINE='al init | source'

# Check if already added
if [ -f "$SHELL_CONFIG" ]; then
    if grep -q "al init" "$SHELL_CONFIG" 2>/dev/null; then
        echo "Shell integration already configured in $SHELL_CONFIG"
    else
        echo "Adding shell integration to $SHELL_CONFIG..."
        echo "" >> "$SHELL_CONFIG"
        echo "# Aliasly - command alias manager" >> "$SHELL_CONFIG"
        if [ "$SHELL_NAME" = "fish" ]; then
            echo "$FISH_INIT_LINE" >> "$SHELL_CONFIG"
        else
            echo "$INIT_LINE" >> "$SHELL_CONFIG"
        fi
        echo -e "${GREEN}Shell integration added!${NC}"
    fi
else
    echo "Creating $SHELL_CONFIG with shell integration..."
    if [ "$SHELL_NAME" = "fish" ]; then
        echo "$FISH_INIT_LINE" > "$SHELL_CONFIG"
    else
        echo "# Aliasly - command alias manager" > "$SHELL_CONFIG"
        echo "$INIT_LINE" >> "$SHELL_CONFIG"
    fi
    echo -e "${GREEN}Shell integration added!${NC}"
fi

# =============================================
# Done!
# =============================================

echo ""
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}  Aliasly installed successfully!${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo ""
"${INSTALL_DIR}/${BINARY_NAME}" --version
echo ""
echo "To activate now, run:"
echo -e "  ${BLUE}source $SHELL_CONFIG${NC}"
echo ""
echo "Or just open a new terminal window."
echo ""
echo "Quick start:"
echo "  al list        # See default aliases"
echo "  al add         # Add a new alias"
echo "  al config      # Open web UI"
echo ""
echo "After activation, use aliases directly:"
echo "  gs             # instead of: al gs"
echo "  gc \"message\"   # instead of: al gc \"message\""
echo ""
