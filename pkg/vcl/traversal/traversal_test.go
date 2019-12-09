package traversal

import (
	"testing"

	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/lexer"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/parser"
)

func TestContents(t *testing.T) {
	testCases := []struct {
		input              string
		expectedAttrCount  int
		expectedBlockCount int
	}{
		{
			`x = 10`,
			1,
			0,
		},
		{
			`acl type name {
				"local"	
			}`,
			0,
			1,
		}, {
			`sub pipe_if_local { x }`,
			0,
			1,
		},
	}

	for n, tc := range testCases {
		l := lexer.NewLexer(tc.input)
		p := parser.NewParser(l)

		program := p.ParseProgram()
		contents := Contents(program)
		if len(contents.Body.Attributes) != tc.expectedAttrCount {
			t.Fatalf("contents.Attributes length failed[testcase:%d], got:%d, want:%d", n, len(contents.Body.Attributes), tc.expectedAttrCount)
		}

		if len(contents.Body.Blocks) != tc.expectedBlockCount {
			t.Fatalf("contents.Blocks length failed[testcase:%d], got:%d, want:%d", n, len(contents.Body.Blocks), tc.expectedBlockCount)
		}
	}
}

func TestConvertBody(t *testing.T) {
	testCases := []struct {
		input              string
		expectedAttrCount  int
		expectedBlockCount int
	}{
		{
			`x = 10`,
			1,
			0,
		},
		{
			`acl type name {
				"local"	
			}`,
			0,
			1,
		}, {
			`sub pipe_if_local { x }`,
			0,
			1,
		},
	}

	for n, tc := range testCases {
		l := lexer.NewLexer(tc.input)
		p := parser.NewParser(l)

		program := p.ParseProgram()

		body := convertBody(program.Statements)
		if len(body.Attributes) != tc.expectedAttrCount {
			t.Fatalf("contents.Attributes length failed[testcase:%d], got:%d, want:%d", n, len(body.Attributes), tc.expectedAttrCount)
		}

		if len(body.Blocks) != tc.expectedBlockCount {
			t.Fatalf("contents.Blocks length failed[testcase:%d], got:%d, want:%d", n, len(body.Blocks), tc.expectedBlockCount)
		}
	}
}
