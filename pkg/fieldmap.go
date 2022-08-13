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

func (f Field) ToAstField() *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(f.Name)},
		Type:  f.Type,
		Tag:   f.Tag,
	}
}

func FieldListToMap(fl *ast.FieldList) (map[string]Field, error) {
	if fl == nil {
		return nil, nil
	}
	fields := map[string]Field{}
	for _, f := range fl.List {
		for _, n := range f.Names {
			if _, exists := fields[n.Name]; exists {
				return nil, fmt.Errorf("duplicate field %v in one field list", n.Name)
			}
			fields[n.Name] = Field{Name: n.Name, Tag: f.Tag, Type: f.Type}
		}
	}
	return fields, nil
}
