package schema

// Schema ...
type Schema struct {
	Attributes []AttributeSchema
	Blocks     []BlockHeaderSchema
}

// AttributeSchema ...
type AttributeSchema struct {
	Name string
}

// BlockHeaderSchema ...
type BlockHeaderSchema struct {
	Type       string
	LabelNames []string
}
