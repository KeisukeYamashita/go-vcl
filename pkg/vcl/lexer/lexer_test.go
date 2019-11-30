package lexer

import (
	"testing"

	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/token"
)

func TestNextToken(t *testing.T) {
	testCases := []struct {
		input          string
		expectedTokens []struct {
			expectedType    token.Type
			expectedLiteral string
		}
	}{
		{
			"=~,; call == && || 10",
			[]struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.ASSIGN, "="},
				{token.MATCH, "~"},
				{token.COMMA, ","},
				{token.SEMICOLON, ";"},
				{token.CALL, "call"},
				{token.EQUAL, "=="},
				{token.AND, "&&"},
				{token.OR, "||"},
				{token.INT, "10"},
			},
		},
		{
			`sub pipe_if_local {
	if (client.ip ~ local) {
		return (pipe);
	}
}
`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.SUBROUTINE, "sub"},
				{token.INDENT, "pipe_if_local"},
				{token.LBRACE, "{"},
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.INDENT, "client.ip"},
				{token.MATCH, "~"},
				{token.INDENT, "local"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.LPAREN, "("},
				{token.INDENT, "pipe"},
				{token.RPAREN, ")"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.RBRACE, "}"},
			},
		},
		{
			"import directors; # load the directors",
			[]struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.IMPORT, "import"},
				{token.INDENT, "directors"},
				{token.SEMICOLON, ";"},
				{token.COMMENT, "#"},
				{token.INDENT, "load"},
				{token.INDENT, "the"},
				{token.INDENT, "directors"},
			},
		},
	}

	for i, tc := range testCases {
		l := New(tc.input)

		for j, expectedToken := range tc.expectedTokens {
			tok := l.NextToken()
			if tok.Type != expectedToken.expectedType {
				t.Fatalf("failed[testCase:%d:%d] - wrong tokenType, want: %s, got: %s", i+1, j+1, expectedToken.expectedType, tok.Type)
			}

			if tok.Literal != expectedToken.expectedLiteral {
				t.Fatalf("failed[testCase:%d:%d] - wrong literal, want%s, got:%s", i+1, j+1, expectedToken.expectedLiteral, tok.Literal)
			}
		}
	}
}
