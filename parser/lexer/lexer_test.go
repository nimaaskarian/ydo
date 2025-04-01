package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/nimaaskarian/ydo/parser/token"
)

func TestNextToken(t *testing.T) {
  input := `(due+ 5d) - dickwhitman + 5dondraper - 2025-12-12/now + 2025-12-30/8:12:12==5d != 2025-12-30 = !due #`
  tests := []struct {
    expected_type token.Type
    expected_literal string
  } {
    {token.LPAREN, "("},
    {token.FIELD, "due"},
    {token.PLUS, "+"},
    {token.DUE, "5d"},
    {token.RPAREN, ")"},
    {token.MINUS, "-"},
    {token.FIELD, "dickwhitman"},
    {token.PLUS, "+"},
    {token.DUE, "5dondraper"},
    {token.MINUS, "-"},
    {token.DUE, "2025-12-12/now"},
    {token.PLUS, "+"},
    {token.DUE, "2025-12-30/8:12:12"},
    {token.EQ, "=="},
    {token.DUE, "5d"},
    {token.NOT_EQ, "!="},
    {token.DUE, "2025-12-30"},
    {token.ASSIGN, "="},
    {token.BANG, "!"},
    {token.FIELD, "due"},
    {token.ILLEGAL, "#"},
    {token.EOF, ""},
  }
  l := New(input)
  for _, test := range tests {
    tok := l.NextToken()
    assert.Equal(t, test.expected_type, tok.Type)
    assert.Equal(t, test.expected_literal, tok.Literal)
  }
}
