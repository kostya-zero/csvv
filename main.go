package main

import (
	"fmt"
	"os"
)

func main() {
	cmd := BuildCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
