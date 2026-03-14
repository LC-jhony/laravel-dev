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
	composerTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("14")).
				Width(50).
				Align(lipgloss.Center)

	composerInfoStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Width(50).
				Align(lipgloss.Center)

	composerErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("9")).
				Width(50).
				Align(lipgloss.Center)

	composerSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("10")).
				Width(50).
				Align(lipgloss.Center)
)

func ShowComposerSelector() error {
	if !isInteractive() {
		fmt.Println()
		fmt.Println(composerInfoStyle.Render("Ejecuta el programa en una terminal interactiva."))
		fmt.Println(composerInfoStyle.Render("O usa: laravel-dev install composer"))
		return nil
	}

	var installConfirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("¿Instalar Composer?").
				Affirmative("Sí, instalar Composer").
				Negative("No, cancelar").
				Value(&installConfirm).
				Description("Administrador de dependencias para PHP"),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el formulario: %w", err)
	}

	if !installConfirm {
		fmt.Println()
		fmt.Println(composerInfoStyle.Render("Instalación cancelada por el usuario."))
		return nil
	}

	return InstallComposer()
}

func InstallComposer() error {
	fmt.Println()
	fmt.Println(composerTitleStyle.Render("Instalando Composer..."))
	fmt.Println()

	if _, err := exec.LookPath("php"); err != nil {
		return fmt.Errorf("PHP no está instalado. Instala PHP primero con: laravel-dev install php")
	}

	fmt.Println(composerInfoStyle.Render("1/4 Descargando instalador..."))
	if err := runPHPCommand("-r", "copy('https://getcomposer.org/installer', 'composer-setup.php');"); err != nil {
		return fmt.Errorf("error al descargar instalador: %w", err)
	}

	fmt.Println(composerInfoStyle.Render("2/4 Verificando instalador..."))
	hashVerify := `if (hash_file('sha384', 'composer-setup.php') === 'c8b085408188070d5f52bcfe4ecfbee5f727afa458b2573b8eaaf77b3419b0bf2768dc67c86944da1544f06fa544fd47') { echo 'Installer verified'.PHP_EOL; } else { echo 'Installer corrupt'.PHP_EOL; unlink('composer-setup.php'); exit(1); }`
	if err := runPHPCommand("-r", hashVerify); err != nil {
		return fmt.Errorf("error al verificar instalador: %w", err)
	}

	fmt.Println(composerInfoStyle.Render("3/4 Ejecutando instalador..."))
	if err := runPHPCommand("composer-setup.php"); err != nil {
		return fmt.Errorf("error al ejecutar instalador: %w", err)
	}

	fmt.Println(composerInfoStyle.Render("4/4 Instalando globalmente..."))
	if err := runCommand("sudo", "mv", "composer.phar", "/usr/local/bin/composer"); err != nil {
		return fmt.Errorf("error al mover composer: %w", err)
	}

	fmt.Println(composerInfoStyle.Render("Limpiando archivos temporales..."))
	if err := runPHPCommand("-r", "unlink('composer-setup.php');"); err != nil {
		return fmt.Errorf("error al limpiar: %w", err)
	}

	fmt.Println()
	fmt.Println(composerSuccessStyle.Render("✓ Composer instalado correctamente!"))
	fmt.Println()
	fmt.Println(composerInfoStyle.Render("Para verificar la instalación, ejecuta: composer --version"))

	return nil
}

func runPHPCommand(args ...string) error {
	cmd := exec.Command("php", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func GetComposerInstallCommand() string {
	return `php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');" && php -r "if (hash_file('sha384', 'composer-setup.php') === 'c8b085408188070d5f52bcfe4ecfbee5f727afa458b2573b8eaaf77b3419b0bf2768dc67c86944da1544f06fa544fd47') { echo 'Installer verified'.PHP_EOL; } else { echo 'Installer corrupt'.PHP_EOL; unlink('composer-setup.php'); exit(1); }" && php composer-setup.php && php -r "unlink('composer-setup.php');" && sudo mv composer.phar /usr/local/bin/composer`
}

func ComposerVersion() error {
	cmd := exec.Command("composer", "--version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(composerInfoStyle.Render("Composer no está instalado."))
		return nil
	}

	fmt.Println(composerSuccessStyle.Render(strings.TrimSpace(string(output))))
	return nil
}
