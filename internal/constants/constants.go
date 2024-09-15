package constants

import (
	"os"
	"path/filepath"
)

var (
	HomeDir, _    = os.UserHomeDir()
	ExtensionsDir = filepath.Join(HomeDir, ".vscode", "extensions")
)

var AnsiColorFromVSCode = map[string][]string{
	"Ansi 0 Color":        {"terminal.ansiBlack"},
	"Ansi 1 Color":        {"terminal.ansiRed"},
	"Ansi 10 Color":       {"terminal.ansiBrightGreen"},
	"Ansi 11 Color":       {"terminal.ansiBrightYellow"},
	"Ansi 12 Color":       {"terminal.ansiBrightBlue"},
	"Ansi 13 Color":       {"terminal.ansiBrightMagenta"},
	"Ansi 14 Color":       {"terminal.ansiBrightCyan"},
	"Ansi 15 Color":       {"terminal.ansiBrightWhite"},
	"Ansi 2 Color":        {"terminal.ansiGreen"},
	"Ansi 3 Color":        {"terminal.ansiYellow"},
	"Ansi 4 Color":        {"terminal.ansiBlue"},
	"Ansi 5 Color":        {"terminal.ansiMagenta"},
	"Ansi 6 Color":        {"terminal.ansiCyan"},
	"Ansi 7 Color":        {"terminal.ansiWhite"},
	"Ansi 8 Color":        {"terminal.ansiBrightBlack"},
	"Ansi 9 Color":        {"terminal.ansiBrightRed"},
	"Background Color":    {"terminal.background", "editor.background"},
	"Bold Color":          {},
	"Cursor Color":        {"terminalCursor.foreground", "editorCursor.foreground"},
	"Cursor Text Color":   {"terminalCursor.foreground", "editorCursor.foreground"},
	"Foreground Color":    {"terminal.foreground", "editor.foreground"},
	"Selected Text Color": {"terminal.background", "editor.background"},
	"Selection Color":     {"terminal.selectionBackground", "terminal.foreground", "editor.foreground"},
	"Link Color":          {"textLink.foreground"},
}

var DefaultFallbackColors = map[string]map[string]string{
	"dark": {
		"Ansi 0 Color":        "#21222c",
		"Ansi 1 Color":        "#ff5555",
		"Ansi 2 Color":        "#50fa7b",
		"Ansi 3 Color":        "#f1fa8c",
		"Ansi 4 Color":        "#bd93f9",
		"Ansi 5 Color":        "#ff79c6",
		"Ansi 6 Color":        "#8be9fd",
		"Ansi 7 Color":        "#f8f8f2",
		"Ansi 8 Color":        "#6272a4",
		"Ansi 9 Color":        "#ff6e6e",
		"Ansi 10 Color":       "#69ff94",
		"Ansi 11 Color":       "#ffffa5",
		"Ansi 12 Color":       "#d6acff",
		"Ansi 13 Color":       "#ff92df",
		"Ansi 14 Color":       "#a4ffff",
		"Ansi 15 Color":       "#ffffff",
		"Background Color":    "#282a36",
		"Bold Color":          "#ffffff",
		"Cursor Color":        "#f8f8f2",
		"Cursor Guide Color":  "#b3ecff30",
		"Cursor Text Color":   "#282a36",
		"Foreground Color":    "#f8f8f2",
		"Link Color":          "#8be9fd",
		"Selected Text Color": "#ffffff",
		"Selection Color":     "#44475a",
		"Tab Color":           "#bf1a12ff",
	},
	"light": {
		"Ansi 0 Color":        "#21222c",
		"Ansi 1 Color":        "#e64747",
		"Ansi 2 Color":        "#50fa7b",
		"Ansi 3 Color":        "#e7c547",
		"Ansi 4 Color":        "#7aa2f7",
		"Ansi 5 Color":        "#ad8ee6",
		"Ansi 6 Color":        "#449dab",
		"Ansi 7 Color":        "#787c99",
		"Ansi 8 Color":        "#444b6a",
		"Ansi 9 Color":        "#ff7a85",
		"Ansi 10 Color":       "#b9f27c",
		"Ansi 11 Color":       "#ff9e64",
		"Ansi 12 Color":       "#7da6ff",
		"Ansi 13 Color":       "#bb9af7",
		"Ansi 14 Color":       "#0db9d7",
		"Ansi 15 Color":       "#acb0d0",
		"Background Color":    "#fafafa",
		"Bold Color":          "#444b6a",
		"Cursor Color":        "#444b6a",
		"Cursor Guide Color":  "#b3ecff30",
		"Cursor Text Color":   "#fafafa",
		"Foreground Color":    "#444b6a",
		"Link Color":          "#005fb8",
		"Selected Text Color": "#fafafa",
		"Selection Color":     "#b4d8fd",
		"Tab Color":           "#bf1a12ff",
	},
}
