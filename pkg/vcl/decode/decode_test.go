package decoder

import (
	"reflect"
	"testing"

	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/ast"
)

func TestDecode(t *testing.T) {
	type TestStruct struct {
		Name string
	}

	testStruct := &TestStruct{}
	prog := ast.Program{}

	testCases := []struct {
		input   interface{}
		program ast.Program
	}{
		{testStruct, prog},
	}

	for n, tc := range testCases {
		Decode(tc.program, tc.input)
		_ = n
		// TODO(KeisukeYamashita): Add test for decode
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
		schema := impliedBodySchema(tc.input)
		if len(schema.Attributes) != 1 {
			t.Fatalf("Attribute length wrong[testCase:%d], got:%d, want:%d", n, len(schema.Attributes), 1)
		}

		if len(schema.Blocks) != 1 {
			t.Fatalf("Block length wrong[testCase:%d], got:%d, want:%d", n, len(schema.Blocks), 1)
		}

		if len(schema.Blocks[0].LabelNames) != 1 {
			t.Fatalf("Block label are not expected[testCase:%d], got:%d, want:%d", n, len(schema.Blocks[0].LabelNames), 1)
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
