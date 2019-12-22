package schema

// File is the root of the data structure
type File struct {
	Body Body
}

// Body contains multiple attributes and blocks
type Body interface{}

// Attributes are like field of objects
type Attributes map[string]*Attribute

// Blocks are block types containing other block
type Blocks []*Block

// Flats are attributes that are not key-value
type Flats []interface{}

// Comments are comment lines
type Comments []string

// Block ais a structure which contains block header, labels and body
type Block struct {
	Type   string
	Labels []string
	Body   Body
}

// BodySchema represents the desired structure of a body.
type BodySchema struct {
	Attributes []AttributeSchema
	Blocks     []BlockHeaderSchema
}

// Attribute are field of the object
type Attribute struct {
	Name  string
	Value interface{}
}

// BodyContent is a content from body
type BodyContent struct {
	Attributes Attributes
	Blocks     Blocks
	Flats      Flats
	Comments   Comments
}

// AttributeSchema is the desired attribute
type AttributeSchema struct {
	Name     string
	Required bool
}

// BlockHeaderSchema is the desired block
type BlockHeaderSchema struct {
	Type       string
	LabelNames []string
	Body       Body
}

// ByType transforms the receiving block sequence into a map from type
// name to block sequences of only that type.
func (bs Blocks) ByType() map[string]Blocks {
	ret := make(map[string]Blocks)
	for _, b := range bs {
		ty := b.Type
		if ret[ty] == nil {
			ret[ty] = make(Blocks, 0, 1)
		}

		ret[ty] = append(ret[ty], b)
	}

	return ret
}
