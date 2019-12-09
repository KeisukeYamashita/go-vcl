package schema

// File is the root of the data structure
type File struct {
	Body *Body
}

// Body contains multiple attributes and blocks
type Body struct {
	Attributes []AttributeSchema
	Blocks     []BlockHeaderSchema
}

// AttributeSchema ...
type AttributeSchema struct {
	Name     string
	Required bool
}

// BlockHeaderSchema ...
type BlockHeaderSchema struct {
	Type       string
	LabelNames []string
	Body       *Body
}
