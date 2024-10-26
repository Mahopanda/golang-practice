package main

import (
	"fmt"
	"strings"
)

type Node struct {
	Value  int
	Left   *Node
	Right  *Node
	Height int
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

	// 根據節點位置選擇適當的符號
	if isLeft {
		sb.WriteString(prefix + "├── ")
	} else {
		sb.WriteString(prefix + "└── ")
	}
	sb.WriteString(fmt.Sprintf("%d\n", node.Value))

	// 計算新的前綴
	newPrefix := prefix + (map[bool]string{true: "│   ", false: "    "})[isLeft]

	// 遞迴處理左右子節點
	if node.Left != nil || node.Right != nil {
		if node.Left != nil {
			sb.WriteString(b.printTree(node.Left, newPrefix, true))
		} else {
			sb.WriteString(newPrefix + "├── (null)\n") // 顯示空節點位置
		}
		if node.Right != nil {
			sb.WriteString(b.printTree(node.Right, newPrefix, false))
		} else {
			sb.WriteString(newPrefix + "└── (null)\n") // 顯示空節點位置
		}
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
	b.Root = b.addNodeBalanced(b.Root, value)
	b.Len++
}

func (b *BinarySearchTree) addNodeBalanced(node *Node, value int) *Node {
	if node == nil {
		return &Node{Value: value, Height: 1}
	}

	if value < node.Value {
		node.Left = b.addNodeBalanced(node.Left, value)
	} else {
		node.Right = b.addNodeBalanced(node.Right, value)
	}

	node.Height = 1 + max(height(node.Left), height(node.Right))

	balance := getBalance(node)

	// 左左情況
	if balance > 1 && value < node.Left.Value {
		return rightRotate(node)
	}

	// 右右情況
	if balance < -1 && value > node.Right.Value {
		return leftRotate(node)
	}

	// 左右情況
	if balance > 1 && value > node.Left.Value {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// 右左情況
	if balance < -1 && value < node.Right.Value {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	return node
}

func height(node *Node) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func getBalance(node *Node) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

func rightRotate(y *Node) *Node {
	x := y.Left
	T2 := x.Right

	x.Right = y
	y.Left = T2

	y.Height = 1 + max(height(y.Left), height(y.Right))
	x.Height = 1 + max(height(x.Left), height(x.Right))

	return x
}

func leftRotate(x *Node) *Node {
	y := x.Right
	T2 := y.Left

	y.Left = x
	x.Right = T2

	x.Height = 1 + max(height(x.Left), height(x.Right))
	y.Height = 1 + max(height(y.Left), height(y.Right))

	return y
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
	bst.addNode(9)
	fmt.Println("After adding nodes:")
	fmt.Println(bst)

	node, found := bst.search(4)
	if found {
		fmt.Printf("搜索 4: 找到了，節點值為 %d\n", node.Value)
	} else {
		fmt.Println("搜索 4: 未找到")
	}

}
