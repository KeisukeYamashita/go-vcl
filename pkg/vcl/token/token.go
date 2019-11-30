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
	MATCH  = "~"
	BANG   = "!"
	EQUAL  = "=="
	AND    = "&&"
	OR     = "||"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	SUBROUTINE = "SUBROUTINE"
	CALL       = "CALL"

	IF = "IF"

	RETURN = "RETURN"
	IMPORT = "IMPORT"
)

// NewToken ...
func NewToken(tokenType Type, char byte) *Token {
	return &Token{
		Type:    tokenType,
		Literal: string(char),
	}
}

var keywords = map[string]Type{
	"sub":    SUBROUTINE,
	"call":   CALL,
	"if":     IF,
	"return": RETURN,
	"import": IMPORT,
}

// LookupIndent ...
func LookupIndent(indent string) Type {
	if tokenType, ok := keywords[indent]; ok {
		return tokenType
	}

	return INDENT
}
