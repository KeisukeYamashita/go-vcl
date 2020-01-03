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

	testCases := map[string]struct {
		input       interface{}
		program     *ast.Program
		shouldError bool
	}{
		"with pointer":     {testStruct, prog, false},
		"with not-pointer": {*testStruct, prog, true},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			if errs := Decode(tc.program, tc.input); len(errs) > 0 {
				if tc.shouldError {
					return
				}

				t.Fatalf("decode failed with error: %v", errs)
			}

			if tc.shouldError {
				t.Fatalf("decode should failed but successed")
			}
		})
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
		Name      string   `vcl:"name,label"`
		Endpoints []string `vcl:",flat"`
	}

	type Sub struct {
		Type      string   `vcl:"type,label"`
		Endpoints []string `vcl:",flat"` // Memo(KeisukeYamashita): Wont test inside of the block
	}

	type SubObj struct {
		Type string `vcl:"type,label"`
		Name string `vcl:"name,label"`
		Host string `vcl:".host"`
		IP   string `vcl:".ip"`
	}

	type RootSub struct {
		Subs []*SubObj `vcl:"sub,block"`
	}

	type Root struct {
		ACLs []*ACL `vcl:"acl,block"`
		Subs []*Sub `vcl:"sub,block"`
	}

	testCases := map[string]struct {
		input    string
		val      interface{}
		expected interface{}
	}{
		"with single block": {
			`acl local {
	"local";
	"localhost";
}`, &Root{}, &Root{Subs: []*Sub{}, ACLs: []*ACL{&ACL{Type: "local", Endpoints: []string{"local", "localhost"}}}},
		},
		"with two same block": {
			`acl local {
	"local";	
	"localhost";
}

acl remote {
	"remote";
}
`, &Root{}, &Root{Subs: []*Sub{}, ACLs: []*ACL{&ACL{Type: "local", Endpoints: []string{"local", "localhost"}}, &ACL{Type: "remote", Endpoints: []string{"remote"}}}},
		},
		"with two mixed block type": {
			`acl local {
	"local";	
	"localhost";
}

sub pipe_something {
	"inside_sub";
	"34.100.0.0"/23;
}
`, &Root{}, &Root{ACLs: []*ACL{&ACL{Type: "local", Endpoints: []string{"local", "localhost"}}}, Subs: []*Sub{&Sub{Type: "pipe_something", Endpoints: []string{"inside_sub", "\"34.100.0.0\"/23"}}}},
		},
		"with sub block": {
			`sub pipe_something {
	.host = "host";
	.ip = "ip";
}
`, &RootSub{}, &RootSub{Subs: []*SubObj{&SubObj{Type: "pipe_something", Host: "host", IP: "ip"}}},
		},
		"with multi label": {
			`sub pipe_something pipe_keke {
	.host = "host";
	.ip = "ip";
}
`, &RootSub{}, &RootSub{Subs: []*SubObj{&SubObj{Type: "pipe_something", Name: "pipe_keke", Host: "host", IP: "ip"}}},
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()
			root := tc.val
			val := reflect.ValueOf(root).Elem()
			errs := decodeProgramToStruct(program, val)

			if len(errs) > 0 {
				t.Fatalf("decodeProgramToStruct_Block has errorr, err:%v", errs)
			}

			if !reflect.DeepEqual(tc.val, tc.expected) {
				t.Fatalf("decodeProgramToStruct_Block got wrong result, got:%#v, want:%#v", tc.val, tc.expected)
			}
		})
	}
}

