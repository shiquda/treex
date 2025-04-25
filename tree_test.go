package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetRelativePath(t *testing.T) {
	testCases := []struct {
		absolute string
		root     string
		expected string
	}{
		{"/home/user/project/file.txt", "/home/user/", "project/file.txt"},
		{"/home/user/project/", "/home/user/", "project"},
		{"/home/user/project", "/home/user", "project"},
		{"/other/path", "/home/user/", "/other/path"},
		{"project/file.txt", "/home/user/", "project/file.txt"},
	}

	for _, tc := range testCases {
		result := getRelativePath(tc.absolute, tc.root)
		if result != tc.expected {
			t.Errorf("For %s and %s, expected %s, but got %s", tc.absolute, tc.root, tc.expected, result)
		}
	}
}

func TestTreeNode(t *testing.T) {
	// Create a test tree
	root := &TreeNode{
		Name:  "root",
		IsDir: true,
		Depth: 0,
	}

	child1 := &TreeNode{
		Name:  "child1",
		IsDir: true,
		Depth: 1,
	}

	child2 := &TreeNode{
		Name:  "child2",
		IsDir: false,
		Depth: 1,
	}

	grandchild1 := &TreeNode{
		Name:  "grandchild1",
		IsDir: false,
		Depth: 2,
	}

	// Build tree structure
	root.Children = append(root.Children, child1, child2)
	child1.Children = append(child1.Children, grandchild1)

	// Test node count
	nodes := root.getAllNodes()
	if len(nodes) != 4 {
		t.Errorf("Expected 4 nodes, but got %d", len(nodes))
	}

	// Test specific node presence in result
	found := false
	for _, node := range nodes {
		if node.Name == "grandchild1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Did not find grandchild1 in all nodes")
	}
}

func TestGetTreeNode(t *testing.T) {
	// Create test directory structure
	testDir := "test_dir_structure"
	defer os.RemoveAll(testDir)

	// Create directory structure
	os.MkdirAll(filepath.Join(testDir, "dir1", "subdir"), 0755)
	os.MkdirAll(filepath.Join(testDir, "dir2"), 0755)
	os.MkdirAll(filepath.Join(testDir, ".hidden_dir"), 0755)

	// Create some files
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, "dir1", "file2.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(testDir, ".hidden_file"), []byte("test"), 0644)

	// Get current working directory as basePath
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("Failed to get current working directory")
	}

	// Test basic tree generation
	filter := NewFilter("", false)
	node, err := getTreeNode(testDir, 1, cwd, 0, filter, false, false)
	if err != nil {
		t.Fatalf("getTreeNode error: %v", err)
	}

	// Check root node
	if node.Name != testDir || !node.IsDir {
		t.Errorf("Root node error: name=%s, isDir=%v", node.Name, node.IsDir)
	}

	// Check child node count (should include 5: dir1, dir2, .hidden_dir, file1.txt, .hidden_file)
	if len(node.Children) != 5 {
		t.Errorf("Expected 5 child nodes, but got %d", len(node.Children))
	}

	// Test hidden file filtering
	node, err = getTreeNode(testDir, 1, cwd, 0, filter, true, false)
	if err != nil {
		t.Fatalf("getTreeNode error: %v", err)
	}

	// Check child node count (should exclude .hidden_dir and .hidden_file)
	if len(node.Children) != 3 {
		t.Errorf("Expected 3 child nodes after filtering hidden files, but got %d", len(node.Children))
	}

	// Test maximum depth
	node, err = getTreeNode(testDir, 1, cwd, 1, filter, false, false)
	if err != nil {
		t.Fatalf("getTreeNode error: %v", err)
	}

	// Find dir1 node
	var dir1Node *TreeNode
	for _, child := range node.Children {
		if child.Name == "dir1" && child.IsDir {
			dir1Node = child
			break
		}
	}

	if dir1Node == nil {
		t.Fatal("Did not find dir1 node")
	}

	// Check if dir1 node has no children (due to depth limit)
	if len(dir1Node.Children) != 0 {
		t.Errorf("With depth limit 1, dir1 should have no children, but got %d", len(dir1Node.Children))
	}

	// Test only directories
	node, err = getTreeNode(testDir, 1, cwd, 0, filter, false, true)
	if err != nil {
		t.Fatalf("getTreeNode error: %v", err)
	}

	// Check if only directories are included
	for _, child := range node.Children {
		if !child.IsDir {
			t.Errorf("In dirs-only mode, node %s should not be a file", child.Name)
		}
	}
}
