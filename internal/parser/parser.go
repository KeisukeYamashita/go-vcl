package parser

import (
	"fmt"
	"strconv"

	"github.com/KeisukeYamashita/go-vcl/internal/ast"
	"github.com/KeisukeYamashita/go-vcl/internal/lexer"
	"github.com/KeisukeYamashita/go-vcl/internal/token"
)

var precedences = map[token.Type]int{
	token.EQUAL: EQUALS,
	token.MATCH: EQUALS,
	token.PLUS:  SUM,
	token.AND:   EQUALS,
	token.OR:    EQUALS,
}

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser ...
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	errors        []error
	prefixParseFn map[token.Type]prefixParseFn
	infixParseFn  map[token.Type]infixParseFn
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

func (p *Parser) init() {
	p.nextToken()
	p.nextToken()

	p.prefixParseFn = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.CIDR, p.parseCIDRLiteral)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.SUBROUTINE, p.parseBlockExpression)
	p.registerPrefix(token.ACL, p.parseBlockExpression)
	p.registerPrefix(token.BACKEND, p.parseBlockExpression)
	p.registerPrefix(token.DIRECTOR, p.parseBlockExpression)
	p.registerPrefix(token.LBRACE, p.parseObjectExpression)

	p.infixParseFn = make(map[token.Type]infixParseFn)
	p.registerInfix(token.MATCH, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{
		Token: p.curToken,
	}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, err)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return lit
}

func (p *Parser) parseCIDRLiteral() ast.Expression {
	lit := &ast.CIDRLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)
	return expr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseBlockExpression() ast.Expression {
	expr := &ast.BlockExpression{
		Token: p.curToken,
	}

	labels := []string{}
	for !p.peekTokenIs(token.LBRACE) && !p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		labels = append(labels, p.curToken.Literal)
	}

	expr.Labels = labels

	if p.peekTokenIs(token.SEMICOLON) {
		return expr
	}

	if !p.peekTokenIs(token.LBRACE) {
		return nil
	}

	p.nextToken()

	expr.Blocks = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.curToken,
	}
	block.Statements = []ast.Statement{}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseObjectExpression() ast.Expression {
	expr := &ast.BlockExpression{
		Token: p.curToken,
	}

	expr.Labels = []string{}

	if p.peekTokenIs(token.SEMICOLON) {
		return expr
	}

	expr.Blocks = p.parseBlockStatement()
	return expr
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFn[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFn[tokenType] = fn
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
		switch p.peekToken.Type {
		case token.ASSIGN:
			return p.parseAssignStatement()
		default:
			return p.parseExpressionStatement()
		}
	case token.RETURN:
		return p.parseReturnStatement()
	case token.CALL:
		return p.parseCallStatement()
	default:
		return p.parseExpressionStatement()
	}
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

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		p.peekError(token.ASSIGN)
		return nil
	}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		p.peekError(token.ASSIGN)
		return nil
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseCallStatement() ast.Statement {
	stmt := &ast.CallStatement{
		Token: p.curToken,
	}

	p.nextToken()

	stmt.CallValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{
		Token: p.curToken,
	}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedentce int) ast.Expression {
	prefix := p.prefixParseFn[p.curToken.Type]
	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedentce < p.peekPrecedence() {
		infix := p.infixParseFn[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
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

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}
