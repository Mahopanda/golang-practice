package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

// Key 表示 B+ 樹中的鍵類型，這裡為簡化起見使用整數型別
type Key int

// Value 表示存儲在 B+ 樹中的值
type Value struct {
	Data interface{} // 可存儲任何類型的數據
}

// Node 表示 B+ 樹中的一個節點
type Node struct {
	IsLeaf   bool     // 指示此節點是否為葉節點 (true 表示葉節點，false 表示內部節點)
	Keys     []Key    // 此節點內的已排序鍵
	Children []*Node  // 指向子節點的引用（僅在 IsLeaf 為 false 時使用）
	Values   []*Value // 與鍵對應的值（僅在 IsLeaf 為 true 時使用）
	Next     *Node    // 指向下一個葉節點，用於支持高效範圍查詢
}

// Record 表示一條完整記錄，包含多個欄位
type Record struct {
	ID   int
	Name string
	Age  int
}

// BPlusTree 表示 B+ 樹的結構

type BPlusTree struct {
	Root  *Node // 指向 B+ 樹的根節點
	Order int   // 每個節點可以容納的最大鍵數
}

// 定義多 B+ 樹來支持多個欄位的查詢
type Database struct {
	ByID   *BPlusTree // ID 索引樹
	ByName *BPlusTree // Name 索引樹
}

// 查詢數據庫，選擇適合的 B+ 樹
func (db *Database) QueryByID(id int) *Value {
	return db.ByID.Search(Key(id))
}

// hashKey 將字符串轉換為 Key，可以使用哈希函數
func hashKey(name string) Key {
	// 簡單的哈希函數示例，可以根據需求替換
	var hash int
	for _, char := range name {
		hash += int(char)
	}
	return Key(hash)
}

func (db *Database) QueryByName(name string) *Value {
	// 遍歷 Name 索引樹，找到匹配的記錄
	// 假設我們可以轉換 Name 為 Key（通過哈希或其他方式）
	return db.ByName.Search(hashKey(name))
}

// NewBPlusTree 初始化並返回具有指定階數的 B+ 樹
func NewBPlusTree(order int) *BPlusTree {
	return &BPlusTree{
		Root:  &Node{IsLeaf: true}, // 初始時，根節點為葉節點
		Order: order,
	}
}

// Insert adds a key-value pair to the B+ tree.
func (tree *BPlusTree) Insert(key Key, value Value) {
	root := tree.Root
	if len(root.Keys) == tree.Order {
		// If root is full, split it and create a new root.
		newRoot := &Node{IsLeaf: false}
		newRoot.Children = append(newRoot.Children, root)
		tree.Root = newRoot
		tree.splitChild(newRoot, 0)
	}
	tree.insertNonFull(tree.Root, key, value)
}

// insertNonFull 將鍵和值插入到非滿的節點中
func (tree *BPlusTree) insertNonFull(node *Node, key Key, value Value) {
	if node.IsLeaf {
		// 插入到葉節點中
		idx := 0
		for idx < len(node.Keys) && key > node.Keys[idx] {
			idx++
		}
		// 使用切片的 append 寫法來插入鍵和值
		node.Keys = append(node.Keys[:idx], append([]Key{key}, node.Keys[idx:]...)...)             // 將鍵按順序插入
		node.Values = append(node.Values[:idx], append([]*Value{&value}, node.Values[idx:]...)...) // 將值插入對應位置
		/*
			內層 append([]Key{key}, node.Keys[idx:]...)

				將新的鍵 key 放在切片 node.Keys[idx:] 的前面。
				這樣可以保持插入後的 Keys 切片仍然是有序的。

			外層 append(node.Keys[:idx], ... )

				將 key 插入後的新切片與 node.Keys 前 idx 個元素組合，形成一個完整的有序切片。
				這種用法允許我們在不創建額外臨時變量的情況下，將一個元素插入到切片的特定位置中。
				這種寫法的主要目的是簡潔地完成中間插入操作，同時保持切片的順序。

			slice Expansion:
				append 函數在處理切片時，如果目標切片容量不足，會自動擴展容量。
				這樣可以避免在每次插入時都重新分配大量內存，提高性能。
				example:
					slice := []int{1, 2, 3}
					slice = append(slice, []int{4, 5, 6}...)
					fmt.Println(slice) // 輸出: [1, 2, 3, 4, 5, 6]
		*/
	} else {
		// 插入到內部節點中
		idx := 0
		for idx < len(node.Keys) && key > node.Keys[idx] {
			idx++
		}
		// 如果子節點已滿，則需要先進行分裂
		if len(node.Children[idx].Keys) == tree.Order {
			tree.splitChild(node, idx) // 當子節點滿了時分裂
			if key > node.Keys[idx] {
				idx++
			}
		}
		// 遞歸地向下插入
		tree.insertNonFull(node.Children[idx], key, value)
	}
}

