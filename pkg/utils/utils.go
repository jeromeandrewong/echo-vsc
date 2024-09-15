package utils

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type RGBA struct {
	Red   float64
	Green float64
	Blue  float64
	Alpha float64
}

func GetDownloadsFolder() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, "Downloads"), nil
}

func HexToRGBA(hex string) (RGBA, error) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}

	values, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return RGBA{}, err
	}

	return RGBA{
		Red:   float64((values>>16)&255) / 255.0,
		Green: float64((values>>8)&255) / 255.0,
		Blue:  float64(values&255) / 255.0,
		Alpha: 1.0,
	}, nil
}

func RemoveCommentsAndTrailingCommas(data []byte) []byte {
	var result []byte
	inString := false
	inLineComment := false
	inBlockComment := false
	lastNonWhitespace := byte(0)

	for i := 0; i < len(data); i++ {
		if inLineComment {
			if data[i] == '\n' {
				inLineComment = false
			}
			continue
		}
		if inBlockComment {
			if i < len(data)-1 && data[i] == '*' && data[i+1] == '/' {
				inBlockComment = false
				i++
			}
			continue
		}
		if inString {
			if data[i] == '"' && data[i-1] != '\\' {
				inString = false
			}
			result = append(result, data[i])
			lastNonWhitespace = data[i]
			continue
		}

		switch {
		case data[i] == '"':
			inString = true
			result = append(result, data[i])
			lastNonWhitespace = data[i]
		case i < len(data)-1 && data[i] == '/' && data[i+1] == '/':
			inLineComment = true
			i++
		case i < len(data)-1 && data[i] == '/' && data[i+1] == '*':
			inBlockComment = true
			i++
		case data[i] == ',' && i < len(data)-1:
			// check if the next non-whitespace character is a closing bracket or brace
			for j := i + 1; j < len(data); j++ {
				if data[j] == ' ' || data[j] == '\t' || data[j] == '\n' || data[j] == '\r' {
					continue
				}
				if data[j] == '}' || data[j] == ']' {
					// skip this comma
					i = j - 1
					goto nextIteration
				}
				break
			}
			result = append(result, data[i])
			lastNonWhitespace = data[i]
		case data[i] == ' ' || data[i] == '\t' || data[i] == '\n' || data[i] == '\r':
			// only append whitespace if the last character wasn't a comma
			if lastNonWhitespace != ',' {
				result = append(result, data[i])
			}
		default:
			result = append(result, data[i])
			lastNonWhitespace = data[i]
		}
	nextIteration:
	}

	return result
}
