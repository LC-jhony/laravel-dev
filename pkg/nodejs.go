package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var nodeVersions = []struct {
	version string
	desc    string
}{
	{"22", "LTS - Versión más reciente"},
	{"20", "LTS - Versión estable actual"},
	{"18", "LTS - Versión anterior"},
	{"21", "Versión actual (latest)"},
	{"19", "Legacy"},
	{"16", "Legacy"},
}

var (
	nodeTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10")).
			Width(50).
			Align(lipgloss.Center)

	nodeInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Width(50).
			Align(lipgloss.Center)

	nodeErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Width(50).
			Align(lipgloss.Center)

	nodeSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("10")).
				Width(50).
				Align(lipgloss.Center)
)

func ShowNodeVersionSelector() error {
	if !isInteractive() {
		fmt.Println()
		fmt.Println(nodeInfoStyle.Render("Ejecuta el programa en una terminal interactiva."))
		fmt.Println(nodeInfoStyle.Render("O usa: laravel-dev install node"))
		return nil
	}

	var selectedVersion string
	var installConfirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Ingresa la versión de Node.js").
				Placeholder("Ej: 20, 22, 18, 21, etc.").
				Value(&selectedVersion).
				Description("Ingresa el número de versión (sin el prefijo 'v')").
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("la versión no puede estar vacía")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("¿Confirmar instalación?").
				Affirmative("Sí, instalar Node.js "+selectedVersion).
				Negative("No, cancelar").
				Value(&installConfirm).
				Description(fmt.Sprintf("Se instalará Node.js %s vía NVM", selectedVersion)),
		),
	).WithWidth(60)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el formulario: %w", err)
	}

	if !installConfirm {
		fmt.Println()
		fmt.Println(nodeInfoStyle.Render("Instalación cancelada por el usuario."))
		return nil
	}

	if err := InstallNVM(); err != nil {
		return err
	}

	return InstallNode(selectedVersion)
}

func InstallNVM() error {
	fmt.Println()
	fmt.Println(nodeTitleStyle.Render("Instalando NVM..."))
	fmt.Println()

	nvmPath := os.Getenv("HOME") + "/.nvm/nvm.sh"

	if _, err := os.Stat(nvmPath); err == nil {
		fmt.Println(nodeInfoStyle.Render("NVM ya está instalado."))
		return nil
	}

	fmt.Println(nodeInfoStyle.Render("Descargando e instalando NVM..."))

	installScript := `curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.4/install.sh | bash`

	cmd := exec.Command("bash", "-c", installScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al instalar NVM: %w", err)
	}

	fmt.Println()
	fmt.Println(nodeSuccessStyle.Render("✓ NVM instalado correctamente!"))

	return nil
}

func InstallNode(version string) error {
	fmt.Println()
	fmt.Println(nodeTitleStyle.Render(fmt.Sprintf("Instalando Node.js %s...", version)))
	fmt.Println()

	fmt.Println(nodeInfoStyle.Render("Cargando NVM e instalando Node.js..."))

	loadNVM := fmt.Sprintf(`export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && nvm install %s`, version)

	cmd := exec.Command("bash", "-c", loadNVM)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al instalar Node.js: %w", err)
	}

	fmt.Println()
	fmt.Println(nodeSuccessStyle.Render(fmt.Sprintf("✓ Node.js %s instalado correctamente!", version)))
	fmt.Println()
	fmt.Println(nodeInfoStyle.Render("Para usar Node.js, ejecuta:"))
	fmt.Println(nodeInfoStyle.Render("  source ~/.nvm/nvm.sh"))
	fmt.Println(nodeInfoStyle.Render(fmt.Sprintf("  nvm use %s", version)))

	return nil
}

func GetNodeInstallCommand() string {
	return `curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.4/install.sh | bash && export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && nvm install <versión>`
}

func NodeVersion() error {
	cmd := exec.Command("bash", "-c", `export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && node -v 2>/dev/null || echo "Node.js no está instalado"`)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(nodeInfoStyle.Render("Node.js no está instalado o NVM no está configurado."))
		return nil
	}

	fmt.Println(nodeSuccessStyle.Render("Node.js: " + strings.TrimSpace(string(output))))
	return nil
}
