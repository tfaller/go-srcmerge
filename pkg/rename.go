package pkg

import (
	"go/ast"
)

// RenameDeclarations renames all occurrences of a declaration and their usage.
func RenameDeclarations(node ast.Node, oldName, newName string) {
	ast.Inspect(node, func(n ast.Node) bool {
		if n, ok := n.(*ast.SelectorExpr); ok {
			// only rename SelectorExpr expression -> not the selector (field) itself
			RenameDeclarations(n.X, oldName, newName)
			return false
		}
		if field, isField := n.(*ast.Field); isField {
			// only rename field types ... not the field name
			RenameDeclarations(field.Type, oldName, newName)
			return false
		}
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if ident.Name == oldName {
			ident.Name = newName
		}
		return true
	})
}
