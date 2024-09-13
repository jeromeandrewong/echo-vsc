package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
	}

	vscDir := homeDir + VSC_EXTENSION_PATH
	fmt.Printf("VSC directory: %s\n", vscDir)

	themes, err := getVSCThemes(vscDir)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, theme := range themes {
		fmt.Printf("Theme: %s\n", theme)
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
			log.Printf("Error processing extension %s: %v", extension.Name(), err)
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
