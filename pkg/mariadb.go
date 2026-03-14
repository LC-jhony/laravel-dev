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
	mariaTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9")).
			Width(50).
			Align(lipgloss.Center)

	mariaInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Width(50).
			Align(lipgloss.Center)

	mariaErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Width(50).
			Align(lipgloss.Center)

	mariaSuccessStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("10")).
				Width(50).
				Align(lipgloss.Center)
)

func ShowMariaDBSelector() error {
	if !isInteractive() {
		fmt.Println()
		fmt.Println(mariaInfoStyle.Render("Ejecuta el programa en una terminal interactiva para usar el selector de MariaDB."))
		fmt.Println(mariaInfoStyle.Render("O usa: laravel-dev install mariadb"))
		return nil
	}

	var installConfirm bool
	var secureConfirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("¿Instalar MariaDB Server?").
				Affirmative("Sí, instalar MariaDB").
				Negative("No, cancelar").
				Value(&installConfirm).
				Description("Se instalará MariaDB Server y Client"),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el formulario: %w", err)
	}

	if !installConfirm {
		fmt.Println()
		fmt.Println(mariaInfoStyle.Render("Instalación cancelada por el usuario."))
		return nil
	}

	if err := InstallMariaDB(); err != nil {
		return err
	}

	secureForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("¿Ejecutar mariadb-secure-installation?").
				Affirmative("Sí, configurar seguridad").
				Negative("No, omitir").
				Value(&secureConfirm).
				Description("Configurar contraseña root y opciones de seguridad"),
		),
	)

	err = secureForm.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el formulario: %w", err)
	}

	if secureConfirm {
		if err := SecureMariaDB(); err != nil {
			return err
		}
	}

	return nil
}

func InstallMariaDB() error {
	fmt.Println()
	fmt.Println(mariaTitleStyle.Render("Instalando MariaDB..."))
	fmt.Println()

	var sudoPassword string
	if os.Geteuid() != 0 {
		fmt.Println(mariaInfoStyle.Render("Este comando requiere permisos de sudo..."))

		pwd, err := askSudoPassword()
		if err != nil {
			return err
		}
		sudoPassword = pwd
	}

	fmt.Println(mariaInfoStyle.Render("1/2 Instalando MariaDB Server y Client..."))
	if err := runMariaDBCommandWithPassword(sudoPassword, "sudo", "apt", "install", "-y", "mariadb-server", "mariadb-client"); err != nil {
		return fmt.Errorf("error al instalar MariaDB: %w", err)
	}

	fmt.Println(mariaInfoStyle.Render("2/2 Iniciando servicio MariaDB..."))
	if err := runMariaDBCommandWithPassword(sudoPassword, "sudo", "systemctl", "start", "mariadb"); err != nil {
		return fmt.Errorf("error al iniciar MariaDB: %w", err)
	}

	if err := runMariaDBCommandWithPassword(sudoPassword, "sudo", "systemctl", "enable", "mariadb"); err != nil {
		return fmt.Errorf("error al habilitar MariaDB: %w", err)
	}

	fmt.Println()
	fmt.Println(mariaSuccessStyle.Render("✓ MariaDB instalado correctamente!"))
	fmt.Println()
	fmt.Println(mariaInfoStyle.Render("Para verificar la instalación, ejecuta: mysql -V"))

	return nil
}

func SecureMariaDB() error {
	fmt.Println()
	fmt.Println(mariaTitleStyle.Render("Configurando MariaDB..."))
	fmt.Println()

	var sudoPassword string
	if os.Geteuid() != 0 {
		pwd, err := askSudoPassword()
		if err != nil {
			return err
		}
		sudoPassword = pwd
	}

	fmt.Println(mariaInfoStyle.Render("Ejecutando mariadb-secure-installation..."))
	fmt.Println()
	fmt.Println(mariaInfoStyle.Render("Sigue las instrucciones en pantalla:"))
	fmt.Println(mariaInfoStyle.Render(" 1. Enter (sin contraseña actual)"))
	fmt.Println(mariaInfoStyle.Render(" 2. Set root password? [Y/n] → Y"))
	fmt.Println(mariaInfoStyle.Render(" 3. Nueva contraseña root"))
	fmt.Println(mariaInfoStyle.Render(" 4. Confirmar contraseña"))
	fmt.Println(mariaInfoStyle.Render(" 5. Remove anonymous users? [Y/n] → Y"))
	fmt.Println(mariaInfoStyle.Render(" 6. Disallow root login remotely? [Y/n] → Y"))
	fmt.Println(mariaInfoStyle.Render(" 7. Remove test database? [Y/n] → Y"))
	fmt.Println(mariaInfoStyle.Render(" 8. Reload privilege tables? [Y/n] → Y"))
	fmt.Println()

	if err := runMariaDBCommandWithPassword(sudoPassword, "sudo", "mariadb-secure-installation"); err != nil {
		return fmt.Errorf("error al ejecutar mariadb-secure-installation: %w", err)
	}

	fmt.Println()
	fmt.Println(mariaSuccessStyle.Render("✓ MariaDB configurado correctamente!"))
	fmt.Println()
	fmt.Println(mariaInfoStyle.Render("Para conectar: mysql -u root -p"))

	return nil
}

func GetMariaDBInstallCommand() string {
	return "sudo apt install mariadb-server mariadb-client && sudo mariadb-secure-installation"
}

func runMariaDBCommand(name string, args ...string) error {
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

func MariaDBStatus() error {
	cmd := exec.Command("systemctl", "is-active", "mariadb")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("MariaDB no está en ejecución")
	}
	fmt.Println(mariaSuccessStyle.Render("✓ MariaDB está en ejecución"))
	return nil
}

func MariaDBVersion() error {
	cmd := exec.Command("mysql", "-V")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error al obtener versión de MariaDB: %w", err)
	}
	fmt.Println(mariaInfoStyle.Render(strings.TrimSpace(string(output))))
	return nil
}

func runMariaDBCommandWithPassword(password string, name string, args ...string) error {
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
