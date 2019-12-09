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
		content := Content(program)
		if len(content.Attributes) != tc.expectedAttrCount {
			t.Fatalf("contents.Attributes length failed[testcase:%d], got:%d, want:%d", n, len(content.Attributes), tc.expectedAttrCount)
		}

		if len(content.Blocks) != tc.expectedBlockCount {
			t.Fatalf("contents.Blocks length failed[testcase:%d], got:%d, want:%d", n, len(content.Blocks), tc.expectedBlockCount)
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
		content := convertBody(program.Statements)
		if len(content.Attributes) != tc.expectedAttrCount {
			t.Fatalf("contents.Attributes length failed[testcase:%d], got:%d, want:%d", n, len(content.Attributes), tc.expectedAttrCount)
		}

		if len(content.Blocks) != tc.expectedBlockCount {
			t.Fatalf("contents.Blocks length failed[testcase:%d], got:%d, want:%d", n, len(content.Blocks), tc.expectedBlockCount)
		}
	}
}
