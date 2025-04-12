package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

type TreeNode struct {
	Name     string
	IsDir    bool
	Children []*TreeNode
	Depth    int
}

type Filter struct {
	dirs      []string
	suffixes  []string
	patterns  []string
	gitignore bool
}

func NewFilter(excludeStr string, useGitIgnore bool) *Filter {
	f := &Filter{
		gitignore: useGitIgnore,
	}

	if excludeStr == "" && !useGitIgnore {
		return f
	}

	if excludeStr != "" {
		rules := strings.Split(excludeStr, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if strings.HasSuffix(rule, "/") {
				f.dirs = append(f.dirs, strings.TrimSuffix(rule, "/"))
			} else if strings.HasPrefix(rule, ".") {
				f.suffixes = append(f.suffixes, rule)
			} else {
				f.patterns = append(f.patterns, rule)
			}
		}
	}

	if useGitIgnore {
		f.loadGitIgnorePatterns()
	}

	return f
}

func (f *Filter) loadGitIgnorePatterns() {
	gitignorePath := ".gitignore"
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		// 如果 .gitignore 不存在，忽略错误
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 处理 .gitignore 规则
		if strings.HasSuffix(line, "/") {
			// 目录规则
			f.dirs = append(f.dirs, strings.TrimSuffix(line, "/"))
		} else if strings.HasPrefix(line, "*.") {
			// 后缀规则，如 *.txt
			f.suffixes = append(f.suffixes, strings.TrimPrefix(line, "*"))
		} else {
			// 其他模式
			f.patterns = append(f.patterns, line)
		}
	}
}

func (f *Filter) shouldExclude(name string, isDir bool, path string) bool {
	if isDir {
		for _, dir := range f.dirs {
			if matchPattern(name, dir) {
				return true
			}
		}
	} else {
		for _, suffix := range f.suffixes {
			if strings.HasSuffix(name, suffix) {
				return true
			}
		}
	}

	// 检查通用模式
	for _, pattern := range f.patterns {
		if matchPattern(path, pattern) || matchPattern(name, pattern) {
			return true
		}
	}

	return false
}

// 简单的通配符匹配
func matchPattern(name, pattern string) bool {
	// 完全匹配
	if pattern == name {
		return true
	}

	// 前缀星号: *suffix
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(name, strings.TrimPrefix(pattern, "*")) {
		return true
	}

	// 后缀星号: prefix*
	if strings.HasSuffix(pattern, "*") && strings.HasPrefix(name, strings.TrimSuffix(pattern, "*")) {
		return true
	}

	return false
}

func getRelativePath(absolute string, root string) string {
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}

	if !strings.HasPrefix(absolute, root) {
		return absolute
	}

	relativePath := strings.TrimPrefix(absolute, root)

	relativePath = strings.TrimSuffix(relativePath, "/")

	return relativePath
}

func getTreeNode(root string, depth int, basePath string, maxDepth int, filter *Filter, hideHidden bool, dirsOnly bool) (*TreeNode, error) {
	// 检查是否超过最大深度
	if maxDepth > 0 && depth > maxDepth {
		// 返回目录本身，但不递归获取其内容
		dirName := filepath.Base(strings.TrimSuffix(root, "/"))
		return &TreeNode{
			Name:  dirName,
			IsDir: true,
			Depth: depth - 1,
		}, nil
	}

	files, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	relativeName := getRelativePath(root, basePath)
	if depth == 1 {
		node := TreeNode{
			Name:  relativeName,
			IsDir: true,
			Depth: depth - 1,
		}
		for _, entry := range files {
			// 检查是否是隐藏文件
			if hideHidden && strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			childPath := relativeName
			if childPath != "" {
				childPath += "/"
			}
			childPath += entry.Name()

			if filter.shouldExclude(entry.Name(), entry.IsDir(), childPath) {
				continue
			}
			if entry.IsDir() {
				fullChildPath := root + "/" + entry.Name() + "/"
				child, e := getTreeNode(fullChildPath, depth+1, basePath, maxDepth, filter, hideHidden, dirsOnly)
				if e != nil {
					return nil, e
				}
				if child != nil {
					node.Children = append(node.Children, child)
				}
			} else if !dirsOnly {
				child := &TreeNode{
					Name:  entry.Name(),
					IsDir: false,
					Depth: depth,
				}
				node.Children = append(node.Children, child)
			}
		}
		return &node, nil
	} else {
		dirName := filepath.Base(strings.TrimSuffix(root, "/"))
		node := TreeNode{
			Name:  dirName,
			IsDir: true,
			Depth: depth - 1,
		}
		for _, entry := range files {
			// 检查是否是隐藏文件
			if hideHidden && strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			childPath := relativeName
			if childPath != "" {
				childPath += "/"
			}
			childPath += entry.Name()

			if filter.shouldExclude(entry.Name(), entry.IsDir(), childPath) {
				continue
			}
			if entry.IsDir() {
				fullChildPath := root + "/" + entry.Name() + "/"
				child, e := getTreeNode(fullChildPath, depth+1, basePath, maxDepth, filter, hideHidden, dirsOnly)
				if e != nil {
					return nil, e
				}
				if child != nil {
					node.Children = append(node.Children, child)
				}
			} else if !dirsOnly {
				child := &TreeNode{
					Name:  entry.Name(),
					IsDir: false,
					Depth: depth,
				}
				node.Children = append(node.Children, child)
			}
		}
		return &node, nil
	}
}

