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

type Tables struct {
	Tables map[string]*bplustree.BPlusTree
}

func (db *Tables) AddTable(name string) {
	db.Tables[name] = bplustree.NewBPlusTree(3) // 假設默認階數為 3
}

func (db *Tables) RemoveTable(name string) {
	delete(db.Tables, name)
}

func (db *Tables) GetTable(name string) *bplustree.BPlusTree {
	return db.Tables[name]
}

func (db *Database) ExecuteQuery(query string) (*models.Value, error) {
	statement, err := parser.Parse(query)
	if err != nil {
		return nil, err
	}

	executor := parser.QueryExecutor{DB: db}
	return executor.Execute(statement)
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

// RangeQueryByID 查詢給定範圍內的記錄（根據 ID）
func (db *Database) RangeQueryByID(minID, maxID int) []*models.Value {
	return db.ByID.RangeQuery(models.Key(minID), models.Key(maxID))
}

// RangeQueryByName 查詢給定範圍內的記錄（根據 Name hash）
func (db *Database) RangeQueryByName(minName, maxName string) []*models.Value {
	minHash := models.Key(HashKey(minName))
	maxHash := models.Key(HashKey(maxName))
	return db.ByName.RangeQuery(minHash, maxHash)
}
