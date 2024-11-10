// 這是 Lexer 模組，用來解析 SQL 語句並產生標記 (tokens)
// -----------------------------------------------------------
//                  讀取 SQL 查詢        *詞法分析*
// SQL 查詢文件 ---------------> 字符流 ---------------> Tokens --> 語法解析 --> 執行/編譯
// -----------------------------------------------------------

// 將輸入的 SQL 語句分解為一系列標記，用於後續的語法解析和編譯。

package parser

import (
	"fmt"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// 定義標記的類型
type tokenType int

const (
	tokenError      tokenType = iota // 表示錯誤標記
	tokenEOF                         // 表示文件結束
	tokenIdentifier                  // 標識符，如表名、欄位名
	tokenInteger                     // 整數數值
	tokenString                      // 字符串
	tokenKeyword                     // 關鍵字，如 SELECT、FROM 等
	tokenLeftParen                   // 左括號 "("
	tokenRightParen                  // 右括號 ")"
	tokenComma                       // 逗號 ","
	tokenAssign                      // 賦值 "="
	tokenEquals                      // 等號 "=="
)

// Token 結構用來表示詞法分析後生成的標記
type Token struct {
	Type    tokenType // 標記的類型
	Literal string    // 標記的字面值
}

// Lexer 是詞法分析器，負責將 SQL 字符流解析成 Tokens
type Lexer struct {
	input    string     // 輸入的 SQL 語句
	start    int        // 當前標記的起始位置
	position int        // 當前解析字符的位置
	width    int        // 最後讀取的字符寬度
	tokens   chan Token // 用來存儲生成的 Tokens
	mu       sync.Mutex // 互斥鎖，保證並發安全
	closed   bool       // 標記是否已經關閉 tokens 通道
}

// keywords 映射，將 SQL 關鍵字對應到標記類型
var keywords = map[string]tokenType{
	"SELECT": tokenKeyword,
	"INSERT": tokenKeyword,
	"DELETE": tokenKeyword,
	"INTO":   tokenKeyword,
	"FROM":   tokenKeyword,
	"WHERE":  tokenKeyword,
	"LIMIT":  tokenKeyword,
	"AND":    tokenKeyword,
	"VALUES": tokenKeyword,
	"UPDATE": tokenKeyword,
	"SET":    tokenKeyword,
	"CREATE": tokenKeyword,
	"DROP":   tokenKeyword,
	"TABLE":  tokenKeyword,
}

// NewLexer 創建並返回一個新的 Lexer 實例
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		tokens: make(chan Token, 100), // 給通道設置緩衝，以提高性能
		closed: false,
	}
	return l
}

// Run 啟動 Lexer 解析並返回生成的 Tokens 列表
func (l *Lexer) Run() []Token {
	go l.run() // 啟動解析過程
	var tokens []Token
	for token := range l.tokens {
		if token.Type == tokenEOF {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens
}

// run 是 Lexer 的核心解析函數，負責管理解析狀態
func (l *Lexer) run() {
	for state := lexStatement; state != nil; {
		state = state(l)
	}
	// 確保 tokens 通道只關閉一次
	l.mu.Lock()
	if !l.closed {
		close(l.tokens)
		l.closed = true
	}
	l.mu.Unlock()
}

// stateFn 定義詞法分析的狀態函數類型
type stateFn func(*Lexer) stateFn

// lexStatement 是解析的入口點，用來根據字符決定下一步的解析動作
func lexStatement(l *Lexer) stateFn {
	switch r := l.next(); {
	case isSpace(r):
		l.ignore() // 忽略空白字符
		return lexStatement
	case r == '*':
		l.emit(tokenIdentifier) // 處理 "*" 符號
		return lexStatement
	case isLetter(r):
		l.backup()
		return lexIdentifier // 轉去解析標識符或關鍵字
	case isDigit(r):
		l.backup()
		return lexNumber // 轉去解析數字
	case r == '"':
		return lexString // 進入字符串解析
	case r == '(':
		l.emit(tokenLeftParen)
		return lexStatement
	case r == ')':
		l.emit(tokenRightParen)
		return lexStatement
	case r == ',':
		l.emit(tokenComma)
		return lexStatement
	case r == '=':
		if l.peek() == '=' {
			l.next()
			l.emit(tokenEquals) // 處理 "=="
		} else {
			l.emit(tokenAssign) // 處理 "="
		}
		return lexStatement
	case r == eof:
		l.emit(tokenEOF) // 發出 EOF 標記，結束解析
		return nil
	default:
		return l.errorf("unexpected character: %q", r) // 處理意外字符
	}
}

// 輔助方法

// next 讀取下一個字符
func (l *Lexer) next() rune {
	if l.position >= len(l.input) {
		l.width = 0
		return eof // 返回 EOF 標記
	}
	r, w := utf8.DecodeRuneInString(l.input[l.position:])
	l.width = w
	l.position += w
	return r
}

// backup 回退一個字符
func (l *Lexer) backup() {
	l.position -= l.width
}

// peek 預覽下一個字符而不移動位置
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// emit 發送標記到 tokens 通道
func (l *Lexer) emit(t tokenType) {
	l.mu.Lock()
	if !l.closed {
		l.tokens <- Token{t, l.input[l.start:l.position]}
	}
	l.mu.Unlock()
	l.start = l.position
}

// ignore 忽略當前字符並移動 start
func (l *Lexer) ignore() {
	l.start = l.position
}

// errorf 發送錯誤標記並結束解析
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.mu.Lock()
	if !l.closed {
		l.tokens <- Token{
			Type:    tokenError,
			Literal: fmt.Sprintf(format, args...),
		}
		close(l.tokens)
		l.closed = true
	}
	l.mu.Unlock()
	return nil
}

// 字符判斷函數

// isSpace 判斷是否為空白字符
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// isLetter 判斷是否為字母或底線
func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

// isDigit 判斷是否為數字
func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

const eof = -1 // 定義 EOF 標記

// 詞法分析器狀態函數

// lexIdentifier 解析標識符或關鍵字
func lexIdentifier(l *Lexer) stateFn {
	for {
		r := l.next()
		if !isLetter(r) && !isDigit(r) {
			l.backup() // 回退非標識符字符
			break
		}
	}
	word := l.input[l.start:l.position]
	if t, ok := keywords[strings.ToUpper(word)]; ok {
		l.emit(t) // 發送關鍵字標記
	} else {
		l.emit(tokenIdentifier) // 發送標識符標記
	}
	return lexStatement
}

// lexNumber 解析數字
func lexNumber(l *Lexer) stateFn {
	for isDigit(l.next()) {
	}
	l.backup()
	l.emit(tokenInteger)
	return lexStatement
}

// lexString 解析用引號括起來的字符串
func lexString(l *Lexer) stateFn {
	l.next() // 跳過起始引號
	for {
		switch r := l.next(); {
		case r == eof:
			return l.errorf("unterminated string") // 未結束的字符串錯誤
		case r == '"':
			l.emit(tokenString) // 發送字符串標記
			return lexStatement
		}
	}
}
