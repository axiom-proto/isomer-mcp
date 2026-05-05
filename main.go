package main

import (
	"fmt"
	"isomer/internal/cmd"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}
}