func (t *TreeNode) ToIndentString(spaces int) string {
	var result string
	for i := 0; i < t.Depth*spaces; i++ {
		result += " "
	}
	result += t.Name + "\n"
	for _, child := range t.Children {
		result += child.ToIndentString(spaces)
	}
	return result
}

func (t *TreeNode) ToTreeString(isLast bool, prefix string) string {
	var result string
	// current node prefix
	currentPrefix := prefix
	if t.Depth > 0 {
		if isLast {
			currentPrefix += "└── "
		} else {
			currentPrefix += "├── "
		}
	}
	result += currentPrefix + t.Name + "\n"

	// sub node prefix
	childPrefix := prefix
	if t.Depth > 0 {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	for i, child := range t.Children {
		isLastChild := i == len(t.Children)-1
		result += child.ToTreeString(isLastChild, childPrefix)
	}
	return result
}

func (t *TreeNode) ToMarkdownString(level int) string {
	var result string
	// 添加缩进
	for i := 0; i < level; i++ {
		result += "  "
	}
	// 添加列表标记
	result += "- "
	// 如果是目录，添加目录名和斜杠
	if t.IsDir {
		result += t.Name + "/"
	} else {
		result += t.Name
	}
	result += "\n"

	// 处理子节点
	for _, child := range t.Children {
		result += child.ToMarkdownString(level + 1)
	}
	return result
}

func (t *TreeNode) ToMermaidString() string {
	var result string
	result += "graph TD\n"
	result += t.toMermaidNodes("", 1)
	return result
}

func (t *TreeNode) toMermaidNodes(parentID string, nodeID int) string {
	var result string
	currentID := fmt.Sprintf("N%d", nodeID)

	// 添加当前节点
	if t.IsDir {
		result += fmt.Sprintf("    %s[%s/]\n", currentID, t.Name)
	} else {
		result += fmt.Sprintf("    %s[%s]\n", currentID, t.Name)
	}

	// 如果不是根节点，添加与父节点的连接
	if parentID != "" {
		result += fmt.Sprintf("    %s --> %s\n", parentID, currentID)
	}

	// 处理子节点
	childID := nodeID + 1
	for _, child := range t.Children {
		result += child.toMermaidNodes(currentID, childID)
		childID += len(child.getAllNodes())
	}
	return result
}

func (t *TreeNode) getAllNodes() []*TreeNode {
	nodes := []*TreeNode{t}
	for _, child := range t.Children {
		nodes = append(nodes, child.getAllNodes()...)
	}
	return nodes
}

func main() {
	dir := flag.StringP("dir", "d", ".", "directory to scan")
	format := flag.StringP("format", "f", "tree", "output format. allowed: [indent, tree, md, mermaid]")
	maxDepth := flag.IntP("max-depth", "m", 0, "maximum directory depth (0 for unlimited)")
	output := flag.StringP("output", "o", "", "output file path (default: stdout)")
	exclude := flag.StringP("exclude", "e", "", "exclude rules (comma-separated, e.g. 'dir/, .txt')")
	hideHidden := flag.BoolP("hide-hidden", "H", false, "hide hidden files and directories (default: false)")
	dirsOnly := flag.BoolP("dirs-only", "D", false, "show directories only (default: false)")
	useGitIgnore := flag.BoolP("use-gitignore", "I", false, "use .gitignore patterns to exclude files/directories (default: false)")
	flag.Parse()

	// 获取绝对路径并确保以"/"结尾
	absPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}
	if !strings.HasSuffix(absPath, "/") {
		absPath += "/"
	}

	filter := NewFilter(*exclude, *useGitIgnore)
	node, err := getTreeNode(*dir, 1, absPath, *maxDepth, filter, *hideHidden, *dirsOnly)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		flag.Usage()
		return
	}

	var outputStr string
	switch *format {
	case "tree":
		outputStr = node.ToTreeString(true, "")
	case "indent":
		outputStr = node.ToIndentString(4)
	case "md":
		outputStr = node.ToMarkdownString(0)
	case "mermaid":
		outputStr = node.ToMermaidString()
	default:
		fmt.Fprintf(os.Stderr, "error: unknown format '%s'\n", *format)
		flag.Usage()
		return
	}

	if *output != "" {
		// create output dir (if not exist)
		outputDir := filepath.Dir(*output)
		if outputDir != "." {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "error creating output directory: %s\n", err)
				return
			}
		}
		// write
		if err := os.WriteFile(*output, []byte(outputStr), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing to file: %s\n", err)
			return
		}
		fmt.Printf("Output written to: %s\n", *output)
	} else {
		fmt.Print(outputStr)
	}
}
