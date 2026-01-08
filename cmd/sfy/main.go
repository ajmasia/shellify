// Package main is the entry point for the sfy CLI.
package main

import (
	"os"

	"github.com/ajmasia/shellify/internal/interfaces/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
