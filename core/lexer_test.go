package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
  input := `due + 5d - dickwhitman + 5dondraper`
  tests := []struct {
    expected_type TokenType
    expected_literal string
  } {
    {FIELD, "due"},
    {PLUS, "+"},
    {DURATION, "5d"},
    {MINUS, "-"},
    {ILLEGAL, "dickwhitman"},
    {PLUS, "+"},
    {ILLEGAL, "5dondraper"},
    {EOF, ""},
  }
  l := New(input)
  for _, test := range tests {
    tok := l.NextToken()
    assert.Equal(t, test.expected_type, tok.Type)
    assert.Equal(t, test.expected_literal, tok.Literal)
  }
}
