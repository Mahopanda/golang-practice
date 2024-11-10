package parser

import (
	"testing"
)

// 修改 collectTokens 函數
func collectTokens(l *Lexer) []Token {
	return l.Run()  // 直接使用新的 Run 方法
}

// 輔助函數：驗證token序列是否符合預期
func validateTokens(t *testing.T, got []Token, want []Token) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("token數量不匹配：想要 %d 個，得到 %d 個", len(want), len(got))
		return
	}

	for i, token := range got {
		if token.Type != want[i].Type || token.Literal != want[i].Literal {
			t.Errorf("token[%d] 不匹配\n想要：{Type: %v, Literal: %q}\n得到：{Type: %v, Literal: %q}",
				i, want[i].Type, want[i].Literal, token.Type, token.Literal)
		}
	}
}

func TestLexer_CreateTable(t *testing.T) {
	input := "CREATE TABLE table1 (c1 INTEGER, c2 STRING)"
	l := NewLexer(input)
	got := collectTokens(l)

	want := []Token{
		{Type: tokenKeyword, Literal: "CREATE"},
		{Type: tokenKeyword, Literal: "TABLE"},
		{Type: tokenIdentifier, Literal: "table1"},
		{Type: tokenLeftParen, Literal: "("},
		{Type: tokenIdentifier, Literal: "c1"},
		{Type: tokenIdentifier, Literal: "INTEGER"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenIdentifier, Literal: "c2"},
		{Type: tokenIdentifier, Literal: "STRING"},
		{Type: tokenRightParen, Literal: ")"},
	}

	validateTokens(t, got, want)
}

func TestLexer_DropTable(t *testing.T) {
	input := "DROP TABLE table1"
	l := NewLexer(input)
	got := collectTokens(l)

	want := []Token{
		{Type: tokenKeyword, Literal: "DROP"},
		{Type: tokenKeyword, Literal: "TABLE"},
		{Type: tokenIdentifier, Literal: "table1"},
	}

	validateTokens(t, got, want)
}

func TestLexer_Select(t *testing.T) {
	input := `SELECT c1, c2 FROM table1 WHERE c3 == c4 AND c5 == c6`
	l := NewLexer(input)
	got := collectTokens(l)

	want := []Token{
		{Type: tokenKeyword, Literal: "SELECT"},
		{Type: tokenIdentifier, Literal: "c1"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenIdentifier, Literal: "c2"},
		{Type: tokenKeyword, Literal: "FROM"},
		{Type: tokenIdentifier, Literal: "table1"},
		{Type: tokenKeyword, Literal: "WHERE"},
		{Type: tokenIdentifier, Literal: "c3"},
		{Type: tokenEquals, Literal: "=="},
		{Type: tokenIdentifier, Literal: "c4"},
		{Type: tokenKeyword, Literal: "AND"},
		{Type: tokenIdentifier, Literal: "c5"},
		{Type: tokenEquals, Literal: "=="},
		{Type: tokenIdentifier, Literal: "c6"},
	}

	validateTokens(t, got, want)
}

func TestLexer_Insert(t *testing.T) {
	input := `INSERT INTO table1 (c1, c2, c3) VALUES (5, "some string", 10)`
	l := NewLexer(input)
	got := collectTokens(l)

	want := []Token{
		{Type: tokenKeyword, Literal: "INSERT"},
		{Type: tokenKeyword, Literal: "INTO"},
		{Type: tokenIdentifier, Literal: "table1"},
		{Type: tokenLeftParen, Literal: "("},
		{Type: tokenIdentifier, Literal: "c1"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenIdentifier, Literal: "c2"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenIdentifier, Literal: "c3"},
		{Type: tokenRightParen, Literal: ")"},
		{Type: tokenKeyword, Literal: "VALUES"},
		{Type: tokenLeftParen, Literal: "("},
		{Type: tokenInteger, Literal: "5"},
		{Type: tokenComma, Literal: ","},
		{Type: tokenString, Literal: `"some string"`},
		{Type: tokenComma, Literal: ","},
		{Type: tokenInteger, Literal: "10"},
		{Type: tokenRightParen, Literal: ")"},
	}

	validateTokens(t, got, want)
}

func TestLexer_Update(t *testing.T) {
	input := `UPDATE table1 SET c1 = 10 WHERE c1 == 5 AND c3 == "quoted string"`
	l := NewLexer(input)
	got := collectTokens(l)

	want := []Token{
		{Type: tokenKeyword, Literal: "UPDATE"},
		{Type: tokenIdentifier, Literal: "table1"},
		{Type: tokenKeyword, Literal: "SET"},
		{Type: tokenIdentifier, Literal: "c1"},
		{Type: tokenAssign, Literal: "="},
		{Type: tokenInteger, Literal: "10"},
		{Type: tokenKeyword, Literal: "WHERE"},
		{Type: tokenIdentifier, Literal: "c1"},
		{Type: tokenEquals, Literal: "=="},
		{Type: tokenInteger, Literal: "5"},
		{Type: tokenKeyword, Literal: "AND"},
		{Type: tokenIdentifier, Literal: "c3"},
		{Type: tokenEquals, Literal: "=="},
		{Type: tokenString, Literal: `"quoted string"`},
	}

	validateTokens(t, got, want)
}

func TestLexer_Delete(t *testing.T) {
	input := `DELETE FROM table1 WHERE c1 == 5 AND c3 == "quoted string"`
	l := NewLexer(input)
	got := collectTokens(l)

	want := []Token{
		{Type: tokenKeyword, Literal: "DELETE"},
		{Type: tokenKeyword, Literal: "FROM"},
		{Type: tokenIdentifier, Literal: "table1"},
		{Type: tokenKeyword, Literal: "WHERE"},
		{Type: tokenIdentifier, Literal: "c1"},
		{Type: tokenEquals, Literal: "=="},
		{Type: tokenInteger, Literal: "5"},
		{Type: tokenKeyword, Literal: "AND"},
		{Type: tokenIdentifier, Literal: "c3"},
		{Type: tokenEquals, Literal: "=="},
		{Type: tokenString, Literal: `"quoted string"`},
	}

	validateTokens(t, got, want)
}

func TestLexer_ErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Token
	}{
		{
			name:  "未結束的字符串",
			input: `SELECT * FROM table WHERE name = "unclosed string`,
			want:  Token{Type: tokenError, Literal: "unterminated string"},
		},
		{
			name:  "非法字符",
			input: "SELECT * FROM table WHERE name @ value",
			want:  Token{Type: tokenError, Literal: "unexpected character: '@'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tokens := collectTokens(l)

			// 檢查是否有錯誤token
			for _, token := range tokens {
				if token.Type == tokenError {
					if token.Literal != tt.want.Literal {
						t.Errorf("錯誤消息不匹配\n想要：%q\n得到：%q",
							tt.want.Literal, token.Literal)
					}
					return
				}
			}
			t.Error("預期出現錯誤token，但沒有找到")
		})
	}
}
