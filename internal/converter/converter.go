package converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jeromeandrewong/echo-vsc/internal/constants"
	log "github.com/jeromeandrewong/echo-vsc/internal/logger"
	"github.com/jeromeandrewong/echo-vsc/internal/theme"
	"github.com/jeromeandrewong/echo-vsc/pkg/utils"
)

type ThemeOptions struct {
	Theme       theme.Theme
	Directory   string
	ShouldWrite bool
}

type vscodeTheme struct {
	Colors map[string]interface{} `json:"colors"`
	Type   string                 `json:"type"`
}

func GenerateTheme(options ThemeOptions) (string, error) {
	if options.Directory == "" {
		options.Directory = ""
	}

	if !options.ShouldWrite {
		options.ShouldWrite = true
	}

	fileName := fmt.Sprintf("%s-%d.itermcolors", options.Theme.Label, time.Now().Unix())
	filePath := filepath.Join(options.Directory, fileName)

	iTermTheme, err := convertTheme(options.Theme)
	if err != nil {
		return "", err
	}

	if options.ShouldWrite {
		err = os.WriteFile(filePath, []byte(iTermTheme), 0644)
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}

func convertTheme(selectedTheme theme.Theme) (string, error) {
	vscodeTheme, err := readTheme(selectedTheme.Path)
	if err != nil {
		log.Error("🚨 Failed to read theme file", "path", selectedTheme.Path, "error", err)
		return "", fmt.Errorf("error reading theme file: %v", err)
	}

	if vscodeTheme.Colors == nil {
		log.Error("🚨 Invalid theme format", "path", selectedTheme.Path)
		return "", fmt.Errorf("invalid theme format: colors not found or not a map")
	}

	themeType := vscodeTheme.Type
	if themeType == "" {
		themeType, err = theme.GetThemeType()
		if err != nil {
			return "", err
		}
		if themeType == "" {
			return "", fmt.Errorf("🚨 theme type selection cancelled")
		}
	}

	var itermColors []map[string]interface{}

	for name := range constants.AnsiColorFromVSCode {
		colorHex := getItermColor(themeType, name, vscodeTheme.Colors)
		colorRGBA, err := utils.HexToRGBA(colorHex)
		if err != nil {
			return "", err
		}

		itermColors = append(itermColors, map[string]interface{}{
			"key":   name,
			"red":   colorRGBA.Red,
			"green": colorRGBA.Green,
			"blue":  colorRGBA.Blue,
			"alpha": colorRGBA.Alpha,
		})
	}
	return getThemeXML(itermColors), nil
}

func readTheme(path string) (vscodeTheme, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return vscodeTheme{}, fmt.Errorf("error reading file: %v", err)
	}

	cleanContents := utils.RemoveCommentsAndTrailingCommas(contents)

	var themeData vscodeTheme
	err = json.Unmarshal(cleanContents, &themeData)
	if err != nil {
		return vscodeTheme{}, fmt.Errorf("error parsing theme JSON: %v", err)
	}

	return themeData, nil
}

func getItermColor(themeType string, name string, vscodeTheme map[string]interface{}) string {
	possibleKeys := constants.AnsiColorFromVSCode[name]
	for _, color := range possibleKeys {
		if val, ok := vscodeTheme[color]; ok {
			if strVal, ok := val.(string); ok {
				return strVal
			}
		}
	}

	fallback := constants.DefaultFallbackColors[themeType][name]
	userMessage := fmt.Sprintf("🔧 Color '%s' is missing for this %s theme, using default fallback color.", name, themeType)
	log.Info(userMessage)

	// Structured logging for debugging
	log.Debug("Using fallback color",
		"colorName", name,
		"themeType", themeType,
		"fallback", fallback)

	return fallback
}

func getThemeXML(colors []map[string]interface{}) string {
	var buffer bytes.Buffer

	buffer.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
`)

	for _, color := range colors {
		buffer.WriteString(getItermColorComponent(color))
	}

	buffer.WriteString(`</dict>
</plist>
`)

	return buffer.String()
}

func getItermColorComponent(color map[string]interface{}) string {
	return fmt.Sprintf(`  <key>%s</key>
  <dict>
    <key>Alpha Component</key>
    <real>%f</real>
    <key>Blue Component</key>
    <real>%f</real>
    <key>Color Space</key>
    <string>sRGB</string>
    <key>Green Component</key>
    <real>%f</real>
    <key>Red Component</key>
    <real>%f</real>
  </dict>
`, color["key"], color["alpha"], color["blue"], color["green"], color["red"])
}
