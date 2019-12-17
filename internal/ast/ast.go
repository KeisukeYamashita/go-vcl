package ast

import "github.com/KeisukeYamashita/go-vcl/internal/token"

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

// PrefixExpression ...
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (exp *PrefixExpression) expressionNode() {}
func (exp *PrefixExpression) TokenLiteral() string {
	return exp.Token.Literal
}

// InfixExpression ...
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

// IfExpression ...
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (exp *IfExpression) expressionNode() {}
func (exp *IfExpression) TokenLiteral() string {
	return exp.Token.Literal
}

// BlockExpression ...
type BlockExpression struct {
	Token  token.Token
	Labels []string
	Blocks *BlockStatement
}

func (exp *BlockExpression) expressionNode() {}
func (exp *BlockExpression) TokenLiteral() string {
	return exp.Token.Literal
}

// BlockStatement ...
type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// AssignStatement holds the Name for the Identifier and its value
type AssignStatement struct {
	Token token.Token // token.ASSIGN
	Name  *Identifier
	Value Expression
}

// AssignFieldStatement holds the Name for the Identifier and its value
type AssignFieldStatement struct {
	Token token.Token // token.ASSIGN_FIELD
	Name  *Identifier
	Value Expression
}

func (as *AssignFieldStatement) statementNode() {}
func (as *AssignFieldStatement) TokenLiteral() string {
	return as.Token.Literal
}

func (as *AssignStatement) statementNode() {}
func (as *AssignStatement) TokenLiteral() string {
	return as.Token.Literal
}

// ReturnStatement holds the Name for the Identifier and its value
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (as *ReturnStatement) statementNode() {}
func (as *ReturnStatement) TokenLiteral() string {
	return as.Token.Literal
}

type CommentStatement struct {
	Token token.Token
	Value string
}

func (as *CommentStatement) statementNode() {}
func (as *CommentStatement) TokenLiteral() string {
	return as.Token.Literal
}

// CallStatement holds the Name for the Identifier and its value
type CallStatement struct {
	Token     token.Token // token.ASSIGN
	CallValue Expression
}

func (as *CallStatement) statementNode() {}
func (as *CallStatement) TokenLiteral() string {
	return as.Token.Literal
}

// ExpressionStatement holds the Name for the Identifier and its value
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

// IntegerLiteral ...
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

// BooleanLiteral ...
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (i *BooleanLiteral) expressionNode() {}
func (i *BooleanLiteral) TokenLiteral() string {
	return i.Token.Literal
}

// StringLiteral ...
type StringLiteral struct {
	Token token.Token
	Value string
}

func (i *StringLiteral) expressionNode() {}
func (i *StringLiteral) TokenLiteral() string {
	return i.Token.Literal
}

type CIDRLiteral struct {
	Token token.Token
	Value string
}

func (i *CIDRLiteral) expressionNode() {}
func (i *CIDRLiteral) TokenLiteral() string {
	return i.Token.Literal
}

// PercentageLiteral ...
type PercentageLiteral struct {
	Token token.Token
	Value string
}

func (i *PercentageLiteral) expressionNode() {}
func (i *PercentageLiteral) TokenLiteral() string {
	return i.Token.Literal
}
