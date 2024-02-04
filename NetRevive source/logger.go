package main

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

const (
	logFilePath       = "NetRevive.log"
	logTruncateLength = 1000
	timeFormat        = "2006-01-02 15:04:05" // Timestamp format
)

// LogInfo logs a message with an [INFO] prefix
func LogInfo(msg string) {
	log.Println(fmt.Sprintf("[INFO] %s", msg))
}

// LogWarning logs a message with a [WARN] prefix
func LogWarning(msg string) {
	logMessage := fmt.Sprintf("[WARN] %s", msg)
	log.Println(logMessage)
	logToFile(logFilePath, logMessage)
}

// LogError logs a message with an [ERRO] prefix, along with the error message and stack trace if an error is present
func LogError(msg string, err error) {
	logMessage := fmt.Sprintf("[ERRO] %s", msg)
	log.Println(logMessage)
	if err != nil {
		logMessage += fmt.Sprintf("\n%+v", errors.WithStack(err))
	}
	logToFile(logFilePath, logMessage)
}

// LogFatal logs a message with an [FATL] prefix and exits the program with a non-zero status code, along with the error message and stack trace if an error is present
func LogFatal(msg string, err error) {
	logMessage := fmt.Sprintf("[FATL] %s", msg)
	log.Println(logMessage)
	if err != nil {
		logMessage += fmt.Sprintf("\n%+v", errors.WithStack(err))
	}
	logToFile(logFilePath, logMessage)
	panic(2)
}

func logToFile(filePath, message string) error {
	// Generate timestamp
	timestamp := time.Now().Format(timeFormat)

	// Format log message with timestamp
	logMessage := fmt.Sprintf("[%s] %s", timestamp, message)

	// Open the file in append mode
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	// Write the formatted log message to the file
	if _, err := file.WriteString(logMessage + "\n"); err != nil {
		return err
	}

	// Close the file after writing
	file.Close()

	// Re-open the file in read mode
	file, err = os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read existing contents and keep only the last logTruncateLength lines
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Close the file after reading
	file.Close()

	// Re-open the file in write mode
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write filtered lines back to the file
	startIndex := len(lines) - logTruncateLength
	if startIndex < 0 {
		startIndex = 0
	}
	for _, line := range lines[startIndex:] {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}
