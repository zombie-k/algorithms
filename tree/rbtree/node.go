package rbtree

import "fmt"

type color bool

const black, red color = true, false

var NIL = &Node{nil, nil, black, nil, nil, nil}

type Node struct {
	key   interface{}
	value interface{}
	color color

	left   *Node
	right  *Node
	parent *Node
}

func NewNode(key, value interface{}) *Node {
	return &Node{
		key:    key,
		value:  value,
		color:  red,
		left:   NIL,
		right:  NIL,
		parent: NIL,
	}
}

func (node *Node) String() string {
	//return fmt.Sprintf("{key:%v, value:%v]", node.key, node.value)
	if node.color == red {
		return fmt.Sprintf("*%v", node.key)
	}
	return fmt.Sprintf("%v", node.key)
}

// Size 返回subtree节点个数
func (node *Node) Size() int {
	if node == nil {
		return 0
	}
	size := 1
	if node.left != nil {
		size += node.left.Size()
	}
	if node.right != nil {
		size += node.right.Size()
	}
	return size
}

// 祖父节点
func (node *Node) grandParent() *Node {
	if node != nil && node.parent != nil {
		return node.parent.parent
	}
	return nil
}

// 兄弟节点
func (node *Node) brother() *Node {
	if node == nil || node.parent == nil {
		return nil
	}
	if node == node.parent.left {
		return node.parent.right
	}
	return node.parent.left
}

// 叔伯节点
func (node *Node) uncle() *Node {
	if node == nil || node.parent == nil || node.grandParent() == nil {
		return nil
	}
	return node.parent.brother()
}

func (node *Node) maxNode() *Node {
	if node == nil {
		return nil
	}
	n := node.right
	for n != nil {
		n = node.right
	}
	return n
}

func (node *Node) Color() color {
	if node == nil {
		return black
	}
	return node.color
}

func (node *Node) output(prefix string, isTail bool, str *string) {
	if node.right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "|    "
		} else {
			newPrefix = "    "
		}
		node.right.output(newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "|    "
		}
		node.left.output(newPrefix, true, str)
	}
}

func nodeColor(node *Node) color {
	return node.color
}
