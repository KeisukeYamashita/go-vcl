package lexer

import (
	"testing"

	"github.com/KeisukeYamashita/go-vcl/internal/token"
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
			`=~,; call == && || 10 "keke" false ! "35.0.0.0"/23; server1 K_backend1 50% table`,
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
				{token.STRING, "keke"},
				{token.FALSE, "false"},
				{token.BANG, "!"},
				{token.CIDR, "\"35.0.0.0\"/23"},
				{token.IDENT, "server1"},
				{token.IDENT, "K_backend1"},
				{token.PERCENTAGE, "50%"},
				{token.TABLE, "table"},
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
				{token.IDENT, "pipe_if_local"},
				{token.LBRACE, "{"},
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.IDENT, "client.ip"},
				{token.MATCH, "~"},
				{token.IDENT, "local"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.LPAREN, "("},
				{token.IDENT, "pipe"},
				{token.RPAREN, ")"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.RBRACE, "}"},
			},
		},
		{
			`director my_dir random {
				.retries = 3;
			}`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.DIRECTOR, "director"},
				{token.IDENT, "my_dir"},
				{token.IDENT, "random"},
				{token.LBRACE, "{"},
				{token.IDENT, ".retries"},
				{token.ASSIGN, "="},
				{token.INT, "3"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
			},
		},
		{
			`table my_id {
				"key1": "value 1",
			}`,
			[]struct {
				expectedType    token.Type
				expectedLiteral string
			}{
				{token.TABLE, "table"},
				{token.IDENT, "my_id"},
				{token.LBRACE, "{"},
				{token.STRING, "key1"},
				{token.COLON, ":"},
				{token.STRING, "value 1"},
				{token.COMMA, ","},
				{token.RBRACE, "}"},
			},
		},
	}

	for i, tc := range testCases {
		l := NewLexer(tc.input)

		for j, expectedToken := range tc.expectedTokens {
			tok := l.NextToken()
			if tok.Type != expectedToken.expectedType {
				t.Fatalf("failed[testCase:%d:%d] - wrong tokenType, want: %s(literal:%s), got: %s(literal:%s)", i+1, j+1, expectedToken.expectedType, expectedToken.expectedLiteral, tok.Type, tok.Literal)
			}

			if tok.Literal != expectedToken.expectedLiteral {
				t.Fatalf("failed[testCase:%d:%d] - wrong literal, want: %s, got: %s", i+1, j+1, expectedToken.expectedLiteral, tok.Literal)
			}
		}
	}
}
