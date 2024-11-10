// parser_test.go
package parser

import (
	"testing"

	"github.com/Mahopanda/mini-project/bplustree/types"
	"github.com/stretchr/testify/assert"
)

func TestParseSelect(t *testing.T) {
	tokens := []Token{
		{Type: tokenKeyword, Literal: "SELECT"},
		{Type: tokenIdentifier, Literal: "name"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenIdentifier, Literal: "age"},
		{Type: tokenKeyword, Literal: "FROM"},
		{Type: tokenIdentifier, Literal: "users"},
		{Type: tokenKeyword, Literal: "WHERE"},
		{Type: tokenIdentifier, Literal: "age"},
		{Type: tokenEquals, Literal: "=="},
		{Type: tokenInteger, Literal: "30"},
	}

	parser := NewParser(tokens)
	stmt, err := parser.Parse()

	assert.NoError(t, err)
	selectStmt, ok := stmt.(*types.SelectStatement)
	assert.True(t, ok)

	// 驗證解析結果
	assert.Equal(t, []string{"name", "age"}, selectStmt.Columns)
	assert.Equal(t, "users", selectStmt.Table)
	assert.NotNil(t, selectStmt.Where)
	if selectStmt.Where != nil {
		assert.Equal(t, "age", selectStmt.Where.Expr.Left)
		assert.Equal(t, "==", selectStmt.Where.Expr.Operator)
		assert.Equal(t, "30", selectStmt.Where.Expr.Right)
	}
}

// 測試SQL插入的語句
func TestParseInsert(t *testing.T) {
	tokens := []Token{
		{Type: tokenKeyword, Literal: "INSERT"},
		{Type: tokenKeyword, Literal: "INTO"},
		{Type: tokenIdentifier, Literal: "users"},
		{Type: tokenLeftParen, Literal: "("},
		{Type: tokenIdentifier, Literal: "name"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenIdentifier, Literal: "age"},
		{Type: tokenRightParen, Literal: ")"},
		{Type: tokenKeyword, Literal: "VALUES"},
		{Type: tokenLeftParen, Literal: "("},
		{Type: tokenString, Literal: "\"John\""},
		{Type: tokenComma, Literal: ","},
		{Type: tokenInteger, Literal: "25"},
		{Type: tokenRightParen, Literal: ")"},
	}

	parser := NewParser(tokens)
	stmt, err := parser.Parse()

	assert.NoError(t, err)
	insertStmt, ok := stmt.(*types.InsertStatement)
	assert.True(t, ok)
	assert.Equal(t, "users", insertStmt.Table)
	assert.Equal(t, []string{"name", "age"}, insertStmt.Columns)
	assert.Equal(t, []types.Value{{Literal: "\"John\""}, {Literal: "25"}}, insertStmt.Values)
}

// 測試SQL更新的語句
func TestParseUpdate(t *testing.T) {
	tokens := []Token{
		{Type: tokenKeyword, Literal: "UPDATE"},
		{Type: tokenIdentifier, Literal: "users"},
		{Type: tokenKeyword, Literal: "SET"},
		{Type: tokenIdentifier, Literal: "name"},
		{Type: tokenAssign, Literal: "="},
		{Type: tokenString, Literal: "\"Jane\""},
		{Type: tokenComma, Literal: ","},
		{Type: tokenIdentifier, Literal: "age"},
		{Type: tokenAssign, Literal: "="},
		{Type: tokenInteger, Literal: "32"},
		{Type: tokenKeyword, Literal: "WHERE"},
		{Type: tokenIdentifier, Literal: "id"},
		{Type: tokenEquals, Literal: "="},
		{Type: tokenInteger, Literal: "1"},
	}

	parser := NewParser(tokens)
	stmt, err := parser.Parse()

	assert.NoError(t, err)
	updateStmt, ok := stmt.(*types.UpdateStatement)
	assert.True(t, ok)
	assert.Equal(t, "users", updateStmt.Table)
	assert.Equal(t, []string{"name", "age"}, updateStmt.Columns)
	assert.Equal(t, []types.Value{{Literal: "\"Jane\""}, {Literal: "32"}}, updateStmt.Values)
	assert.NotNil(t, updateStmt.Where)
	assert.Equal(t, "id", updateStmt.Where.Expr.Left)
	assert.Equal(t, "=", updateStmt.Where.Expr.Operator)
	assert.Equal(t, "1", updateStmt.Where.Expr.Right)
}

