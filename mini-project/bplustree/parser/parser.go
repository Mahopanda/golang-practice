// parser.go
// 此程式碼用於解析 SQL 查詢語句，支援 SELECT、INSERT、UPDATE、DELETE、CREATE TABLE 和 DROP TABLE 語句。

package parser

import (
	"errors"
	"fmt"

	"github.com/Mahopanda/mini-project/bplustree/types"
)

var (
	// ErrUnsupportedStatement 表示不支援的 SQL 語句錯誤
	ErrUnsupportedStatement = errors.New("不支援的 SQL 語句")
)

// Parser 表示一個 SQL 語句解析器
type Parser struct {
	tokens []Token // 輸入的標記列表
	pos    int     // 當前解析的位置
}

// NewParser 建立一個新的 Parser 實例
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// Parse 開始解析標記，返回 SQL 語句結構或錯誤
func (p *Parser) Parse() (types.Statement, error) {
	if len(p.tokens) == 0 {
		return nil, errors.New("查詢為空")
	}

	switch p.tokens[p.pos].Type {
	case tokenKeyword:
		keyword := p.tokens[p.pos].Literal
		switch keyword {
		case "SELECT":
			return p.parseSelect()
		case "INSERT":
			return p.parseInsert()
		case "UPDATE":
			return p.parseUpdate()
		case "DELETE":
			return p.parseDelete()
		case "CREATE":
			return p.parseCreateTable()
		case "DROP":
			return p.parseDropTable()
		default:
			return nil, ErrUnsupportedStatement
		}
	default:
		return nil, ErrUnsupportedStatement
	}
}

// parseSelect 解析 SELECT 語句
func (p *Parser) parseSelect() (*types.SelectStatement, error) {
	p.advance() // 略過 SELECT 關鍵字
	columns := []string{}

	// 解析第一個欄位
	if !p.match(tokenIdentifier) {
		return nil, errors.New("語法錯誤：SELECT 語句中缺少欄位名稱")
	}
	columns = append(columns, p.tokens[p.pos-1].Literal)

	// 解析後續欄位
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == tokenComma {
		p.advance() // 略過逗號
		if !p.match(tokenIdentifier) {
			return nil, errors.New("語法錯誤：逗號後缺少欄位名稱")
		}
		columns = append(columns, p.tokens[p.pos-1].Literal)
	}

	// 檢查 FROM 子句
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != tokenKeyword || p.tokens[p.pos].Literal != "FROM" {
		return nil, errors.New("語法錯誤：缺少 FROM 子句")
	}
	p.advance() // 略過 FROM

	// 解析表名
	if !p.match(tokenIdentifier) {
		return nil, errors.New("語法錯誤：缺少表名")
	}
	table := p.tokens[p.pos-1].Literal

	// 解析 WHERE 子句（可選）
	var where *types.WhereCondition
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == tokenKeyword && p.tokens[p.pos].Literal == "WHERE" {
		p.advance() // 略過 WHERE
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		where = &types.WhereCondition{Expr: expr}
	}

	return &types.SelectStatement{
		Columns: columns,
		Table:   table,
		Where:   where,
	}, nil
}

// parseInsert 解析 INSERT 語句
func (p *Parser) parseInsert() (types.Statement, error) {
	p.advance() // 略過 INSERT 關鍵字

	// 檢查 INTO 關鍵字（可選）
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == tokenKeyword && p.tokens[p.pos].Literal == "INTO" {
		p.advance()
	}

	// 解析表名
	if !p.match(tokenIdentifier) {
		return nil, errors.New("語法錯誤：缺少表名")
	}
	table := p.tokens[p.pos-1].Literal

	// 解析欄位列表
	columns := []string{}
	if !p.match(tokenLeftParen) {
		return nil, errors.New("語法錯誤：缺少左括號")
	}

	for {
		if !p.match(tokenIdentifier) {
			return nil, errors.New("語法錯誤：缺少欄位名稱")
		}
		columns = append(columns, p.tokens[p.pos-1].Literal)

		if p.pos >= len(p.tokens) {
			return nil, errors.New("語法錯誤：意外的語句結束")
		}

		if p.tokens[p.pos].Type == tokenRightParen {
			p.advance() // 略過右括號
			break
		}

		if p.tokens[p.pos].Type != tokenComma {
			return nil, errors.New("語法錯誤：缺少逗號或右括號")
		}
		p.advance() // 略過逗號
	}

	// 檢查 VALUES 關鍵字
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != tokenKeyword || p.tokens[p.pos].Literal != "VALUES" {
		return nil, errors.New("語法錯誤：缺少 VALUES 關鍵字")
	}
	p.advance() // 略過 VALUES

	// 解析多行值
	var values [][]types.Value
	for {
		if !p.match(tokenLeftParen) {
			return nil, errors.New("語法錯誤：缺少左括號")
		}

		rowValues := []types.Value{}
		for {
			if p.pos >= len(p.tokens) {
				return nil, errors.New("語法錯誤：意外的語句結束")
			}

			if p.tokens[p.pos].Type == tokenString || p.tokens[p.pos].Type == tokenInteger {
				rowValues = append(rowValues, types.Value{Literal: p.tokens[p.pos].Literal})
				p.advance()
			} else {
				return nil, errors.New("語法錯誤：缺少值")
			}

			if p.pos >= len(p.tokens) {
				return nil, errors.New("語法錯誤：意外的語句結束")
			}

			if p.tokens[p.pos].Type == tokenRightParen {
				p.advance() // 略過右括號
				break
			}

			if p.tokens[p.pos].Type != tokenComma {
				return nil, errors.New("語法錯誤：缺少逗號或右括號")
			}
			p.advance() // 略過逗號
		}

		values = append(values, rowValues)

		// 檢查是否還有更多行
		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != tokenComma {
			break
		}
		p.advance() // 略過行之間的逗號
	}

	// 如果只有一行數據，返回普通的 InsertStatement
	if len(values) == 1 {
		return &types.InsertStatement{
			Table:   table,
			Columns: columns,
			Values:  values[0],
		}, nil
	}

	// 如果有多行數據，返回 BulkInsertStatement
	return &types.BulkInsertStatement{
		Table:   table,
		Columns: columns,
		Values:  values,
	}, nil
}

