package parser

import (
	"fmt"
	"testing"

	"github.com/KeisukeYamashita/go-vcl/internal/ast"
	"github.com/KeisukeYamashita/go-vcl/internal/lexer"
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

func TestAssignStatement_Object(t *testing.T) {
	testCases := []struct {
		input               string
		expectedIdentifiers []struct {
			expectedIdentifier string
			expectedStmtName   string
			expectedStmtValue  int
		}
	}{
		{
			`x = {
	y = 10;
};`,
			[]struct {
				expectedIdentifier string
				expectedStmtName   string
				expectedStmtValue  int
			}{
				{"x", "y", 10},
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

			asStmt, ok := stmt.(*ast.AssignStatement)
			if !ok {
				t.Fatalf("stmt was not ast.AssignStatement, got:%T", stmt)
			}

			expr, ok := asStmt.Value.(*ast.BlockExpression)
			if !ok {
				t.Fatalf("stmt.Value was not ast.BlockExpression, got:%T", asStmt.Value)
			}

			if len(expr.Blocks.Statements) != 1 {
				t.Fatalf("expr.Blocks.Statements wrong number returned, got:%d, want%d", len(expr.Blocks.Statements), 1)
			}

			asStmt, ok = expr.Blocks.Statements[0].(*ast.AssignStatement)
			if !ok {
				t.Fatalf("stmt of block statement was not ast.AssignStatement, got:%T", stmt)
			}

			if !testAssignStatement(t, asStmt, expectedIdentifiers.expectedStmtName) {
				t.Fatalf("parse assigntStatement failed")
			}
		}
	}
}

func testAssignStatement(t *testing.T, s ast.Statement, name string) bool {
	asStmt, ok := s.(*ast.AssignStatement)
	if !ok {
		t.Errorf("s not *ast.AssignStatement, got:%T", s)
		return false
	}
	if asStmt.Name.Value != name {
		t.Errorf("asStmt.Name.Value(=Identifier) wrong, got: '%s', want: %s", asStmt.Name.Value, name)
		return false
	}

	return true
}

func TestCommentStatement(t *testing.T) {
	testCases := map[string]struct {
		input           string
		expectedComment string
	}{
		"with comment line by hash":         {`# keke`, "keke"},
		"with comment line by double slash": {"// keke", "keke"},
		"with single comment by multi line": {"/* keke */", "keke"},
		"with long comment by multi line":   {"/* keke is happy */", "keke is happy"},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			l := lexer.NewLexer(tc.input)
			p := NewParser(l)

			program := p.ParseProgram()
			if program == nil {
				t.Fatalf("ParseProgram() failed got nil program")
			}

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements wrong number returned, got:%d, want:%d", len(program.Statements), 1)
			}

			stmt, ok := program.Statements[0].(*ast.CommentStatement)
			if !ok {
				t.Fatalf("stmt was not ast.CommentStatement, got:%T", program.Statements[0])
			}

			if stmt.Value != tc.expectedComment {
				t.Fatalf("stmt.Value got wrong value got:%s, want:%s", stmt.Value, tc.expectedComment)
			}
		})
	}
}

