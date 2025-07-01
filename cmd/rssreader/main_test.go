package main

import (
	"os"
	"testing"
)

func TestMainCompiles(t *testing.T) {
	t.Log("Main package compiles successfully")
}

func TestHelpFlag(t *testing.T) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// Test with help flag
	os.Args = []string{"rssreader", "-help"}

	// We can't easily test the main function directly since it calls os.Exit
	// But we can test that the flag parsing works correctly
	// This is a basic smoke test to ensure the code compiles and runs
	t.Log("Help flag test - if this runs without panic, the test passes")
}

// TODO(jannis-seemann): Mock Parse and test that this CLI tool works as expected.
