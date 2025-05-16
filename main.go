package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

func main() {
	cfg := NewDefaultConfig()

	// parse flags
	cfg.Dir = *flag.StringP("dir", "d", ".", "directory to scan")
	cfg.OutputFormat = *flag.StringP("format", "f", "tree", "output format. allowed: [indent, tree, md, mermaid]")
	cfg.MaxDepth = *flag.IntP("max-depth", "m", 0, "maximum directory depth (0 for unlimited)")
	cfg.OutputFilePath = *flag.StringP("output", "o", "", "output file path (default: stdout)")
	cfg.ExcludeRuleStr = *flag.StringP("exclude", "e", "", "exclude rules (comma-separated, e.g. 'dir/, .txt')")
	cfg.HideHidden = *flag.BoolP("hide-hidden", "H", false, "hide hidden files and directories (default: false)")
	cfg.DirsOnly = *flag.BoolP("dirs-only", "D", false, "show directories only (default: false)")
	cfg.UseGitIgnore = *flag.BoolP("use-gitignore", "I", false, "use .gitignore patterns to exclude files/directories (default: false)")
	cfg.UseIcons = *flag.BoolP("icons", "C", false, "display file type icons (default: false)")

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
	filter := NewFilter(cfg.ExcludeRuleStr, cfg.UseGitIgnore)
	node, err := getTreeNode(cfg.Dir, 1, absolutePath, *maxDepth, filter, *hideHidden, *dirsOnly)
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