func TestAssignFieldStatement(t *testing.T) {
	testCases := []struct {
		input               string
		expectedIdentifiers []struct {
			expectedStmtName  string
			expectedStmtValue string
		}
	}{
		{
			`"key": "value"`,
			[]struct {
				expectedStmtName  string
				expectedStmtValue string
			}{
				{"key", "value"},
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

			if !testAssignFieldStatement(t, stmt, expectedIdentifiers.expectedStmtName) {
				t.Fatalf("parse assigntStatement failed")
			}

			asStmt, ok := stmt.(*ast.AssignFieldStatement)
			if !ok {
				t.Fatalf("stmt was not ast.AssignStatement, got:%T", stmt)
			}

			expr, ok := asStmt.Value.(*ast.StringLiteral)
			if !ok {
				t.Fatalf("stmt.Value was not ast.StringLiteral, got:%T", expr.Value)
			}

			if expr.Value != expectedIdentifiers.expectedStmtValue {
				t.Fatalf("stmt.Value was not currect, got:%s, want:%s", expr.Value, expectedIdentifiers.expectedStmtValue)
			}
		}
	}
}

func testAssignFieldStatement(t *testing.T, s ast.Statement, name string) bool {
	asStmt, ok := s.(*ast.AssignFieldStatement)
	if !ok {
		t.Errorf("s not *ast.AssignFieldStatement, got:%T", s)
		return false
	}
	if asStmt.Name.Value != name {
		t.Errorf("asStmt.Name.Value(=Identifier) wrong, got: '%s', want: %s", asStmt.Name.Value, name)
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

func TestCallStatement(t *testing.T) {
	testCases := []struct {
		input string

		expectedIdentifier string
	}{
		{`call pipe_if_local;`, "pipe_if_local"},
	}

	for n, tc := range testCases {
		l := lexer.NewLexer(tc.input)
		p := NewParser(l)

		program := p.ParseProgram()
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements wrong number returned testCase[%d], got:%d, want%d", n, len(program.Statements), 3)
		}

		for _, stmt := range program.Statements {
			callStmt, ok := stmt.(*ast.CallStatement)
			if !ok {
				t.Fatalf("stmt not *ast.CallStatement testCase[%d], got:%T", n, stmt)
			}

			if callStmt.TokenLiteral() != "call" {
				t.Fatalf("returnStmt.TokenLiteral not 'return' testCase[%d], got:%q", n, callStmt.TokenLiteral())
			}

			if callStmt.CallValue.(*ast.Identifier).Value != tc.expectedIdentifier {
				t.Fatalf("callStmt callValue wrong in testCase[%d], got:%s, want:%s", n, callStmt.CallValue, tc.expectedIdentifier)
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
	testCases := map[string]struct {
		input       string
		expected    int64
		shouldError bool
	}{
		"with single integer":        {"5;", 5, false},
		"with invalid float integer": {"5.0;", 5, true},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			l := lexer.NewLexer(tc.input)
			p := NewParser(l)
			program := p.ParseProgram()

			if len(program.Statements) != 1 {
				if tc.shouldError {
					return
				}
				t.Fatalf("program.Statements length is not expected, got:%d, want:%d", len(program.Statements), 1)
			}

			if tc.shouldError {
				t.Fatalf("test should fail but successed")
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got:%T", program.Statements[0])
			}

			ident, ok := stmt.Expression.(*ast.IntegerLiteral)
			if !ok {
				t.Fatalf("exp not *ast.Identifier, got:%T", stmt.Expression)
			}

			if ident.Value != tc.expected {
				t.Fatalf("ident.Value wrong, got:%d, want:%d", ident.Value, 5)
			}

			if ident.TokenLiteral() != "5" {
				t.Errorf("ident.TokenLiteral wrong, got:%s, want:%s", ident.TokenLiteral(), "5")
			}
		})
	}
}

func TestPercentageLiteralExpression(t *testing.T) {
	input := "5%;"

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

	ident, ok := stmt.Expression.(*ast.PercentageLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Identifier, got:%T", stmt.Expression)
	}

	if ident.Value != "5%" {
		t.Fatalf("ident.Value wrong, got:%s, want:%s", ident.Value, "5%")
	}

	if ident.TokenLiteral() != "5%" {
		t.Errorf("ident.TokenLiteral wrong, got:%s, want:%s", ident.TokenLiteral(), "5%")
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

		if !testIntegerLiter(exp.Right, tc.rightValue) {
			t.Fatalf("parsingInfixExpression failed to get right value[testCase:%d], want:%d", n, tc.rightValue)
		}

		if !testIntegerLiter(exp.Right, tc.leftValue) {
			t.Fatalf("parsingInfixExpression failed to getleft value[testCase:%d], want:%d", n, tc.leftValue)
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

func TestOperatorPrecedenceParsing(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"1+(2+3)+4", "((1+(2+3)) + 4)"},
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
		_ = stmt
	}
}

func TestIfStatement(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{"if (x ~ y) { x }"},
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

		expr, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("program.Statement[0] is not ast.IfExpression[testCase:%d], got:%T", n, stmt.Expression)
		}

		if !testInfixExpression(t, expr.Condition, "x", "~", "y") {
			t.Fatalf("ifStatement test failed to parse condition[testCase:%d]", n)
		}

		if len(expr.Consequence.Statements) != 1 {
			t.Fatalf("consequence is not 1 statements[testCase:%d], got:%d", n, len(expr.Consequence.Statements))
		}

		consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("statement[0] in if consequence is not ast.ExpressionStatement[testCase:%d], got:%T", n, expr.Consequence.Statements[0])
		}

		if !testIdentifier(t, consequence.Expression, "x") {
			t.Fatalf("ifStatement failed to test identifier[testCase:%d]", n)
		}
	}
}

func TestIfElseStatement(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{`if (x ~ y) { 
			x 
		} else {
			y
		}`},
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

		expr, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("program.Statement[0] is not ast.IfExpression[testCase:%d], got:%T", n, stmt.Expression)
		}

		if !testInfixExpression(t, expr.Condition, "x", "~", "y") {
			t.Fatalf("ifStatement test failed to parse condition[testCase:%d]", n)
		}

		if len(expr.Consequence.Statements) != 1 {
			t.Fatalf("consequence is not 1 statements[testCase:%d], got:%d", n, len(expr.Consequence.Statements))
		}

		consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("statement[0] in if consequence is not ast.ExpressionStatement[testCase:%d], got:%T", n, expr.Consequence.Statements[0])
		}

		if !testIdentifier(t, consequence.Expression, "x") {
			t.Fatalf("ifStatement failed to test identifier[testCase:%d]", n)
		}
	}
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("expr is not ast.Identifier, got:%T", expr)
		return false
	}

	if ident.Value != value {
		t.Errorf("identifier value wrong, want:%s, got:%s", ident.Value, value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("identifier token literal wrong, want:%s, got:%s", ident.TokenLiteral(), value)
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiter(expr, int64(v))
	case int64:
		return testIntegerLiter(expr, v)
	case string:
		return testIdentifier(t, expr, v)
	}

	return false
}

