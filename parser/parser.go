package parser

import (
	"fmt"
	"strings"
)

type StatementType int

const (
	StatementSelect StatementType = iota
	StatementInsert
	StatementUpdate
	StatementDelete
	StatementCreateTable
	StatementDropTable
)

// 簡化的 Statement 接口和結構體
type Statement interface {
	GetType() StatementType
}

type Select struct {
	Table   string
	Columns []string
	Where   string // 簡單條件處理
}

type Insert struct {
	Table   string
	Columns []string
	Values  []string
}

type Update struct {
	Table   string
	Columns []string
	Values  []string
}

type Delete struct {
	Table string
	Where string
}

type CreateTable struct {
	Table   string
	Columns []string
}

type DropTable struct {
	Table string
}

type AlterTable struct {
	Table   string
	Columns []string
}

type ColumnDefinition struct {
	Table   string
	Columns []string
}

type DropColumn struct {
	Table   string
	Columns []string
}

type AddColumn struct {
	Table   string
	Columns []string
}

type ModifyColumn struct {
	Table   string
	Columns []string
}


func (*Select) GetType() StatementType { return StatementSelect }

// GetType returns the statement type.
func (*Insert) GetType() StatementType { return StatementInsert }

// GetType returns the statement type.
func (*Update) GetType() StatementType { return StatementUpdate }

// GetType returns the statement type.
func (*Delete) GetType() StatementType { return StatementDelete }

// GetType returns the statement type.
func (*CreateTable) GetType() StatementType { return StatementCreateTable }

// GetType returns the statement type.
func (*DropTable) GetType() StatementType { return StatementDropTable }

// 簡化的 Parse 函數
func Parse(input string) (Statement, error) {
	tokens := strings.Fields(input)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	switch strings.ToUpper(tokens[0]) {
	case "SELECT":
		return parseSelect(tokens)
	case "INSERT":
		// 簡化版的 INSERT 解析邏輯
	case "CREATE":
		// 簡化版的 CREATE TABLE 解析邏輯
	case "DROP":
		// 簡化版的 DROP TABLE 解析邏輯
	default:
		return nil, fmt.Errorf("unsupported statement type")
	}
	return nil, fmt.Errorf("parsing not implemented")
}

func parseSelect(tokens []string) (*Select, error) {
	if len(tokens) < 4 || strings.ToUpper(tokens[2]) != "FROM" {
		return nil, fmt.Errorf("invalid SELECT syntax")
	}
	return &Select{
		Columns: []string{tokens[1]},
		Table:   tokens[3],
		Where:   "", // 可根據需要擴展 WHERE 子句
	}, nil
}
