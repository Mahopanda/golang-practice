package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

type Key int
type Value struct {
	Data interface{}
}

type Node struct {
	IsLeaf   bool
	Keys     []Key
	Children []*Node
	Values   []*Value
	Next     *Node
}

type BPlusTree struct {
	Root  *Node
	Order int
}

func NewBPlusTree(order int) *BPlusTree {
	return &BPlusTree{
		Root:  &Node{IsLeaf: true},
		Order: order,
	}
}

func (tree *BPlusTree) Insert(key Key, value Value) {
	root := tree.Root
	if len(root.Keys) == tree.Order {
		// 根節點滿時，創建新根
		newRoot := &Node{IsLeaf: false}
		newRoot.Children = append(newRoot.Children, root)
		tree.Root = newRoot
		tree.splitChild(newRoot, 0)
	}
	tree.insertNonFull(tree.Root, key, value)
}

func (tree *BPlusTree) insertNonFull(node *Node, key Key, value Value) {
	if node.IsLeaf {
		// 插入葉節點
		idx := 0
		for idx < len(node.Keys) && key > node.Keys[idx] {
			idx++
		}
		node.Keys = append(node.Keys[:idx], append([]Key{key}, node.Keys[idx:]...)...)
		node.Values = append(node.Values[:idx], append([]*Value{&value}, node.Values[idx:]...)...)
	} else {
		// 插入內部節點
		idx := 0
		for idx < len(node.Keys) && key > node.Keys[idx] {
			idx++
		}
		if len(node.Children[idx].Keys) == tree.Order {
			tree.splitChild(node, idx)
			if key > node.Keys[idx] {
				idx++
			}
		}
		tree.insertNonFull(node.Children[idx], key, value)
	}
}

func (tree *BPlusTree) splitChild(parent *Node, index int) {
	fullNode := parent.Children[index]
	newNode := &Node{IsLeaf: fullNode.IsLeaf}
	mid := len(fullNode.Keys) / 2

	if fullNode.IsLeaf {
		// 分裂葉子節點
		newNode.Keys = append(newNode.Keys, fullNode.Keys[mid:]...)
		newNode.Values = append(newNode.Values, fullNode.Values[mid:]...)
		fullNode.Keys = fullNode.Keys[:mid]
		fullNode.Values = fullNode.Values[:mid]
		newNode.Next = fullNode.Next
		fullNode.Next = newNode
	} else {
		// 分裂內部節點
		newNode.Keys = append(newNode.Keys, fullNode.Keys[mid+1:]...)
		newNode.Children = append(newNode.Children, fullNode.Children[mid+1:]...)
		fullNode.Keys = fullNode.Keys[:mid]
		fullNode.Children = fullNode.Children[:mid+1]
	}

	// 在父節點中插入分裂點
	parent.Keys = append(parent.Keys[:index], append([]Key{fullNode.Keys[mid]}, parent.Keys[index:]...)...)
	parent.Children = append(parent.Children[:index+1], append([]*Node{newNode}, parent.Children[index+1:]...)...)
}

func (tree *BPlusTree) Search(key Key) *Value {
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

func (tree *BPlusTree) searchNode(node *Node, key Key) *Node {
	idx := 0
	for idx < len(node.Keys) && key > node.Keys[idx] {
		idx++
	}
	if node.IsLeaf {
		return node
	}
	if idx < len(node.Keys) && key == node.Keys[idx] {
		idx++
	}
	return tree.searchNode(node.Children[idx], key)
}

func (tree *BPlusTree) Update(key Key, newValue Value) bool {
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

func (tree *BPlusTree) Delete(key Key) bool {
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

// SaveTree 將 B+ 樹存儲到文件
func (tree *BPlusTree) SaveTree(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(tree)
}

// LoadTree 從文件加載 B+ 樹
func LoadTree(filename string) (*BPlusTree, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var tree BPlusTree
	if err := decoder.Decode(&tree); err != nil {
		return nil, err
	}
	return &tree, nil
}

func main() {
	tree := NewBPlusTree(3)

	// 插入數據
	tree.Insert(1, Value{Data: "Alice"})
	tree.Insert(2, Value{Data: "Bob"})
	tree.Insert(3, Value{Data: "Charlie"})

	// 查詢數據
	result := tree.Search(2)
	if result != nil {
		fmt.Printf("Search Result for key 2: %v\n", result.Data)
	} else {
		fmt.Println("Key 2 not found.")
	}

	// 更新數據
	success := tree.Update(2, Value{Data: "Bob Updated"})
	if success {
		fmt.Println("Update successful for key 2.")
	}

	// 刪除數據
	deleted := tree.Delete(2)
	if deleted {
		fmt.Println("Delete successful for key 2.")
	}

	// 保存樹到文件
	err := tree.SaveTree("bptree_data.gob")
	if err != nil {
		fmt.Println("Error saving B+ tree:", err)
	} else {
		fmt.Println("B+ tree saved successfully.")
	}

	// 從文件加載樹
	loadedTree, err := LoadTree("bptree_data.gob")
	if err != nil {
		fmt.Println("Error loading B+ tree:", err)
	} else {
		fmt.Println("B+ tree loaded successfully.")
	}

	// 查詢加載的樹中的數據
	result2 := loadedTree.Search(2)
	if result2 != nil {
		fmt.Printf("Search Result for key 2 in loaded tree: %v\n", result.Data)
	} else {
		fmt.Println("Key 2 not found in loaded tree.")
	}
}