func TestDecodeProgramToStruct_DirectorBlock(t *testing.T) {
	type Backend struct {
		Backend string `vcl:".backend"`
		Weight  int64  `vcl:".weight"`
	}

	type Director struct {
		Type     string     `vcl:"type,label"`
		Name     string     `vcl:"name,label"`
		Quorum   string     `vcl:".quorum"`
		Retries  int64      `vcl:".retries"`
		Backends []*Backend `vcl:",flat"`
	}

	type Root struct {
		Directors []*Director `vcl:"director,block"`
	}

	testCases := map[string]struct {
		input    string
		val      interface{}
		expected interface{}
	}{
		"with single director block": {
			`director my_dir random {
				.quorum = 50%;
				.retries = 3;
			}`, &Root{}, &Root{Directors: []*Director{&Director{Type: "my_dir", Name: "random", Quorum: "50%", Retries: 3, Backends: []*Backend{}}}},
		},
		"with deep director block": {
			`director my_dir random {
				.quorum = 50%;
				.retries = 3;
				{ .backend = K_backend1; .weight = 1; }
			}`, &Root{}, &Root{Directors: []*Director{&Director{Type: "my_dir", Name: "random", Quorum: "50%", Retries: 3, Backends: []*Backend{&Backend{Backend: "K_backend1", Weight: 1}}}}},
		},
		"with multiple deep director block": {
			`director my_dir random {
				.quorum = 50%;
				.retries = 3;
				{ .backend = K_backend1; .weight = 1; }
				{ .backend = E_backend1; .weight = 3; }
			}`, &Root{}, &Root{Directors: []*Director{&Director{Type: "my_dir", Name: "random", Quorum: "50%", Retries: 3, Backends: []*Backend{&Backend{Backend: "K_backend1", Weight: 1}, &Backend{Backend: "E_backend1", Weight: 3}}}}},
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()
			root := tc.val
			val := reflect.ValueOf(root).Elem()
			errs := decodeProgramToStruct(program, val)

			if len(errs) > 0 {
				t.Fatalf("decodeProgramToStruct_Block has errorr, err:%v", errs)
			}

			if !reflect.DeepEqual(tc.val, tc.expected) {
				t.Fatalf("decodeProgramToStruct_Block got wrong result, got:%#v, want:%#v", tc.val, tc.expected)
			}
		})
	}
}

func TestDecodeProgramToStruct_TableBlock(t *testing.T) {
	type Table struct {
		Type     string `vcl:"type,label"`
		Username string `vcl:"username"`
	}

	type Root struct {
		Tables []*Table `vcl:"table,block"`
	}

	testCases := map[string]struct {
		input    string
		val      interface{}
		expected interface{}
	}{
		"with single table block": {
			`table my_id {
	"username": "keke"
}`, &Root{}, &Root{[]*Table{&Table{Type: "my_id", Username: "keke"}}},
		},
		"with multiple table block": {
			`table my_id {
	"username": "keke"
}

table my_keke {
	"username": "kekekun",
}`, &Root{}, &Root{[]*Table{&Table{Type: "my_id", Username: "keke"}, &Table{Type: "my_keke", Username: "kekekun"}}},
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()
			root := tc.val
			val := reflect.ValueOf(root).Elem()
			errs := decodeProgramToStruct(program, val)

			if len(errs) > 0 {
				t.Fatalf("decodeProgramToStruct_Block has errorr, err:%v", errs)
			}

			if !reflect.DeepEqual(tc.val, tc.expected) {
				t.Fatalf("decodeProgramToStruct_Block got wrong result, got:%#v, want:%#v", tc.val, tc.expected)
			}
		})
	}
}

func TestDecodeProgramToStruct_NestedBlock(t *testing.T) {
	type Probe struct {
		X int64 `vcl:"x"`
	}

	type Backend struct {
		Type  string `vcl:"type,label"`
		IP    string `vcl:".ip"`
		Probe *Probe `vcl:".probe,block"`
	}

	type Root struct {
		Backends []*Backend `vcl:"backend,block"`
	}

	testCases := map[string]struct {
		input    string
		val      interface{}
		expected interface{}
	}{
		"with nested simple block": {
			`backend remote {
	.ip = "localhost";
	.probe = {
		x = 10;
	};
}`, &Root{}, &Root{Backends: []*Backend{&Backend{Type: "remote", IP: "localhost", Probe: &Probe{X: 10}}}},
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()
			root := tc.val
			val := reflect.ValueOf(root).Elem()
			errs := decodeProgramToStruct(program, val)

			if len(errs) > 0 {
				t.Fatalf("decodeProgramToStruct_Block has errorr, err:%v", errs)
			}

			if !reflect.DeepEqual(tc.val, tc.expected) {
				t.Fatalf("decodeProgramToStruct_Block got wrong result, got:%#v", tc.val)
			}
		})
	}
}

