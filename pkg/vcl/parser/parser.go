package parser

import (
	"fmt"

	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/ast"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/lexer"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/token"
)

// Parser ...
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	errors []error
}

// NewParser ...
func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []error{},
	}
	p.init()
	return p
}

// Errors return the parse errors
func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram ...
func (p *Parser) ParseProgram() *ast.Program {
	program := new(ast.Program)
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.IDENT:
		return p.parseAssignStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) init() {
	p.nextToken()
	p.nextToken()
}

func (p *Parser) parseAssignStatement() ast.Statement {
	stmt := &ast.AssignStatement{
		Token: p.curToken,
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		p.peekError(token.ASSIGN)
		return nil
	}

	if !p.expectPeek(token.INT) && !p.expectPeek(token.STRING) && !p.expectPeek(token.CIDR) {
		return nil
	}

	//TOOD(KeisukeYamashita): Add Expression right-hand
	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	p.nextToken()

	if !p.expectPeek(token.IDENT) {
		p.peekError(token.IDENT)
		return nil
	}

	// TODO(KeisukeYamashita): Add return expression

	if !p.expectPeek(token.RPAREN) {
		p.peekError(token.RPAREN)
		return nil
	}

	return stmt
}

func (p *Parser) peekError(t token.Type) {
	err := fmt.Errorf("expected to be token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, err)
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}
