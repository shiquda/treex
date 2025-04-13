package main

import (
	"os"
	"path/filepath"
	"strings"
)

type TreeNode struct {
	Name     string
	IsDir    bool
	Children []*TreeNode
	Depth    int
}

func getRelativePath(absolute string, root string) string {
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}

	if !strings.HasPrefix(absolute, root) {
		return absolute
	}

	relativePath := strings.TrimPrefix(absolute, root)
	relativePath = strings.TrimSuffix(relativePath, "/")

	return relativePath
}

func getTreeNode(root string, depth int, basePath string, maxDepth int, filter *Filter, hideHidden bool, dirsOnly bool) (*TreeNode, error) {
	// 检查是否超过最大深度
	if maxDepth > 0 && depth > maxDepth {
		// 返回目录本身，但不递归获取其内容
		dirName := filepath.Base(strings.TrimSuffix(root, "/"))
		return &TreeNode{
			Name:  dirName,
			IsDir: true,
			Depth: depth - 1,
		}, nil
	}

	files, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	relativeName := getRelativePath(root, basePath)

	// 决定节点名称
	nodeName := relativeName
	if depth > 1 {
		nodeName = filepath.Base(strings.TrimSuffix(root, "/"))
	}

	// 创建节点
	node := TreeNode{
		Name:  nodeName,
		IsDir: true,
		Depth: depth - 1,
	}

	// 处理子条目
	for _, entry := range files {
		// 检查是否是隐藏文件
		if hideHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		childPath := relativeName
		if childPath != "" {
			childPath += "/"
		}
		childPath += entry.Name()

		if filter.shouldExclude(entry.Name(), entry.IsDir(), childPath) {
			continue
		}

		if entry.IsDir() {
			fullChildPath := root + "/" + entry.Name() + "/"
			child, e := getTreeNode(fullChildPath, depth+1, basePath, maxDepth, filter, hideHidden, dirsOnly)
			if e != nil {
				return nil, e
			}
			if child != nil {
				node.Children = append(node.Children, child)
			}
		} else if !dirsOnly {
			child := &TreeNode{
				Name:  entry.Name(),
				IsDir: false,
				Depth: depth,
			}
			node.Children = append(node.Children, child)
		}
	}

	return &node, nil
}

func (t *TreeNode) getAllNodes() []*TreeNode {
	nodes := []*TreeNode{t}
	for _, child := range t.Children {
		nodes = append(nodes, child.getAllNodes()...)
	}
	return nodes
}
