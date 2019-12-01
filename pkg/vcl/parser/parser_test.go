package parser

import (
	"testing"

	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/ast"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/lexer"
)

func TestAssignStatement(t *testing.T) {
	testCases := []struct {
		input               string
		expectedIdentifiers []struct {
			expectedIdentifier string
		}
	}{
		{
			`x = 10
y = 20
keke = 20
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
		if len(program.Statements) != len(tc.expectedIdentifiers) {
			t.Fatalf("program.Statements wrong number returned testCase[%d], got:%d, want%d", n, len(program.Statements), len(tc.expectedIdentifiers))
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
