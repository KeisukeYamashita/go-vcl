package traversal

import (
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/ast"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/schema"
)

// Contents retrives from ast.Program
func Contents(prog *ast.Program) *schema.File {
	b := convertBody(prog.Statements)
	return &schema.File{
		Body: b,
	}
}

// Contents will ast.Program to schema
func convertBody(stmts []ast.Statement) *schema.Body {
	var attrs []schema.AttributeSchema
	var blocks []schema.BlockHeaderSchema

	for _, stmt := range stmts {
		switch v := stmt.(type) {
		case *ast.AssignStatement:
			attrs = append(attrs, schema.AttributeSchema{
				Name: v.Name.Value,
			})
		case *ast.ExpressionStatement:
			switch expr := v.Expression.(type) {
			case *ast.BlockExpression:
				body := convertBody(expr.Blocks.Statements)
				block := schema.BlockHeaderSchema{
					Body: body,
				}

				var blockType string
				var labels []string
				if len(expr.Labels) > 0 {
					blockType = expr.Labels[0]
					if len(expr.Labels) > 1 {
						labels = expr.Labels[1:]
					}
				}

				block.Type = blockType
				block.LabelNames = labels
				blocks = append(blocks, block)
			}
		}
	}

	body := &schema.Body{
		Attributes: attrs,
		Blocks:     blocks,
	}

	return body
}
