package main

import (
	"fmt"
	"os"

	"github.com/HanSoBored/agtop/internal/ui"
)

func main() {
	if err := ui.Start(); err != nil {
		fmt.Printf("Error running Android GPU Monitor: %v\n", err)
		os.Exit(1)
	}
}
