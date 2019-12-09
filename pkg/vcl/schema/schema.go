package schema

// Schema ...
type Schema struct {
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
}
