package main

import (
	"fmt"
)

func main() {
	cmd := BuildCmd()
	if err := cmd.Execute(); err != nil {
		PrintFatal(fmt.Sprintf("cli failed: %v", err))
	}
}
