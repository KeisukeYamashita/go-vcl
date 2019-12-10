package vcl

import (
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type ACL struct {
		Type      string   `vcl:"type,label"`
		Endpoints []string `vcl:"endpoints,flat"`
	}

	type Root struct {
		ACLs []*ACL `vcl:"acl,block"`
	}

	testCases := []struct {
		input       []byte
		val         interface{}
		expectedVal interface{}
	}{
		{
			[]byte(`acl local {
	"localhost"			
}`),
			&Root{},
			&Root{ACLs: []*ACL{&ACL{Type: "local", Endpoints: []string{"localhost"}}}},
		},
	}

	for n, tc := range testCases {
		err := Decode(tc.input, tc.val)
		if err != nil {
			t.Fatalf("decode failed with error[testcase:%d], error:%v", n, err)
		}

		if !reflect.DeepEqual(tc.val, tc.expectedVal) {
			t.Fatalf("decode got wrong value, got:%v, want:%v", tc.val, tc.expectedVal)
		}
	}
}
