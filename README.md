# laravel-dev

Herramienta CLI interactiva para instalar y configurar tu entorno de desarrollo Laravel en Linux/macOS.

## Descripción

`laravel-dev` es una interfaz de terminal (TUI) que te permite instalar de forma sencilla:

- **PHP** - Múltiples versiones con extensiones seleccionables
- **MariaDB** - Servidor de base de datos MySQL
- **Node.js** - Entorno de JavaScript del lado del servidor
- **Composer** - Gestor de dependencias para PHP
- **Laravel Valet** - Entorno de desarrollo local para Laravel

## Requisitos

- Go 1.25.8 o superior
- Sistema operativo: Linux o macOS
- Acceso sudo (para instalación de paquetes)
- Terminal interactiva

## Instalación

### Método rápido (curl)

```bash
curl -sSL https://raw.githubusercontent.com/LC-jhony/laravel-dev/main/install.sh | bash
```

### Método manual (go install)

```bash
go install github.com/LC-jhony/laravel-dev@latest
```

Asegúrate de que `$HOME/go/bin` esté en tu PATH, o agrega:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Uso

Simplemente ejecuta:

```bash
laravel-dev
```

Esto abrirá una interfaz interactiva donde podrás:

1. Ver el estado de tu sistema
2. Seleccionar qué componentes instalar
3. Elegir versiones específicas
4. Configurar opciones adicionales

### Opciones de línea de comandos

```bash
laravel-dev --help    # Muestra la ayuda
```

## Características

### PHP
- Selección de versión (PHP 7.4 - PHP 8.4)
- Instalación de extensiones comunes:
  - pdo, mysql, mysqli, sqlite3
  - gd, imagick, zip
  - mbstring, curl, xml
  - redis, memcached
  - xdebug

### MariaDB
- Instalación de MariaDB Server
- Configuración segura interactiva
- Configuración de contraseña root

### Node.js
- Instalación de Node.js (versión personalizada)
- npm incluido

### Composer
- Instalación global de Composer
- Configuración PATH automática

### Laravel Valet
- Instalación de Valet
- Verificación de requisitos del sistema
- Configuración de dominio local (.test)

## Desinstalación

```bash
go uninstall github.com/LC-jhony/laravel-dev@latest
```

## Desarrollo

### Compilar desde código fuente

```bash
git clone https://github.com/LC-jhony/laravel-dev.git
cd laravel-dev
go build -o laravel-dev .
```

### Ejecutar en modo desarrollo

```bash
go run .
```

### Ejecutar pruebas

```bash
go test ./...
```

## Estructura del proyecto

```
laravel-dev/
├── Main.go              # Punto de entrada
├── cmd/
│   └── Welcome.go       # Pantalla de bienvenida
├── pkg/
│   ├── php.go           # Instalador de PHP
│   ├── mariadb.go       # Instalador de MariaDB
│   ├── nodejs.go        # Instalador de Node.js
│   ├── composer.go      # Instalador de Composer
│   └── valet.go         # Instalador de Laravel Valet
└── install.sh           # Script de instalación
```

## Tecnologías usadas

- [Go](https://go.dev/) - Lenguaje de programación
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Estilos de terminal
- [Huh](https://github.com/charmbracelet/huh) - Formularios interactivos
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Framework TUI

## Licencia

MIT License

## Contribuciones

Las contribuciones son bienvenidas. Por favor, abre un issue o pull request en [GitHub](https://github.com/LC-jhony/laravel-dev).
