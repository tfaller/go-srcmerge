package pkg

import (
	"fmt"
	"go/ast"
)

type Field struct {
	Name string
	Tag  *ast.BasicLit
	Type ast.Expr
}

func FieldListToMap(fl *ast.FieldList) (map[string]Field, error) {
	fields := map[string]Field{}
	for _, f := range fl.List {
		for _, n := range f.Names {
			if _, exists := fields[n.Name]; exists {
				return nil, fmt.Errorf("duplicate field %v in one field list", n.Name)
			}
			fields[n.Name] = Field{Tag: f.Tag, Type: f.Type}
		}
	}
	return fields, nil
}
