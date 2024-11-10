package database

import (
	"hash/fnv"

	"github.com/Mahopanda/mini-project/bplustree"
	"github.com/Mahopanda/mini-project/bplustree/models"
	"github.com/Mahopanda/mini-project/bplustree/parser"
	"github.com/Mahopanda/mini-project/bplustree/types"
)

// 定義多 B+ 樹來支持多個欄位的查詢
type Database struct {
	ByID   *bplustree.BPlusTree // ID 索引樹
	ByName *bplustree.BPlusTree // Name 索引樹
}

type Tables struct {
	Tables map[string]*bplustree.BPlusTree
}

func NewTables() *Tables {
	return &Tables{
		Tables: make(map[string]*bplustree.BPlusTree),
	}
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

// ExecuteQuery 執行SQL查詢
func (db *Database) ExecuteQuery(query string) (*models.Value, error) {
	lexer := parser.NewLexer(query)
	tokens := lexer.Run()

	p := parser.NewParser(tokens)
	stmt, err := p.Parse()
	if err != nil {
		return nil, err
	}

	executor := parser.NewQueryExecutor(db)
	return executor.Execute(stmt)
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
	return db.ByName.Search(models.Key(HashKey(name)))
}

func (db *Database) QueryByID(id int) *models.Value {
	return db.ByID.Search(models.Key(id))
}

func (db *Database) RangeQueryByID(minID, maxID int) []*models.Value {
	return db.ByID.RangeQuery(models.Key(minID), models.Key(maxID))
}

func (db *Database) RangeQueryByName(minName, maxName string) []*models.Value {
	minHash := models.Key(HashKey(minName))
	maxHash := models.Key(HashKey(maxName))
	return db.ByName.RangeQuery(minHash, maxHash)
}
