package decoder

import (
	"reflect"
	"testing"

	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/ast"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/lexer"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/parser"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/schema"
)

func TestDecode(t *testing.T) {
	type TestStruct struct {
		Name string
	}

	testStruct := &TestStruct{}
	prog := &ast.Program{}

	testCases := []struct {
		input   interface{}
		program *ast.Program
	}{
		{testStruct, prog},
	}

	for n, tc := range testCases {
		Decode(tc.program, tc.input)
		_ = n
		// TODO(KeisukeYamashita): Add test for decode
	}
}

func TestDecodeProgramToStruct(t *testing.T) {
	type Root struct {
		X   int64  `vcl:"x"`
		API string `vcl:"api"`
	}

	testCases := []struct {
		input    string
		val      interface{}
		expected interface{}
	}{
		{`x = 1`, &Root{}, &Root{X: 1}},
		{`api = "localhost"`, &Root{}, &Root{API: "localhost"}},
	}

	for n, tc := range testCases {
		l := lexer.NewLexer(tc.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		root := tc.val
		val := reflect.ValueOf(root).Elem()
		errs := decodeProgramToStruct(program, val)

		if len(errs) > 0 {
			t.Fatalf("decodeProgramToStruct has errors[testCase:%d], err:%v", n, errs)
		}
	}
}

func TestImpliedBodySchema(t *testing.T) {
	type testBlock struct {
		Type       string `vcl:"type,label"`
		MiddelName string `vcl:"middelname"`
	}

	type testStruct struct {
		Type     string     `vcl:"type,label"`
		Name     string     `vcl:"name"`
		Resource *testBlock `vcl:"resource,block"`
	}

	input := &testStruct{
		Type: "my-type",
		Name: "keke",
		Resource: &testBlock{
			MiddelName: "middelName",
		},
	}

	testCases := []struct {
		input interface{}
	}{
		{input},
	}

	for n, tc := range testCases {
		file := impliedBodySchema(tc.input)
		bs := file.Body.(*schema.BodySchema)
		if len(bs.Attributes) != 1 {
			t.Fatalf("Attribute length wrong[testCase:%d], got:%d, want:%d", n, len(bs.Attributes), 1)
		}

		if len(bs.Blocks) != 1 {
			t.Fatalf("Block length wrong[testCase:%d], got:%d, want:%d", n, len(bs.Blocks), 1)
		}

		if len(bs.Blocks[0].LabelNames) != 1 {
			t.Fatalf("Block label are not expected[testCase:%d], got:%d, want:%d", n, len(bs.Blocks[0].LabelNames), 1)
		}
	}
}

func TestGetFieldTags(t *testing.T) {
	type testStruct struct {
		Type     string      `vcl:"type,label"`
		Name     string      `vcl:"name"` // implied attribute
		Resource interface{} `vcl:"resource,block"`
	}

	input := &testStruct{
		Type:     "my-type",
		Name:     "keke",
		Resource: "",
	}

	testCases := []struct {
		input *testStruct
	}{
		{input},
	}

	for n, tc := range testCases {
		ty := reflect.TypeOf(*tc.input)
		tags := getFieldTags(ty)

		if len(tags.Attributes) != 1 {
			t.Fatalf("Attribute length wrong[testCase:%d], got:%d, want:%d", n, len(tags.Attributes), 1)
		}

		if len(tags.Labels) != 1 {
			t.Fatalf("Labels length wrong[testCase:%d], got:%d, want:%d", n, len(tags.Labels), 1)
		}

		if len(tags.Blocks) != 1 {
			t.Fatalf("Blocks length wrong[testCase:%d], got:%d, want:%d", n, len(tags.Blocks), 1)
		}

	}
}
