package rbtree

import "github.com/zombie-k/algorithms/tree/util"

// 插入操作
// CASE 1：如果关注节点是 a（待插入节点），它的叔叔节点 d 是红色，(红色 带*号)
// 依次执行下面的操作:
//  1. 将关注节点 a 的父节点 b，叔叔节点 d 的顔色都设置成黑色
//  2. 将关注节点 a 的祖父节点 c 的顔色设置成红色
//  3. 关注节点变成 a 的祖父节点 c
//  4. 跳到 CASE2 或者 CASE3 继续处理
//     c                        c*
//     /   \                    /   \
//     b*     d*   --------->   b      d
//     / \    / \               / \    / \
//     z  a*  p  q              z  a*  p  q
//     / \					  / \
//     x   y                    x   y
//
// CASE 2:
// 1. 如果关注节点是 a ， 他的叔叔节点 d 是黑色，关注节点 a 是其父节点 b 的右子节点, 其父节点 b 是其祖父节点 e 的左子节点
// 依次执行下面的操作：
//  1. 关注节点变成节点 a 的父节点 b
//  2. 围绕新的关注节点 b 左旋
//
// 2. 如果关注节点是 a ， 他的叔叔节点 d 是黑色，关注节点 a 是其父节点 b 的左子节点, 其父节点 b 是其祖父节点 e 的右子节点
// 依次执行下面的操作：
//
//  1. 关注节点变成节点 a 的父节点 b
//
//  2. 围绕新的关注节点 b 右旋
//
//  3. 跳到 CASE 3
//     c                        c
//     /   \     b左旋转         /   \
//     b*     d    --------->   a*     d
//     / \    / \               / \    / \
//     z  a*  p  q              b*  y  p  q
//     / \				   / \
//     x   y                 z  x
//
//     c                        c
//     /   \     b右旋转         /   \
//     d      b*   --------->   d      a*
//     / \    / \               / \    / \
//     p  q   a*  z             p  q   x   b
//     / \				          / \
//     x  y                         y  z
//
// CASE 3: 如果关注节点是 a, 他的叔叔节点 d 是黑色，关注节点 a 是其父节点 b 的左子节点
// 依次执行下面的操作：
//  1. 围绕关注节点 a 的祖父节点 c 右旋
//  2. 将关注节点 a 的父节点 b，兄弟节点 c 的顔色互换
//  3. 调整结束，形成红黑树
//     c                        b*							b
//     /   \     c右旋转         /   \       顔色互换		      /   \
//     b*     d    --------->   a*     c     --------->       a*     c*
//     / \    / \               / \    / \                    / \    / \
//     a*  y    p  q             x  z   y  d                   x  z   y  d
//     / \				                 / \                           / \
//     x  z                                p  q                          p  q

type Tree struct {
	root       *Node
	size       int
	comparator util.Comparator
}

func NewTree(comparator util.Comparator) *Tree {
	return &Tree{root: NIL, comparator: comparator}
}

func (tree *Tree) Empty() bool {
	return tree.size == 0
}

func (tree *Tree) Size() int {
	return tree.size
}

func (tree *Tree) String() string {
	str := "RedBlackTree\n"
	if !tree.Empty() {
		output(tree.root, "", true, &str)
	}
	return str
}

func (tree *Tree) Get(key interface{}) (interface{}, bool) {
	node := tree.lookup(key)
	if node != nil {
		return node.value, true
	}
	return nil, false
}

// Insert 插入
func (tree *Tree) Insert(key interface{}, value interface{}) {
	var insertNode *Node

	if tree.root == NIL {
		tree.comparator(key, key)
		tree.root = &Node{key: key, value: value, color: red, left: NIL, right: NIL, parent: NIL}
		insertNode = tree.root
		insertNode.parent = NIL
		tree.insertFixup(insertNode)
		tree.size++
		return
	}

	node := tree.root
	loop := true
	for loop {
		compare := tree.comparator(key, node.key)
		switch {
		case compare == 0:
			node.key = key
			node.value = value
			return
		case compare < 0:
			if node.left == NIL {
				insertNode = &Node{key: key, value: value, color: red, left: NIL, right: NIL, parent: NIL}
				node.left = insertNode
				loop = false
			} else {
				node = node.left
			}
		case compare > 0:
			if node.right == NIL {
				insertNode = &Node{key: key, value: value, color: red, left: NIL, right: NIL, parent: NIL}
				node.right = insertNode
				loop = false
			} else {
				node = node.right
			}
		}
	}
	insertNode.parent = node

	tree.insertFixup(insertNode)
	tree.size++
}

// Delete
// z 是待删除的目标节点
// y 是z的子节点或后继节点, 用于替换z
// x 是y的右子树, 用于替换y
func (tree *Tree) Delete(key interface{}) {
	z := tree.lookup(key)
	if z == NIL {
		return
	}
	tree.delete(z)
}

func (tree *Tree) Maximum() (interface{}, interface{}){
	if tree.root == NIL {
		return nil, nil
	}
	maxi := tree.maximum(tree.root)
	return maxi.key, maxi.value
}

func (tree *Tree) Minimum() (interface{}, interface{}){
	if tree.root == NIL {
		return nil, nil
	}
	maxi := tree.minimum(tree.root)
	return maxi.key, maxi.value
}

func (tree *Tree) Floor(key interface{}) (interface{}, interface{}, bool) {
	node, found := tree.FloorNode(key)
	if !found {
		return nil, nil, found
	}
	return node.key, node.value, found
}

func (tree *Tree) FloorNode(key interface{}) (floor *Node, found bool) {
	found = false
	node := tree.root
	for node != NIL {
		compare := tree.comparator(key, node.key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			node = node.left
		case compare > 0:
			floor, found = node, true
			node = node.right
		}
	}
	if found {
		return floor, true
	}
	return nil, false
}

