package decoder

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/ast"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/schema"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/traversal"
)

// Decode is a function for mapping the program of parser output to your custom struct.
func Decode(program *ast.Program, val interface{}) []error {
	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("target value must be a pointer, not: %s", rv.Type().String()))
	}

	return decodeProgramToValue(program, rv.Elem())
}

func decodeProgramToValue(program *ast.Program, val reflect.Value) []error {
	et := val.Type()
	switch et.Kind() {
	case reflect.Struct:
		return decodeProgramToStruct(program, val)
	// TODO(KeisukeYamashita): implement map[string]interface{}
	default:
		panic(fmt.Sprintf("target value must be a pointer to struct, not: %s", et.String()))
	}
}

var attrType = reflect.TypeOf((*schema.Attribute)(nil))

func decodeProgramToStruct(program *ast.Program, val reflect.Value) []error {
	errs := []error{}
	content := traversal.Content(program)

	tags := getFieldTags(val.Type())
	for name, fieldIdx := range tags.Attributes {
		attr := content.Attributes[name]
		field := val.Type().Field(fieldIdx)
		fieldTy := field.Type
		fieldV := val.Field(fieldIdx)

		if attr == nil {
			fieldV.Set(reflect.Zero(field.Type))
			continue
		}

		switch {
		case attrType.AssignableTo(field.Type):
			fieldV.Set(reflect.ValueOf(attr))
		case fieldTy.AssignableTo(reflect.ValueOf(attr.Value).Type()):
			fieldV.Set(reflect.ValueOf(attr.Value))
		}
	}

	blocksByType := content.Blocks.ByType()

	for typeName, fieldIdx := range tags.Blocks {
		blocks := blocksByType[typeName]
		field := val.Type().Field(fieldIdx)

		ty := field.Type

		var isSlice bool
		var isPtr bool
		if ty.Kind() == reflect.Slice {
			isSlice = true
			ty = ty.Elem()
		}

		if ty.Kind() == reflect.Ptr {
			isPtr = true
			ty = ty.Elem()
		}

		if len(blocks) > 1 && !isSlice {
			errs = append(errs, errors.New("more than one block but the field type is not slice"))
		}

		if len(blocks) == 0 {
			if isSlice || isPtr {
				val.Field(fieldIdx).Set(reflect.Zero(field.Type))
			} else {
				errs = append(errs, errors.New("no block"))
			}
		}

		switch {
		case isSlice:
			elemType := ty
			if isPtr {
				elemType = reflect.PtrTo(ty)
			}

			sli := reflect.MakeSlice(reflect.SliceOf(elemType), len(blocks), len(blocks))

			for i, block := range blocks {
				if isPtr {
					v := reflect.New(ty)
					decodeBlockToStruct(block, v.Elem())
					sli.Index(i).Set(v)
				} else {
					errs = append(errs, errors.New("block is not a pointer"))
				}
			}

			val.Field(fieldIdx).Set(sli)
		default:
			if isPtr {
				v := reflect.New(ty)
				val.Field(fieldIdx).Set(v)
			} else {
				errs = append(errs, errors.New("block is not a pointer"))
			}
		}
	}

	return nil
}

func decodeBlockToStruct(block *schema.Block, val reflect.Value) {
	var errs []error
	tags := getFieldTags(val.Type())
	content := traversal.BodyContent(block.Body)

	for name, fieldIdx := range tags.Attributes {
		attr := content.Attributes[name]
		field := val.Type().Field(fieldIdx)
		fieldTy := field.Type
		fieldV := val.Field(fieldIdx)

		if attr == nil {
			fieldV.Set(reflect.Zero(field.Type))
			continue
		}

		switch {
		case attrType.AssignableTo(field.Type):
			fieldV.Set(reflect.ValueOf(attr))
		case fieldTy.AssignableTo(reflect.ValueOf(attr.Value).Type()):
			fieldV.Set(reflect.ValueOf(attr.Value))
		}
	}

	blocksByType := content.Blocks.ByType()

	for typeName, fieldIdx := range tags.Blocks {
		blocks := blocksByType[typeName]
		field := val.Type().Field(fieldIdx)

		ty := field.Type

		var isSlice bool
		var isPtr bool
		if ty.Kind() == reflect.Slice {
			isSlice = true
			ty = ty.Elem()
		}

		if ty.Kind() == reflect.Ptr {
			isPtr = true
			ty = ty.Elem()
		}

		if len(blocks) > 1 && !isSlice {
			errs = append(errs, errors.New("more than one block but the field type is not slice"))
		}

		if len(blocks) == 0 {
			if isSlice || isPtr {
				val.Field(fieldIdx).Set(reflect.Zero(field.Type))
			} else {
				errs = append(errs, errors.New("no block"))
			}
		}

		switch {
		case isSlice:
			elemType := ty
			if isPtr {
				elemType = reflect.PtrTo(ty)
			}

			sli := reflect.MakeSlice(reflect.SliceOf(elemType), len(blocks), len(blocks))

			for i, block := range blocks {
				if isPtr {
					v := reflect.New(ty)
					decodeBlockToStruct(block, v.Elem())
					sli.Index(i).Set(v)
				} else {
					errs = append(errs, errors.New("block is not a pointer"))
				}
			}

			val.Field(fieldIdx).Set(sli)
		default:
			if isPtr {
				v := reflect.New(ty)
				val.Field(fieldIdx).Set(v)
			} else {
				errs = append(errs, errors.New("block is not a pointer"))
			}
		}
	}

	return
}

// imipliedBodySchema will retrieves the root body schema from the given val.
// For Varnish & Fastly usecases, there will be only blocks in the root. But as a configuration language,
// the root schema can contain attribute as HCL. Therefore, I left the attributes slice for that.
func impliedBodySchema(val interface{}) *schema.File {
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
		attr := tags.Attributes[n]
		field := ty.Field(attr)
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

	file := &schema.File{
		Body: &schema.BodySchema{
			Attributes: attrSchemas,
			Blocks:     blockSchemas,
		},
	}

	return file
}

// fieldTags is a struct that represents info about the field of the passed val.
type fieldTags struct {
	Attributes map[string]int
	Blocks     map[string]int
	Labels     []labelField
	Flats      []flatField
}

// labelField is a struct that represents info about the struct tags of "vcl".
type labelField struct {
	FieldIndex int
	Name       string
}
type flatField struct {
	FieldIndix int
	Name       string
}

// getFieldTags retrieves the "vcl" tags of the given struct type.
func getFieldTags(ty reflect.Type) *fieldTags {
	ret := &fieldTags{
		Attributes: map[string]int{},
		Blocks:     map[string]int{},
		Labels:     []labelField{},
		Flats:      []flatField{},
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
		case "flat":
			ret.Flats = append(ret.Flats, flatField{
				FieldIndix: i,
				Name:       name,
			})
		default:
			panic(fmt.Sprintf("invalid vcl field tag kind %q on %s %q", kind, field.Type.String(), field.Name))
		}
	}

	return ret
}
