package parser

import (
    "fmt"
    "github.com/Mahopanda/mini-project/bplustree/database"
    "github.com/Mahopanda/mini-project/bplustree/models"
)

type QueryExecutor struct {
    DB *database.Database
}

func (qe *QueryExecutor) Execute(statement Statement) (*models.Value, error) {
    switch stmt := statement.(type) {
    case *Select:
        return qe.executeSelect(stmt), nil
    case *Insert:
        return qe.executeInsert(stmt), nil
    case *Update:
        return qe.executeUpdate(stmt), nil
    case *Delete:
        return qe.executeDelete(stmt), nil
    case *CreateTable:
        return nil, qe.executeCreateTable(stmt)
    case *DropTable:
        return nil, qe.executeDropTable(stmt)
    default:
        return nil, fmt.Errorf("unsupported statement type")
    }
}

func (qe *QueryExecutor) executeSelect(selectStmt *Select) *models.Value {
    return qe.DB.QueryByID(models.Key(selectStmt.Table))
}

func (qe *QueryExecutor) executeInsert(insertStmt *Insert) *models.Value {
    // 插入資料到 B+ 樹
    qe.DB.ByID.Insert(models.Key(insertStmt.Table), models.Value{Data: insertStmt.Values})
    return nil
}

func (qe *QueryExecutor) executeUpdate(updateStmt *Update) *models.Value {
    // 實現更新邏輯
    return nil
}

func (qe *QueryExecutor) executeDelete(deleteStmt *Delete) *models.Value {
    // 實現刪除邏輯
    return nil
}

func (qe *QueryExecutor) executeCreateTable(createStmt *CreateTable) error {
    // 動態創建 B+ 樹或其他結構
    qe.DB.AddTable(createStmt.Name)
    return nil
}

func (qe *QueryExecutor) executeDropTable(dropStmt *DropTable) error {
    // 刪除指定的表
    qe.DB.RemoveTable(dropStmt.Table)
    return nil
}
