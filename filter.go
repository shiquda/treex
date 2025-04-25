package main

import (
	"os"
	"strings"
)

type Filter struct {
	dirNames []string
	suffixes []string
	patterns []string
}

func NewFilter(excludeRuleStr string, useGitIgnore bool) *Filter {
	f := &Filter{}

	if excludeRuleStr == "" && !useGitIgnore {
		return f
	}

	if excludeRuleStr != "" {
		rules := strings.Split(excludeRuleStr, ",")
		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if strings.HasSuffix(rule, "/") {
				f.dirNames = append(f.dirNames, strings.TrimSuffix(rule, "/"))
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
		// If .gitignore doesn't exist, ignore the error
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Process .gitignore rules
		if strings.HasSuffix(line, "/") {
			// Directory rule
			f.dirNames = append(f.dirNames, strings.TrimSuffix(line, "/"))
		} else if strings.HasPrefix(line, "*.") {
			// Suffix rule, e.g. *.txt
			f.suffixes = append(f.suffixes, strings.TrimPrefix(line, "*"))
		} else {
			// Other patterns
			f.patterns = append(f.patterns, line)
		}
	}
}

func (f *Filter) shouldExclude(name string, isDir bool, path string) bool {
	if isDir {
		for _, dir := range f.dirNames {
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

	// Check general patterns
	for _, pattern := range f.patterns {
		if matchPattern(path, pattern) || matchPattern(name, pattern) {
			return true
		}
	}

	return false
}

// Simple wildcard matching
func matchPattern(name, pattern string) bool {
	// Exact match
	if pattern == name {
		return true
	}

	// Prefix star: *suffix
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(name, strings.TrimPrefix(pattern, "*")) {
		return true
	}

	// Suffix star: prefix*
	if strings.HasSuffix(pattern, "*") && strings.HasPrefix(name, strings.TrimSuffix(pattern, "*")) {
		return true
	}

	return false
}
