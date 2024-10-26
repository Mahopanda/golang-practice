package main

import (
	"fmt"

	"github.com/Mahopanda/mini-project/bplustree"
	"github.com/Mahopanda/mini-project/bplustree/database"
	"github.com/Mahopanda/mini-project/bplustree/models"
)

func main() {
	tree := bplustree.NewBPlusTree(3)

	// Insert data into the B+ tree.
	tree.Insert(models.Key(1), models.Value{Data: "Alice"})
	tree.Insert(models.Key(2), models.Value{Data: "Bob"})
	tree.Insert(models.Key(3), models.Value{Data: "Charlie"})

	// Search for a key in the tree.
	result := tree.Search(models.Key(2))
	if result != nil {
		fmt.Printf("Search Result for key 2: %v\n", result.Data)
	} else {
		fmt.Println("Key 2 not found.")
	}

	// Update an existing key in the tree.
	success := tree.Update(models.Key(2), models.Value{Data: "Bob Updated"})
	if success {
		fmt.Println("Update successful for key 2.")
	}

	// Delete a key from the tree.
	deleted := tree.Delete(models.Key(2))
	if deleted {
		fmt.Println("Delete successful for key 2.")
	}

	// 初始化數據庫
	db := &database.Database{
		ByID:   bplustree.NewBPlusTree(3),
		ByName: bplustree.NewBPlusTree(3),
	}

	// 插入數據
	db.ByID.Insert(models.Key(1), models.Value{Data: database.Record{ID: 1, Name: "Alice", Age: 25}})
	db.ByID.Insert(models.Key(2), models.Value{Data: database.Record{ID: 2, Name: "Bob", Age: 30}})
	db.ByName.Insert(models.Key(database.HashKey("Alice")), models.Value{Data: database.Record{ID: 1, Name: "Alice", Age: 25}})
	db.ByName.Insert(models.Key(database.HashKey("Bob")), models.Value{Data: database.Record{ID: 2, Name: "Bob", Age: 30}})

	// 查詢 ID 為 1 的記錄
	result3 := db.QueryByID(1)
	if result3 != nil {
		record := result3.Data.(database.Record)
		fmt.Printf("QueryByID Result: ID=%d, Name=%v, Age=%d\n", record.ID, record.Name, record.Age)
	}

	// 查詢 Name 為 "Alice" 的記錄
	result4 := db.QueryByName("Alice")
	if result4 != nil {
		record := result4.Data.(database.Record)
		fmt.Printf("QueryByName Result: ID=%d, Name=%v, Age=%d\n", record.ID, record.Name, record.Age)
	}
}
