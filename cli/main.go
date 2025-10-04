package main

import (
	"fmt"
	"os"
)

var (
	// Version information set during build
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("wild-cli version %s\n", Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Build time: %s\n", BuildTime)
		return
	}

	fmt.Println("Wild Cloud CLI")
	fmt.Println("Usage: wild-cli [command]")
	fmt.Println("")
	fmt.Println("Available commands:")
	fmt.Println("  version    Show version information")
	fmt.Println("")
	fmt.Println("More commands coming soon...")
}