// splitChild 將滿載的子節點分裂為兩個節點，並更新父節點
func (tree *BPlusTree) splitChild(parent *Node, index int) {
	fullNode := parent.Children[index]        // 獲取要分裂的滿載子節點
	newNode := &Node{IsLeaf: fullNode.IsLeaf} // 創建一個新節點，用於存儲分裂後的一半數據
	mid := len(fullNode.Keys) / 2             // 獲取分裂點的索引，將節點一分為二

	if fullNode.IsLeaf {
		// 如果是葉子節點的分裂
		newNode.Keys = append(newNode.Keys, fullNode.Keys[mid:]...)       // 將 fullNode 的後半部分鍵移動到 newNode 中
		newNode.Values = append(newNode.Values, fullNode.Values[mid:]...) // 將 fullNode 的後半部分值移動到 newNode 中
		fullNode.Keys = fullNode.Keys[:mid]                               // 保留 fullNode 的前半部分鍵
		fullNode.Values = fullNode.Values[:mid]                           // 保留 fullNode 的前半部分值
		newNode.Next = fullNode.Next                                      // 更新 newNode 的 Next 指針，使其指向 fullNode 原本指向的下一個節點
		fullNode.Next = newNode                                           // 將 fullNode 的 Next 指針設置為 newNode，以便保持葉節點鏈接
	} else {
		// 如果是內部節點的分裂
		newNode.Keys = append(newNode.Keys, fullNode.Keys[mid+1:]...)             // 將 fullNode 的後半部分鍵移動到 newNode 中
		newNode.Children = append(newNode.Children, fullNode.Children[mid+1:]...) // 將 fullNode 的後半部分子節點移動到 newNode 中
		fullNode.Keys = fullNode.Keys[:mid]                                       // 保留 fullNode 的前半部分鍵
		fullNode.Children = fullNode.Children[:mid+1]                             // 保留 fullNode 的前半部分子節點
	}

	// 將分裂點的鍵插入到父節點中
	parent.Keys = append(parent.Keys[:index], append([]Key{fullNode.Keys[mid]}, parent.Keys[index:]...)...)
	// 將新節點插入到父節點的 Children 中
	parent.Children = append(parent.Children[:index+1], append([]*Node{newNode}, parent.Children[index+1:]...)...)
}

// Search finds the value associated with a key.
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

// searchNode traverses the tree to find the leaf node containing the key.
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

// Update modifies the value associated with a given key, if it exists.
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

// Delete removes a key-value pair from the tree.
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

// SaveTree serializes the B+ tree to a file.
func (tree *BPlusTree) SaveTree(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(tree)
}

// LoadTree deserializes a B+ tree from a file.
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

	// Insert data into the B+ tree.
	tree.Insert(1, Value{Data: "Alice"})
	tree.Insert(2, Value{Data: "Bob"})
	tree.Insert(3, Value{Data: "Charlie"})

	// Search for a key in the tree.
	result := tree.Search(2)
	if result != nil {
		fmt.Printf("Search Result for key 2: %v\n", result.Data)
	} else {
		fmt.Println("Key 2 not found.")
	}

	// Update an existing key in the tree.
	success := tree.Update(2, Value{Data: "Bob Updated"})
	if success {
		fmt.Println("Update successful for key 2.")
	}

	// Delete a key from the tree.
	deleted := tree.Delete(2)
	if deleted {
		fmt.Println("Delete successful for key 2.")
	}

	// Save the tree to a file.
	err := tree.SaveTree("bptree_data.gob")
	if err != nil {
		fmt.Println("Error saving B+ tree:", err)
	} else {
		fmt.Println("B+ tree saved successfully.")
	}

	// Load the tree from a file.
	loadedTree, err := LoadTree("bptree_data.gob")
	if err != nil {
		fmt.Println("Error loading B+ tree:", err)
	} else {
		fmt.Println("B+ tree loaded successfully.")
	}

	// Search in the loaded tree.
	result2 := loadedTree.Search(2)
	if result2 != nil {
		fmt.Printf("Search Result for key 2 in loaded tree: %v\n", result2.Data)
	} else {
		fmt.Println("Key 2 not found in loaded tree.")
	}

	// 初始化數據庫
	db := &Database{
		ByID:   NewBPlusTree(3),
		ByName: NewBPlusTree(3),
	}

	// 插入數據
	db.ByID.Insert(1, Value{Data: Record{ID: 1, Name: "Alice", Age: 25}})
	db.ByID.Insert(2, Value{Data: Record{ID: 2, Name: "Bob", Age: 30}})
	db.ByName.Insert(hashKey("Alice"), Value{Data: Record{ID: 1, Name: "Alice", Age: 25}})
	db.ByName.Insert(hashKey("Bob"), Value{Data: Record{ID: 2, Name: "Bob", Age: 30}})

	// 查詢 ID 為 1 的記錄
	result3 := db.QueryByID(1)
	if result3 != nil {
		fmt.Printf("QueryByID Result: ID=%d, Name=%v, Age=%d\n", result3.Data.(Record).ID, result3.Data.(Record).Name, result3.Data.(Record).Age)
	}

	// 查詢 Name 為 "Alice" 的記錄
	result4 := db.QueryByName("Alice")
	if result4 != nil {
		fmt.Printf("QueryByName Result: ID=%d, Name=%v, Age=%d\n", result4.Data.(Record).ID, result4.Data.(Record).Name, result4.Data.(Record).Age)
	}

}
