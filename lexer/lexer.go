package lexer

import (
	"fmt"
	"github.com/nimaaskarian/ydo/lexer/token"
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
  fmt.Println(string(l.ch), l.ch)
  l.skipWhitespace()
  switch l.ch {
  case '-':
    tok = newToken(token.MINUS, l.ch)
  case '+':
    tok = newToken(token.PLUS, l.ch)
  case '(':
    tok = newToken(token.LPAREN, l.ch)
  case ')':
    fmt.Println("is token.RPAREN")
    tok = newToken(token.RPAREN, l.ch)
  case ';':
    tok = newToken(token.SEMICOLON, l.ch)
  case ',':
    tok = newToken(token.COMMA, l.ch)
  case '=':
    if l.peekChar() == '=' {
      ch := l.ch
      l.readChar()
      literal := string(ch) + string(l.ch)
      tok = token.Token {Type: token.EQ, Literal: literal}
    } else {
      tok = newToken(token.ASSIGN, l.ch)
    }
  case '!':
    if l.peekChar() == '=' {
      ch := l.ch
      l.readChar()
      literal := string(ch) + string(l.ch)
      tok = token.Token {Type: token.NOT_EQ, Literal: literal}
    } else {
      tok = newToken(token.BANG, l.ch)
    }
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
    fmt.Println(string(l.ch))
    l.readChar()
  }
  fmt.Printf("finished read field %q\n", string(l.ch))
  return l.input[position:l.position]
}

func (l *Lexer) readDue() string {
  position := l.position
  for isDue(l.ch) {
    fmt.Println(string(l.ch))
    l.readChar()
  }
  fmt.Println("finished read due", string(l.ch))
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

