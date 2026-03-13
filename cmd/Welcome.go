package cmd

import (
	"fmt"
	"runtime"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Width(50).
			Foreground(lipgloss.Color("9")).
			Align(lipgloss.Center).
			PaddingTop(1).
			PaddingBottom(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Width(50).
			Align(lipgloss.Center)

	sectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("12")).
				Width(50).
				Align(lipgloss.Left).
				PaddingTop(1)

	infoKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Width(18).
			Align(lipgloss.Left)

	infoValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Width(32).
			Align(lipgloss.Left)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Width(15).
			Align(lipgloss.Left)

	commandDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Width(35).
				Align(lipgloss.Left)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(50)

	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(50).
			Align(lipgloss.Center)
)

const AppVersion = "v1.0.0"

func ShowWelcome() {
	infoBlock := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Left, infoKeyStyle.Render("Sistema:"), infoValueStyle.Render(runtime.GOOS)),
		lipgloss.JoinHorizontal(lipgloss.Left, infoKeyStyle.Render("Arquitectura:"), infoValueStyle.Render(runtime.GOARCH)),
		lipgloss.JoinHorizontal(lipgloss.Left, infoKeyStyle.Render("Nucleos CPU:"), infoValueStyle.Render(fmt.Sprintf("%d", runtime.NumCPU()))),
	//	lipgloss.JoinHorizontal(lipgloss.Left, infoKeyStyle.Render("Version de Go:"), infoValueStyle.Render(runtime.Version())),
	)

	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("LARAVEL DEV TOOLS"),
		subtitleStyle.Render("Herramienta de Desarrollo Laravel"),
		subtitleStyle.Render("Instalador de PHP, MariaDB, NodeJS, composer, Laravel Valet y Laravel installer desde codigo fuente"),
		dividerStyle.Render(""),
		sectionTitleStyle.Render("Informacion del Sistema"),
		infoBlock,
		dividerStyle.Render(""),
	)

	fmt.Println(mainContent)
}
