package bplustree

import (
	"github.com/Mahopanda/mini-project/bplustree/models"
)

type BPlusTree struct {
	Root  *Node // 指向 B+ 樹的根節點
	Order int   // 每個節點可以容納的最大鍵數
}

// NewBPlusTree 初始化並返回具有指定階數的 B+ 樹
func NewBPlusTree(order int) *BPlusTree {
	return &BPlusTree{
		Root:  &Node{IsLeaf: true}, // 初始時，根節點為葉節點
		Order: order,
	}
}

func (tree *BPlusTree) Insert(key models.Key, value models.Value) {
	root := tree.Root
	if len(root.Keys) == tree.Order {
		// If root is full, split it and create a new root.
		newRoot := &Node{IsLeaf: false}
		newRoot.Children = append(newRoot.Children, root) // 初始化子節點
		tree.Root = newRoot
		tree.splitChild(newRoot, 0) // 在新根節點中分裂
	}
	tree.insertNonFull(tree.Root, key, value)
}


// Search finds the value associated with a key.
func (tree *BPlusTree) Search(key models.Key) *models.Value {
	node := tree.searchNode(tree.Root, key)
	if node != nil {
		for i, k := range node.Keys {
			if k == key {
				return node.Values[i]
			}
		}
	}
	return nil
}

// Update modifies the value associated with a given key, if it exists.
func (tree *BPlusTree) Update(key models.Key, newValue models.Value) bool {
	node := tree.searchNode(tree.Root, key)
	if node != nil {
		for i, k := range node.Keys {
			if k == key {
				node.Values[i] = &newValue
				return true
			}
		}
	}
	return false
}

// Delete removes a key-value pair from the tree.
func (tree *BPlusTree) Delete(key models.Key) bool {
	node := tree.searchNode(tree.Root, key)
	if node != nil {
		for i, k := range node.Keys {
			if k == key {
				node.Keys = append(node.Keys[:i], node.Keys[i+1:]...)
				node.Values = append(node.Values[:i], node.Values[i+1:]...)
				return true
			}
		}
	}
	return false
}

// RangeQuery 查詢範圍內的所有鍵值對
func (tree *BPlusTree) RangeQuery(minKey, maxKey models.Key) []*models.Value {
	var result []*models.Value
	node := tree.searchRangeStartNode(tree.Root, minKey)
	if node == nil {
		return result
	}

	// 從找到的節點開始，遍歷葉節點，收集範圍內的值
	for node != nil {
		for i, key := range node.Keys {
			if key > maxKey {
				return result
			}
			if key >= minKey {
				result = append(result, node.Values[i])
			}
		}
		node = node.Next
	}
	return result
}

// searchRangeStartNode 查找範圍查詢的起始葉節點
func (tree *BPlusTree) searchRangeStartNode(node *Node, minKey models.Key) *Node {
	idx := 0
	for idx < len(node.Keys) && minKey > node.Keys[idx] {
		idx++
	}
	if node.IsLeaf {
		return node
	}
	if idx < len(node.Keys) && minKey == node.Keys[idx] {
		idx++
	}
	return tree.searchRangeStartNode(node.Children[idx], minKey)
}