// parseUpdate 解析 UPDATE 語句
func (p *Parser) parseUpdate() (*types.UpdateStatement, error) {
	p.advance() // 略過 UPDATE 關鍵字

	// 解析表名
	if !p.match(tokenIdentifier) {
		return nil, errors.New("語法錯誤：缺少表名")
	}
	table := p.tokens[p.pos-1].Literal

	// 解析 SET 關鍵字
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != tokenKeyword || p.tokens[p.pos].Literal != "SET" {
		return nil, errors.New("語法錯誤：缺少 SET 關鍵字")
	}
	p.advance() // 略過 SET

	// 解析欄位和值
	columns := []string{}
	values := []types.Value{}
	for {
		// 解析欄位名
		if !p.match(tokenIdentifier) {
			return nil, errors.New("語法錯誤：缺少欄位名稱")
		}
		columns = append(columns, p.tokens[p.pos-1].Literal)

		// 解析賦值符號
		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != tokenAssign {
			return nil, errors.New("語法錯誤：缺少賦值符號")
		}
		p.advance() // 略過賦值符號

		// 解析值
		if !p.match(tokenString) && !p.match(tokenInteger) {
			return nil, errors.New("語法錯誤：缺少值")
		}
		values = append(values, types.Value{Literal: p.tokens[p.pos-1].Literal})

		// 檢查是否還有更多的欄位值對
		if p.pos >= len(p.tokens) {
			return nil, errors.New("語法錯誤：意外的語句結束")
		}

		if p.tokens[p.pos].Type != tokenComma {
			break
		}
		p.advance() // 略過逗號
	}

	// 解析 WHERE 子句（可選）
	var where *types.WhereCondition
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == tokenKeyword && p.tokens[p.pos].Literal == "WHERE" {
		p.advance() // 略過 WHERE
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		where = &types.WhereCondition{Expr: expr}
	}

	return &types.UpdateStatement{
		Table:   table,
		Columns: columns,
		Values:  values,
		Where:   where,
	}, nil
}

// parseDelete 解析 DELETE 語句
func (p *Parser) parseDelete() (*types.DeleteStatement, error) {
	p.advance() // 略過 DELETE 關鍵字

	// 解析 FROM 關鍵字
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != tokenKeyword || p.tokens[p.pos].Literal != "FROM" {
		return nil, errors.New("語法錯誤：少 FROM 關鍵字")
	}
	p.advance() // 略過 FROM

	// 解析表名
	if !p.match(tokenIdentifier) {
		return nil, errors.New("語法錯誤：缺少表名")
	}
	table := p.tokens[p.pos-1].Literal

	// 解析 WHERE 子句（可選）
	var where *types.WhereCondition
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == tokenKeyword && p.tokens[p.pos].Literal == "WHERE" {
		p.advance() // 略過 WHERE
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		where = &types.WhereCondition{Expr: expr}
	}

	return &types.DeleteStatement{
		Table: table,
		Where: where,
	}, nil
}

