package main

import (
	"fmt"
	"strings"
)

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

type BinarySearchTree struct {
	Root *Node
	Len  int
}

func (b *BinarySearchTree) String() string {
	sb := strings.Builder{}
	b.inOrderTraversal(&sb)
	return b.printTree(b.Root, "", true)
}

func (b *BinarySearchTree) printTree(node *Node, prefix string, isLeft bool) string {
	if node == nil {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(prefix)
	if isLeft {
		sb.WriteString("├── ")
	} else {
		sb.WriteString("└── ")
	}
	sb.WriteString(fmt.Sprintf("%d\n", node.Value))

	childPrefix := prefix + (map[bool]string{true: "│   ", false: "    "})[isLeft]
	if node.Left != nil {
		sb.WriteString(b.printTree(node.Left, childPrefix, true))
	}
	if node.Right != nil {
		sb.WriteString(b.printTree(node.Right, childPrefix, false))
	}

	return sb.String()
}

func (b *BinarySearchTree) inOrderTraversal(sb *strings.Builder) {
	b.inOrderTraversalByNode(b.Root, sb)
}

func (b *BinarySearchTree) inOrderTraversalByNode(root *Node, sb *strings.Builder) {
	if root == nil {
		return
	}

	b.inOrderTraversalByNode(root.Left, sb)
	sb.WriteString(fmt.Sprintf("%d ", root.Value))
	b.inOrderTraversalByNode(root.Right, sb)
}

func (b *BinarySearchTree) addNode(value int) {
	b.Root = b.addNodeByNode(b.Root, value)
	b.Len++
}

func (b *BinarySearchTree) addNodeByNode(root *Node, value int) *Node {
	if root == nil {
		return &Node{Value: value}
	}

	if value < root.Value {
		root.Left = b.addNodeByNode(root.Left, value)
	} else {
		root.Right = b.addNodeByNode(root.Right, value)
	}

	return root
}

func (b *BinarySearchTree) search(value int) (*Node, bool) {
	return b.searchByNode(b.Root, value)
}

func (b *BinarySearchTree) searchByNode(root *Node, value int) (*Node, bool) {
	if root == nil {
		return nil, false
	}

	if root.Value == value {
		return root, true
	}

	if value < root.Value {
		return b.searchByNode(root.Left, value)
	}

	return b.searchByNode(root.Right, value)
}

func main() {
	n := &Node{Value: 2, Left: nil, Right: nil}
	n.Left = &Node{Value: 1, Left: nil, Right: nil}
	n.Right = &Node{Value: 3, Left: nil, Right: nil}

	bst := &BinarySearchTree{Root: n, Len: 3}
	fmt.Println("Initial tree structure:")

	fmt.Println(bst)

	bst.addNode(4)
	bst.addNode(5)
	bst.addNode(6)
	bst.addNode(7)
	bst.addNode(8)

	fmt.Println("After adding nodes:")
	fmt.Println(bst)

	node, found := bst.search(4)
	if found {
		fmt.Printf("搜索 4: 找到了，節點值為 %d\n", node.Value)
	} else {
		fmt.Println("搜索 4: 未找到")
	}

}
