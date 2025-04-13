package main

import (
	"os"
	"path/filepath"
	"testing"
)

// 用于测试的临时目录结构
func setupTestDir(t *testing.T) string {
	tempDir := filepath.Join(os.TempDir(), "treex_test")
	os.RemoveAll(tempDir) // 清理之前的测试目录

	// 创建目录结构
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

	// 创建一些文件
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

	// 为gitignore测试创建文件
	ignoreContent := "*.log\nbuild/\n"
	if err := os.WriteFile(filepath.Join(tempDir, ".gitignore"), []byte(ignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore file: %v", err)
	}

	// 创建应该被排除的文件
	os.WriteFile(filepath.Join(tempDir, "log.log"), []byte("log"), 0644)
	os.MkdirAll(filepath.Join(tempDir, "build"), 0755)
	os.WriteFile(filepath.Join(tempDir, "build", "app.js"), []byte("js"), 0644)

	return tempDir
}

// 通过程序逻辑集成测试不同类型的输出
func TestIntegration(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	// 确保工作目录正确
	oldWd, _ := os.Getwd()
	os.Chdir(testDir)
	defer os.Chdir(oldWd)

	// 测试场景
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
				// 检查所有文件和目录都包括在内
				allNodes := node.getAllNodes()
				return len(allNodes) >= 10 // 根 + 至少9个子节点
			},
		},
		{
			name:       "Hidden Files Filter",
			filter:     NewFilter("", false),
			hideHidden: true,
			dirsOnly:   false,
			maxDepth:   0,
			checkFunc: func(node *TreeNode) bool {
				// 检查隐藏文件是否被过滤
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
				// 检查是否只包含目录
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
				// 检查深度是否受限
				maxFoundDepth := 0
				for _, n := range allChildrenRecursive(node) {
					if n.Depth > maxFoundDepth {
						maxFoundDepth = n.Depth
					}
				}
				return maxFoundDepth <= 1 // 深度应该不超过1
			},
		},
		{
			name:       "Exclude Rules",
			filter:     NewFilter(".log, build/", false),
			hideHidden: false,
			dirsOnly:   false,
			maxDepth:   0,
			checkFunc: func(node *TreeNode) bool {
				// 检查排除规则是否生效
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
				// 检查.gitignore规则是否生效
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

			// 测试各种输出格式
			_ = node.ToTreeString(true, "", false)
			_ = node.ToIndentString(2, false)
			_ = node.ToMarkdownString(0, false)
			_ = node.ToMermaidString()
		})
	}
}

// 递归获取所有子节点
func allChildrenRecursive(node *TreeNode) []*TreeNode {
	var result []*TreeNode
	for _, child := range node.Children {
		result = append(result, child)
		result = append(result, allChildrenRecursive(child)...)
	}
	return result
}
