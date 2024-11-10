@startuml
actor User
participant Lexer
participant "Token Channel" as TokenChannel

User -> Lexer: Run()
activate Lexer
Lexer -> Lexer: run() -> lexStatement()
loop 每個字符解析
    Lexer -> Lexer: next() 讀取字符
    alt 字符是空白
        Lexer -> Lexer: ignore() 忽略字符
    else 字符是字母
        Lexer -> Lexer: lexIdentifier()
        loop 解析標識符
            Lexer -> Lexer: next() 讀取字符
        end
        Lexer -> TokenChannel: emit() 標識符 Token
    else 字符是數字
        Lexer -> Lexer: lexNumber()
        loop 解析數字
            Lexer -> Lexer: next() 讀取字符
        end
        Lexer -> TokenChannel: emit() 數字 Token
    else 字符是引號
        Lexer -> Lexer: lexString()
        loop 解析字符串
            Lexer -> Lexer: next() 讀取字符
            alt 字符是 EOF
                Lexer -> Lexer: errorf() 發送錯誤標記
                break
            else 字符是引號結尾
                Lexer -> TokenChannel: emit() 字符串 Token
                break
            end
        end
    else 特殊符號 (=, *, (, ), ,)
        Lexer -> Lexer: emit() 處理符號並發送 Token
    else 未知字符
        Lexer -> Lexer: errorf() 發送錯誤標記
    end
end
Lexer -> TokenChannel: emit() EOF Token
deactivate Lexer
TokenChannel -> User: 返回 Tokens
@enduml
