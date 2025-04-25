package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

func main() {
	// parse flags
	dir := flag.StringP("dir", "d", ".", "directory to scan")
	outputFormat := flag.StringP("format", "f", "tree", "output format. allowed: [indent, tree, md, mermaid]")
	maxDepth := flag.IntP("max-depth", "m", 0, "maximum directory depth (0 for unlimited)")
	outputFilePath := flag.StringP("output", "o", "", "output file path (default: stdout)")
	excludeRuleStr := flag.StringP("exclude", "e", "", "exclude rules (comma-separated, e.g. 'dir/, .txt')")
	hideHidden := flag.BoolP("hide-hidden", "H", false, "hide hidden files and directories (default: false)")
	dirsOnly := flag.BoolP("dirs-only", "D", false, "show directories only (default: false)")
	useGitIgnore := flag.BoolP("use-gitignore", "I", false, "use .gitignore patterns to exclude files/directories (default: false)")
	useIcons := flag.BoolP("icons", "C", false, "display file type icons (default: false)")
	flag.Parse()

	// get the absolute path and ensure it ends with "/"
	absolutePath, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}
	if !strings.HasSuffix(absolutePath, "/") {
		absolutePath += "/"
	}

	// filters
	filter := NewFilter(*excludeRuleStr, *useGitIgnore)
	node, err := getTreeNode(*dir, 1, absolutePath, *maxDepth, filter, *hideHidden, *dirsOnly)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		flag.Usage()
		return
	}

	// output
	var outputStr string
	switch *outputFormat {
	case "tree":
		outputStr = node.ToTreeString(true, "", *useIcons)
	case "indent":
		outputStr = node.ToIndentString(4, *useIcons)
	case "md":
		outputStr = node.ToMarkdownString(0, *useIcons)
	case "mermaid":
		outputStr = node.ToMermaidString()
	default:
		fmt.Fprintf(os.Stderr, "error: unknown outputFormat '%s'\n", *outputFormat)
		flag.Usage()
		return
	}

	// write to file
	if *outputFilePath != "" {
		// create output dir (if not exist)
		outputDirPath := filepath.Dir(*outputFilePath)
		if outputDirPath != "." {
			if err := os.MkdirAll(outputDirPath, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "error creating output directory: %s\n", err)
				return
			}
		}

		if err := os.WriteFile(*outputFilePath, []byte(outputStr), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing to file: %s\n", err)
			return
		}
		fmt.Printf("Output written to: %s\n", *outputFilePath)
	} else {
		fmt.Print(outputStr)
	}
}
