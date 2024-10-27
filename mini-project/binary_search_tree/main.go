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

func (b *BinarySearchTree) deleteNode(value int) {
	b.Root = b.deleteNodeByNode(b.Root, value)
	b.Len--
}

func (b *BinarySearchTree) deleteNodeByNode(node *Node, value int) *Node {
	if node == nil {
		return nil
	}

	// 樹中查找要刪除的節點
	if value < node.Value {
		node.Left = b.deleteNodeByNode(node.Left, value)
	} else if value > node.Value {
		node.Right = b.deleteNodeByNode(node.Right, value)
	} else {
		// 找到要刪除的節點
		if node.Left == nil {
			return node.Right
		} else if node.Right == nil {
			return node.Left
		}

		// 情況 3：節點有兩個子節點，找到右子樹的最小值節點
		minRight := minValueNode(node.Right)
		node.Value = minRight.Value
		node.Right = b.deleteNodeByNode(node.Right, minRight.Value)
	}

	// 更新節點高度
	node.Height = 1 + max(height(node.Left), height(node.Right))

	// 檢查並恢復 AVL 樹的平衡
	balance := getBalance(node)

	// 左左情況
	if balance > 1 && getBalance(node.Left) >= 0 {
		return rightRotate(node)
	}

	// 左右情況
	if balance > 1 && getBalance(node.Left) < 0 {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// 右右情況
	if balance < -1 && getBalance(node.Right) <= 0 {
		return leftRotate(node)
	}

	// 右左情況
	if balance < -1 && getBalance(node.Right) > 0 {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	return node
}

// 尋找最小值節點
func minValueNode(node *Node) *Node {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
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

// RangeQuery 查詢指定範圍內的所有節點值
func (b *BinarySearchTree) RangeQuery(lower, upper int) []int {
	var result []int
	b.rangeQueryByNode(b.Root, lower, upper, &result)
	return result
}

// rangeQueryByNode 遞迴查詢每個節點，找出範圍內的值
func (b *BinarySearchTree) rangeQueryByNode(node *Node, lower, upper int, result *[]int) {
	if node == nil {
		return
	}

	// 範圍查詢：若當前節點的值在範圍內，則將其加入結果
	if node.Value >= lower && node.Value <= upper {
		*result = append(*result, node.Value)
	}

	// 若當前節點的值大於下限，則向左子樹查找
	if node.Value > lower {
		b.rangeQueryByNode(node.Left, lower, upper, result)
	}

	// 若當前節點的值小於上限，則向右子樹查找
	if node.Value < upper {
		b.rangeQueryByNode(node.Right, lower, upper, result)
	}
}

func main() {
	n := &Node{Value: 2, Left: nil, Right: nil}
	n.Left = &Node{Value: 1, Left: nil, Right: nil}
	n.Right = &Node{Value: 3, Left: nil, Right: nil}

	bst := &BinarySearchTree{Root: n, Len: 3}
	fmt.Println("Initial tree structure:")

	fmt.Println(bst)

	fmt.Println("--------------------------------")

	bst.addNode(4)
	bst.addNode(5)
	bst.addNode(6)
	bst.addNode(7)
	bst.addNode(8)
	bst.addNode(9)
	fmt.Println("新增節點 4, 5, 6, 7, 8, 9")
	fmt.Println(bst)

	fmt.Println("--------------------------------")

	bst.addNode(12)
	bst.addNode(18)
	bst.addNode(15)
	bst.addNode(10)

	fmt.Println("範圍查詢 6 到 15:")
	result := bst.RangeQuery(6, 15)
	fmt.Println(result)

	fmt.Println("--------------------------------")
	node, found := bst.search(4)
	if found {
		fmt.Printf("搜索 4: 找到了，節點值為 %d\n", node.Value)
	} else {
		fmt.Println("搜索 4: 未找到")
	}

	fmt.Println("--------------------------------")

	bst.deleteNode(5)
	fmt.Println("刪除節點 5")
	fmt.Println(bst)

}
