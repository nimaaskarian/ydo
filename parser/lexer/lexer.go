package lexer

import (
	"github.com/nimaaskarian/ydo/parser/token"
)

type Lexer struct {
  input           string
  // points to current char
  position        int
  // points to the char to read (after current char)
  read_position   int
  // char under examination
  ch              byte
}

func New(input string) *Lexer {
  l := &Lexer{input: input}
  l.readChar()
  return l
}

func (l *Lexer) readChar() {
   if l.read_position >= len(l.input) {
    l.ch = 0
  } else {
    l.ch = l.input[l.read_position]
  }
  l.position = l.read_position
  l.read_position += 1
}

func (l *Lexer) NextToken() token.Token {
  var tok token.Token
  l.skipWhitespace()
  switch l.ch {
  case '-':
    tok = newToken(token.MINUS, l.ch)
  case '+':
    tok = newToken(token.PLUS, l.ch)
  case '(':
    tok = newToken(token.LPAREN, l.ch)
  case ')':
    tok = newToken(token.RPAREN, l.ch)
  case '&':
    tok = l.peekToken('&', token.BIT_AND, token.AND)
  case '|':
    tok = l.peekToken('|', token.BIT_OR, token.OR)
  case '=':
    tok = l.peekToken('=', token.ASSIGN, token.EQ)
  case '!':
    tok = l.peekToken('=', token.BANG, token.NOT_EQ)
  case 0:
    tok.Literal = ""
    tok.Type = token.EOF
  default:
    if isLetter(l.ch) {
      tok.Literal = l.readField()
      tok.Type = token.FIELD
      return tok
    } else if isDigit(l.ch) {
      tok.Type = token.DUE
      tok.Literal = l.readDue()
      return tok
    } else {
      tok = newToken(token.ILLEGAL, l.ch)
    }
  }
  l.readChar()
  return tok
}

// a helper for peekChar and assign
func (l *Lexer) peekToken(next byte, no_peek_type, peek_type token.Type) token.Token {
  if l.peekChar() == next {
    ch := l.ch
    l.readChar()
    literal := string(ch) + string(l.ch)
    return token.Token {Type: peek_type, Literal: literal}
  } else {
    return newToken(no_peek_type, l.ch)
  }
}

func (l *Lexer) peekChar() byte {
  if l.read_position >= len(l.input) {
    return 0
  } else {
    return l.input[l.read_position]
  }
}

func (l *Lexer) skipWhitespace() {
  for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
    l.readChar()
  }
}

func (l *Lexer) readField() string {
  position := l.position
  for isField(l.ch) {
    l.readChar()
  }
  return l.input[position:l.position]
}

func (l *Lexer) readDue() string {
  position := l.position
  for isDue(l.ch) {
    l.readChar()
  }
  return l.input[position:l.position]
}

func isLetter(ch byte) bool {
  return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isField(ch byte) bool {
  return isLetter(ch) || ch == '-'
}

func isDigit(ch byte) bool {
  return ch >= '0' && ch <= '9'
}

func isDue(ch byte) bool {
  return isDigit(ch) || isField(ch) || ch == '/' || ch == ':'
}

func newToken(token_type token.Type, ch byte) token.Token {
  return token.Token {Type: token_type, Literal: string(ch)}
}

