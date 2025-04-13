package main

import (
	"os"
	"testing"
)

func TestNewFilter(t *testing.T) {
	// 测试空过滤器
	f := NewFilter("", false)
	if len(f.dirs) != 0 || len(f.suffixes) != 0 || len(f.patterns) != 0 {
		t.Error("Empty filter should have no rules")
	}

	// 测试各种规则
	f = NewFilter("dir/, .txt, pattern*", false)
	if len(f.dirs) != 1 || f.dirs[0] != "dir" {
		t.Error("Filter should have one directory rule")
	}
	if len(f.suffixes) != 1 || f.suffixes[0] != ".txt" {
		t.Error("Filter should have one suffix rule")
	}
	if len(f.patterns) != 1 || f.patterns[0] != "pattern*" {
		t.Error("Filter should have one pattern rule")
	}
}

func TestShouldExclude(t *testing.T) {
	// 测试目录过滤
	f := NewFilter("node_modules/, .git/", false)
	if !f.shouldExclude("node_modules", true, "node_modules") {
		t.Error("Should exclude node_modules directory")
	}
	if f.shouldExclude("src", true, "src") {
		t.Error("Should not exclude src directory")
	}

	// 测试后缀过滤
	f = NewFilter(".txt, .log", false)
	if !f.shouldExclude("file.txt", false, "file.txt") {
		t.Error("Should exclude .txt files")
	}
	if !f.shouldExclude("file.log", false, "file.log") {
		t.Error("Should exclude .log files")
	}
	if f.shouldExclude("file.go", false, "file.go") {
		t.Error("Should not exclude .go files")
	}

	// 测试模式匹配
	f = NewFilter("test*, *temp", false)
	if !f.shouldExclude("test_file", false, "test_file") {
		t.Error("Should exclude files starting with test")
	}
	if !f.shouldExclude("file_temp", false, "file_temp") {
		t.Error("Should exclude files ending with temp")
	}
	if f.shouldExclude("source.go", false, "source.go") {
		t.Error("Should not exclude normal go files")
	}
}

func TestGitIgnoreIntegration(t *testing.T) {
	// 创建临时的.gitignore文件用于测试
	content := []byte("*.log\nbuild/\n# comment\ntemp*\n")
	err := os.WriteFile(".gitignore.test", content, 0644)
	if err != nil {
		t.Fatal("Failed to create test .gitignore file")
	}
	defer os.Remove(".gitignore.test")

	// 重命名现有的.gitignore文件（如果有）
	if _, err := os.Stat(".gitignore"); err == nil {
		os.Rename(".gitignore", ".gitignore.bak")
		defer os.Rename(".gitignore.bak", ".gitignore")
	}
	os.Rename(".gitignore.test", ".gitignore")
	defer os.Rename(".gitignore", ".gitignore.test")

	// 测试gitignore集成
	f := NewFilter("", true)

	// 检查是否正确读取了规则
	if len(f.suffixes) == 0 || f.suffixes[0] != ".log" {
		t.Error("Should read .log suffix rule from .gitignore")
	}

	if len(f.dirs) == 0 || f.dirs[0] != "build" {
		t.Error("Should read build directory rule from .gitignore")
	}

	// 测试匹配
	if !f.shouldExclude("server.log", false, "server.log") {
		t.Error("Should exclude .log files")
	}
	if !f.shouldExclude("build", true, "build") {
		t.Error("Should exclude build directory")
	}
	if !f.shouldExclude("temp_file", false, "temp_file") {
		t.Error("Should exclude files starting with temp")
	}
}
