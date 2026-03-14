#!/bin/bash

# Colores ANSI
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Función para imprimir con formato
print_header() {
    echo -e "${CYAN}┌─────────────────────────────────────────────────────────┐${NC}"
    echo -e "${CYAN}│${NC}                                                         ${CYAN}│${NC}"
    echo -e "${CYAN}│${NC}         ${MAGENTA}🚀 laravel-dev Installer${NC}                        ${CYAN}│${NC}"
    echo -e "${CYAN}│${NC}                                                         ${CYAN}│${NC}"
    echo -e "${CYAN}└─────────────────────────────────────────────────────────┘${NC}"
}

print_step() {
    echo -e "${BLUE}→${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Limpiar pantalla
clear

# Mostrar cabecera
print_header
echo ""

# Verificar si Go está instalado
print_step "Verificando instalación de Go..."
if ! command -v go &> /dev/null; then
    print_error "Go no está instalado. Instálalo primero."
    exit 1
fi
print_success "Go detectado"

# Verificar versión de Go
GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+\.[0-9]+')
print_step "Verificando versión de Go ($GO_VERSION)..."
if [[ "$GO_VERSION" < "1.22.2" ]]; then
    print_error "Versión $GO_VERSION detectada, se requiere >= 1.22.2"
    exit 1
fi
print_success "Versión de Go compatible"

print_step "Preparando instalación de laravel-dev..."
echo ""

# Instalar el paquete
print_step "Descargando e instalando laravel-dev..."
go install github.com/LC-jhony/laravel-dev@latest

if [ $? -eq 0 ]; then
    print_success "laravel-dev instalado correctamente"
    echo ""
    print_step "Iniciando aplicación..."
    echo ""
    sleep 1
    
    # Ejecutar el binario instalado
    GOPATH=$(go env GOPATH)
    LARAVEL_DEV_BIN="$GOPATH/bin/laravel-dev"
    
    if [ -f "$LARAVEL_DEV_BIN" ]; then
        "$LARAVEL_DEV_BIN"
    else
        # Fallback: intentar con go run si el binario no existe
        if [ -d "/home/crack/Sites/laravel-dev" ]; then
            cd /home/crack/Sites/laravel-dev && go run .
        else
            print_error "No se pudo encontrar la aplicación"
            exit 1
        fi
    fi
else
    print_error "Error al instalar laravel-dev"
    exit 1
fi
