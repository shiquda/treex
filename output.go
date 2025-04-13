package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// 获取文件类型图标
func getFileIcon(name string, isDir bool) string {
	if isDir {
		return "📁 " // 文件夹图标
	}

	extension := strings.ToLower(filepath.Ext(name))
	switch extension {
	case ".go":
		return "🔹 " // Go 文件
	case ".py":
		return "🐍 " // Python 文件
	case ".js", ".jsx", ".ts", ".tsx":
		return "📜 " // JavaScript/TypeScript 文件
	case ".html", ".htm":
		return "🌐 " // HTML 文件
	case ".css":
		return "🎨 " // CSS 文件
	case ".md":
		return "📝 " // Markdown 文件
	case ".json":
		return "📋 " // JSON 文件
	case ".xml":
		return "📋 " // XML 文件
	case ".yml", ".yaml":
		return "⚙️ " // YAML 文件
	case ".txt":
		return "📄 " // 纯文本文件
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg":
		return "🖼️ " // 图片文件
	case ".mp3", ".wav", ".ogg":
		return "🎵 " // 音频文件
	case ".mp4", ".avi", ".mkv", ".mov":
		return "🎬 " // 视频文件
	case ".pdf":
		return "📕 " // PDF 文件
	case ".zip", ".tar", ".gz", ".7z", ".rar":
		return "📦 " // 压缩文件
	case ".exe", ".dll":
		return "⚡ " // 可执行文件
	case ".sh", ".bash", ".zsh", ".ps1":
		return "⚙️ " // 脚本文件
	case ".c", ".cpp", ".h", ".hpp":
		return "🔧 " // C/C++ 文件
	case ".java":
		return "☕ " // Java 文件
	case ".rb":
		return "💎 " // Ruby 文件
	case ".php":
		return "🐘 " // PHP 文件
	case ".rs":
		return "🦀 " // Rust 文件
	case ".sql":
		return "🗄️ " // SQL 文件
	case ".gitignore", ".dockerignore":
		return "🔒 " // 忽略文件
	default:
		return "📄 " // 默认为普通文件
	}
}

func (t *TreeNode) getEntryString(useIcons bool) string {
	s := t.Name
	if t.IsDir {
		s += "/"
	}
	if useIcons {
		s = getFileIcon(t.Name, t.IsDir) + s
	}
	return s
}

// 缩进格式输出
func (t *TreeNode) ToIndentString(spaces int, useIcons bool) string {
	var result string
	for i := 0; i < t.Depth*spaces; i++ {
		result += " "
	}

	result += t.getEntryString(useIcons) + "\n"

	for _, child := range t.Children {
		result += child.ToIndentString(spaces, useIcons)
	}
	return result
}

// 树形结构输出
func (t *TreeNode) ToTreeString(isLast bool, prefix string, useIcons bool) string {
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

	nodeName := t.getEntryString(useIcons)

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
		result += child.ToTreeString(isLastChild, childPrefix, useIcons)
	}
	return result
}

// Markdown格式输出
func (t *TreeNode) ToMarkdownString(level int, useIcons bool) string {
	var result string
	// 添加缩进
	for i := 0; i < level; i++ {
		result += "  "
	}
	// 添加列表标记
	result += "- "
	result += t.getEntryString(useIcons)
	result += "\n"

	// 处理子节点
	for _, child := range t.Children {
		result += child.ToMarkdownString(level+1, useIcons)
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
