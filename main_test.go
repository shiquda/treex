package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Setup temporary directory structure for testing
func setupTestDir(t *testing.T) string {
	tempDir := filepath.Join(os.TempDir(), "treex_test")
	os.RemoveAll(tempDir) // Clean up previous test directory

	// Create directory structure
	dirs := []string{
		filepath.Join(tempDir, "dir1"),
		filepath.Join(tempDir, "dir2", "subdir"),
		filepath.Join(tempDir, ".hidden_dir"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	// Create some files
	files := map[string]string{
		filepath.Join(tempDir, "file1.txt"):         "test content",
		filepath.Join(tempDir, "dir1", "file2.go"):  "go code",
		filepath.Join(tempDir, ".hidden_file"):      "hidden",
		filepath.Join(tempDir, "dir2", "config.md"): "# Config",
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Create files for gitignore test
	ignoreContent := "*.log\nbuild/\n"
	if err := os.WriteFile(filepath.Join(tempDir, ".gitignore"), []byte(ignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore file: %v", err)
	}

	// Create files that should be excluded
	os.WriteFile(filepath.Join(tempDir, "log.log"), []byte("log"), 0644)
	os.MkdirAll(filepath.Join(tempDir, "build"), 0755)
	os.WriteFile(filepath.Join(tempDir, "build", "app.js"), []byte("js"), 0644)

	return tempDir
}

// Integration test for different output types
func TestIntegration(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	// Ensure working directory is correct
	oldWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(oldWd)

	// Test scenarios
	testCases := []struct {
		name       string
		filter     *Filter
		hideHidden bool
		dirsOnly   bool
		maxDepth   int
		checkFunc  func(node *TreeNode) bool
	}{
		{
			name:       "Basic Tree Generation",
			filter:     NewFilter("", false),
			hideHidden: false,
			dirsOnly:   false,
			maxDepth:   0,
			checkFunc: func(node *TreeNode) bool {
				// Check if all files and directories are included
				allNodes := node.getAllNodes()
				return len(allNodes) >= 10 // root + at least 9 child nodes
			},
		},
		{
			name:       "Hidden Files Filter",
			filter:     NewFilter("", false),
			hideHidden: true,
			dirsOnly:   false,
			maxDepth:   0,
			checkFunc: func(node *TreeNode) bool {
				// Check if hidden files are filtered
				for _, child := range node.Children {
					if child.Name == ".hidden_dir" || child.Name == ".hidden_file" {
						return false
					}
				}
				return true
			},
		},
		{
			name:       "Dirs Only",
			filter:     NewFilter("", false),
			hideHidden: false,
			dirsOnly:   true,
			maxDepth:   0,
			checkFunc: func(node *TreeNode) bool {
				// Check if only directories are included
				for _, child := range allChildrenRecursive(node) {
					if !child.IsDir {
						return false
					}
				}
				return true
			},
		},
		{
			name:       "Max Depth Limit",
			filter:     NewFilter("", false),
			hideHidden: false,
			dirsOnly:   false,
			maxDepth:   1,
			checkFunc: func(node *TreeNode) bool {
				// Check if depth is limited
				maxFoundDepth := 0
				for _, n := range allChildrenRecursive(node) {
					if n.Depth > maxFoundDepth {
						maxFoundDepth = n.Depth
					}
				}
				return maxFoundDepth <= 1 // depth should not exceed 1
			},
		},
		{
			name:       "Exclude Rules",
			filter:     NewFilter(".log, build/", false),
			hideHidden: false,
			dirsOnly:   false,
			maxDepth:   0,
			checkFunc: func(node *TreeNode) bool {
				// Check if exclude rules are effective
				for _, child := range node.Children {
					if child.Name == "log.log" || child.Name == "build" {
						return false
					}
				}
				return true
			},
		},
		{
			name:       "GitIgnore Integration",
			filter:     NewFilter("", true),
			hideHidden: false,
			dirsOnly:   false,
			maxDepth:   0,
			checkFunc: func(node *TreeNode) bool {
				// Check if .gitignore rules are effective
				for _, child := range node.Children {
					if child.Name == "log.log" || child.Name == "build" {
						return false
					}
				}
				return true
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node, err := getTreeNode(".", 1, testDir, tc.maxDepth, tc.filter, tc.hideHidden, tc.dirsOnly)
			if err != nil {
				t.Fatalf("%s: getTreeNode error: %v", tc.name, err)
			}

			if !tc.checkFunc(node) {
				t.Errorf("%s: Check failed", tc.name)
			}

			// Test various output formats
			_ = node.ToTreeString(true, "", false)
			_ = node.ToIndentString(2, false)
			_ = node.ToMarkdownString(0, false)
			_ = node.ToMermaidString()
		})
	}
}

// Get all child nodes recursively
func allChildrenRecursive(node *TreeNode) []*TreeNode {
	var result []*TreeNode
	for _, child := range node.Children {
		result = append(result, child)
		result = append(result, allChildrenRecursive(child)...)
	}
	return result
}
