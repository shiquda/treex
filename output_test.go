package main

import (
	"strings"
	"testing"
)

func createTestTree() *TreeNode {
	// Create a simple test tree structure
	root := &TreeNode{
		Name:  "root",
		IsDir: true,
		Depth: 0,
	}

	dir1 := &TreeNode{
		Name:  "dir1",
		IsDir: true,
		Depth: 1,
	}

	file1 := &TreeNode{
		Name:  "file1.txt",
		IsDir: false,
		Depth: 1,
	}

	file2 := &TreeNode{
		Name:  "file2.go",
		IsDir: false,
		Depth: 2,
	}

	// Build tree structure
	root.Children = append(root.Children, dir1, file1)
	dir1.Children = append(dir1.Children, file2)

	return root
}

func TestToIndentString(t *testing.T) {
	// Create a simple test tree structure
	tree := createTestTree()

	// Test without icons
	result := tree.ToIndentString(2, false) // No icons

	// Check if output contains expected content
	expectedLines := []string{
		"root/",
		"  dir1/",
		"    file2.go",
		"  file1.txt",
	}

	for _, line := range expectedLines {
		if !strings.Contains(result, line) {
			t.Errorf("Indented output missing expected line: %s", line)
		}
	}

	// Check indentation
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if !strings.HasPrefix(lines[1], "  ") {
		t.Error("dir1 line should have 2 spaces indentation")
	}
	if !strings.HasPrefix(lines[2], "    ") {
		t.Error("file2.go line should have 4 spaces indentation")
	}

	// Check if directories have trailing slash
	if !strings.HasSuffix(lines[0], "/") {
		t.Error("Root directory should have a trailing slash")
	}

	if !strings.HasSuffix(lines[1], "/") {
		t.Error("Directory 'dir1' should have a trailing slash")
	}

	if strings.HasSuffix(lines[2], "/") {
		t.Error("File 'file2.go' should not have a trailing slash")
	}

	if strings.HasSuffix(lines[3], "/") {
		t.Error("File 'file1.txt' should not have a trailing slash")
	}

	//
	iconResult := tree.ToIndentString(2, true)
	if !strings.Contains(iconResult, "📁") {
		t.Error("Icon mode should display folder icon for directories")
	}
}

func TestToTreeString(t *testing.T) {
	tree := createTestTree()
	result := tree.ToTreeString(true, "", false) // No icons

	// Check if output contains expected content and tree symbols
	expectedPatterns := []string{
		"root/",
		"├── dir1/",
		"│   └── file2.go",
		"└── file1.txt",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(result, pattern) {
			t.Errorf("Tree output missing expected pattern: %s", pattern)
		}
	}

	// Ensure empty depth (root node) doesn't show prefix symbols
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if strings.HasPrefix(lines[0], "└── ") {
		t.Error("Root node should not display prefix symbol")
	}

	// Check if directories have trailing slash
	for i, line := range lines {
		if i == 0 && !strings.HasSuffix(line, "/") {
			t.Error("Root directory should have a trailing slash")
		}

		if strings.Contains(line, "dir1") && !strings.HasSuffix(line, "/") {
			t.Error("Directory 'dir1' should have a trailing slash")
		}

		if strings.Contains(line, "file1.txt") && strings.HasSuffix(line, "/") {
			t.Error("File 'file1.txt' should not have a trailing slash")
		}

		if strings.Contains(line, "file2.go") && strings.HasSuffix(line, "/") {
			t.Error("File 'file2.go' should not have a trailing slash")
		}
	}

	// Test icon mode
	iconResult := tree.ToTreeString(true, "", true)
	if !strings.Contains(iconResult, "📁") {
		t.Error("Icon mode should display folder icon for directories")
	}
	if !strings.Contains(iconResult, "🔹") {
		t.Error("Icon mode should display Go file icon for .go files")
	}
}

func TestToMarkdownString(t *testing.T) {
	tree := createTestTree()
	result := tree.ToMarkdownString(0, false)

	// Check if output contains expected markdown list format
	expectedPatterns := []string{
		"- root/",
		"  - dir1/",
		"    - file2.go",
		"  - file1.txt",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(result, pattern) {
			t.Errorf("Markdown output missing expected pattern: %s", pattern)
		}
	}

	// Ensure directories have trailing slash
	if !strings.Contains(result, "root/") || !strings.Contains(result, "dir1/") {
		t.Error("Directory names should have trailing slashes")
	}

	// Ensure files don't have trailing slash
	if strings.Contains(result, "file1.txt/") || strings.Contains(result, "file2.go/") {
		t.Error("File names should not have trailing slashes")
	}
}

func TestToMermaidString(t *testing.T) {
	tree := createTestTree()
	result := tree.ToMermaidString()

	// Check if mermaid output format is correct
	expectedPatterns := []string{
		"graph TD",
		"N1[root/]",
		"N2[dir1/]",
		"N1 --> N2",
		"N3[file2.go]",
		"N2 --> N3",
		"N4[file1.txt]",
		"N1 --> N4",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(result, pattern) {
			t.Errorf("Mermaid output missing expected pattern: %s\n%s", pattern, result)
		}
	}

	// Ensure node IDs are incremental
	if !strings.Contains(result, "N1") || !strings.Contains(result, "N2") ||
		!strings.Contains(result, "N3") || !strings.Contains(result, "N4") {
		t.Error("Mermaid output should contain incremental node IDs")
	}
}