func TestParseDelete(t *testing.T) {
	tokens := []Token{
		{Type: tokenKeyword, Literal: "DELETE"},
		{Type: tokenKeyword, Literal: "FROM"},
		{Type: tokenIdentifier, Literal: "users"},
		{Type: tokenKeyword, Literal: "WHERE"},
		{Type: tokenIdentifier, Literal: "age"},
		{Type: tokenEquals, Literal: "="},
		{Type: tokenInteger, Literal: "25"},
	}

	parser := NewParser(tokens)
	stmt, err := parser.Parse()

	assert.NoError(t, err)
	deleteStmt, ok := stmt.(*types.DeleteStatement)
	assert.True(t, ok)
	assert.Equal(t, "users", deleteStmt.Table)
	assert.NotNil(t, deleteStmt.Where)
	assert.Equal(t, "age", deleteStmt.Where.Expr.Left)
	assert.Equal(t, "=", deleteStmt.Where.Expr.Operator)
	assert.Equal(t, "25", deleteStmt.Where.Expr.Right)
}

func TestParseCreateTable(t *testing.T) {
	tokens := []Token{
		{Type: tokenKeyword, Literal: "CREATE"},
		{Type: tokenKeyword, Literal: "TABLE"},
		{Type: tokenIdentifier, Literal: "users"},
		{Type: tokenLeftParen, Literal: "("},
		{Type: tokenIdentifier, Literal: "id"},
		{Type: tokenKeyword, Literal: "INTEGER"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenIdentifier, Literal: "name"},
		{Type: tokenKeyword, Literal: "STRING"},
		{Type: tokenRightParen, Literal: ")"},
	}

	parser := NewParser(tokens)
	stmt, err := parser.Parse()

	assert.NoError(t, err)
	createStmt, ok := stmt.(*types.CreateTableStatement)
	assert.True(t, ok)
	assert.Equal(t, "users", createStmt.Table)
	assert.Equal(t, []types.ColumnDefinition{
		{Name: "id", Type: "INTEGER"},
		{Name: "name", Type: "STRING"},
	}, createStmt.Columns)
}

func TestParseDropTable(t *testing.T) {
	tokens := []Token{
		{Type: tokenKeyword, Literal: "DROP"},
		{Type: tokenKeyword, Literal: "TABLE"},
		{Type: tokenIdentifier, Literal: "users"},
	}

	parser := NewParser(tokens)
	stmt, err := parser.Parse()

	assert.NoError(t, err)
	dropStmt, ok := stmt.(*types.DropTableStatement)
	assert.True(t, ok)
	assert.Equal(t, "users", dropStmt.Table)
}

func TestParseBulkInsert(t *testing.T) {
	tokens := []Token{
		{Type: tokenKeyword, Literal: "INSERT"},
		{Type: tokenKeyword, Literal: "INTO"},
		{Type: tokenIdentifier, Literal: "employees"},
		{Type: tokenLeftParen, Literal: "("},
		{Type: tokenIdentifier, Literal: "id"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenIdentifier, Literal: "name"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenIdentifier, Literal: "age"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenIdentifier, Literal: "department"},
		{Type: tokenRightParen, Literal: ")"},
		{Type: tokenKeyword, Literal: "VALUES"},
		{Type: tokenLeftParen, Literal: "("},
		{Type: tokenInteger, Literal: "1"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenString, Literal: "\"Peter Parker\""},
		{Type: tokenComma, Literal: ","},
		{Type: tokenInteger, Literal: "19"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenString, Literal: "\"IT\""},
		{Type: tokenRightParen, Literal: ")"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenLeftParen, Literal: "("},
		{Type: tokenInteger, Literal: "2"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenString, Literal: "\"Bruce Wayne\""},
		{Type: tokenComma, Literal: ","},
		{Type: tokenInteger, Literal: "29"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenString, Literal: "\"Automobile\""},
		{Type: tokenRightParen, Literal: ")"},
	}

	parser := NewParser(tokens)
	stmt, err := parser.Parse()

	assert.NoError(t, err)
	bulkInsertStmt, ok := stmt.(*types.BulkInsertStatement)
	assert.True(t, ok)
	assert.Equal(t, "employees", bulkInsertStmt.Table)
	assert.Equal(t, []string{"id", "name", "age", "department"}, bulkInsertStmt.Columns)
	
	expectedValues := [][]types.Value{
		{
			{Literal: "1"},
			{Literal: "\"Peter Parker\""},
			{Literal: "19"},
			{Literal: "\"IT\""},
		},
		{
			{Literal: "2"},
			{Literal: "\"Bruce Wayne\""},
			{Literal: "29"},
			{Literal: "\"Automobile\""},
		},
	}
	assert.Equal(t, expectedValues, bulkInsertStmt.Values)
}
