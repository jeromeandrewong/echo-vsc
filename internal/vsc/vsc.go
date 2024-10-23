package vsc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/jeromeandrewong/echo-vsc/internal/logger"
	"github.com/jeromeandrewong/echo-vsc/internal/theme"
)

// job to process a single extension
type ExtensionJob struct {
	Dir     string
	ExtInfo os.DirEntry
}

// result of processing an extension
type ExtensionResult struct {
	Themes []theme.Theme
	Err    error
}

func processExtensionWorker(
	jobs <-chan ExtensionJob,
	results chan<- ExtensionResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for job := range jobs {
		themes, err := getThemesFromExtension(job.Dir, job.ExtInfo)
		results <- ExtensionResult{
			Themes: themes,
			Err:    err,
		}
	}

}
func GetVSCThemes(vscDir string) ([]theme.Theme, error) {
	startTime := time.Now()

	extensions, err := os.ReadDir(vscDir)
	if err != nil {
		return nil, fmt.Errorf("error reading VSC directory: %v", err)
	}

	// count actual directories to process
	dirCount := 0
	for _, ext := range extensions {
		if ext.IsDir() {
			dirCount++
		}
	}

	if dirCount == 0 {
		return nil, fmt.Errorf("no extensions found in directory")
	}

	// buffered channels for jobs and results
	const numWorkers = 5
	jobs := make(chan ExtensionJob, dirCount)
	results := make(chan ExtensionResult, dirCount)

	// start worker pool
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go processExtensionWorker(jobs, results, &wg)
	}

	// send jobs to workers
	jobCount := 0
	for _, extension := range extensions {
		if !extension.IsDir() {
			continue
		}
		jobs <- ExtensionJob{
			Dir:     vscDir,
			ExtInfo: extension,
		}
		jobCount++
	}
	close(jobs) // all jobs are sent

	// start a goroutine to close results channel after all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// collect results with timeout
	var allThemes []theme.Theme
	processedCount := 0
	const timeout = 30 * time.Second
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case result, ok := <-results:
			if !ok {
				return allThemes, nil
			}
			processedCount++

			if result.Err != nil {
				log.Error("Warning: %v", result.Err)
				continue
			}

			allThemes = append(allThemes, result.Themes...)

			// reset timer for next result
			timer.Reset(timeout)

		case <-timer.C:
			return allThemes, fmt.Errorf("timeout waiting for theme processing after %v", timeout)
		}

		// Check if we've processed all extensions
		if processedCount == jobCount {
			processingTime := time.Since(startTime)
			fmt.Printf("GetVSCThemes processing time: %v\n", processingTime)
			fmt.Printf("Total themes found: %d\n", len(allThemes))
			return allThemes, nil
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
			Themes []theme.Theme `json:"themes"`
		} `json:"contributes"`
	}

	if err := json.Unmarshal(packageJSON, &packageData); err != nil {
		return nil, fmt.Errorf("error parsing package.json: %v", err)
	}

	var themes []theme.Theme

	for _, t := range packageData.Contributes.Themes {
		themePath := filepath.Join(extensionPath, t.Path)

		themes = append(themes, theme.Theme{
			Label: t.Label,
			Path:  themePath,
		})
	}

	return themes, nil
}
