package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("ASDP Centralized Validation Tool")
	fmt.Println("================================")

	// In a real tool, we might implement a logic here without relying on 'go test'.
	// But since we want to run the suite exactly as verified:

	wd, _ := os.Getwd()
	fmt.Printf("Running from: %s\n", wd)

	// Ensure we are in the right place or find the root
	// For simplicity, we assume this tool is run via 'go run tools/validate/main.go' from root
	// or built and installed.

	cmd := exec.Command("go", "test", "-v", "./tools/validate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n[FAIL] Validation suite failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n[SUCCESS] All functional and compliance checks passed.")
}
