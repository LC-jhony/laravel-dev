package main

import (
	"fmt"
	"laravel-dev/pkg"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

var (
	uiTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9")).
			Width(70).
			Align(lipgloss.Center).
			PaddingTop(1).
			PaddingBottom(1)

	uiSubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Width(70).
			Align(lipgloss.Center)

	uiInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Width(70).
			Align(lipgloss.Center)

	uiSuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Width(70).
			Align(lipgloss.Center).
			Bold(true)

	uiErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Width(70).
			Align(lipgloss.Center)

	uiBoxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("236")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9")).
			Width(72).
			Padding(1, 2)
)

type ComponentStatus struct {
	Name      string
	Icon      string
	Installed bool
	Version   string
}

func main() {
	showWelcomeUI()

	if !isInteractiveTerminal() {
		showNonInteractiveMenu()
		return
	}

	showSystemStatus()
	showInstallerSelector()
}

func showWelcomeUI() {
	title := uiTitleStyle.Render("Laravel Dev Tools")
	subtitle := uiSubtitleStyle.Render("Instalador de entorno de desarrollo")

	sysInfo := lipgloss.JoinVertical(
		lipgloss.Center,
		uiInfoStyle.Render(fmt.Sprintf("Sistema: %s | Arquitectura: %s", runtime.GOOS, runtime.GOARCH)),
	)

	box := uiBoxStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			"",
			subtitle,
			"",
			sysInfo,
		),
	)

	fmt.Println(box)
}

func showSystemStatus() {
	components := checkAllComponents()

	if !hasAnyComponent(components) {
		fmt.Println()
		fmt.Println(uiTitleStyle.Render(" Estado del Sistema "))
		fmt.Println()
		fmt.Println(uiInfoStyle.Render("No se encontró ningún componente instalado."))
		fmt.Println(uiInfoStyle.Render("Selecciona los componentes a instalar."))
		return
	}

	fmt.Println()
	fmt.Println(uiTitleStyle.Render(" Estado del Sistema "))
	fmt.Println()

	renderTable(components)
}

func checkAllComponents() []ComponentStatus {
	return []ComponentStatus{
		checkPHP(),
		checkMariaDB(),
		checkNode(),
		checkComposer(),
		checkValet(),
	}
}

func checkPHP() ComponentStatus {
	cmd := exec.Command("bash", "-c", "php -v 2>/dev/null")
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	if err != nil || outputStr == "" {
		return ComponentStatus{Name: "PHP", Icon: "🐘", Installed: false, Version: "No instalado"}
	}
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "php") {
			return ComponentStatus{Name: "PHP", Icon: "🐘", Installed: true, Version: extractVersion(line)}
		}
	}
	return ComponentStatus{Name: "PHP", Icon: "🐘", Installed: true, Version: outputStr}
}

func checkMariaDB() ComponentStatus {
	cmd := exec.Command("bash", "-c", "mysql -V 2>/dev/null || mariadb -V 2>/dev/null")
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	if err != nil || outputStr == "" {
		return ComponentStatus{Name: "MariaDB", Icon: "🗄️", Installed: false, Version: "No instalado"}
	}
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "maria") || strings.Contains(strings.ToLower(line), "mysql") {
			return ComponentStatus{Name: "MariaDB", Icon: "🗄️", Installed: true, Version: extractVersion(line)}
		}
	}
	return ComponentStatus{Name: "MariaDB", Icon: "🗄️", Installed: true, Version: outputStr}
}

func checkNode() ComponentStatus {
	cmd := exec.Command("bash", "-c", `export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" && node -v 2>/dev/null`)
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))
	if err != nil || outputStr == "" || strings.Contains(outputStr, "command not found") {
		return ComponentStatus{Name: "Node.js", Icon: "🟢", Installed: false, Version: "No instalado"}
	}
	return ComponentStatus{Name: "Node.js", Icon: "🟢", Installed: true, Version: extractVersion(outputStr)}
}

func checkComposer() ComponentStatus {
	cmd := exec.Command("bash", "-c", "composer --version 2>/dev/null")
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	if err != nil || outputStr == "" {
		return ComponentStatus{Name: "Composer", Icon: "📦", Installed: false, Version: "No instalado"}
	}
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "composer") {
			return ComponentStatus{Name: "Composer", Icon: "📦", Installed: true, Version: extractVersion(line)}
		}
	}
	return ComponentStatus{Name: "Composer", Icon: "📦", Installed: true, Version: outputStr}
}

func checkValet() ComponentStatus {
	cmd := exec.Command("bash", "-c", `export PATH="$HOME/.composer/vendor/bin:$PATH" && valet --version 2>/dev/null`)
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))
	if err != nil || outputStr == "" || strings.Contains(outputStr, "command not found") {
		return ComponentStatus{Name: "Valet", Icon: "🚀", Installed: false, Version: "No instalado"}
	}
	return ComponentStatus{Name: "Valet", Icon: "🚀", Installed: true, Version: extractVersion(outputStr)}
}

