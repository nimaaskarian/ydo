package core

import (
	"time"

	"github.com/nimaaskarian/ydo/utils"
)

type TokenType string
type Token struct {
  Type TokenType
  Literal string
}

const (
  // special
  ILLEGAL      = "ILLEGAL"
  EOF          = "EOF"
  // ident/literal
  FIELD        = "FIELD"
  DURATION     = "DURATION"
  // delimitares
  COMMA        = ","
  SEMICOLON    = ";"
  // operators
  EQUAL        = "=="
  MINUS        = "-"
  PLUS         = "+"

  LPAREN       = "("
  RPAREN       = ")"
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

func (l *Lexer) NextToken() Token {
  var tok Token
  l.skipWhitespace()
  switch l.ch {
  case '-':
    tok = newToken(MINUS, l.ch)
  case '+':
    tok = newToken(PLUS, l.ch)
  case '(':
    tok = newToken(LPAREN, l.ch)
  case ')':
    tok = newToken(RPAREN, l.ch)
  case ';':
    tok = newToken(SEMICOLON, l.ch)
  case ',':
    tok = newToken(COMMA, l.ch)
  case 0:
    tok.Literal = ""
    tok.Type = EOF
  default:
    if isField(l.ch) {
      tok.Literal = l.readField()
      tok.Type = l.lookUpField(tok.Literal)
      return tok
    } else if isDigit(l.ch) {
      tok.Literal = l.readDuration()
      tok.Type = l.lookUpDuration(tok.Literal)
    } else {
      tok = newToken(ILLEGAL, l.ch)
    }
  }
  l.readChar()
  return tok
}

func (l *Lexer) skipWhitespace() {
  for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
    l.readChar()
  }
}
func (l *Lexer) lookUpField(field string) TokenType {
  task := &Task{}
  _, err := task.ReflectAccessField(field)
  if err != nil {
    return ILLEGAL
  }
  return FIELD
}

func (l *Lexer) lookUpDuration(duration string) TokenType {
  _, err := utils.ParseDuration(duration, time.Time{});
  if err != nil {
    return ILLEGAL
  }
  return DURATION
}

func (l *Lexer) readField() string {
  position := l.position
  for isField(l.ch) {
    l.readChar()
  }
  return l.input[position:l.position]
}

func (l *Lexer) readDuration() string {
  position := l.position
  for isDuration(l.ch) {
    l.readChar()
  }
  return l.input[position:l.position]
}

func isLetter(ch byte) bool {
  return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isField(ch byte) bool {
  return isLetter(ch) || ch == '-'
}

func isDigit(ch byte) bool {
  return ch >= '0' && ch <= '9'
}

func isDuration(ch byte) bool {
  return isDigit(ch) || isLetter(ch)
}

func newToken(token_type TokenType, ch byte) Token {
  return Token {Type: token_type, Literal: string(ch)}
}

