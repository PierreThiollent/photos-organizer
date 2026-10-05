package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	executablePath, err := os.Executable()

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get executable path: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🏞️ Photos/videos organizer")
	fmt.Println("Processing files...")
	fmt.Println()

	startedAt := time.Now()

	workingDir := filepath.Dir(executablePath)
	entries, err := os.ReadDir(workingDir)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read directory %s: %v\n", workingDir, err)
		os.Exit(1)
	}

	movedCount := 0
	errorCount := 0
	folderCount := 0

	for _, entry := range entries {
		fileName := entry.Name()
		path := filepath.Join(workingDir, fileName)

		if path == executablePath ||
			strings.HasPrefix(fileName, ".") ||
			filepath.Ext(fileName) == "" {
			continue
		}

		fileInfo, err := os.Stat(path)
		if err != nil {
			errorCount++
			fmt.Fprintf(os.Stderr, "failed to get metadata for %s: %v\n", path, err)
			continue
		}

		if !fileInfo.Mode().IsRegular() {
			continue
		}

		modificationDate := fileInfo.ModTime().Format("2006-01-02")

		newDir := filepath.Join(workingDir, modificationDate)

		if _, err := os.Stat(newDir); err != nil {
			if !os.IsNotExist(err) {
				errorCount++
				fmt.Fprintf(os.Stderr, "failed to check directory %s: %v\n", newDir, err)
				continue
			}

			if err := os.Mkdir(newDir, 0o755); err != nil {
				errorCount++
				fmt.Fprintf(os.Stderr, "failed to create directory %s: %v\n", newDir, err)
				continue
			}

			folderCount++
		}

		newPath := filepath.Join(newDir, fileName)

		if err := os.Rename(path, newPath); err != nil {
			errorCount++
			fmt.Fprintf(os.Stderr, "failed to move file %s to %s: %v\n", path, newPath, err)
			continue
		}

		movedCount++
	}

	elapsed := time.Since(startedAt)
	elapsedSeconds := elapsed.Seconds()

	filesPerSecond := 0.0
	if elapsedSeconds > 0 {
		filesPerSecond = float64(movedCount) / elapsedSeconds
	}

	status := "Success"
	if errorCount > 0 {
		status = "Completed with errors"
	}

	duration := fmt.Sprintf("%.3f s", elapsedSeconds)
	throughput := fmt.Sprintf("%.0f files/s", filesPerSecond)

	fmt.Println("┌──────────────────────────────────────────┐")
	fmt.Printf("│ %-40s │\n", "Processing summary")
	fmt.Println("├──────────────────────────────────────────┤")
	fmt.Printf("│ %-15s : %22s │\n", "Status", status)
	fmt.Printf("│ %-15s : %22d │\n", "Files moved", movedCount)
	fmt.Printf("│ %-15s : %22d │\n", "Folders created", folderCount)
	fmt.Printf("│ %-15s : %22d │\n", "Errors", errorCount)
	fmt.Printf("│ %-15s : %22s │\n", "Duration", duration)
	fmt.Printf("│ %-15s : %22s │\n", "Throughput", throughput)
	fmt.Println("└──────────────────────────────────────────┘")
}
