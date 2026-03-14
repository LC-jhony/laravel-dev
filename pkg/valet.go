package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	valetTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("13")).
			Width(50).
			Align(lipgloss.Center)

	valetInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Width(50).
			Align(lipgloss.Center)

	valetErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Width(50).
			Align(lipgloss.Center)

	valetSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("10")).
				Width(50).
				Align(lipgloss.Center)

	valetWarningStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("11")).
				Width(50).
				Align(lipgloss.Center)
)

func ShowValetSelector() error {
	if !isInteractive() {
		fmt.Println()
		fmt.Println(valetInfoStyle.Render("Ejecuta el programa en una terminal interactiva."))
		fmt.Println(valetInfoStyle.Render("O usa: laravel-dev install valet"))
		return nil
	}

	var installConfirm bool
	var installLaravel bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("¿Instalar Laravel Valet?").
				Affirmative("Sí, instalar Valet").
				Negative("No, cancelar").
				Value(&installConfirm).
				Description("Entorno de desarrollo Laravel para Linux"),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("¿Instalar Laravel Installer?").
				Affirmative("Sí, instalar Laravel Installer").
				Negative("No, omitir").
				Value(&installLaravel).
				Description("Para crear nuevos proyectos Laravel"),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el formulario: %w", err)
	}

	if !installConfirm {
		fmt.Println()
		fmt.Println(valetInfoStyle.Render("Instalación cancelada por el usuario."))
		return nil
	}

	if err := InstallValetRequirements(); err != nil {
		return err
	}

	if err := InstallValet(); err != nil {
		return err
	}

	if err := SetupValet(); err != nil {
		return err
	}

	if installLaravel {
		if err := InstallLaravelInstaller(); err != nil {
			return err
		}
	}

	return nil
}

func InstallValetRequirements() error {
	fmt.Println()
	fmt.Println(valetTitleStyle.Render("Instalando requerimientos de Valet..."))
	fmt.Println()

	fmt.Println(valetInfoStyle.Render("Requerimientos: network-manager, libnss3-tools, jq, xsel"))

	if err := runValetCommand("sudo", "apt-get", "update"); err != nil {
		return fmt.Errorf("error al actualizar repositorios: %w", err)
	}

	if err := runValetCommand("sudo", "apt-get", "install", "-y", "network-manager", "libnss3-tools", "jq", "xsel"); err != nil {
		return fmt.Errorf("error al instalar requerimientos: %w", err)
	}

	fmt.Println()
	fmt.Println(valetSuccessStyle.Render("✓ Requerimientos instalados correctamente!"))

	return nil
}

func InstallValet() error {
	fmt.Println()
	fmt.Println(valetTitleStyle.Render("Instalando Valet..."))
	fmt.Println()

	if _, err := exec.LookPath("composer"); err != nil {
		fmt.Println(valetWarningStyle.Render("Composer no encontrado. Instalando Composer primero..."))
		if err := InstallComposer(); err != nil {
			return fmt.Errorf("error al instalar Composer: %w", err)
		}
	}

	fmt.Println(valetInfoStyle.Render("Instalando valet-linux globalmente..."))
	fmt.Println(valetInfoStyle.Render("Esto puede tomar unos minutos..."))

	installCmd := `export PATH="$HOME/.composer/vendor/bin:$PATH" && composer global require cpriego/valet-linux`

	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al instalar valet-linux: %w", err)
	}

	fmt.Println()
	fmt.Println(valetSuccessStyle.Render("✓ Valet instalado correctamente!"))

	return nil
}

func SetupValet() error {
	fmt.Println()
	fmt.Println(valetTitleStyle.Render("Configurando Valet..."))
	fmt.Println()

	fmt.Println(valetInfoStyle.Render("Ejecutando valet install..."))
	fmt.Println(valetWarningStyle.Render("Esto puede pedir tu contraseña de sudo..."))

	setupCmd := `export PATH="$HOME/.composer/vendor/bin:$PATH" && valet install`

	cmd := exec.Command("bash", "-c", setupCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al ejecutar valet install: %w", err)
	}

	fmt.Println()
	fmt.Println(valetInfoStyle.Render("Creando directorio ~/Sites..."))

	sitesPath := os.Getenv("HOME") + "/Sites"
	if _, err := os.Stat(sitesPath); os.IsNotExist(err) {
		if err := os.MkdirAll(sitesPath, 0755); err != nil {
			return fmt.Errorf("error al crear directorio Sites: %w", err)
		}
	}

	fmt.Println(valetInfoStyle.Render("Ejecutando valet park..."))

	parkCmd := `export PATH="$HOME/.composer/vendor/bin:$PATH" && cd $HOME/Sites && valet park`

	cmd = exec.Command("bash", "-c", parkCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al ejecutar valet park: %w", err)
	}

	fmt.Println()
	fmt.Println(valetSuccessStyle.Render("✓ Valet configurado correctamente!"))

	return nil
}

func InstallLaravelInstaller() error {
	fmt.Println()
	fmt.Println(valetTitleStyle.Render("Instalando Laravel Installer..."))
	fmt.Println()

	fmt.Println(valetInfoStyle.Render("Instalando laravel/installer globalmente..."))

	installCmd := `export PATH="$HOME/.composer/vendor/bin:$PATH" && composer global require laravel/installer`

	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error al instalar Laravel Installer: %w", err)
	}

	fmt.Println()
	fmt.Println(valetSuccessStyle.Render("✓ Laravel Installer instalado correctamente!"))
	fmt.Println()
	fmt.Println(valetInfoStyle.Render("Para crear un nuevo proyecto Laravel:"))
	fmt.Println(valetInfoStyle.Render("  laravel new mi-proyecto"))

	return nil
}

func GetValetInstallCommand() string {
	return `sudo apt-get install network-manager libnss3-tools jq xsel && composer global require cpriego/valet-linux && valet install && mkdir -p ~/Sites && cd ~/Sites && valet park && composer global require laravel/installer`
}

func ValetStatus() error {
	cmd := exec.Command("bash", "-c", `export PATH="$HOME/.composer/vendor/bin:$PATH" && valet status 2>/dev/null || echo "Valet no está configurado"`)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(valetInfoStyle.Render("Valet no está instalado."))
		return nil
	}

	fmt.Println(valetSuccessStyle.Render(strings.TrimSpace(string(output))))
	return nil
}

func runValetCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 0 {
				return fmt.Errorf("el comando salió con código de error: %d", exitErr.ExitCode())
			}
		}
		return err
	}
	return nil
}