func testInfixExpression(t *testing.T, expr ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expr is not ast.InfixExpression, got:%T", expr)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		t.Errorf("operationExpression.Right is not expected literal expression, right:%s, got:%s", opExp.Left, left)
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("expr.Operator is wrong, got:%s, want:%s", opExp.Operator, operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		t.Errorf("operationExpression.Right is not expected literal expression, right:%s, got:%s", opExp.Right, right)
		return false
	}

	return true
}

func TestBlockStatement(t *testing.T) {
	testCases := map[string]struct {
		input           string
		expectedLabels  []string
		blockType       string
		blockIdentifier []string
	}{
		"with single block sub":  {"sub pipe_if_local { x }", []string{"pipe_if_local"}, "sub", []string{"x"}},
		"with single block acl":  {"acl local { \"localhost\"; }", []string{"local"}, "acl", []string{"localhost"}},
		"with two statement acl": {"acl local { \"local\"; \"localhost\"}", []string{"local"}, "acl", []string{"local", "localhost"}},
		"with backend statement": {"backend server1 { .host = \"localhost\"}", []string{"server"}, "backend", []string{}},
		"with none backend":      {"backend default none;", []string{"default", "none"}, "backend", []string{}},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			l := lexer.NewLexer(tc.input)
			p := NewParser(l)
			program := p.ParseProgram()

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements wrong length got:%d, want:%d", len(program.Statements), 1)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statement[0] is not ast.ExpressionStatement, got:%T", program.Statements[0])
			}

			expr, ok := stmt.Expression.(*ast.BlockExpression)
			if !ok {
				t.Fatalf("program.Statement[0] is not ast.BlockExpression, got:%T", stmt.Expression)
			}

			if len(expr.Labels) != len(tc.expectedLabels) {
				t.Fatalf("blockExpression labels length does not match, got:%d, want:%d", len(expr.Labels), len(tc.expectedLabels))
			}

			for idx, identifier := range tc.blockIdentifier {
				block, ok := expr.Blocks.Statements[idx].(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("statement[%d] in if consequence is not ast.ExpressionStatement, got:%T", idx, expr.Blocks.Statements[0])
				}

				switch block.Expression.(type) {
				case *ast.Identifier:
					if !testIdentifier(t, block.Expression, identifier) {
						t.Fatalf("blockExpression failed to test identifier")
					}
				case *ast.StringLiteral:
					if !testStringLiteral(t, block.Expression, identifier) {
						t.Fatalf("blockExpression failed to test stringLiteral")
					}
				}
			}
		})
	}
}

func testStringLiteral(t *testing.T, expr ast.Expression, value string) bool {
	opExp, ok := expr.(*ast.StringLiteral)
	if !ok {
		t.Errorf("expr is not ast.InfixExpression, got:%T", expr)
		return false
	}

	if opExp.Value != value {
		t.Errorf("value of string literal wrong, got:%s, want:%s", opExp.Value, value)
		return false
	}

	return true
}
