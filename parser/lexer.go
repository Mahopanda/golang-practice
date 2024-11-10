package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer 是解析器的前端，負責將 SQL 語句轉換為語法樹。
type Lexer struct {
	input    string     // 輸入的SQL語句
	start    int        // 詞法單元的開始位置
	position int        // 詞法單元的結束位置
	width    int        // 詞法單元的寬度
	tokens   chan Token // 語法解析器之間的通信管道
}

type Token struct {
	Type    tokenType // 詞法單元的類型
	Literal string    // 詞法單元的值
}

// 定義語言的詞法單元
type tokenType int

// 定義詞法單元的類型
const (
	tokenError tokenType = iota
	tokenIdentifier
	tokenInteger
	tokenString
	tokenKeyword // 用於表示關鍵字的通用類型
	tokenLeftParenthesis
	tokenRightParenthesis
	tokenDelimeter
	tokenAssign
	tokenEquals
	tokenEnd
)

var keywords = map[string]tokenType{
	"SELECT": tokenSelect,
	"INSERT": tokenInsert,
	"DELETE": tokenDelete,
	"INTO":   tokenInto,
	"FROM":   tokenFrom,
	"WHERE":  tokenWhere,
	"LIMIT":  tokenLimit,
	"AND":    tokenAnd,
	"VALUES": tokenValues,
	"UPDATE": tokenUpdate,
	"SET":    tokenSet,
	"CREATE": tokenCreate,
	"DROP":   tokenDrop,
	"TABLE":  tokenTable,
	"ENGINE": tokenEngine,
	"LSM":    tokenLSM,
	"BPTREE": tokenBPTree,
}

// token convert to string
func (t tokenType) String() string {
	s, ok := tokenStringMap[t]
	if !ok {
		return fmt.Sprintf("unknown token: %d", t)
	}
	return s
}

const eof = -1

type lexFunc func(*Lexer) lexFunc

// 建立以及回傳一個Lexer
func NewLexer(input string) *Lexer {
	return &Lexer{input: input,
		tokens: make(chan Token),
	}
}

// 啟動詞法分析器
func (l *Lexer) Run() {
	go l.lexer()
}

// 詞法分析器
func (l *Lexer) lexer() {
	for state := lexText; state != nil; {
		state = state(l)
	}
}

// 負責處理SQL語句的詞法單元，解析輸入字符串，狀態機模式
func lexText(l *Lexer) lexFunc {
	r := l.next()
	switch true { //當為true的時候就可以讓每個case子句的條件進行比較
	case unicode.IsDigit(r):
		return lexInteger
	case r == eof:
		return lexEnd
	}
	return lexText
}

// 處理整數
func lexInteger(l *Lexer) lexFunc {
	r := l.peek()
	if !unicode.IsDigit(r) {
		l.next()
		return lexInteger
	}
	l.produce(tokenInteger)
	return lexText
}

// 生成一個新的token，並發送到tokens管道，當lexr解析到一個完整的標記時，將標記發送到語法解析器
// 下列各個方法用來讓lexer 在解析輸入 SQL 查詢字符串時的運作，使其能夠有效地管理解析狀態、字符讀取和標記生成。
func (l *Lexer) produce(t tokenType) {
	l.tokens <- Token{Type: t, Literal: l.input[l.start:l.position]}
	l.start = l.position
}

// 讀取並回傳下一個字符，用於解析過程中
func (l *Lexer) next() rune {
	if l.position >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.position:])
	l.position += w
	return r
}

// 將position後移一位，並退回到上個位置，當lexer發現當前字符不符合預期時，需要退回來重新解析
func (l *Lexer) revert() {
	l.position -= l.width
}

// 查看輸入字符串的下一個字符，但不會移動position，用於在不改當變當前位置在不改變變當前位置的情況下檢查下一個字符，
// 例如:判斷是否需要產生復合詞法單元
func (l *Lexer) peek() rune {
	r, _ := utf8.DecodeRuneInString(l.input[l.position:])
	return r
}

// 清空tokens管道
func (l *Lexer) drain() {
	// 使用 for range l.tokens 循環來消耗 tokens 通道中的所有剩餘標記，直到通道關閉。
	for range l.tokens {
	}
}

// 產生錯誤
func (l *Lexer) errorf(format string, args ...interface{}) lexFunc {
	l.tokens <- Token{Type: tokenError, Literal: fmt.Sprintf(format, args...)}
	return nil
}

// 處理結束
func lexEnd(l *Lexer) lexFunc {
	l.produce(tokenEnd)
	return nil
}

// 判斷字母、底線與數字
func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// 解析並產生用引號括起的字符串，例如select "column" from table
func lexString(l *Lexer) lexFunc {
	r := l.next()
	switch true {
	case r == '"':
		l.produce(tokenString)
		return lexText
	case r == eof:
		return l.errorf("unterminated string")
	}

	return lexString
}

// 詞法分析的入口點，當 lexer 開始解析輸入的 SQL 查詢時，會從這個函式開始，
// 處理lexer的狀態,解析整體SQL查並根據輸入字符串決定下一步的解析羅輯， 例如: 解析到字符串、整數、空白等
func lexStatement(l *Lexer) lexFunc {
	r := l.next()
	switch true {
	case unicode.IsSpace(r):
		return lexStatement
	case isAlphaNumeric(r):
		return lexIdentifier
	case r == eof:
		return lexEnd
	}
	return lexStatement
}

// 判斷SQL中的關鍵字的狀態
func lexIdentifier(l *Lexer) lexFunc {
	for {
		r := l.next()
		if !isAlphaNumeric(r) {
			l.revert()
			break
		}
	}

	word := l.input[l.start:l.position]
	if t, ok := keywords[strings.ToUpper(word)]; ok {
		l.produce(t) // 生成對應的關鍵字 token
	} else {
		l.produce(tokenIdentifier) // 生成標識符 token
	}

	return lexStatement
}
