package bplustree

import (
	"github.com/Mahopanda/mini-project/bplustree/models"
)

const MaxKeysPerNode = 4

// Node 表示 B+ 樹中的一個節點
type Node struct {
	IsLeaf   bool            // 指示此節點是否為葉節點 (true 表示葉節點，false 表示內部節點)
	Keys     []models.Key    // 此節點內的已排序鍵
	Children []*Node         // 指向子節點的引用（僅在 IsLeaf 為 false 時使用）
	Values   []*models.Value // 與鍵對應的值（僅在 IsLeaf 為 true 時使用）
	Next     *Node           // 指向下一個葉節點，用於支持高效範圍查詢
}

// insertNonFull 將鍵和值插入到非滿的節點中
func (tree *BPlusTree) insertNonFull(node *Node, key models.Key, value models.Value) {
	if node.IsLeaf {
		// 插入到葉節點中
		idx := 0
		for idx < len(node.Keys) && key > node.Keys[idx] {
			idx++
		}
		// 使用切片的 append 寫法來插入鍵和值
		node.Keys = append(node.Keys[:idx], append([]models.Key{key}, node.Keys[idx:]...)...)             // 將鍵按順序插入
		node.Values = append(node.Values[:idx], append([]*models.Value{&value}, node.Values[idx:]...)...) // 將值插入對應位置
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

func (tree *BPlusTree) splitChild(parent *Node, index int) {
	fullNode := parent.Children[index]        // 獲取要分裂的滿載子節點
	newNode := &Node{IsLeaf: fullNode.IsLeaf} // 創建新節點，用於存儲分裂後的一半數據
	mid := len(fullNode.Keys) / 2             // 獲取分裂點的索引

	// 初始化和分配葉節點或內部節點的鍵值
	if fullNode.IsLeaf {
		// 分裂葉節點
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

	// 更新父節點的鍵值和子節點
	parent.Keys = append(parent.Keys[:index], append([]models.Key{fullNode.Keys[mid]}, parent.Keys[index:]...)...)
	parent.Children = append(parent.Children[:index+1], append([]*Node{newNode}, parent.Children[index+1:]...)...)
}

// searchNode traverses the tree to find the leaf node containing the key.
func (tree *BPlusTree) searchNode(node *Node, key models.Key) *Node {
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
