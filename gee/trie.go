package gee

import "strings"

type node struct {
	full_path   string // /p/hello
	part        string // hello
	children    []*node
	isprecesion bool // if * or :(sth) -> true else false
}

// For insertion
func (n *node) find_child(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

// For search
func (n *node) find_children(part string) []*node {
	result := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isprecesion {
			result = append(result, child)
		}
	}
	return result
}

func (n *node) insert(full_path string, parts []string, height int) {

	//因为根节点那一层是个dummy node 必然匹配不到，所以 len(parts) == height
	if len(parts) == height {
		n.full_path = full_path
		return
	}

	part := parts[height]
	child := n.find_child(part)

	// 例如：匹配 /hello/a， hello 下面没有 a 的node，创建一个并插入
	if child == nil {
		child = &node{part: part, isprecesion: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	child.insert(full_path, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {

	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 注意：inser当中对于full_path的赋值是发生在最后一个点
		// 假如当前点是由于 /p/*/m 这种形式得来，当我试图匹配 /p/q时应该返回失败
		// 如果该节点是由于 /p/* 得来，那么full_path应不为空 而是 /p/*
		if n.full_path == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.find_children(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
