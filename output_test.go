package main

import (
	"strings"
	"testing"
)

func createTestTree() *TreeNode {
	// 创建一个简单的测试树结构
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

	// 构建树结构
	root.Children = append(root.Children, dir1, file1)
	dir1.Children = append(dir1.Children, file2)

	return root
}

func TestToIndentString(t *testing.T) {
	tree := createTestTree()
	result := tree.ToIndentString(2)

	// 检查输出是否包含预期的内容
	expectedLines := []string{
		"root",
		"  dir1",
		"    file2.go",
		"  file1.txt",
	}

	for _, line := range expectedLines {
		if !strings.Contains(result, line) {
			t.Errorf("Indented output missing expected line: %s", line)
		}
	}

	// 检查缩进是否正确
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if !strings.HasPrefix(lines[1], "  ") {
		t.Error("dir1 line should have 2 spaces indentation")
	}
	if !strings.HasPrefix(lines[2], "    ") {
		t.Error("file2.go line should have 4 spaces indentation")
	}
}

func TestToTreeString(t *testing.T) {
	tree := createTestTree()
	result := tree.ToTreeString(true, "")

	// 检查输出是否包含预期的内容和树形符号
	expectedPatterns := []string{
		"root",
		"├── dir1",
		"│   └── file2.go",
		"└── file1.txt",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(result, pattern) {
			t.Errorf("Tree output missing expected pattern: %s", pattern)
		}
	}

	// 确保空深度（根节点）不显示前缀符号
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if strings.HasPrefix(lines[0], "└── ") {
		t.Error("Root node should not display prefix symbol")
	}
}

func TestToMarkdownString(t *testing.T) {
	tree := createTestTree()
	result := tree.ToMarkdownString(0)

	// 检查输出是否包含预期的markdown列表格式
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

	// 确保目录有斜杠
	if !strings.Contains(result, "root/") || !strings.Contains(result, "dir1/") {
		t.Error("Directory names should have trailing slashes")
	}

	// 确保文件没有斜杠
	if strings.Contains(result, "file1.txt/") || strings.Contains(result, "file2.go/") {
		t.Error("File names should not have trailing slashes")
	}
}

func TestToMermaidString(t *testing.T) {
	tree := createTestTree()
	result := tree.ToMermaidString()

	// 检查mermaid输出格式是否正确
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

	// 确保节点ID递增
	if !strings.Contains(result, "N1") || !strings.Contains(result, "N2") ||
		!strings.Contains(result, "N3") || !strings.Contains(result, "N4") {
		t.Error("Mermaid output should contain incremental node IDs")
	}
}