// parseCreateTable 解析 CREATE TABLE 語句
func (p *Parser) parseCreateTable() (*types.CreateTableStatement, error) {
	p.advance() // 略過 CREATE 關鍵字

	// 檢查 TABLE 關鍵字
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != tokenKeyword || p.tokens[p.pos].Literal != "TABLE" {
		return nil, errors.New("語法錯誤：缺少 TABLE 關鍵字")
	}
	p.advance() // 略過 TABLE

	// 解析表名
	if !p.match(tokenIdentifier) {
		return nil, errors.New("語法錯誤：缺少表名")
	}
	tableName := p.tokens[p.pos-1].Literal

	// 解析欄位定義
	if !p.match(tokenLeftParen) {
		return nil, errors.New("語法錯誤：缺少左括號")
	}

	columns := []types.ColumnDefinition{}
	for {
		if !p.match(tokenIdentifier) {
			return nil, errors.New("語法錯誤：缺少欄位名稱")
		}
		columnName := p.tokens[p.pos-1].Literal

		if !p.match(tokenKeyword) {
			return nil, errors.New("語法錯誤：缺少欄位類型")
		}
		columnType := p.tokens[p.pos-1].Literal

		columns = append(columns, types.ColumnDefinition{
			Name: columnName,
			Type: columnType,
		})

		if p.pos >= len(p.tokens) {
			return nil, errors.New("語法錯誤：意外的語句結束")
		}

		if p.tokens[p.pos].Type == tokenRightParen {
			p.advance() // 略過右括號
			break
		}

		if p.tokens[p.pos].Type != tokenComma {
			return nil, errors.New("語法錯誤：缺少逗號或右括號")
		}
		p.advance() // 略過逗號
	}

	return &types.CreateTableStatement{
		Table:   tableName,
		Columns: columns,
	}, nil
}

// parseDropTable 解析 DROP TABLE 語句
func (p *Parser) parseDropTable() (*types.DropTableStatement, error) {
	p.advance() // 略過 DROP 關鍵字

	// 檢查 TABLE 關鍵字
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != tokenKeyword || p.tokens[p.pos].Literal != "TABLE" {
		return nil, errors.New("語法錯誤：缺少 TABLE 關鍵字")
	}
	p.advance() // 略過 TABLE

	// 解析表名
	if !p.match(tokenIdentifier) {
		return nil, errors.New("語法錯誤：缺少表名")
	}
	tableName := p.tokens[p.pos-1].Literal

	return &types.DropTableStatement{
		Table: tableName,
	}, nil
}

// match 判斷當前標記是否符合預期，並前進到下一個標記
func (p *Parser) match(expectedType tokenType) bool {
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == expectedType {
		p.pos++
		return true
	}
	return false
}

// current 返回當前的標記
func (p *Parser) current() Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return Token{Type: tokenEOF}
}

// advance 將當前位置前進一個標記
func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

// expect 確保當前標記符合預期，否則返回錯誤
func (p *Parser) expect(expectedType tokenType) error {
	if !p.match(expectedType) {
		return errors.New("語法錯誤：預期 " + expectedType.String() + "，但找到了 " + p.current().Literal)
	}
	return nil
}

// parseExpression 解析表達式
func (p *Parser) parseExpression() (types.Expression, error) {
	// 解析左側
	if !p.match(tokenIdentifier) {
		return types.Expression{}, errors.New("語法錯誤：表達式缺少左側標識符")
	}
	left := p.tokens[p.pos-1].Literal

	// 解析運算符
	if !p.match(tokenEquals) {
		return types.Expression{}, errors.New("語法錯誤：表達式缺少運算符")
	}
	operator := p.tokens[p.pos-1].Literal

	// 解析右側
	if !p.match(tokenIdentifier) && !p.match(tokenString) && !p.match(tokenInteger) {
		return types.Expression{}, errors.New("語法錯誤：表達式缺少右側值")
	}
	right := p.tokens[p.pos-1].Literal

	return types.Expression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}, nil
}

// Token 結構添加 IntValue 方法
func (t Token) IntValue() int {
	value := 0
	fmt.Sscanf(t.Literal, "%d", &value)
	return value
}

// 添加 tokenType 的 String 方法
func (t tokenType) String() string {
	switch t {
	case tokenError:
		return "ERROR"
	case tokenEOF:
		return "EOF"
	case tokenIdentifier:
		return "IDENTIFIER"
	case tokenInteger:
		return "INTEGER"
	case tokenString:
		return "STRING"
	case tokenKeyword:
		return "KEYWORD"
	case tokenLeftParen:
		return "LEFT_PAREN"
	case tokenRightParen:
		return "RIGHT_PAREN"
	case tokenComma:
		return "COMMA"
	case tokenAssign:
		return "ASSIGN"
	case tokenEquals:
		return "EQUALS"
	default:
		return "UNKNOWN"
	}
}

// peek 查看當前的 token，但不移動位置
func (p *Parser) peek() Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return Token{Type: tokenEOF}
}

// isAtEnd 檢查是否已到達所有 tokens 的末尾
func (p *Parser) isAtEnd() bool {
	return p.pos >= len(p.tokens)
}

// previous 獲取前一個 token
func (p *Parser) previous() Token {
	if p.pos > 0 {
		return p.tokens[p.pos-1]
	}
	return Token{Type: tokenEOF}
}
