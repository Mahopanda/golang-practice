package types

// Statement 是所有 SQL 語句的通用接口
type Statement interface{}

// Expression 表示條件表達式
type Expression struct {
    Left     string
    Operator string
    Right    string
}

// WhereCondition 表示 WHERE 子句
type WhereCondition struct {
    Expr Expression
}

// Value 表示一個值
type Value struct {
    Literal string
}

// ColumnDefinition 表示列定義
type ColumnDefinition struct {
    Name string
    Type string
}

// SelectStatement 表示 SELECT 語句
type SelectStatement struct {
    Columns []string
    Table   string
    Where   *WhereCondition
    Limit   int
}

// InsertStatement 表示 INSERT 語句
type InsertStatement struct {
    Table   string
    Columns []string
    Values  []Value
}

// UpdateStatement 表示 UPDATE 語句
type UpdateStatement struct {
    Table   string
    Columns []string
    Values  []Value
    Where   *WhereCondition
}

// DeleteStatement 表示 DELETE 語句
type DeleteStatement struct {
    Table string
    Where *WhereCondition
}

// CreateTableStatement 表示 CREATE TABLE 語句
type CreateTableStatement struct {
    Table   string
    Columns []ColumnDefinition
}

// DropTableStatement 表示 DROP TABLE 語句
type DropTableStatement struct {
    Table string
}

// BulkInsertStatement 表示批量插入語句
type BulkInsertStatement struct {
    Table   string
    Columns []string
    Values  [][]Value // 每行的值
} 