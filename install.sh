#!/bin/bash

set -e

echo "Installing laravel-dev..."

if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed."
    echo "Please install Go 1.25.8 or higher from https://go.dev/dl/"
    exit 1
fi

GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+(\.[0-9]+)?')
REQUIRED_VERSION="1.25.8"

compare_versions() {
    printf '%s\n%s\n' "$1" "$2" | sort -V -r | head -n1
}

if [ "$(compare_versions "$GO_VERSION" "$REQUIRED_VERSION")" != "$REQUIRED_VERSION" ]; then
    echo "Warning: Go version $GO_VERSION is installed."
    echo "laravel-dev requires Go $REQUIRED_VERSION or higher."
fi

echo "Installing laravel-dev with go install..."
go install github.com/LC-jhony/laravel-dev@latest

INSTALL_DIR=$(go env GOPATH)/bin

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo "Warning: $INSTALL_DIR is not in your PATH."
    echo "Add the following line to your shell profile (.bashrc, .zshrc, etc.):"
    echo ""
    echo "    export PATH=\$PATH:$INSTALL_DIR"
    echo ""
fi

echo ""
echo "Installation complete!"
echo "Run 'laravel-dev' to start the installer."
