package ast

import "github.com/KeisukeYamashita/go-vcl/pkg/vcl/token"

// Program ...
type Program struct {
	Statements []Statement
}

// Node ...
type Node interface {
	TokenLiteral() string
}

// Statement ...
type Statement interface {
	Node
	statementNode()
}

// Expression ...
type Expression interface {
	Node
	expressionNode()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (exp *PrefixExpression) expressionNode() {}
func (exp *PrefixExpression) TokenLiteral() string {
	return exp.Token.Literal
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (exp *InfixExpression) expressionNode() {}
func (exp *InfixExpression) TokenLiteral() string {
	return exp.Token.Literal
}

// AssignStatement holds the Name for the Identifier and its value
type AssignStatement struct {
	Token token.Token // token.ASSIGN
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode() {}
func (as *AssignStatement) TokenLiteral() string {
	return as.Token.Literal
}

// ReturnStatement holds the Name for the Identifier and its value
type ReturnStatement struct {
	Token       token.Token // token.ASSIGN
	ReturnValue Expression
}

func (as *ReturnStatement) statementNode() {}
func (as *ReturnStatement) TokenLiteral() string {
	return as.Token.Literal
}

// ReturnStatement holds the Name for the Identifier and its value
type ExpressionStatement struct {
	Token      token.Token // token.ASSIGN
	Expression Expression
}

func (as *ExpressionStatement) statementNode() {}
func (as *ExpressionStatement) TokenLiteral() string {
	return as.Token.Literal
}

// Identifier ...
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (i *BooleanLiteral) expressionNode() {}
func (i *BooleanLiteral) TokenLiteral() string {
	return i.Token.Literal
}
