package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

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
