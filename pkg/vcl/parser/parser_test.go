package parser

import (
	"fmt"
	"testing"

	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/ast"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/lexer"
)

func TestAssignStatement(t *testing.T) {
	testCases := []struct {
		input               string
		expectedIdentifiers []struct {
			expectedIdentifier string
			expectedValue      interface{}
		}
	}{
		{
			`x = 10;
y = "kekesan";
keke = true;
			`,
			[]struct {
				expectedIdentifier string
				expectedValue      interface{}
			}{
				{"x", 10},
				{"y", "kekesan"},
				{"keke", true},
			},
		},
	}

	for n, tc := range testCases {
		l := lexer.NewLexer(tc.input)
		p := NewParser(l)

		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram() failed testCase[%d] got nil program", n)
		}

		if len(program.Statements) != len(tc.expectedIdentifiers) {
			t.Fatalf("program.Statements wrong number returned, got:%d, want%d", len(program.Statements), len(tc.expectedIdentifiers))
		}

		for i, expectedIdentifiers := range tc.expectedIdentifiers {
			stmt := program.Statements[i]
			if !testAssignStatement(t, stmt, expectedIdentifiers.expectedIdentifier) {
				t.Fatalf("parse assigntStatement failed")
			}
		}
	}
}

func testAssignStatement(t *testing.T, s ast.Statement, name string) bool {
	asStmt, ok := s.(*ast.AssignStatement)
	if !ok {
		t.Errorf("s not *ast.AssignStatement, got=%T", s)
		return false
	}

	if asStmt.Name.Value != name {
		t.Errorf("asStmt.NAme.Value(=Identifier) wrong, got: '%s', want: %s", asStmt.Name.Value, name)
		return false
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	testCases := []struct {
		input               string
		expectedIdentifiers []struct {
			expectedIdentifier string
		}
	}{
		{
			`return (pass);
return (pipe);
return (cache);
			`,
			[]struct {
				expectedIdentifier string
			}{
				{"x"},
				{"y"},
				{"keke"},
			},
		},
	}

	for n, tc := range testCases {
		l := lexer.NewLexer(tc.input)
		p := NewParser(l)

		program := p.ParseProgram()
		if len(program.Statements) != 3 {
			t.Fatalf("program.Statements wrong number returned testCase[%d], got:%d, want%d", n, len(program.Statements), 3)
		}

		for _, stmt := range program.Statements {
			returnStmt, ok := stmt.(*ast.ReturnStatement)
			if !ok {
				t.Fatalf("stmt not *ast.ReturnStatement testCase[%d], got:%T", n, stmt)
			}

			if returnStmt.TokenLiteral() != "return" {
				t.Fatalf("returnStmt.TokenLiteral not 'return' testCase[%d], got:%q", n, returnStmt.TokenLiteral())
			}
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "keke;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements length is not expected, got:%d, want:%d", len(program.Statements), 1)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got:%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier, got:%T", stmt.Expression)
	}

	if ident.Value != "keke" {
		t.Errorf("ident.Value wrong, got:%s, want:%s", ident.Value, "keke")
	}

	if ident.TokenLiteral() != "keke" {
		t.Errorf("ident.TokenLiteral wrong, got:%s, want:%s", ident.TokenLiteral(), "keke")
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements length is not expected, got:%d, want:%d", len(program.Statements), 1)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got:%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Identifier, got:%T", stmt.Expression)
	}

	if ident.Value != 5 {
		t.Fatalf("ident.Value wrong, got:%d, want:%d", ident.Value, 5)
	}

	if ident.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral wrong, got:%s, want:%s", ident.TokenLiteral(), "5")
	}
}

func TestBooleanExpression(t *testing.T) {
	input := `true;
false;`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements length is not expected, got:%d, want:%d", len(program.Statements), 2)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got:%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.BooleanLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Identifier, got:%T", stmt.Expression)
	}

	if ident.Value != true {
		t.Fatalf("ident.Value wrong, got:%t, want:%t", ident.Value, true)
	}

	if ident.TokenLiteral() != "true" {
		t.Errorf("ident.TokenLiteral wrong, got:%s, want:%s", ident.TokenLiteral(), "true")
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	testCases := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 ~ 5", 5, "~", 5},
	}

	for n, tc := range testCases {
		l := lexer.NewLexer(tc.input)
		p := NewParser(l)
		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements length is not expected[testCase:%d], got:%d, want:%d", n, len(program.Statements), 1)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not ast.ExpressionStatement[testCase:%d], got:%T", n, program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("program.Statement[0]'s Expression is not ast.InfixExpression[testCase:%d], got:%T", n, program.Statements[0])
		}

		if exp.Operator != tc.operator {
			t.Fatalf("parsingInfixExpression failed to get Operator[testCase:%d], got:%s, want:%s", n, exp.Operator, tc.operator)
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	testCases := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
	}

	for n, tc := range testCases {
		l := lexer.NewLexer(tc.input)
		p := NewParser(l)
		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements wrong length got:%d, want:%d", len(program.Statements), 1)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not ast.ExpressionStatement[testCase:%d], got:%T", n, program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression, got:%T", stmt.Expression)
		}

		if exp.Operator != tc.operator {
			t.Fatalf("prefixExpression operator does not match, got:%s, want:%s", exp.Operator, tc.operator)
		}

		if !testIntegerLiter(exp.Right, tc.integerValue) {
			t.Fatalf("prefixExpression integerValue does not match want:%d", tc.integerValue)
		}
	}
}

func testIntegerLiter(il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		return false
	}

	if integ.Value != value {
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		return false
	}

	return true
}
