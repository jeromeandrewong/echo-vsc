package main

import (
	"fmt"
	"os"

	log "echo/internal/logger"
	"echo/internal/themepicker"
	"echo/internal/vsc"

	tea "github.com/charmbracelet/bubbletea"
)

const VSC_EXTENSION_PATH = "/.vscode/extensions"

type PackageData struct {
	Contributes struct {
		Themes []themepicker.Theme `json:"themes"`
	} `json:"contributes"`
}

func main() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error("Error getting home directory", "error", err)
	}

	vscDir := homeDir + VSC_EXTENSION_PATH

	themes, err := vsc.GetVSCThemes(vscDir)
	if err != nil {
		log.Fatal("Failed to get VSC themes", "error", err)
	}

	p := tea.NewProgram(themepicker.New(themes), tea.WithAltScreen())

	m, err := p.Run()
	if err != nil {
		log.Fatal("Error running program", "error", err)
		os.Exit(1)
	}

	if m, ok := m.(themepicker.Model); ok && m.Choice.Label != "" {
		fmt.Printf("Selected theme: %s\nPath: %s\n", m.Choice.Label, m.Choice.Path)
	} else {
		fmt.Println("No theme selected")
	}
}
