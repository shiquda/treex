package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// Format file size to human readable format
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

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

func (t *TreeNode) getEntryString(useIcons bool, showSize, showOwner, showModTime bool) string {
	s := t.Name
	if t.IsDir {
		s += "/"
	}
	if useIcons {
		s = getFileIcon(t.Name, t.IsDir) + s
	}
	
	// Add metadata information
	var metadata []string
	if showSize {
		metadata = append(metadata, formatFileSize(t.Size))
	}
	if showOwner && t.Owner != "" {
		metadata = append(metadata, t.Owner)
	}
	if showModTime {
		metadata = append(metadata, t.ModTime.Format("2006-01-02 15:04:05"))
	}
	
	if len(metadata) > 0 {
		s += " [" + strings.Join(metadata, ", ") + "]"
	}
	
	return s
}

func (t *TreeNode) ToIndentString(spaces int, useIcons bool, showSize, showOwner, showModTime bool) string {
	var result string
	for i := 0; i < t.Depth*spaces; i++ {
		result += " "
	}

	result += t.getEntryString(useIcons, showSize, showOwner, showModTime) + "\n"

	for _, child := range t.Children {
		result += child.ToIndentString(spaces, useIcons, showSize, showOwner, showModTime)
	}
	return result
}

func (t *TreeNode) ToTreeString(isLast bool, prefix string, useIcons bool, showSize, showOwner, showModTime bool) string {
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

	nodeName := t.getEntryString(useIcons, showSize, showOwner, showModTime)

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
		result += child.ToTreeString(isLastChild, childPrefix, useIcons, showSize, showOwner, showModTime)
	}
	return result
}

func (t *TreeNode) ToMarkdownString(level int, useIcons bool, showSize, showOwner, showModTime bool) string {
	var result string
	result += strings.Repeat("  ", level)

	result += "- "
	result += t.getEntryString(useIcons, showSize, showOwner, showModTime)
	result += "\n"

	// Process child nodes
	for _, child := range t.Children {
		result += child.ToMarkdownString(level+1, useIcons, showSize, showOwner, showModTime)
	}
	return result
}

func (t *TreeNode) ToMermaidString(showSize, showOwner, showModTime bool) string {
	var result string
	result += "graph TD\n" // Mermaid graph directive
	result += t.toMermaidNodes("", 1, showSize, showOwner, showModTime)
	return result
}

func (t *TreeNode) toMermaidNodes(parentID string, nodeID int, showSize, showOwner, showModTime bool) string {
	var result string
	currentID := fmt.Sprintf("N%d", nodeID)

	// Add current node
	var nodeLabel string
	if t.IsDir {
		nodeLabel = t.Name + "/"
	} else {
		nodeLabel = t.Name
	}

	// Add metadata to label
	var metadata []string
	if showSize {
		metadata = append(metadata, formatFileSize(t.Size))
	}
	if showOwner && t.Owner != "" {
		metadata = append(metadata, t.Owner)
	}
	if showModTime {
		metadata = append(metadata, t.ModTime.Format("2006-01-02 15:04:05"))
	}
	
	if len(metadata) > 0 {
		nodeLabel += "\\n[" + strings.Join(metadata, ", ") + "]"
	}

	result += fmt.Sprintf("    %s[\"%s\"]\n", currentID, nodeLabel)

	if parentID != "" {
		result += fmt.Sprintf("    %s --> %s\n", parentID, currentID)
	}

	// Process child nodes
	childID := nodeID + 1
	for _, child := range t.Children {
		result += child.toMermaidNodes(currentID, childID, showSize, showOwner, showModTime)
		childID += len(child.getAllNodes())
	}
	return result
}

// ToJSONString converts the tree to JSON format
func (t *TreeNode) ToJSONString(showSize, showOwner, showModTime bool) string {
	data := t.toJSONNode(showSize, showOwner, showModTime)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error creating JSON: %v", err)
	}
	return string(jsonData)
}

// toJSONNode creates a JSON representation of the node
func (t *TreeNode) toJSONNode(showSize, showOwner, showModTime bool) map[string]interface{} {
	node := map[string]interface{}{
		"name":  t.Name,
		"isDir": t.IsDir,
		"depth": t.Depth,
	}

	if showSize {
		node["size"] = t.Size
	}
	if showOwner && t.Owner != "" {
		node["owner"] = t.Owner
	}
	if showModTime {
		node["modTime"] = t.ModTime.Format("2006-01-02 15:04:05")
	}

	if len(t.Children) > 0 {
		children := make([]map[string]interface{}, len(t.Children))
		for i, child := range t.Children {
			children[i] = child.toJSONNode(showSize, showOwner, showModTime)
		}
		node["children"] = children
	}

	return node
}

// ToCSVString converts the tree to CSV format
func (t *TreeNode) ToCSVString(showSize, showOwner, showModTime bool) string {
	var result strings.Builder
	
	// Write header
	header := []string{"Name", "Type", "Depth"}
	if showSize {
		header = append(header, "Size")
	}
	if showOwner {
		header = append(header, "Owner")
	}
	if showModTime {
		header = append(header, "ModTime")
	}
	
	result.WriteString(strings.Join(header, ",") + "\n")
	
	// Write all nodes
	nodes := t.getAllNodes()
	for _, node := range nodes {
		row := []string{node.Name, getTypeString(node.IsDir), fmt.Sprintf("%d", node.Depth)}
		if showSize {
			row = append(row, formatFileSize(node.Size))
		}
		if showOwner {
			row = append(row, node.Owner)
		}
		if showModTime {
			row = append(row, node.ModTime.Format("2006-01-02 15:04:05"))
		}
		
		// Escape quotes and wrap in quotes if necessary
		for i, field := range row {
			if strings.Contains(field, ",") || strings.Contains(field, "\"") {
				row[i] = "\"" + strings.ReplaceAll(field, "\"", "\"\"") + "\""
			}
		}
		
		result.WriteString(strings.Join(row, ",") + "\n")
	}
	
	return result.String()
}

// getTypeString returns "Directory" or "File" based on isDir
func getTypeString(isDir bool) string {
	if isDir {
		return "Directory"
	}
	return "File"
}
