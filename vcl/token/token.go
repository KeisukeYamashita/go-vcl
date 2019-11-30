package token

// Token defineds a single VCL token
type Token struct {
	Type    Type
	Literal string
}

// Type is a set of lexical tokens of the VCL
type Type string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	COMMENT = "COMMENT"

	INDENT = "INDENT"
	INT    = "INT"

	ASSIGN = "="
	PLUS   = "+"
)

// NewToken ...
func NewToken(tokenType Type, char byte) *Token {
	return &Token{
		Type:    tokenType,
		Literal: string(char),
	}
}
