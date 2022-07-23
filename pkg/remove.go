package pkg

import (
	"go/ast"
	"go/token"
)

// RemoveGenDeclByName removes a const, var or type by its declaration name.
// References which use the declared thing are unchanged.
func RemoveGenDeclByName(declarations []ast.Decl, name string) []ast.Decl {
	newDecls := make([]ast.Decl, 0, len(declarations))
	for _, d := range declarations {
		genDecl, ok := d.(*ast.GenDecl)
		if !ok {
			newDecls = append(newDecls, d)
			continue
		}
		newSpecs := genDecl.Specs[:0]
		for i, spec := range genDecl.Specs {
			// prevent memory leaks if we deleted an item earlier
			genDecl.Specs[i] = nil

			switch spec := spec.(type) {
			case *ast.TypeSpec:
				if spec.Name.Name != name {
					newSpecs = append(newSpecs, spec)
				}
			case *ast.ValueSpec:
				names := spec.Names[:0]
				values := spec.Values[:0]
				for i, ident := range spec.Names {
					if ident.Name != name {
						names = append(names, ident)
						values = append(values, spec.Values[i])
					}
				}
				if len(names) > 0 {
					spec.Names = names
					spec.Values = values
					newSpecs = append(newSpecs, spec)
				}
			default:
				newSpecs = append(newSpecs, spec)
			}
		}
		if len(newSpecs) > 0 {
			// if there are still spec entries, add the remaining
			genDecl.Specs = newSpecs
			newDecls = append(newDecls, d)
		}
	}
	return newDecls
}

func RemoveImports(file *ast.File) {
	newDecls := file.Decls[:0]
	for i, decl := range file.Decls {
		// prevent memory leaks if we deleted an item earlier
		file.Decls[i] = nil

		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			newDecls = append(newDecls, decl)
		}
	}
	file.Decls = newDecls
}
