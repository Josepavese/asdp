#!/bin/bash
set -e

APP_NAME="asdp"
INSTALL_BASE="$HOME/.asdp"
INSTALL_DIR="$INSTALL_BASE/bin"
REPO="Josepavese/asdp"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Authentication for private repos
AUTH_HEADER=()
if [ -n "$GITHUB_TOKEN" ]; then
    AUTH_HEADER=(-H "Authorization: token $GITHUB_TOKEN")
fi

echo -e "${GREEN}Starting ASDP Installer...${NC}"

# 1. Detect OS and Arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$OS" == "darwin" ]; then
    PLATFORM="darwin"
    if [ "$ARCH" == "arm64" ]; then
        BINARY_ARCH="arm64"
    else
        BINARY_ARCH="amd64"
    fi
elif [ "$OS" == "linux" ]; then
    PLATFORM="linux"
    if [ "$ARCH" == "aarch64" ]; then
        BINARY_ARCH="arm64"
    else
        BINARY_ARCH="amd64"
    fi
else
    echo -e "${RED}Unsupported OS: $OS${NC}"
    exit 1
fi

BINARY_NAME="asdp-${PLATFORM}-${BINARY_ARCH}.zst"
echo "Detected: ${PLATFORM} / ${BINARY_ARCH}"

# 2. Check Prerequisites (Universal Ctags & zstd)
echo "Checking prerequisites..."
PREREQ_MISSING=0
if ! command -v ctags &> /dev/null; then
    echo -e "${YELLOW}Warning: 'ctags' not found.${NC}"
    echo "ASDP requires 'universal-ctags' for multi-language support."
    PREREQ_MISSING=1
fi

if ! command -v zstd &> /dev/null; then
    echo -e "${YELLOW}Warning: 'zstd' not found.${NC}"
    echo "Required for decompressing the asdp binary."
    PREREQ_MISSING=1
fi

if [ $PREREQ_MISSING -eq 1 ]; then
    if [ "$OS" == "linux" ]; then
       echo -e "Please run: ${GREEN}sudo apt-get install universal-ctags zstd${NC}"
    elif [ "$OS" == "darwin" ]; then
       echo -e "Please run: ${GREEN}brew install universal-ctags zstd${NC}"
    fi
    # If zstd is missing, we must exit
    if ! command -v zstd &> /dev/null; then exit 1; fi
fi

# 3. Create Directories
echo "Setting up directories..."
mkdir -p "$INSTALL_DIR"

# 4. Download Binary
echo "Downloading latest release..."
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}"
CORE_URL="https://github.com/${REPO}/releases/latest/download/asdp-core.zip"

echo "Fetching Binary from: $DOWNLOAD_URL"
# Use -L to follow redirects and -f to fail on 404
if curl "${AUTH_HEADER[@]}" -L -f -o "${INSTALL_DIR}/${APP_NAME}.zst" "$DOWNLOAD_URL"; then
    echo -e "${GREEN}Binary download successful.${NC}"
    echo "Decompressing binary..."
    zstd -d --rm "${INSTALL_DIR}/${APP_NAME}.zst" -o "${INSTALL_DIR}/${APP_NAME}"
else
    echo -e "${RED}Binary download failed (404 Not Found).${NC}"
    echo -e "${YELLOW}Note: If the repository is private, ensure you are authenticated.${NC}"
    echo -e "Tip: Try setting the token: ${GREEN}GITHUB_TOKEN=your_token ./install.sh${NC}"
    exit 1
fi
chmod +x "${INSTALL_DIR}/${APP_NAME}"

echo "Fetching Core Assets from: $CORE_URL"
# Use temp file for zip
if curl "${AUTH_HEADER[@]}" -L -f -o "/tmp/asdp-core.zip" "$CORE_URL"; then
    echo "Unzipping core assets..."
    # Unzip into ~/.asdp, overwrite existing
    unzip -o /tmp/asdp-core.zip -d "$HOME/.asdp/"
    rm /tmp/asdp-core.zip
    echo -e "${GREEN}Core assets installed.${NC}"
