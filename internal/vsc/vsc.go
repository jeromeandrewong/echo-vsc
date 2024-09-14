package vsc

import (
	"echo/internal/logger"
	"echo/internal/themepicker"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type PackageData struct {
	Contributes struct {
		Themes []themepicker.Theme `json:"themes"`
	} `json:"contributes"`
}

func GetVSCThemes(vscDir string) ([]themepicker.Theme, error) {
	var themes []themepicker.Theme

	extensions, err := os.ReadDir(vscDir)
	if err != nil {
		return nil, fmt.Errorf("error reading VSC directory: %v", err)
	}

	for _, extension := range extensions {
		if !extension.IsDir() {
			continue
		}

		extensionThemes, err := getThemesFromExtension(vscDir, extension)
		if err != nil {
			logger.Warn("Error processing extension", "extension", extension.Name(), "error", err)
			continue
		}

		themes = append(themes, extensionThemes...)
	}

	return themes, nil
}

func getThemesFromExtension(vscDir string, extension os.DirEntry) ([]themepicker.Theme, error) {
	packageJSONPath := filepath.Join(vscDir, extension.Name(), "package.json")
	packageJSON, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return nil, fmt.Errorf("error reading package.json: %v", err)
	}

	var packageData PackageData
	if err := json.Unmarshal(packageJSON, &packageData); err != nil {
		return nil, fmt.Errorf("error parsing package.json: %v", err)
	}

	var themes []themepicker.Theme
	for _, t := range packageData.Contributes.Themes {
		themes = append(themes, themepicker.Theme{
			Label: t.Label,
			Path:  filepath.Join(vscDir, extension.Name(), t.Path),
		})
	}

	return themes, nil
}
