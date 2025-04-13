package main

import (
	"os"
	"strings"
)

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
