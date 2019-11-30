package lexer

import (
	"testing"

	"github.com/KeisukeYamashita/go-vcl/vcl/token"
)

func TestNextToken(t *testing.T) {
	input := `=+,;`

	testCases := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
	}

	l := New(input)

	for i, tc := range testCases {
		tok := l.NextChar()
		if tok.Type != tc.expectedType {
			t.Fatalf("failed[testCase:%d] - wrong tokenType, want: %s, got: %s", i, tc.expectedType, tok.Type)
		}

		if tok.Literal != tc.expectedLiteral {
			t.Fatalf("failed[testCase:%d] - wrong literal, want%s, got:%s", i, tc.expectedLiteral, tok.Literal)
		}
	}
}
