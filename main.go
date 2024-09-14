package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "echo/internal/logger"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const VSC_EXTENSION_PATH = "/.vscode/extensions"

type Theme struct {
	Label string
	Path  string
}

type PackageData struct {
	Contributes struct {
		Themes []Theme `json:"themes"`
	} `json:"contributes"`
}
type model struct {
	list   list.Model
	choice Theme
}

func (t Theme) Title() string       { return t.Label }
func (t Theme) Description() string { return t.Path }
func (t Theme) FilterValue() string { return t.Label }

func initialModel(themes []Theme) model {
	items := make([]list.Item, len(themes))
	for i, theme := range themes {
		items[i] = theme
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "VSCode Themes"

	return model{list: l}
}
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			t, ok := m.list.SelectedItem().(Theme)
			if ok {
				m.choice = t
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func main() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error("Error getting home directory", "error", err)
	}

	vscDir := homeDir + VSC_EXTENSION_PATH
	log.Info("VSC directory", "path", vscDir)

	themes, err := getVSCThemes(vscDir)
	if err != nil {
		log.Fatal("Failed to get VSC themes", "error", err)
	}

	p := tea.NewProgram(initialModel(themes), tea.WithAltScreen())

	m, err := p.Run()
	if err != nil {
		log.Fatal("Error running program", "error", err)
		os.Exit(1)
	}

	if m, ok := m.(model); ok && m.choice.Label != "" {
		fmt.Printf("Selected theme: %s\nPath: %s\n", m.choice.Label, m.choice.Path)
	} else {
		fmt.Println("No theme selected")
	}
}

func getVSCThemes(vscDir string) ([]Theme, error) {
	var themes []Theme

	extensions, err := os.ReadDir(vscDir)
	if err != nil {
		return nil, fmt.Errorf("error reading VSC directory: %v", err)
	}

	// check if extension is a theme by checking package.json.contributes.theme
	for _, extension := range extensions {
		if !extension.IsDir() {
			continue
		}

		extensionThemes, err := getThemesFromExtension(vscDir, extension)
		if err != nil {
			log.Warn("Error processing extension", "extension", extension.Name(), "error", err)
			continue
		}

		themes = append(themes, extensionThemes...)
	}

	return themes, nil
}

func getThemesFromExtension(vscDir string, extension os.DirEntry) ([]Theme, error) {
	packageJSONPath := vscDir + "/" + extension.Name() + "/package.json"
	packageJSON, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return nil, fmt.Errorf("error reading package.json: %v", err)
	}

	// unmarshal package.json, handles absent filed (contributes.themes) gracefully
	var packageData PackageData
	if err := json.Unmarshal(packageJSON, &packageData); err != nil {
		return nil, fmt.Errorf("error parsing package.json: %v", err)
	}

	var themes []Theme
	for _, t := range packageData.Contributes.Themes {
		themes = append(themes, Theme{
			Label: t.Label,
			Path:  filepath.Join(vscDir, extension.Name(), t.Path),
		})
	}

	return themes, nil
}
