package pkg

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func LoadAstFile(file string) (*ast.File, error) {
	fs := token.NewFileSet()
	return parser.ParseFile(fs, file, nil, 0)
}

func WriteAstFile(file string, ast *ast.File) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return format.Node(f, token.NewFileSet(), ast)
}
