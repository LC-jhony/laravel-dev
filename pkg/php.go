package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func init() {
	_ = os.Stdin
}

func isInteractive() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return false
	}

	if _, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err != nil {
		return false
	}

	return true
}

var phpVersions = []struct {
	version string
	desc    string
}{
	{"7.4", "Versión legacy (soporte extendido)"},
	{"8.0", "Versión estable"},
	{"8.1", "Versión estable"},
	{"8.2", "Versión estable"},
	{"8.3", "Versión estable"},
	{"8.4", "Versión estable (recomendada)"},
	{"8.5", "Versión más reciente"},
}

var phpExtensions = []string{
	"cli",
	"common",
	"curl",
	"pgsql",
	"fpm",
	"gd",
	"imap",
	"intl",
	"mbstring",
	"mysql",
	"opcache",
	"soap",
	"xml",
	"zip",
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9")).
			Width(50).
			Align(lipgloss.Center)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Width(50).
			Align(lipgloss.Center)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Width(50).
			Align(lipgloss.Center)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Width(50).
			Align(lipgloss.Center)
)

func ShowVersionSelector() error {
	if !isInteractive() {
		fmt.Println()
		fmt.Println(infoStyle.Render("Ejecuta el programa en una terminal interactiva para usar el selector de versión."))
		fmt.Println(infoStyle.Render("O usa: laravel-dev install --version=8.4"))
		return nil
	}

	var selectedVersion string
	var installConfirm bool
	var selectedExtensions []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Selecciona la versión de PHP").
				Options(
					huh.NewOption("PHP 7.4  - Versión legacy (soporte extendido)", "7.4"),
					huh.NewOption("PHP 8.0  - Versión estable", "8.0"),
					huh.NewOption("PHP 8.1  - Versión estable", "8.1"),
					huh.NewOption("PHP 8.2  - Versión estable", "8.2"),
					huh.NewOption("PHP 8.3  - Versión estable", "8.3"),
					huh.NewOption("PHP 8.4  - Versión estable (recomendada)", "8.4"),
					huh.NewOption("PHP 8.5  - Versión más reciente", "8.5"),
				).
				Value(&selectedVersion).
				Description("Usa las flechas ↑↓ para navegar, Enter para seleccionar"),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Selecciona las extensiones de PHP (Enter para confirmar)").
				Options(
					huh.NewOption("CLI", "cli").Selected(true),
					huh.NewOption("Common", "common").Selected(true),
					huh.NewOption("CURL", "curl").Selected(true),
					huh.NewOption("PostgreSQL", "pgsql").Selected(true),
					huh.NewOption("FPM", "fpm").Selected(true),
					huh.NewOption("GD", "gd").Selected(true),
					huh.NewOption("IMAP", "imap").Selected(true),
					huh.NewOption("Intl", "intl").Selected(true),
					huh.NewOption("Mbstring", "mbstring").Selected(true),
					huh.NewOption("MySQL", "mysql").Selected(true),
					huh.NewOption("OPcache", "opcache").Selected(true),
					huh.NewOption("SOAP", "soap").Selected(true),
					huh.NewOption("XML", "xml").Selected(true),
					huh.NewOption("ZIP", "zip").Selected(true),
				).
				Value(&selectedExtensions),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("¿Confirmar instalación?").
				Affirmative("Sí, instalar PHP "+selectedVersion).
				Negative("No, cancelar").
				Value(&installConfirm).
				Description(fmt.Sprintf("Se instalará PHP %s con las extensiones seleccionadas", selectedVersion)),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el formulario: %w", err)
	}

	if !installConfirm {
		fmt.Println()
		fmt.Println(infoStyle.Render("Instalación cancelada por el usuario."))
		return nil
	}

	return InstallPHP(selectedVersion, selectedExtensions)
}

func InstallPHP(version string, extensions []string) error {
	fmt.Println()
	fmt.Println(titleStyle.Render(fmt.Sprintf("Instalando PHP %s...", version)))
	fmt.Println()

	var sudoPassword string
	if os.Geteuid() != 0 {
		fmt.Println(infoStyle.Render("Este comando requiere permisos de sudo..."))

		pwd, err := askSudoPassword()
		if err != nil {
			return err
		}
		sudoPassword = pwd
	}

	fmt.Println(infoStyle.Render("1/3 Agregando PPA ondrej/php..."))
	if err := runCommandWithPassword(sudoPassword, "sudo", "add-apt-repository", "-y", "ppa:ondrej/php"); err != nil {
		return fmt.Errorf("error al agregar PPA: %w", err)
	}

	fmt.Println(infoStyle.Render("2/3 Actualizando repositorios..."))
	if err := runCommandWithPassword(sudoPassword, "sudo", "apt-get", "update"); err != nil {
		return fmt.Errorf("error al actualizar repositorios: %w", err)
	}

	fmt.Println(infoStyle.Render("3/3 Instalando PHP y extensiones..."))
	packages := buildPackageList(version, extensions)
	if err := runCommandWithPassword(sudoPassword, "sudo", append([]string{"apt-get", "install", "-y"}, packages...)...); err != nil {
		return fmt.Errorf("error al instalar PHP: %w", err)
	}

	fmt.Println()
	fmt.Println(successStyle.Render(fmt.Sprintf("✓ PHP %s instalado correctamente!", version)))
	fmt.Println()
	fmt.Println(infoStyle.Render("Para verificar la instalación, ejecuta: php -v"))

	return nil
}

func buildPackageList(version string, extensions []string) []string {
	packages := []string{fmt.Sprintf("php%s", version)}
	for _, ext := range extensions {
		packages = append(packages, fmt.Sprintf("php%s-%s", version, ext))
	}
	return packages
}

func runCommand(name string, args ...string) error {
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

func GetInstallCommand(version string) string {
	var pkgNames []string
	for _, ext := range phpExtensions {
		pkgNames = append(pkgNames, fmt.Sprintf("php%s-%s", version, ext))
	}
	return fmt.Sprintf("sudo add-apt-repository -y ppa:ondrej/php && sudo apt-get install -y php%s %s",
		version, strings.Join(pkgNames, " "))
}

func askSudoPassword() (string, error) {
	var password string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Ingresa tu contraseña de sudo").
				Placeholder("Password").
				Value(&password).
				Password(true).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("la contraseña no puede estar vacía")
					}
					return nil
				}),
		),
	).WithWidth(50)

	err := form.Run()
	if err != nil {
		return "", fmt.Errorf("error al pedir contraseña: %w", err)
	}

	return password, nil
}

func runCommandWithPassword(password string, name string, args ...string) error {
	if password != "" {
		args = append([]string{"-S", "-k"}, args...)
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if password != "" {
		cmd.Stdin = strings.NewReader(password + "\n")
	} else {
		cmd.Stdin = os.Stdin
	}

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
