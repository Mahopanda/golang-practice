@startuml
actor User
participant Parser
participant "Lexer" as Lexer
participant "Types" as Types

User -> Parser: Parse()
activate Parser
Parser -> Lexer: 獲取 tokens
alt tokens 為空
    Parser --> User: 返回錯誤「查詢為空」
else tokens 非空
    Parser -> Parser: 檢查第一個 token 類型
    alt 是 keyword
        Parser -> Parser: 檢查 token Literal
        opt Literal 為 SELECT
            Parser -> Parser: parseSelect()
            activate Parser
            Parser -> Types: 創建 SelectStatement
            Types --> Parser: 返回 SelectStatement
            deactivate Parser
            Parser --> User: 返回解析結果
        end
        opt Literal 為 INSERT
            Parser -> Parser: parseInsert()
            activate Parser
            Parser -> Types: 創建 InsertStatement
            Types --> Parser: 返回 InsertStatement
            deactivate Parser
            Parser --> User: 返回解析結果
        end
        opt Literal 為 UPDATE
            Parser -> Parser: parseUpdate()
            activate Parser
            Parser -> Types: 創建 UpdateStatement
            Types --> Parser: 返回 UpdateStatement
            deactivate Parser
            Parser --> User: 返回解析結果
        end
        opt Literal 為 DELETE
            Parser -> Parser: parseDelete()
            activate Parser
            Parser -> Types: 創建 DeleteStatement
            Types --> Parser: 返回 DeleteStatement
            deactivate Parser
            Parser --> User: 返回解析結果
        end
        opt Literal 為 CREATE
            Parser -> Parser: parseCreateTable()
            activate Parser
            Parser -> Types: 創建 CreateTableStatement
            Types --> Parser: 返回 CreateTableStatement
            deactivate Parser
            Parser --> User: 返回解析結果
        end
        opt Literal 為 DROP
            Parser -> Parser: parseDropTable()
            activate Parser
            Parser -> Types: 創建 DropTableStatement
            Types --> Parser: 返回 DropTableStatement
            deactivate Parser
            Parser --> User: 返回解析結果
        end
    else 非支援的 keyword
        Parser --> User: 返回錯誤「不支援的 SQL 語句」
    end
else token 類型非 keyword
    Parser --> User: 返回錯誤「不支援的 SQL 語句」
end
deactivate Parser
@enduml
