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
		ByID:   bplustree.NewBPlusTree(10),
		ByName: bplustree.NewBPlusTree(10),
	}

	// 插入數據
	db.ByID.Insert(models.Key(1), models.Value{Data: database.Record{ID: 1, Name: "Alice", Age: 25}})
	db.ByID.Insert(models.Key(2), models.Value{Data: database.Record{ID: 2, Name: "Bob", Age: 30}})
	db.ByID.Insert(models.Key(3), models.Value{Data: database.Record{ID: 3, Name: "Aken", Age: 35}})
	db.ByID.Insert(models.Key(4), models.Value{Data: database.Record{ID: 4, Name: "Banana", Age: 35}})
	db.ByID.Insert(models.Key(5), models.Value{Data: database.Record{ID: 5, Name: "Charlie", Age: 35}})
	db.ByID.Insert(models.Key(6), models.Value{Data: database.Record{ID: 6, Name: "Dan", Age: 35}})
	db.ByID.Insert(models.Key(7), models.Value{Data: database.Record{ID: 7, Name: "Eric", Age: 35}})
	db.ByID.Insert(models.Key(8), models.Value{Data: database.Record{ID: 8, Name: "Frank", Age: 35}})
	db.ByID.Insert(models.Key(9), models.Value{Data: database.Record{ID: 9, Name: "Grace", Age: 35}})
	db.ByID.Insert(models.Key(10), models.Value{Data: database.Record{ID: 10, Name: "Helen", Age: 35}})

	db.ByName.Insert(models.Key(database.HashKey("Alice")), models.Value{Data: database.Record{ID: 1, Name: "Alice", Age: 25}})
	db.ByName.Insert(models.Key(database.HashKey("Bob")), models.Value{Data: database.Record{ID: 2, Name: "Bob", Age: 30}})
	db.ByName.Insert(models.Key(database.HashKey("Aken")), models.Value{Data: database.Record{ID: 3, Name: "Aken", Age: 35}})
	db.ByName.Insert(models.Key(database.HashKey("Banana")), models.Value{Data: database.Record{ID: 4, Name: "Banana", Age: 35}})
	db.ByName.Insert(models.Key(database.HashKey("Charlie")), models.Value{Data: database.Record{ID: 5, Name: "Charlie", Age: 35}})
	db.ByName.Insert(models.Key(database.HashKey("Dan")), models.Value{Data: database.Record{ID: 6, Name: "Dan", Age: 35}})
	db.ByName.Insert(models.Key(database.HashKey("Eric")), models.Value{Data: database.Record{ID: 7, Name: "Eric", Age: 35}})
	db.ByName.Insert(models.Key(database.HashKey("Frank")), models.Value{Data: database.Record{ID: 8, Name: "Frank", Age: 35}})
	db.ByName.Insert(models.Key(database.HashKey("Grace")), models.Value{Data: database.Record{ID: 9, Name: "Grace", Age: 35}})
	db.ByName.Insert(models.Key(database.HashKey("Helen")), models.Value{Data: database.Record{ID: 10, Name: "Helen", Age: 35}})
	// 印出哈希值
	fmt.Println(database.HashKey("Alice"))
	fmt.Println(database.HashKey("Bob"))
	fmt.Println(database.HashKey("Aken"))
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

	// // 查詢 ID 在範圍 2 到 7 的所有記錄
	idResults := db.RangeQueryByID(2, 7)
	fmt.Println("ID 範圍 2 到 7 的記錄")
	for _, result := range idResults {
		if result != nil {
			record := result.Data.(database.Record)
			fmt.Printf("ID=%d, Name=%v, Age=%d\n", record.ID, record.Name, record.Age)
		}
	}

}
