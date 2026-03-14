package main

import (
	"fmt"
	"os"
)

// PrintError prints formatted error message to stdout.
func PrintError(msg string) {
	fmt.Printf("\x1b[91m\x1b[1merror\x1b[0m: %s\n", msg)
}

// PrintFatal prints formatted error message to stdout and exit with error code "1".
func PrintFatal(msg string) {
	fmt.Printf("\x1b[91m\x1b[1merror\x1b[0m: %s\n", msg)
	os.Exit(1)
}
