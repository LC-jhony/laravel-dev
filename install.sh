#!/bin/bash

set -e

REQUIRED_GO_VERSION="1.25.8"
GO_INSTALL_VERSION="1.25.8"
TEMP_DIR="/tmp/laravel-dev-install"
INSTALL_DIR="/usr/local/go"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

detect_os() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if [ -f /etc/debian_version ]; then
            echo "debian"
        elif [ -f /etc/fedora-release ]; then
            echo "fedora"
        elif [ -f /etc/arch-release ]; then
            echo "arch"
        elif [ -f /etc/alpine-release ]; then
            echo "alpine"
        else
            echo "linux"
        fi
    else
        echo "unsupported"
    fi
}

detect_architecture() {
    case $(uname -m) in
        x86_64|amd64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        armv7l)
            echo "armv6l"
            ;;
        *)
            echo "amd64"
            ;;
    esac
}

compare_versions() {
    printf '%s\n%s\n' "$1" "$2" | sort -V -r | head -n1
}

check_go_installed() {
    if command -v go &> /dev/null; then
        local current_version=$(go version 2>/dev/null | grep -oP 'go\K[0-9]+\.[0-9]+(\.[0-9]+)?' || echo "0")
        if [ "$(compare_versions "$current_version" "$REQUIRED_GO_VERSION")" = "$REQUIRED_GO_VERSION" ]; then
            return 0
        else
            log_warn "Go $current_version is installed but $REQUIRED_GO_VERSION+ is required"
            return 1
        fi
    fi
    return 2
}

install_go_linux() {
    local os=$1
    local arch=$(detect_architecture)
    local download_url="https://go.dev/dl/go${GO_INSTALL_VERSION}.linux-${arch}.tar.gz"

    log_info "Installing Go ${GO_INSTALL_VERSION} for Linux (${arch})..."

    case "$os" in
        debian)
            log_info "Using apt to install Go..."
            sudo apt update
            sudo apt install -y curl wget
            ;;
        fedora)
            log_info "Using dnf to install Go..."
            sudo dnf install -y curl wget
            ;;
        arch)
            log_info "Using pacman to install Go..."
            sudo pacman -S --noconfirm go
            return 0
            ;;
        alpine)
            log_info "Using apk to install Go..."
            sudo apk add --no-cache curl wget
            ;;
    esac

    log_info "Downloading Go from $download_url..."
    mkdir -p "$TEMP_DIR"
    curl -fsSL "$download_url" -o "$TEMP_DIR/go.tar.gz"

    log_info "Extracting Go to $INSTALL_DIR..."
    sudo rm -rf "$INSTALL_DIR"
    sudo mkdir -p "$INSTALL_DIR"
    sudo tar -C "$INSTALL_DIR" --strip-components=1 -xzf "$TEMP_DIR/go.tar.gz"

    rm -rf "$TEMP_DIR"

    export PATH="$INSTALL_DIR/bin:$PATH"
    echo "export PATH=$INSTALL_DIR/bin:\$PATH" >> ~/.bashrc
    if [ -f ~/.zshrc ]; then
        echo "export PATH=$INSTALL_DIR/bin:\$PATH" >> ~/.zshrc
    fi

    log_info "Go installed successfully!"
}

install_go_macos() {
    local arch=$(detect_architecture)
    local pkg_name="go${GO_INSTALL_VERSION}.darwin-${arch}.pkg"
    local archive_name="go${GO_INSTALL_VERSION}.darwin-${arch}.tar.gz"

    log_info "Installing Go ${GO_INSTALL_VERSION} for macOS (${arch})..."

    if command -v brew &> /dev/null; then
        log_info "Using Homebrew to install Go..."
        brew install go
        return 0
    fi

    local download_url="https://go.dev/dl/${archive_name}"

    log_info "Downloading Go from $download_url..."
    mkdir -p "$TEMP_DIR"
    curl -fsSL "$download_url" -o "$TEMP_DIR/go.tar.gz"

    log_info "Extracting Go to /usr/local..."
    sudo rm -rf "$INSTALL_DIR"
    sudo tar -C "/usr/local" --strip-components=1 -xzf "$TEMP_DIR/go.tar.gz"

    rm -rf "$TEMP_DIR"

    export PATH="$INSTALL_DIR/bin:$PATH"
    echo "export PATH=$INSTALL_DIR/bin:\$PATH" >> ~/.bashrc
    if [ -f ~/.zshrc ]; then
        echo "export PATH=$INSTALL_DIR/bin:\$PATH" >> ~/.zshrc
    fi

    log_info "Go installed successfully!"
}

install_go() {
    local os=$(detect_os)

    log_info "Detected OS: $os"

    case "$os" in
        macos)
            install_go_macos
            ;;
        debian|fedora|arch|alpine|linux)
            install_go_linux "$os"
            ;;
        *)
            log_error "Unsupported operating system: $OSTYPE"
            log_error "Please install Go manually from https://go.dev/dl/"
            exit 1
            ;;
    esac
}

main() {
    echo ""
    echo "============================================"
    echo "       laravel-dev Installer"
    echo "============================================"
    echo ""

    local go_status=$(check_go_installed)

    if [ $? -eq 0 ]; then
        log_info "Go is already installed with the required version"
    else
        log_warn "Go is not installed or version is insufficient"
        log_info "Installing Go ${REQUIRED_GO_VERSION}..."
        install_go

        export PATH="$INSTALL_DIR/bin:$PATH"
    fi

    echo ""
    log_info "Installing laravel-dev..."
    go install github.com/LC-jhony/laravel-dev@latest

    local bin_dir=$(go env GOPATH)/bin

    echo ""
    if [[ ":$PATH:" != *":$bin_dir:"* ]]; then
        log_warn "$bin_dir is not in your PATH"
        echo ""
        echo "Add the following line to your shell profile (.bashrc, .zshrc, etc.):"
        echo ""
        echo -e "  ${GREEN}export PATH=\$PATH:$bin_dir${NC}"
        echo ""
    fi

    echo ""
    echo "============================================"
    log_info "Installation complete!"
    echo "============================================"
    echo ""
    echo "Run the following command to start:"
    echo ""
    echo -e "  ${GREEN}laravel-dev${NC}"
    echo ""
}

main "$@"
