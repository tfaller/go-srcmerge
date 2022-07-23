package pkg

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"path"
	"strconv"
)

type Merger struct {
	File ast.File

	declares    map[string]ast.Node
	imports     map[string]string
	importsDecl ast.GenDecl
}

func NewMerger(pkgName string) *Merger {
	return &Merger{
		File: ast.File{
			Name:    &ast.Ident{Name: pkgName},
			Imports: []*ast.ImportSpec{},
			Decls:   []ast.Decl{},
		},
		declares: map[string]ast.Node{},
		imports:  map[string]string{},
		importsDecl: ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: []ast.Spec{},
		},
	}
}

func (m *Merger) Merge(b *ast.File, duplicatePostfix string) error {
	bDeclares := findDeclarations(b)

	// handle imports
	imps, err := findImports(b)
	if err != nil {
		return nil
	}
	for name, iPath := range imps {

		if len(m.imports) == 0 {
			// add imports to the file ...
			m.File.Decls = append(m.File.Decls, &m.importsDecl)
		}

		mImportPath := m.imports[name]
		if iPath == mImportPath {
			continue
		}

		if mImportPath != "" {
			newName := name + duplicatePostfix
			log.Printf("import name conflict %q with paths %q != %q", name, iPath, mImportPath)
			log.Printf("rename %q -> %q", name, newName)
			RenameDeclarations(b, name, newName)
			name = newName
		}

		var impSpecName *ast.Ident
		if path.Base(iPath) != name {
			impSpecName = &ast.Ident{Name: name}
		}

		impSpec := &ast.ImportSpec{
			Name: impSpecName,
			Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(iPath)},
		}

		m.imports[name] = iPath
		m.importsDecl.Specs = append(m.importsDecl.Specs, impSpec)
		m.File.Imports = append(m.File.Imports, impSpec)
	}
	RemoveImports(b)

	// find and remove duplicate declarations
	for name, dec := range bDeclares {
		dup := m.declares[name]
		if dup == nil {
			m.declares[name] = dec
			continue
		}
		if err := NodeEqual(dec, dup); err != nil {
			// ups ... name conflict
			newName := name + duplicatePostfix
			log.Printf("name conflict of declaration %q: %v", name, err)
			log.Printf("rename %q -> %q", name, newName)
			RenameDeclarations(b, name, newName)
			m.declares[newName] = dec
		} else {
			// remove instance of duplicate declaration
			b.Decls = RemoveGenDeclByName(b.Decls, name)
			log.Printf("removed duplicate %q", name)
		}
	}
	m.File.Decls = append(m.File.Decls, b.Decls...)

	return nil
}

func findDeclarations(file *ast.File) map[string]ast.Node {
	declares := map[string]ast.Node{}

	ast.Inspect(file, func(n ast.Node) bool {
		if _, isStmtBlock := n.(*ast.BlockStmt); isStmtBlock {
			// don't check for things which are declared in a func
			return false
		}
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			declares[typeSpec.Name.Name] = typeSpec.Type
			return false
		}
		if valSpec, ok := n.(*ast.ValueSpec); ok {
			for i, name := range valSpec.Names {
				declares[name.Name] = &ast.ValueSpec{
					Names:  []*ast.Ident{name},
					Type:   valSpec.Type,
					Values: []ast.Expr{valSpec.Values[i]},
				}
			}
			return false
		}
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if funcDecl.Recv == nil {
				declares[funcDecl.Name.Name] = funcDecl
			}
			return false
		}
		return true
	})

	return declares
}

func findImports(file *ast.File) (map[string]string, error) {
	imports := map[string]string{}
	for _, impSpec := range file.Imports {
		impName := ""
		impPath, err := strconv.Unquote(impSpec.Path.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid import path %v", err)
		}

		if impSpec.Name != nil {
			impName = impSpec.Name.Name
		} else {
			impName = path.Base(impPath)
		}

		imports[impName] = impPath
	}
	return imports, nil
}