func extractVersion(output string) string {
	output = strings.TrimSpace(output)

	// Para MariaDB: "mysql  Ver 15.1 Distrib 10.11.14-MariaDB, for debian-linux-gnu..."
	if strings.Contains(output, "MariaDB") {
		parts := strings.Split(output, "Distrib ")
		if len(parts) >= 2 {
			version := strings.Split(parts[1], ",")[0]
			return "v" + strings.TrimSpace(version)
		}
	}

	// Para PHP: "PHP 8.4.18 (cli)..."
	if strings.Contains(output, "PHP") {
		parts := strings.Fields(output)
		if len(parts) >= 2 {
			return "v" + parts[1]
		}
	}

	// Para Composer: "Composer version 2.9.5..."
	if strings.Contains(output, "Composer") {
		parts := strings.Fields(output)
		for i, part := range parts {
			if part == "version" && i+1 < len(parts) {
				return "v" + parts[i+1]
			}
		}
	}

	// Para Node.js: "v20.10.0"
	if strings.HasPrefix(output, "v") {
		return output
	}

	// Buscar cualquier número de versión (x.y o x.y.z)
	re := regexp.MustCompile(`(\d+\.\d+(\.\d+)?)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return "v" + matches[1]
	}

	return output
}

func hasAnyComponent(components []ComponentStatus) bool {
	for _, c := range components {
		if c.Installed {
			return true
		}
	}
	return false
}

func renderTable(components []ComponentStatus) {
	renderStatusList(components)
}

func renderStatusList(components []ComponentStatus) {
	for _, c := range components {
		renderComponentRow(c)
	}
}

func renderComponentRow(c ComponentStatus) {
	iconStyle := lipgloss.NewStyle().Width(3)
	nameStyle := lipgloss.NewStyle().Width(15).Bold(true)
	versionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 2).
		Width(50)

	if c.Installed {
		// Componente instalado - estilo verde
		iconStyle = iconStyle.Foreground(lipgloss.Color("10"))
		nameStyle = nameStyle.Foreground(lipgloss.Color("10"))
		boxStyle = boxStyle.BorderForeground(lipgloss.Color("10"))

		content := lipgloss.JoinHorizontal(
			lipgloss.Left,
			iconStyle.Render("✓"),
			nameStyle.Render(c.Name),
			versionStyle.Render(c.Version),
		)
		fmt.Println(boxStyle.Render(content))
	} else {
		// Componente no instalado - estilo rojo/gris
		iconStyle = iconStyle.Foreground(lipgloss.Color("9"))
		nameStyle = nameStyle.Foreground(lipgloss.Color("9"))
		boxStyle = boxStyle.BorderForeground(lipgloss.Color("9"))

		content := lipgloss.JoinHorizontal(
			lipgloss.Left,
			iconStyle.Render("✗"),
			nameStyle.Render(c.Name),
			versionStyle.Render("No instalado"),
		)
		fmt.Println(boxStyle.Render(content))
	}
}

func cleanVersion(version string) string {
	// Extraer solo la versión numérica
	re := regexp.MustCompile(`(\d+\.\d+(\.\d+)?)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) > 0 {
		return "v" + matches[1]
	}
	// Si no encuentra patrón, limpiar manually
	version = strings.ReplaceAll(version, "Ver ", "")
	version = strings.ReplaceAll(version, "Distrib ", "")
	parts := strings.Split(version, ",")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return version
}

func showNonInteractiveMenu() {
	fmt.Println()
	fmt.Println(uiInfoStyle.Render("Ejecuta el programa en una terminal interactiva."))
	fmt.Println()
	fmt.Println(uiBoxStyle.Render("Comandos disponibles:"))
	fmt.Println()
	fmt.Println("  " + uiInfoStyle.Render("laravel-dev install php      - Instalar PHP"))
	fmt.Println("  " + uiInfoStyle.Render("laravel-dev install mysql   - Instalar MariaDB"))
	fmt.Println("  " + uiInfoStyle.Render("laravel-dev install node    - Instalar Node.js"))
	fmt.Println("  " + uiInfoStyle.Render("laravel-dev install composer - Instalar Composer"))
	fmt.Println("  " + uiInfoStyle.Render("laravel-dev install valet   - Instalar Laravel Valet"))
	fmt.Println("  " + uiInfoStyle.Render("laravel-dev install all     - Instalar todo"))
}

