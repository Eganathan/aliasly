#!/bin/bash
# Aliasly Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/yourusername/aliasly/main/scripts/install.sh | bash

set -e

# Configuration
REPO="Eganathan/aliasly"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="al"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo ""
echo "  ___  _ _           _       "
echo " / _ \| (_)         | |      "
echo "/ /_\ \ |_  __ _ ___| |_   _ "
echo "|  _  | | |/ _\` / __| | | | |"
echo "| | | | | | (_| \__ \ | |_| |"
echo "\_| |_/_|_|\__,_|___/_|\__, |"
echo "                        __/ |"
echo "                       |___/ "
echo ""
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
LATEST_VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${YELLOW}Warning: Could not fetch latest version, using v0.1.0${NC}"
    LATEST_VERSION="v0.1.0"
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
if ! curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${ARCHIVE}"; then
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

# Install
echo "Installing to ${INSTALL_DIR}..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo -e "${YELLOW}Need sudo permission to install to ${INSTALL_DIR}${NC}"
    sudo mv "$BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
fi

chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

# Verify installation
if command -v al &> /dev/null; then
    echo ""
    echo -e "${GREEN}Aliasly installed successfully!${NC}"
    echo ""
    al --version
    echo ""
    echo "Get started:"
    echo "  al list        # See default aliases"
    echo "  al add         # Add a new alias"
    echo "  al config      # Open web UI"
    echo ""
    echo "Run 'al --help' for more information."
else
    echo ""
    echo -e "${GREEN}Installation complete!${NC}"
    echo ""
    echo "Make sure ${INSTALL_DIR} is in your PATH, then run:"
    echo "  al --help"
fi
