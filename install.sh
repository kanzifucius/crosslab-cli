#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print step information
print_step() {
    echo -e "${BLUE}==>${NC} $1"
}

# Print success message
print_success() {
    echo -e "${GREEN}==>${NC} $1"
}

# Print error message and exit
print_error() {
    echo -e "${RED}Error:${NC} $1" >&2
    exit 1
}

# Detect the operating system and architecture
detect_platform() {
    local os
    local arch

    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)

    # Convert architecture names
    case "$arch" in
        x86_64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            ;;
    esac

    # Convert OS names
    case "$os" in
        darwin|linux)
            : # These are fine as-is
            ;;
        *)
            print_error "Unsupported operating system: $os"
            ;;
    esac

    echo "$os-$arch"
}

# Main installation function
install_crosslocal() {
    local platform
    platform=$(detect_platform)
    
    print_step "Detected platform: $platform"

    # Create installation directory if it doesn't exist
    local install_dir="$HOME/.crosslocal/bin"
    mkdir -p "$install_dir"

    # Download latest release
    print_step "Downloading crosslocal binary..."
    local download_url="https://github.com/YOUR_ORG/crosslocal/releases/latest/download/crosslocal-${platform}"
    
    if ! curl -fsSL -o "$install_dir/crosslocal" "$download_url"; then
        print_error "Failed to download crosslocal binary"
    fi

    # Make binary executable
    chmod +x "$install_dir/crosslocal"

    # Add to PATH if not already there
    local shell_config="$HOME/.$(basename "$SHELL")rc"
    if ! grep -q "$install_dir" "$shell_config" 2>/dev/null; then
        echo "export PATH=\"\$PATH:$install_dir\"" >> "$shell_config"
    fi

    print_success "crosslocal has been installed successfully to $install_dir/crosslocal!"
    print_success "Please restart your terminal or run: source $shell_config"
    print_step "You can now use 'crosslocal' from your terminal"
}

# Run the installation
install_crosslocal 