func showInstallerSelector() {
	components := checkAllComponents()
	var options []string

	for _, c := range components {
		if !c.Installed {
			options = append(options, c.Name)
		}
	}

	if len(options) == 0 {
		fmt.Println()
		fmt.Println(uiSuccessStyle.Render("✨ Todos los componentes ya están instalados!"))
		fmt.Println()
		fmt.Println(uiInfoStyle.Render("Para actualizar, desinstala y vuelve a instalar."))
		return
	}

	var selectedOptions []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Selecciona los componentes a instalar:").
				Options(
					getInstallOptions(components)...,
				).
				Value(&selectedOptions).
				Description("Usa: Espacio=seleccionar, Enter=continuar, q=cancelar").
				Limit(5).
				Height(10),
		),
	).WithWidth(75).WithKeyMap(&huh.KeyMap{
		Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "cancelar")),
	})

	err := form.Run()
	if err != nil {
		if err.Error() == "user aborted" {
			fmt.Println()
			fmt.Println(uiInfoStyle.Render("Operación cancelada."))
			return
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	if len(selectedOptions) == 0 {
		fmt.Println()
		fmt.Println(uiInfoStyle.Render("No seleccionaste ningún componente."))
		return
	}

	fmt.Println()
	fmt.Println(uiSuccessStyle.Render("⚡ Instalando componentes seleccionados..."))
	fmt.Println()

	installAll(selectedOptions)
}

func getInstallOptions(components []ComponentStatus) []huh.Option[string] {
	var opts []huh.Option[string]

	optionMap := map[string]struct {
		key  string
		desc string
	}{
		"PHP":      {"php", "Intérprete de PHP"},
		"MariaDB":  {"mariadb", "Base de datos MySQL"},
		"Node.js":  {"node", "Runtime de JavaScript"},
		"Composer": {"composer", "Gestor de dependencias PHP"},
		"Valet":    {"valet", "Entorno desarrollo Laravel"},
	}

	for _, c := range components {
		if !c.Installed {
			if opt, ok := optionMap[c.Name]; ok {
				selected := false
				if c.Name == "PHP" || c.Name == "MariaDB" {
					selected = true
				}
				opts = append(opts, huh.NewOption(
					fmt.Sprintf("%s %s - %s", c.Icon, c.Name, opt.desc),
					opt.key,
				).Selected(selected))
			}
		}
	}

	if len(opts) == 0 {
		opts = append(opts, huh.NewOption("Todos instalados", "none"))
	}

	return opts
}

func installAll(options []string) {
	order := []string{"php", "mariadb", "composer", "node", "valet"}
	installed := make(map[string]bool)

	for _, item := range order {
		if !contains(options, item) {
			continue
		}

		if item == "none" {
			continue
		}

		showProgress(item, len(installed)+1, len(options))

		var err error
		switch item {
		case "php":
			err = pkg.ShowVersionSelector()
		case "mariadb":
			err = pkg.ShowMariaDBSelector()
		case "node":
			err = pkg.ShowNodeVersionSelector()
		case "composer":
			err = pkg.ShowComposerSelector()
		case "valet":
			err = pkg.ShowValetSelector()
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s Error al instalar %s: %v\n", uiErrorStyle.Render(""), getComponentName(item), err)
			continue
		}

		installed[item] = true
	}

	showSummary(installed, options, order)
}

func showProgress(item string, current, total int) {
	icon := getComponentIcon(item)
	name := getComponentName(item)
	bar := createProgressBar(current, total)

	fmt.Println()
	fmt.Println(uiBoxStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			uiTitleStyle.Render(fmt.Sprintf("Instalando (%d/%d)", current, total)),
			"",
			uiSuccessStyle.Render(fmt.Sprintf("%s %s", icon, name)),
			"",
			uiInfoStyle.Render(bar),
		),
	))
}

func createProgressBar(current, total int) string {
	width := 30
	filled := (width * current) / total
	empty := width - filled

	result := "[" + repeat("▓", filled) + repeat("░", empty) + "]"
	return result
}

func showSummary(installed map[string]bool, options []string, order []string) {
	fmt.Println()
	fmt.Println(uiSuccessStyle.Render("✅ Instalación completada!"))
	fmt.Println()
	fmt.Println(uiBoxStyle.Render("Resumen:"))

	successCount := 0
	for _, item := range order {
		if contains(options, item) {
			if item == "none" {
				continue
			}
			icon := getComponentIcon(item)
			name := getComponentName(item)
			status := "❌"
			if installed[item] {
				status = "✅"
				successCount++
			}
			fmt.Printf("  %s %s %s\n", status, icon, name)
		}
	}

	fmt.Println()
	fmt.Println(uiInfoStyle.Render(fmt.Sprintf("Total: %d/%d instalados correctamente", successCount, len(options))))

	showSystemStatus()
}

func getComponentIcon(item string) string {
	switch item {
	case "php":
		return "🐘"
	case "mariadb":
		return "🗄️"
	case "node":
		return "🟢"
	case "composer":
		return "📦"
	case "valet":
		return "🚀"
	default:
		return "•"
	}
}

func getComponentName(item string) string {
	switch item {
	case "php":
		return "PHP"
	case "mariadb":
		return "MariaDB"
	case "node":
		return "Node.js"
	case "composer":
		return "Composer"
	case "valet":
		return "Laravel Valet"
	default:
		return item
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func isInteractiveTerminal() bool {
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
