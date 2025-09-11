package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TreeNode struct {
	Name     string
	IsDir    bool
	Children []*TreeNode
	Depth    int
	Size     int64
	Owner    string
	ModTime  time.Time
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

func getFileMetadata(root string, entry os.DirEntry) (int64, string, time.Time, error) {
	var size int64
	var modTime time.Time
	var owner string

	info, err := entry.Info()
	if err != nil {
		return 0, "", time.Time{}, err
	}

	size = info.Size()
	modTime = info.ModTime()

	// Owner information is complex across platforms, so we'll leave it empty for now
	owner = ""

	return size, owner, modTime, nil
}

func getTreeNode(root string, depth int, basePath string, maxDepth int, filter *Filter, hideHidden bool, dirsOnly bool) (*TreeNode, error) {
	// Check if max depth is exceeded
	if maxDepth > 0 && depth > maxDepth {
		// Return directory itself without recursively getting its contents
		dirName := filepath.Base(strings.TrimSuffix(root, "/"))
		
		// Get directory metadata for depth-limited directories
		dirInfo, err := os.Stat(root)
		if err != nil {
			return nil, err
		}
		
		return &TreeNode{
			Name:    dirName,
			IsDir:   true,
			Depth:   depth - 1,
			Size:    dirInfo.Size(),
			ModTime: dirInfo.ModTime(),
		}, nil
	}

	files, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	relativeName := getRelativePath(root, basePath)

	// Determine node name
	nodeName := relativeName
	if depth > 1 {
		nodeName = filepath.Base(strings.TrimSuffix(root, "/"))
	}

	// Create node for directory
	// Get directory metadata
	dirInfo, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	node := TreeNode{
		Name:    nodeName,
		IsDir:   true,
		Depth:   depth - 1,
		Size:    dirInfo.Size(),
		ModTime: dirInfo.ModTime(),
	}

	// Process child entries
	for _, entry := range files {
		// Check if it's a hidden file
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
			fullChildPath := root + "/" + entry.Name()
			size, owner, modTime, err := getFileMetadata(fullChildPath, entry)
			if err != nil {
				return nil, err
			}
			
			child := &TreeNode{
				Name:    entry.Name(),
				IsDir:   false,
				Depth:   depth,
				Size:    size,
				Owner:   owner,
				ModTime: modTime,
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
