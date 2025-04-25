package main

import (
	"os"
	"testing"
)

func TestNewFilter(t *testing.T) {
	// Test empty filter
	f := NewFilter("", false)
	if len(f.dirNames) != 0 || len(f.suffixes) != 0 || len(f.patterns) != 0 {
		t.Error("Empty filter should have no rules")
	}

	// Test various rules
	f = NewFilter("dir/, .txt, pattern*", false)
	if len(f.dirNames) != 1 || f.dirNames[0] != "dir" {
		t.Error("Filter should have one directory rule")
	}
	if len(f.suffixes) != 1 || f.suffixes[0] != ".txt" {
		t.Error("Filter should have one suffix rule")
	}
	if len(f.patterns) != 1 || f.patterns[0] != "pattern*" {
		t.Error("Filter should have one pattern rule")
	}
}

func TestShouldExclude(t *testing.T) {
	// Test directory filtering
	f := NewFilter("node_modules/, .git/", false)
	if !f.shouldExclude("node_modules", true, "node_modules") {
		t.Error("Should exclude node_modules directory")
	}
	if f.shouldExclude("src", true, "src") {
		t.Error("Should not exclude src directory")
	}

	// Test suffix filtering
	f = NewFilter(".txt, .log", false)
	if !f.shouldExclude("file.txt", false, "file.txt") {
		t.Error("Should exclude .txt files")
	}
	if !f.shouldExclude("file.log", false, "file.log") {
		t.Error("Should exclude .log files")
	}
	if f.shouldExclude("file.go", false, "file.go") {
		t.Error("Should not exclude .go files")
	}

	// Test pattern matching
	f = NewFilter("test*, *temp", false)
	if !f.shouldExclude("test_file", false, "test_file") {
		t.Error("Should exclude files starting with test")
	}
	if !f.shouldExclude("file_temp", false, "file_temp") {
		t.Error("Should exclude files ending with temp")
	}
	if f.shouldExclude("source.go", false, "source.go") {
		t.Error("Should not exclude normal go files")
	}
}

func TestGitIgnoreIntegration(t *testing.T) {
	// Create temporary .gitignore file for testing
	content := []byte("*.log\nbuild/\n# comment\ntemp*\n")
	err := os.WriteFile(".gitignore.test", content, 0644)
	if err != nil {
		t.Fatal("Failed to create test .gitignore file")
	}
	defer os.Remove(".gitignore.test")

	// Rename existing .gitignore file if present
	if _, err := os.Stat(".gitignore"); err == nil {
		os.Rename(".gitignore", ".gitignore.bak")
		defer os.Rename(".gitignore.bak", ".gitignore")
	}
	os.Rename(".gitignore.test", ".gitignore")
	defer os.Rename(".gitignore", ".gitignore.test")

	// Test gitignore integration
	f := NewFilter("", true)

	// Check if rules are read correctly
	if len(f.suffixes) == 0 || f.suffixes[0] != ".log" {
		t.Error("Should read .log suffix rule from .gitignore")
	}

	if len(f.dirNames) == 0 || f.dirNames[0] != "build" {
		t.Error("Should read build directory rule from .gitignore")
	}

	// Test matching
	if !f.shouldExclude("server.log", false, "server.log") {
		t.Error("Should exclude .log files")
	}
	if !f.shouldExclude("build", true, "build") {
		t.Error("Should exclude build directory")
	}
	if !f.shouldExclude("temp_file", false, "temp_file") {
		t.Error("Should exclude files starting with temp")
	}
}
