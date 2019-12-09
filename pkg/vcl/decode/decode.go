package decoder

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/ast"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/schema"
)

// Decode is a function for mapping the program of parser output to your custom struct.
func Decode(program ast.Program, val interface{}) error {
	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("target value must be a pointer, not: %s", rv.Type().String()))
	}

	return decodeProgramToValue(program, rv.Elem())
}

func decodeProgramToValue(program ast.Program, val reflect.Value) error {
	et := val.Type()
	switch et.Kind() {
	case reflect.Struct:
		return decodeProgramToStruct(program, val)
	default:
		panic(fmt.Sprintf("target value must be a pointer to struct, not: %s", et.String()))
	}
}

func decodeProgramToStruct(program ast.Program, val reflect.Value) error {
	schema := impliedBodySchema(val.Interface())
	_ = schema
	return nil
}

// imipliedBodySchema will retrieves the root body schema from the given val.
// For Varnish & Fastly usecases, there will be only blocks in the root. But as a configuration language,
// the root schema can contain attribute as HCL. Therefore, I left the attributes slice for that.
func impliedBodySchema(val interface{}) *schema.Schema {
	ty := reflect.TypeOf(val)
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}

	if ty.Kind() != reflect.Struct {
		panic(fmt.Sprintf("target value must be a struct, not: %T", val))
	}

	var attrSchemas []schema.AttributeSchema
	var blockSchemas []schema.BlockHeaderSchema

	tags := getFieldTags(ty)
	attrNames := make([]string, 0, len(tags.Attributes))
	for n := range tags.Attributes {
		attrNames = append(attrNames, n)
	}

	sort.Strings(attrNames)
	for _, n := range attrNames {
		idx := tags.Attributes[n]
		field := ty.Field(idx)
		var required bool

		switch {
		case field.Type.Kind() != reflect.Ptr:
			required = true
		}

		attrSchemas = append(attrSchemas, schema.AttributeSchema{
			Name:     n,
			Required: required,
		})
	}

	blockNames := make([]string, 0, len(tags.Blocks))
	for n := range tags.Blocks {
		blockNames = append(blockNames, n)
	}

	sort.Strings(blockNames)
	for _, n := range blockNames {
		idx := tags.Blocks[n]
		field := ty.Field(idx)
		fty := field.Type
		if fty.Kind() == reflect.Ptr {
			fty = fty.Elem()
		}

		if fty.Kind() != reflect.Struct {
			panic(fmt.Sprintf("hcl 'block' tag kind cannot be applied to %s field %s: struct required", field.Type.String(), field.Name))
		}

		ftags := getFieldTags(fty)
		var labelNames []string
		if len(ftags.Labels) > 0 {
			labelNames = make([]string, len(ftags.Labels))
			for i, l := range ftags.Labels {
				labelNames[i] = l.Name
			}
		}

		blockSchemas = append(blockSchemas, schema.BlockHeaderSchema{
			Type:       n,
			LabelNames: labelNames,
		})
	}

	schema := &schema.Schema{
		Attributes: attrSchemas,
		Blocks:     blockSchemas,
	}

	return schema
}

// fieldTags is a struct that represents info about the field of the passed val.
type fieldTags struct {
	Attributes map[string]int
	Blocks     map[string]int
	Labels     []labelField
}

// labelField is a struct that represents info about the struct tags of "vcl".
type labelField struct {
	FieldIndex int
	Name       string
}

// getFieldTags retrieves the "vcl" tags of the given struct type.
func getFieldTags(ty reflect.Type) *fieldTags {
	ret := &fieldTags{
		Attributes: map[string]int{},
		Blocks:     map[string]int{},
		Labels:     []labelField{},
	}

	ct := ty.NumField()
	for i := 0; i < ct; i++ {
		field := ty.Field(i)
		tag := field.Tag.Get("vcl")
		if tag == "" {
			continue
		}

		comma := strings.Index(tag, ",")
		var name, kind string
		if comma != -1 {
			name = tag[:comma]
			kind = tag[comma+1:]
		} else {
			name = tag
			kind = "attr"
		}

		switch kind {
		case "attr":
			ret.Attributes[name] = i
		case "block":
			ret.Blocks[name] = i
		case "label":
			ret.Labels = append(ret.Labels, labelField{
				FieldIndex: i,
				Name:       name,
			})
		default:
			panic(fmt.Sprintf("invalid vcl field tag kind %q on %s %q", kind, field.Type.String(), field.Name))
		}
	}

	return ret
}
