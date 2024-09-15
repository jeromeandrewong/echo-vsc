package main

import (
	"fmt"
	"os"
	"path/filepath"

	"echo/internal/converter"
	log "echo/internal/logger"
	"echo/internal/theme"
	"echo/internal/vsc"
	"echo/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
)

const VSC_EXTENSION_PATH = "/.vscode/extensions"

type PackageData struct {
	Contributes struct {
		Themes []theme.Theme `json:"themes"`
	} `json:"contributes"`
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("⚠️ Error getting home directory", "error", err)
	}

	vscDir := homeDir + VSC_EXTENSION_PATH

	themes, err := vsc.GetVSCThemes(vscDir)
	if err != nil {
		log.Fatal("⚠️ Failed to get VSC themes", "error", err)
	}

	p := tea.NewProgram(theme.New(themes), tea.WithAltScreen())

	m, err := p.Run()
	if err != nil {
		log.Fatal("⚠️ Error running program", "error", err)
		os.Exit(1)
	}

	if m, ok := m.(theme.Model); ok && m.Choice.Label != "" {
		downloadsDir, err := utils.GetDownloadsFolder()
		if err != nil {
			log.Error("🚨 Failed to get Downloads folder", "error", err)
			os.Exit(1)
		}

		var filePath string

		options := converter.ThemeOptions{
			Theme:       m.Choice,
			Directory:   downloadsDir,
			ShouldWrite: true,
		}
		filePath, err = converter.GenerateTheme(options)

		if err != nil {
			log.Error("🚨 Failed to generate iTerm theme", "error", err)
		} else {
			fileURL := fmt.Sprintf("file://%s", filepath.ToSlash(filePath))

			successMessage := fmt.Sprintf("🎉 Theme generated @%s", fileURL)
			log.Info(successMessage)
		}
	} else {
		fmt.Println("😿 No theme selected, quitting echo!")
	}
}
