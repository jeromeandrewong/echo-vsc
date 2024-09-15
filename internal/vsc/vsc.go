package vsc

import (
	log "echo/internal/logger"
	"echo/internal/theme"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type PackageData struct {
	Contributes struct {
		Themes []theme.Theme `json:"themes"`
	} `json:"contributes"`
}

func GetVSCThemes(vscDir string) ([]theme.Theme, error) {
	var allThemes []theme.Theme
	extensions, err := os.ReadDir(vscDir)
	if err != nil {
		return nil, fmt.Errorf("error reading VSC directory: %v", err)
	}

	themeChan := make(chan []theme.Theme)
	errorChan := make(chan error)
	var wg sync.WaitGroup

	for _, extension := range extensions {
		if !extension.IsDir() {
			continue
		}
		wg.Add(1)
		// start a goroutine for each extension dir
		// which processes an extension and sends its themes through a channel
		go func(ext os.DirEntry) {
			defer wg.Done()
			themes, err := getThemesFromExtension(vscDir, ext)
			if err != nil {
				errorChan <- fmt.Errorf("error processing extension %s: %v", ext.Name(), err)
				return
			}
			themeChan <- themes
		}(extension)
	}

	go func() {
		wg.Wait()
		close(themeChan)
		close(errorChan)
	}()

	for {
		select {
		// main goroutine collects themes from the theme channel and appends them to the allThemes slice.
		case themes, ok := <-themeChan:
			if !ok {
				return allThemes, nil

			}
			allThemes = append(allThemes, themes...)
		case err, ok := <-errorChan:
			if !ok {
				continue
			}
			log.Error("Warning: %v\n", err)
		}
	}
}

func getThemesFromExtension(vscDir string, extension os.DirEntry) ([]theme.Theme, error) {
	extensionPath := filepath.Join(vscDir, extension.Name())

	packageJSONPath := filepath.Join(extensionPath, "package.json")
	packageJSON, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return nil, fmt.Errorf("error reading package.json: %v", err)
	}

	var packageData struct {
		DisplayName string `json:"displayName"`
		Contributes struct {
			Themes []struct {
				Label string `json:"label"`
				Path  string `json:"path"`
			} `json:"themes"`
		} `json:"contributes"`
	}

	if err := json.Unmarshal(packageJSON, &packageData); err != nil {
		return nil, fmt.Errorf("error parsing package.json: %v", err)
	}

	var themes []theme.Theme

	for _, t := range packageData.Contributes.Themes {
		themePath := filepath.Join(extensionPath, t.Path)

		// Check if there's a corresponding .itermcolors file
		itermPath := strings.TrimSuffix(themePath, filepath.Ext(themePath)) + ".itermcolors"
		if _, err := os.Stat(itermPath); err == nil {
			themePath = itermPath
		}

		themes = append(themes, theme.Theme{
			Label: t.Label,
			Path:  themePath,
		})
	}

	return themes, nil
}