else
    echo -e "${YELLOW}Warning: Core assets download failed.${NC}"
    echo -e "If this is a private repo, ensure GITHUB_TOKEN is set."
fi

# 5. Path Setup
echo "Configuring PATH..."
SHELL_RC=""
case "$SHELL" in
    */bash) SHELL_RC="$HOME/.bashrc" ;;
    */zsh) SHELL_RC="$HOME/.zshrc" ;;
esac

if [ -n "$SHELL_RC" ]; then
    if ! grep -q "$INSTALL_DIR" "$SHELL_RC"; then
        echo "" >> "$SHELL_RC"
        echo "# ASDP Path" >> "$SHELL_RC"
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_RC"
        echo -e "${GREEN}Added ${INSTALL_DIR} to ${SHELL_RC}${NC}"
        echo "Please restart your terminal or run: source $SHELL_RC"
    else
        echo "Path already configured in $SHELL_RC"
    fi
fi


# --- MCP Configuration Logic ---

configure_mcp_server() {
    local CONFIG_PATH="$1"
    local SERVER_NAME="asdp"
    local BIN_PATH="${INSTALL_DIR}/${APP_NAME}"

    if [ -f "$CONFIG_PATH" ]; then
        echo "Found MCP config: $CONFIG_PATH"
        
        # Inline Python script to safely update JSON
        python3 -c "
import json
import sys
import os

config_path = '$CONFIG_PATH'
server_name = '$SERVER_NAME'
bin_path = '$BIN_PATH'

try:
    with open(config_path, 'r') as f:
        data = json.load(f)
except Exception as e:
    print(f'Error reading {config_path}: {e}')
    sys.exit(1)

# Ensure 'mcpServers' key exists
if 'mcpServers' not in data:
    data['mcpServers'] = {}

# Update or Add ASDP entry
data['mcpServers'][server_name] = {
    'command': bin_path,
    'args': [],
    'env': {}
}

try:
    with open(config_path, 'w') as f:
        json.dump(data, f, indent=4)
    print(f'Successfully updated {server_name} in {config_path}')
except Exception as e:
    print(f'Error writing {config_path}: {e}')
    sys.exit(1)
"
    fi
}

echo "Configuring IDE Integrations..."

# VS Code (Linux)
configure_mcp_server "$HOME/.config/Code/User/globalStorage/mcp-servers.json"
# VS Code (macOS)
configure_mcp_server "$HOME/Library/Application Support/Code/User/globalStorage/mcp-servers.json"
# Claude Desktop (macOS)
configure_mcp_server "$HOME/Library/Application Support/Claude/claude_desktop_config.json"
# Cursor (Linux)
configure_mcp_server "$HOME/.config/Cursor/User/globalStorage/mcp-servers.json"
# Cursor (macOS)
configure_mcp_server "$HOME/Library/Application Support/Cursor/User/globalStorage/mcp-servers.json"

echo -e "${GREEN}ASDP installed successfully!${NC}"
echo "Run 'asdp' to start."

# 6. Interactive Project Initialization
echo ""
read -p "Do you want to initialize ASDP in the current directory? (y/N): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    AGENT_DIR="./.agent"
    echo "Initializing ASDP in $(pwd)..."
    mkdir -p "$AGENT_DIR"
    
    # Copy from global core/agent if it exists
    SRC_AGENT="$INSTALL_BASE/core/agent"
    if [ -d "$SRC_AGENT" ]; then
        cp -r "$SRC_AGENT/"* "$AGENT_DIR/"
        echo -e "${GREEN}Project initialized successfully in $AGENT_DIR${NC}"
    else
        echo -e "${YELLOW}Warning: Global agent templates not found at $SRC_AGENT${NC}"
    fi
fi
