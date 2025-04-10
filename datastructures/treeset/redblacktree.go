package treeset

type color bool

const (
	black, red color = true, false
)

type Tree[T comparable] struct {
	Root       *Node[T]
	size       int
	Comparator Comparator[T]
}

type Node[T comparable] struct {
	Key    T
	color  color
	Left   *Node[T]
	Right  *Node[T]
	Parent *Node[T]
}

func New[T comparable](comparator Comparator[T]) *Tree[T] {
	return &Tree[T]{Comparator: comparator}
}

func (tree *Tree[T]) Put(key T) {
	var insertedNode *Node[T]
	if tree.Root == nil {
		tree.Root = &Node[T]{Key: key, color: red}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		for loop {
			compare := tree.Comparator(key, node.Key)
			switch {
			case compare == 0:
				node.Key = key
				return
			case compare < 0:
				if node.Left == nil {
					node.Left = &Node[T]{Key: key, color: red}
					insertedNode = node.Left
					loop = false
				} else {
					node = node.Left
				}
			case compare > 0:
				if node.Right == nil {
					node.Right = &Node[T]{Key: key, color: red}
					insertedNode = node.Right
					loop = false
				} else {
					node = node.Right
				}
			}
		}
		insertedNode.Parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++
}

func (tree *Tree[T]) Remove(key T) {
	var child *Node[T]
	node := tree.lookup(key)
	if node == nil {
		return
	}
	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		node.Key = pred.Key
		node = pred
	}
	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}
		if node.color == black {
			node.color = nodeColor(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = black
		}
	}
	tree.size--
}

func (tree *Tree[T]) Keys() []T {
	keys := make([]T, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

func (tree *Tree[T]) Left() *Node[T] {
	var parent *Node[T]
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

func (tree *Tree[T]) GetSize() int {
	return tree.size
}

func (tree *Tree[T]) GetRootKey() T {
	return tree.Root.Key
}

func (tree *Tree[T]) HeadSet(key T) []T {
	keys := make([]T, 0)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		if tree.Comparator(it.Key(), key) >= 0 {
			return keys
		}
		keys = append(keys, it.Key())
	}
	return keys
}

func (tree *Tree[T]) HeadSetSize(key T) int {
	size := 0
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		if tree.Comparator(it.Key(), key) >= 0 {
			return size
		}
		size++
	}
	return size
}

func (tree *Tree[T]) lookup(key T) *Node[T] {
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

func (node *Node[T]) grandparent() *Node[T] {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *Node[T]) uncle() *Node[T] {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *Node[T]) sibling() *Node[T] {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (tree *Tree[T]) rotateLeft(node *Node[T]) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *Tree[T]) rotateRight(node *Node[T]) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (tree *Tree[T]) replaceNode(old *Node[T], new *Node[T]) {
	if old.Parent == nil {
		tree.Root = new
	} else {
		if old == old.Parent.Left {
			old.Parent.Left = new
		} else {
			old.Parent.Right = new
		}
	}
	if new != nil {
		new.Parent = old.Parent
	}
}

func (tree *Tree[T]) insertCase1(node *Node[T]) {
	if node.Parent == nil {
		node.color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *Tree[T]) insertCase2(node *Node[T]) {
	if nodeColor(node.Parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *Tree[T]) insertCase3(node *Node[T]) {
	uncle := node.uncle()
	if nodeColor(uncle) == red {
		node.Parent.color = black
		uncle.color = black
		node.grandparent().color = red
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *Tree[T]) insertCase4(node *Node[T]) {
	grandparent := node.grandparent()
	if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.rotateLeft(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rotateRight(node.Parent)
		node = node.Right
	}
	tree.insertCase5(node)
}

func (tree *Tree[T]) insertCase5(node *Node[T]) {
	node.Parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *Node[T]) maximumNode() *Node[T] {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (tree *Tree[T]) deleteCase1(node *Node[T]) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *Tree[T]) deleteCase2(node *Node[T]) {
	sibling := node.sibling()
	if nodeColor(sibling) == red {
		node.Parent.color = red
		sibling.color = black
		if node == node.Parent.Left {
			tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *Tree[T]) deleteCase3(node *Node[T]) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == black &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *Tree[T]) deleteCase4(node *Node[T]) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		node.Parent.color = black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *Tree[T]) deleteCase5(node *Node[T]) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == red &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		sibling.Left.color = black
		tree.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Right) == red &&
		nodeColor(sibling.Left) == black {
		sibling.color = red
		sibling.Right.color = black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *Tree[T]) deleteCase6(node *Node[T]) {
	sibling := node.sibling()
	sibling.color = nodeColor(node.Parent)
	node.Parent.color = black
	if node == node.Parent.Left && nodeColor(sibling.Right) == red {
		sibling.Right.color = black
		tree.rotateLeft(node.Parent)
	} else if nodeColor(sibling.Left) == red {
		sibling.Left.color = black
		tree.rotateRight(node.Parent)
	}
}

func nodeColor[T comparable](node *Node[T]) color {
	if node == nil {
		return black
	}
	return node.color
}
