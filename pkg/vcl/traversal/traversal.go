package traversal

import (
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/ast"
	"github.com/KeisukeYamashita/go-vcl/pkg/vcl/schema"
)

// Contents retrives from ast.Program
func Content(prog *ast.Program) *schema.BodyContent {
	b := convertBody(prog.Statements)
	return b
}

// Contents will ast.Program to schema
func convertBody(stmts []ast.Statement) *schema.BodyContent {
	attrs := make(map[string]*schema.Attribute)
	var blocks schema.Blocks

	for _, stmt := range stmts {
		switch v := stmt.(type) {
		case *ast.AssignStatement:
			var value interface{}

			switch lit := v.Value.(type) {
			case *ast.StringLiteral:
				value = lit.Value
			case *ast.CIDRLiteral:
				value = lit.Value
			case *ast.BooleanLiteral:
				value = lit.Value
			case *ast.IntegerLiteral:
				value = lit.Value
			default:
				panic("cannot pass invalid argument")
			}

			attrs[v.Name.Value] = &schema.Attribute{
				Name:  v.Name.Value,
				Value: value,
			}
		case *ast.ExpressionStatement:
			switch expr := v.Expression.(type) {
			case *ast.BlockExpression:
				body := convertBody(expr.Blocks.Statements)
				block := &schema.Block{
					Body: body,
				}

				if len(expr.Labels) > 0 {
					block.Type = expr.Labels[0]
					if len(expr.Labels) > 1 {
						labels := expr.Labels[1:]
						block.Labels = labels
					}
				}

				blocks = append(blocks, block)
			}
		}
	}

	body := &schema.BodyContent{
		Attributes: attrs,
		Blocks:     blocks,
	}

	return body
}
