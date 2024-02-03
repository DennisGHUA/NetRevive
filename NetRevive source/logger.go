package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
)

// LogInfo logs a message with an [INFO] prefix
func LogInfo(msg string) {
	log.Printf("[INFO] %s\n", msg)
}

// LogWarning logs a message with a [WARN] prefix
func LogWarning(msg string) {
	log.Printf("[WARN] %s\n", msg)
}

// LogError logs a message with an [ERRO] prefix, along with the error message and stack trace if an error is present
func LogError(msg string, err error) {
	log.Printf("[ERRO] %s\n", msg)
	if err != nil {
		log.Printf("[ERRO] %+v", errors.WithStack(err))
	}
}

// LogFatal logs a message with an [FATL] prefix and exits the program with a non-zero status code, along with the error message and stack trace if an error is present
func LogFatal(msg string, err error) {
	log.Println(fmt.Sprintf("[FATL] %v", msg))
	if err != nil {
		log.Printf("[FATL] %+v", errors.WithStack(err))
	}
	panic(2)
}
