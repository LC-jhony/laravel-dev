#!/bin/bash

echo "============================================"
echo "       laravel-dev Installer"
echo "============================================"

# Verificar si Go está instalado
if ! command -v go &> /dev/null; then
    echo "[ERROR] Go no está instalado. Instálalo primero."
    exit 1
fi

# Verificar versión de Go
GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+\.[0-9]+')
if [[ "$GO_VERSION" < "1.25.8" ]]; then
    echo "[INFO] Go versión $GO_VERSION detectada, requiere >= 1.25.8"
    exit 1
fi

echo "[INFO] Go está instalado con la versión requerida"
echo "[INFO] Instalando laravel-dev..."

# Instalar el paquete
go install github.com/LC-jhony/laravel-dev@latest

if [ $? -eq 0 ]; then
    echo "✅ laravel-dev instalado correctamente"
    echo ""
    echo "Ejecuta:"
    echo "  laravel-dev"
    echo ""
    echo "O:"
    echo "  go run ."
else
    echo "❌ Error al instalar laravel-dev"
    exit 1
fi
