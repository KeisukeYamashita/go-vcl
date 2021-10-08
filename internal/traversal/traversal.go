package traversal

import (
	"github.com/KeisukeYamashita/go-vcl/internal/ast"
	"github.com/KeisukeYamashita/go-vcl/internal/schema"
)

// Content retrives from ast.Program
func Content(prog *ast.Program) *schema.BodyContent {
	b := convertBody(prog.Statements)
	return b
}

// BodyContent retrives body content from body
func BodyContent(body schema.Body) *schema.BodyContent {
	return body.(*schema.BodyContent)
}

// Contents will ast.Program to schema
func convertBody(stmts []ast.Statement) *schema.BodyContent {
	attrs := make(map[string]*schema.Attribute)
	var blocks schema.Blocks
	flats := []interface{}{}
	comments := []string{}

	for _, stmt := range stmts {
		switch v := stmt.(type) {
		case *ast.AssignStatement:
			var isBlock bool
			var value interface{}
			switch lit := v.Value.(type) {
			case *ast.StringLiteral:
				value = lit.Value
			case *ast.CIDRLiteral:
				value = lit.Value
			case *ast.PercentageLiteral:
				value = lit.Value
			case *ast.BooleanLiteral:
				value = lit.Value
			case *ast.IntegerLiteral:
				value = lit.Value
			case *ast.BlockExpression:
				isBlock = true
				body := convertBody(lit.Blocks.Statements)
				block := &schema.Block{
					Body: body,
				}
				block.Type = v.TokenLiteral()
				blocks = append(blocks, block)
			case *ast.Identifier:
				value = lit.Value
			default:
				panic("cannot pass invalid argument which is no a literal")
			}

			if isBlock == false {
				attrs[v.Name.Value] = &schema.Attribute{
					Name:  v.Name.Value,
					Value: value,
				}
			}
		case *ast.AssignFieldStatement:
			var value interface{}
			switch lit := v.Value.(type) {
			case *ast.StringLiteral:
				value = lit.Value
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

				block.Type = expr.TokenLiteral()
				if len(expr.Labels) > 0 {
					block.Labels = expr.Labels
				}

				if block.Type == "{" {
					// this is flatten block
					flats = append(flats, block)
				} else {
					blocks = append(blocks, block)
				}
			case *ast.StringLiteral:
				flats = append(flats, expr.Value)
			case *ast.CIDRLiteral:
				flats = append(flats, expr.Value)
			case *ast.BooleanLiteral:
				flats = append(flats, expr.Value)
			case *ast.IntegerLiteral:
				flats = append(flats, expr.Value)
			}
		case *ast.CommentStatement:
			comments = append(comments, v.TokenLiteral())
		}
	}

	body := &schema.BodyContent{
		Attributes: attrs,
		Blocks:     blocks,
		Flats:      flats,
		Comments:   comments,
	}

	return body
}
