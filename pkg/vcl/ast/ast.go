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

// AssignStatement holds the Name for the Identifier and its value
type AssignStatement struct {
	Token token.Token // token.ASSIGN
	Name *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode(){}
func (as *AssignStatement) TokenLiteral() string {
	return as.Token.Literal
}

// ReturnStatement holds the Name for the Identifier and its value
type ReturnStatement struct {
	Token token.Token // token.ASSIGN
	ReturnValue Expression
}

func (as *ReturnStatement) statementNode(){}
func (as *ReturnStatement) TokenLiteral() string {
	return as.Token.Literal
}

// Identifier ...
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode(){}
func (i *Identifier) TokenLiteral()string {
	return i.Token.Literal
}