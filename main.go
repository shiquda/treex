package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	flag "github.com/spf13/pflag"
)

func main() {
	// parse flags
	dir := flag.StringP("dir", "d", ".", "directory to scan")
	outputFormat := flag.StringP("format", "f", "tree", "output format. allowed: [indent, tree, md, mermaid, json, csv]")
	maxDepth := flag.IntP("max-depth", "m", 0, "maximum directory depth (0 for unlimited)")
	level := flag.IntP("level", "L", 0, "maximum directory depth (0 for unlimited, alias for -m)")
	showSize := flag.BoolP("size", "s", false, "show file sizes")
	showOwner := flag.BoolP("owner", "O", false, "show file owners")
	showModTime := flag.BoolP("modtime", "M", false, "show file modification times")
	outputFilePath := flag.StringP("output", "o", "", "output file path (default: stdout)")
	excludeRuleStr := flag.StringP("exclude", "e", "", "exclude rules (comma-separated, e.g. 'dir/, .txt')")
	hideHidden := flag.BoolP("hide-hidden", "H", false, "hide hidden files and directories (default: false)")
	dirsOnly := flag.BoolP("dirs-only", "D", false, "show directories only (default: false)")
	useGitIgnore := flag.BoolP("use-gitignore", "I", false, "use .gitignore patterns to exclude files/directories (default: false)")
	useIcons := flag.BoolP("icons", "C", false, "display file type icons (default: false)")
	flag.Parse()

	// Use level if max-depth is not explicitly set
	// Note: pflag doesn't have Changed method, so we'll use the default value check
	if *level != 0 && *maxDepth == 0 {
		maxDepth = level
	}

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
		outputStr = node.ToTreeString(true, "", *useIcons, *showSize, *showOwner, *showModTime)
	case "indent":
		outputStr = node.ToIndentString(4, *useIcons, *showSize, *showOwner, *showModTime)
	case "md":
		outputStr = node.ToMarkdownString(0, *useIcons, *showSize, *showOwner, *showModTime)
	case "mermaid":
		outputStr = node.ToMermaidString(*showSize, *showOwner, *showModTime)
	case "json":
		outputStr = node.ToJSONString(*showSize, *showOwner, *showModTime)
	case "csv":
		outputStr = node.ToCSVString(*showSize, *showOwner, *showModTime)
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
		// Handle output encoding for Windows pipe operations
		if runtime.GOOS == "windows" {
			// Check if output is being piped
			if fileInfo, err := os.Stdout.Stat(); err == nil {
				// If it's a character device (console), use normal output
				if (fileInfo.Mode() & os.ModeCharDevice) != 0 {
					fmt.Print(outputStr)
				} else {
					// It's a pipe or file, ensure UTF-8 encoding
					fmt.Fprintf(os.Stdout, "%s", outputStr)
				}
			} else {
				fmt.Print(outputStr)
			}
		} else {
			fmt.Print(outputStr)
		}
	}
}
