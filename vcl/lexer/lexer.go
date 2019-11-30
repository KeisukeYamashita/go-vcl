package lexer

import "github.com/KeisukeYamashita/go-vcl/vcl/token"

// Lexer ...
type Lexer struct {
	input   string
	pos     int
	readPos int
	char    byte
}

// New ...
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.init()
	return l
}

func (l *Lexer) init() {
	l.readChar()
}

// readChar retrieves the byte from readPos
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

// NextChar ...
func (l *Lexer) NextChar() *token.Token {
	var tok *token.Token
	switch l.char {
	case '=':
		tok = token.NewToken(token.ASSIGN, l.char)
	}

	l.readChar()
	return tok
}