func TestDecodeProgramToStruct_Comments(t *testing.T) {
	type ACL struct {
		Type     string   `vcl:"type,label"`
		Comments []string `vcl:",comment"`
	}

	type Root struct {
		ACLs     []*ACL   `vcl:"acl,block"`
		Comments []string `vcl:",comment"`
	}

	testCases := map[string]struct {
		input    string
		val      interface{}
		expected interface{}
	}{
		"with root comment by hash": {
			`# keke`, &Root{}, &Root{Comments: []string{"keke"}, ACLs: []*ACL{}},
		},
		"with root comment by double slash": {
			`// keke`, &Root{}, &Root{Comments: []string{"keke"}, ACLs: []*ACL{}},
		},
		"with root by double slash with block": {
			`// keke
acl "tag" {}		
`, &Root{}, &Root{Comments: []string{"keke"}, ACLs: []*ACL{&ACL{Type: "tag", Comments: []string{}}}},
		},
		"with nested block": {
			`// keke
acl "tag" {
	// internal-keke
	"localhost";
}		
`, &Root{}, &Root{Comments: []string{"keke"}, ACLs: []*ACL{&ACL{Type: "tag", Comments: []string{"internal-keke"}}}},
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()
			root := tc.val
			val := reflect.ValueOf(root).Elem()
			errs := decodeProgramToStruct(program, val)

			if len(errs) > 0 {
				t.Fatalf("decodeProgramToStruct_Block has errorr, err:%v", errs)
			}

			if !reflect.DeepEqual(tc.val, tc.expected) {
				t.Fatalf("decodeProgramToStruct_Block got wrong result, got:%#v", tc.val)
			}
		})
	}
}

func TestDecodeProgramToMap(t *testing.T) {
	testCases := map[string]struct {
		input    string
		val      map[string]interface{}
		expected map[string]interface{}
	}{
		"with single attr": {`x = hello`, map[string]interface{}{}, map[string]interface{}{"x": "hello"}},
		"with multiple attr": {`x = hello;
y = bye`, map[string]interface{}{}, map[string]interface{}{"x": "hello", "y": "bye"}},
		"with single block": {`acl hello {x = "test"}`, map[string]interface{}{}, map[string]interface{}{"acl": map[string]interface{}{"hello": map[string]interface{}{"x": "test"}}}},
		"with multiple block": {`acl hello {
	x = "test";
}

acl bye {
	y = "keke";
}
`, map[string]interface{}{}, map[string]interface{}{"acl": map[string]interface{}{"hello": map[string]interface{}{"x": "test"}, "bye": map[string]interface{}{"y": "keke"}}}},
		"with flat block": {`acl hello {
	"localhost";
	"local";
}`, map[string]interface{}{}, map[string]interface{}{"acl": map[string]interface{}{"hello": []interface{}{"localhost", "local"}}}},
		"with dot attribute block": {`backend default {
	.port = "8080";
}`, map[string]interface{}{}, map[string]interface{}{"backend": map[string]interface{}{"default": map[string]interface{}{"port": "8080"}}}}}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()
			val := reflect.ValueOf(&tc.val).Elem()
			errs := decodeProgramToMap(program, val)

			if len(errs) > 0 {
				t.Fatalf("decodeProgramToStruct has errors, err:%v", errs)
			}

			if !reflect.DeepEqual(&tc.val, &tc.expected) {
				t.Fatalf("decodeProgramToStruct got wrong result got:%v want:%v", tc.val, tc.expected)
			}
		})
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
		Flats    interface{} `vcl:",flat"`
		Comments interface{} `vcl:",comment"`
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

		if len(tags.Flats) != 1 {
			t.Fatalf("Flats length wrong[testCase:%d], got:%d, want:%d", n, len(tags.Flats), 1)
		}

		if len(tags.Comments) != 1 {
			t.Fatalf("Comments length wrong[testCase:%d], got:%d, want:%d", n, len(tags.Comments), 1)
		}
	}
}
