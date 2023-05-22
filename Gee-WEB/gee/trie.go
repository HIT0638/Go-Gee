package gee

import "strings"

type node struct {
	pattern  string  // 待匹配的路由，从根结点到该结点的路径，如 /p/:lang
	part     string  // 该结点存放的部分路由，如 :lang
	children []*node // 该结点的子节点
	isWild   bool    // 是否精准匹配，:lang 或 *filepath 时为true
}

// matchChild 找到该节点的子结点中与目标part匹配成功的结点
// 用于 insert 方法过程中的搜索和确定插入位置。
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		// 如果某个子结点part匹配成功 或者 属于模糊匹配
		if child.part == part || child.isWild {
			return child
		}
	}

	return nil
}

// matchChildren 查找当前结点的子结点中所有匹配成功的结点，记录
// 用于 search 方法过程中每一步的匹配与搜索，相当于添加成功条件的广搜
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)

	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

/*
*
insert 根据给定的路由，在前缀树中添加节点，用于后续的匹配

	@param pattern 接收到的URL
	@param parts pattern 分解之后的每一部分
	@param height 记录当前递归层数（树高度）
*/
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)

	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	// 搜索匹配成功的条件：
	// (len(parts) == height || strings.HasPrefix(n.part, "*")) & n.pattern != ""
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
