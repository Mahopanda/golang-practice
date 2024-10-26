package database

import (
	"hash/fnv"

	"github.com/Mahopanda/mini-project/bplustree"
	"github.com/Mahopanda/mini-project/bplustree/models"
)

// 定義多 B+ 樹來支持多個欄位的查詢
type Database struct {
	ByID   *bplustree.BPlusTree // ID 索引樹
	ByName *bplustree.BPlusTree // Name 索引樹
}

// Record 表示一條完整記錄，包含多個欄位
type Record struct {
	ID   int
	Name string
	Age  int
}

// hashKey 將字符串轉換為uint64類型的哈希值
func HashKey(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func (db *Database) QueryByName(name string) *models.Value {
	// 遍歷 Name 索引樹，找到匹配的記錄
	// 假設我們可以轉換 Name 為 Key（通過哈希或其他方式）
	return db.ByName.Search(models.Key(HashKey(name)))
}

// 查詢數據庫，選擇適合的 B+ 樹
func (db *Database) QueryByID(id int) *models.Value {
	return db.ByID.Search(models.Key(id))
}
