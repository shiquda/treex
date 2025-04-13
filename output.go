package main

import (
	"fmt"
)

// 缩进格式输出
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

// 树形结构输出
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

	// 如果是目录，则在名称后添加斜杠
	nodeName := t.Name
	if t.IsDir {
		nodeName += "/"
	}

	result += currentPrefix + nodeName + "\n"

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

// Markdown格式输出
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

// Mermaid图表格式输出
func (t *TreeNode) ToMermaidString() string {
	var result string
	result += "graph TD\n"
	result += t.toMermaidNodes("", 1)
	return result
}

// 生成Mermaid节点
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