func (tree *Tree) Ceil(key interface{}) (interface{}, interface{}, bool) {
	node, found := tree.CeilNode(key)
	if !found {
		return nil, nil, found
	}
	return node.key, node.value, found
}

func (tree *Tree) CeilNode(key interface{}) (ceil *Node, found bool) {
	found = false
	node := tree.root
	for node != NIL {
		compare := tree.comparator(key, node.key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			ceil, found = node, true
			node = node.left
		case compare > 0:
			node = node.right
		}
	}
	if found {
		return ceil, true
	}
	return nil, false
}

func (tree *Tree) insertFixup(node *Node) {
	defer func() {
		tree.root.color = black
	}()
	for node.parent != nil && node.parent.color == red {
		if node.parent == node.grandParent().left {
			if node.uncle().color == red {
				node.parent.color = black
				node.uncle().color = black
				node.grandParent().color = red
				node = node.grandParent()
			} else {
				if node == node.parent.right {
					node = node.parent
					tree.leftRotate(node)
				}
				node.parent.color = black
				node.grandParent().color = red
				tree.rightRotate(node.grandParent())
			}
		} else {
			if node.uncle().color == red {
				node.parent.color = black
				node.uncle().color = black
				node.grandParent().color = red
				node = node.grandParent()
			} else {
				if node.uncle().color == red {
					node = node.parent
					tree.rightRotate(node)
				}
				node.parent.color = black
				node.grandParent().color = red
				tree.leftRotate(node.grandParent())
			}
		}
	}
}

func (tree *Tree) lookup(key interface{}) *Node {
	node := tree.root
	for node != NIL {
		compare := tree.comparator(key, node.key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.left
		case compare > 0:
			node = node.right
		}
	}
	return NIL
}

//	 node                     x
//	/   \     左旋转         /  \
//
// T1   x   --------->   node   T3
//
//	 / \              /   \
//	T2 T3            T1   T2
func (tree *Tree) leftRotate(node *Node) {
	if node.right == NIL {
		return
	}
	right := node.right
	tree.replaceNode(node, right)
	node.right = right.left
	if right.left != NIL {
		right.left.parent = node
	}
	right.left = node
	node.parent = right
}

//	   node                    x
//	  /   \     右旋转       /  \
//	 x    T2   ------->   y   node
//	/ \                       /  \
//
// y  T1                     T1  T2
func (tree *Tree) rightRotate(node *Node) {
	if node.left == NIL {
		return
	}
	left := node.left
	tree.replaceNode(node, left)
	node.left = left.right
	if left.right != NIL {
		left.right.parent = node
	}
	left.right = node
	node.parent = left
}

func (tree *Tree) replaceNode(old *Node, new *Node) {
	if old.parent == NIL {
		tree.root = new
	} else {
		if old == old.parent.left {
			old.parent.left = new
		} else {
			old.parent.right = new
		}
	}
	if new != NIL {
		new.parent = old.parent
	}
}

func (tree *Tree) maximum(node *Node) *Node {
	current := node
	for current.right != NIL {
		current = current.right
	}
	return current
}

func (tree *Tree) minimum(node *Node) *Node {
	current := node
	for current.left != NIL {
		current = current.left
	}
	return current
}

func (tree *Tree) successor(node *Node) *Node {
	successor := node.right
	for successor.left != NIL {
		successor = successor.left
	}
	return successor
}

func (tree *Tree) delete(z *Node) *Node {
	if z == NIL || z == nil {
		return NIL
	}
	defer func() {
		tree.size--
	}()
	y := z
	x := NIL
	if z.left != NIL && z.right != NIL {
		y = tree.successor(z)
	}

	// 用x替换y
	if y.left != NIL {
		x = y.left
	} else {
		x = y.right
	}
	x.parent = y.parent
	if y.parent == NIL {
		tree.root = x
	} else if y.parent.left == y {
		y.parent.left = x
	} else if y.parent.right == y {
		y.parent.right = x
	}

	// 用y的数据填充z
	if y != z {
		z.key = y.key
		z.value = y.value
	}

	if y.color == black {
		// fixup
		tree.deleteFixup(x)
	}

	return y
}

func (tree *Tree) deleteFixup(x *Node) {
	for x != tree.root && x.color == black {
		if x == x.parent.left {
			w := x.brother()
			if w.color == red {
				w.color = black
				x.parent.color = red
				tree.leftRotate(x.parent)
				w = x.brother()
			}
			if w.left.color == black && w.right.color == black {
				w.color = red
				x = x.parent
			} else {
				if w.right.color == black {
					w.left.color = black
					w.color = red
					tree.rightRotate(w)
					w = x.brother()
				}
				w.color = x.parent.color
				x.parent.color = black
				w.right.color = black
				tree.leftRotate(x.parent)
				x = tree.root
			}
		} else {
			w := x.brother()
			if w.color == red {
				w.color = black
				x.parent.color = red
				tree.rightRotate(x.parent)
				w = x.brother()
			}
			if w.left.color == black && w.right.color == black {
				w.color = red
				x = x.parent
			} else {
				if w.left.color == black {
					w.right.color = black
					w.color = red
					tree.leftRotate(w)
					w = x.brother()
				}
				w.color = x.parent.color
				x.parent.color = black
				w.left.color = black
				tree.rightRotate(x.parent)
				x = tree.root
			}
		}
	}

	x.color = black
}

func output(node *Node, prefix string, isTail bool, str *string) {
	if node.right != NIL {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.left != NIL {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.left, newPrefix, true, str)
	}
}
