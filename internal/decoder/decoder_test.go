package decoder

import (
	"reflect"
	"testing"

	"github.com/KeisukeYamashita/go-vcl/internal/ast"
	"github.com/KeisukeYamashita/go-vcl/internal/lexer"
	"github.com/KeisukeYamashita/go-vcl/internal/parser"
	"github.com/KeisukeYamashita/go-vcl/internal/schema"
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

func TestDecodeProgramToStruct_Attribute(t *testing.T) {
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

		if !reflect.DeepEqual(tc.val, tc.expected) {
			t.Fatalf("decodeProgramToStruct got wrong result[testCase:%d]", n)
		}
	}
}

func TestDecodeProgramToStruct_Block(t *testing.T) {
	type ACL struct {
		Type      string   `vcl:"type,label"`
		Endpoints []string `vcl:"endpoints,flat"`
	}

	type Root struct {
		ACLs []*ACL `vcl:"acl,block"`
	}

	testCases := []struct {
		input    string
		val      interface{}
		expected interface{}
	}{
		{`acl local {
	"local";
	"localhost";
}`, &Root{}, &Root{ACLs: []*ACL{&ACL{Type: "local", Endpoints: []string{"local", "localhost"}}}}},
	}

	for n, tc := range testCases {
		l := lexer.NewLexer(tc.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		root := tc.val
		val := reflect.ValueOf(root).Elem()
		errs := decodeProgramToStruct(program, val)

		if len(errs) > 0 {
			t.Fatalf("decodeProgramToStruct_Block has errors[testCase:%d], err:%v", n, errs)
		}

		if !reflect.DeepEqual(tc.val, tc.expected) {
			t.Fatalf("decodeProgramToStruct_Block got wrong result[testCase:%d]", n)
		}
	}
}

func TestBlockToStruct(t *testing.T) {

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
