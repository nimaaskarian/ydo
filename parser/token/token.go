package token

type Type string
type Token struct {
  Type Type
  Literal string
}

const (
  // special
  ILLEGAL      = "ILLEGAL"
  EOF          = "EOF"
  // ident/literal
  FIELD        = "FIELD"
  DUE          = "DUE"
  // delimiters
  COMMA        = ","
  SEMICOLON    = ";"
  // operators
  EQ           = "=="
  BIT_AND      = "&"
  AND          = "&&"
  BIT_OR       = "|"
  OR           = "||"
  NOT_EQ       = "!="
  BANG         = "!"
  ASSIGN       = "="
  MINUS        = "-"
  PLUS         = "+"

  LPAREN       = "("
  RPAREN       = ")"
)

