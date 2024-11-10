package parser

import (
	"github.com/Mahopanda/mini-project/bplustree/models"
	"github.com/Mahopanda/mini-project/bplustree/types"
)

type QueryExecutor struct {
	DB interface{} // 使用接口而不是具體類型
}

func NewQueryExecutor(db interface{}) *QueryExecutor {
	return &QueryExecutor{DB: db}
}

func (qe *QueryExecutor) Execute(stmt types.Statement) (*models.Value, error) {
	switch stmt := stmt.(type) {
	case *types.SelectStatement:
		return qe.executeSelect(stmt)
	case *types.InsertStatement:
		return qe.executeInsert(stmt)
	case *types.UpdateStatement:
		return qe.executeUpdate(stmt)
	case *types.DeleteStatement:
		return qe.executeDelete(stmt)
	case *types.CreateTableStatement:
		return qe.executeCreateTable(stmt)
	case *types.DropTableStatement:
		return qe.executeDropTable(stmt)
	default:
		return nil, ErrUnsupportedStatement
	}
}

func (qe *QueryExecutor) executeSelect(stmt *types.SelectStatement) (*models.Value, error) {
	return nil, nil
}

func (qe *QueryExecutor) executeInsert(stmt *types.InsertStatement) (*models.Value, error) {
	return nil, nil
}

func (qe *QueryExecutor) executeUpdate(stmt *types.UpdateStatement) (*models.Value, error) {
	return nil, nil
}

func (qe *QueryExecutor) executeDelete(stmt *types.DeleteStatement) (*models.Value, error) {
	return nil, nil
}

func (qe *QueryExecutor) executeCreateTable(stmt *types.CreateTableStatement) (*models.Value, error) {
	return nil, nil
}

func (qe *QueryExecutor) executeDropTable(stmt *types.DropTableStatement) (*models.Value, error) {
	return nil, nil
}
