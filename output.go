package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Get file type icon
func getFileIcon(name string, isDir bool) string {
	if isDir {
		return "📁 " // Folder icon
	}

	extension := strings.ToLower(filepath.Ext(name))
	switch extension {
	case ".go":
		return "🔹 " // Go file
	case ".py":
		return "🐍 " // Python file
	case ".js", ".jsx", ".ts", ".tsx":
		return "📜 " // JavaScript/TypeScript file
	case ".html", ".htm":
		return "🌐 " // HTML file
	case ".css":
		return "🎨 " // CSS file
	case ".md":
		return "📝 " // Markdown file
	case ".json":
		return "📋 " // JSON file
	case ".xml":
		return "📋 " // XML file
	case ".yml", ".yaml":
		return "⚙️ " // YAML file
	case ".txt":
		return "📄 " // Plain text file
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg":
		return "🖼️ " // Image file
	case ".mp3", ".wav", ".ogg":
		return "🎵 " // Audio file
	case ".mp4", ".avi", ".mkv", ".mov":
		return "🎬 " // Video file
	case ".pdf":
		return "📕 " // PDF file
	case ".zip", ".tar", ".gz", ".7z", ".rar":
		return "📦 " // Archive file
	case ".exe", ".dll":
		return "⚡ " // Executable file
	case ".sh", ".bash", ".zsh", ".ps1":
		return "⚙️ " // Script file
	case ".c", ".cpp", ".h", ".hpp":
		return "🔧 " // C/C++ file
	case ".java":
		return "☕ " // Java file
	case ".rb":
		return "💎 " // Ruby file
	case ".php":
		return "🐘 " // PHP file
	case ".rs":
		return "🦀 " // Rust file
	case ".sql":
		return "🗄️ " // SQL file
	case ".gitignore", ".dockerignore":
		return "🔒 " // Ignore file
	default:
		return "📄 " // Default file icon
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

func (t *TreeNode) ToTreeString(isLast bool, prefix string, useIcons bool) string {
	var result string
	currentPrefix := prefix

	// current node prefix
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

func (t *TreeNode) ToMarkdownString(level int, useIcons bool) string {
	var result string
	result += strings.Repeat("  ", level)

	result += "- "
	result += t.getEntryString(useIcons)
	result += "\n"

	// Process child nodes
	for _, child := range t.Children {
		result += child.ToMarkdownString(level+1, useIcons)
	}
	return result
}

func (t *TreeNode) ToMermaidString() string {
	var result string
	result += "graph TD\n" // Mermaid graph directive
	result += t.toMermaidNodes("", 1)
	return result
}

func (t *TreeNode) toMermaidNodes(parentID string, nodeID int) string {
	var result string
	currentID := fmt.Sprintf("N%d", nodeID)

	// Add current node
	if t.IsDir {
		result += fmt.Sprintf("    %s[%s/]\n", currentID, t.Name)
	} else {
		result += fmt.Sprintf("    %s[%s]\n", currentID, t.Name)
	}

	if parentID != "" {
		result += fmt.Sprintf("    %s --> %s\n", parentID, currentID)
	}

	// Process child nodes
	childID := nodeID + 1
	for _, child := range t.Children {
		result += child.toMermaidNodes(currentID, childID)
		childID += len(child.getAllNodes())
	}
	return result
}
