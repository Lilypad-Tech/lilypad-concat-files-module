package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	inputDir := "/inputs"
	outputDir := "/outputs"
	requiredFiles := []string{"one.txt", "two.txt"}
	optionalFiles := []string{"three.txt", "four.txt"}
	outputFile := filepath.Join(outputDir, "output.txt")

	// The separator is passed as a normal input to the job
	separator := os.Getenv("SEPARATOR")

	// Trim quotes that may come with input
	separator = strings.Trim(separator, "\"")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
		return
	}

	filesToProcess, err := buildFileList(inputDir, requiredFiles, optionalFiles)
	if err != nil {
		logError(outputDir, fmt.Sprintf("Error building file list: %v", err))
		return
	}

	if err := concatenateFiles(filesToProcess, separator, outputDir, outputFile); err != nil {
		logError(outputDir, fmt.Sprintf("Error concatenating files: %v", err))
		return
	}
}

func buildFileList(inputDir string, requiredFiles, optionalFiles []string) ([]string, error) {
	filesToProcess := make([]string, len(requiredFiles))

	// Add required files
	for i, filename := range requiredFiles {
		filePath := filepath.Join(inputDir, filename)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return nil, fmt.Errorf("required file not found: %s", filePath)
		} else if err != nil {
			return nil, fmt.Errorf("error checking file %s: %v", filePath, err)
		}

		filesToProcess[i] = filepath.Join(inputDir, filename)
	}

	// Add optional files if they exist
	for _, filename := range optionalFiles {
		filePath := filepath.Join(inputDir, filename)
		if _, err := os.Stat(filePath); err == nil {
			filesToProcess = append(filesToProcess, filePath)
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error checking optional file %s: %v", filePath, err)
		}
	}

	return filesToProcess, nil
}

// Concatenate file contents with optional separator
func concatenateFiles(filesToProcess []string, separator string, outputDir string, outputFile string) error {
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	for i, file := range filesToProcess {
		content, err := os.ReadFile(file)
		if err != nil {
			logError(outputDir, fmt.Sprintf("Error reading file %s: %v", file, err))
			continue
		}

		isLastFile := i == len(filesToProcess)-1
		if err := writeWithSeparator(outFile, content, separator, !isLastFile); err != nil {
			logError(outputDir, fmt.Sprintf("Error processing file %s: %v", file, err))
			continue
		}
	}

	return nil
}

// Write file content with optional separator
func writeWithSeparator(outFile *os.File, content []byte, separator string, addSeparator bool) error {
	_, err := outFile.Write(content)
	if err != nil {
		return fmt.Errorf("error writing file content: %v", err)
	}

	if addSeparator && separator != "" {
		_, err = outFile.WriteString(separator)
		if err != nil {
			return fmt.Errorf("error writing separator: %v", err)
		}

		// Add newline if separator doesn't end with one
		if !strings.HasSuffix(separator, "\n") {
			_, err = outFile.WriteString("\n")
			if err != nil {
				return fmt.Errorf("error writing newline: %v", err)
			}
		}
	}

	return nil
}

func logError(outputDir, message string) {
	errorFile := filepath.Join(outputDir, "error.txt")
	f, err := os.OpenFile(errorFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open error log file: %v\n", err)
		fmt.Fprintf(os.Stderr, "Original error: %s\n", message)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)

	if _, err := f.WriteString(logMessage); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to error log: %v\n", err)
		fmt.Fprintf(os.Stderr, "Original error: %s\n", message)
	}
}
