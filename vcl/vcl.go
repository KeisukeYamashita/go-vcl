package vcl

import (
	"github.com/KeisukeYamashita/go-vcl/internal/decoder"
	"github.com/KeisukeYamashita/go-vcl/internal/lexer"
	"github.com/KeisukeYamashita/go-vcl/internal/parser"
)

// Decode ...
func Decode(bs []byte, val interface{}) []error {
	p := getParser(bs)
	prog := p.ParseProgram()
	return decoder.Decode(prog, val)
}

func getParser(bs []byte) *parser.Parser {
	l := lexer.NewLexer(string(bs))
	return parser.NewParser(l)
}
