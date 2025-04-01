package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
  input := `(due+ 5d) - dickwhitman + 5dondraper - 2025-12-12/now + 2025-12-30/8:12:12==5d != 2025-12-30`
  tests := []struct {
    expected_type TokenType
    expected_literal string
  } {
    {LPAREN, "("},
    {FIELD, "due"},
    {PLUS, "+"},
    {DUE, "5d"},
    {RPAREN, ")"},
    {MINUS, "-"},
    {FIELD, "dickwhitman"},
    {PLUS, "+"},
    {DUE, "5dondraper"},
    {MINUS, "-"},
    {DUE, "2025-12-12/now"},
    {PLUS, "+"},
    {DUE, "2025-12-30/8:12:12"},
    {EQ, "=="},
    {DUE, "5d"},
    {NOT_EQ, "!="},
    {DUE, "2025-12-30"},
    {EOF, ""},
  }
  l := New(input)
  for _, test := range tests {
    tok := l.NextToken()
    assert.Equal(t, test.expected_type, tok.Type)
    assert.Equal(t, test.expected_literal, tok.Literal)
  }
}
