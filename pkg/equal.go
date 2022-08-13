package pkg

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

// ErrAdditionalFields if a field list has additional fields
type ErrAdditionalFields struct {
	A []string
	B []string
}

func (e ErrAdditionalFields) Error() string {
	return fmt.Sprintf(
		"additional fields of a: (%v) and of b: (%v)",
		strings.Join(e.A, ","), strings.Join(e.B, ","))
}

// NodeEqual checks whether two nodes represent the same
// thing. If the are not the same, the first found mismatch will be
// reported. Comments and positions are ignored.
func NodeEqual(a, b ast.Node) error {
	if a == nil && b == nil {
		return nil
	}
	aT := reflect.TypeOf(a)
	if bT := reflect.TypeOf(b); aT != bT {
		return fmt.Errorf("type mismatch %q != %q", aT, bT)
	}
	switch a := a.(type) {
	case *ast.BasicLit:
		return BasicLitEqual(a, b.(*ast.BasicLit))
	case *ast.Ident:
		return IdentEqual(a, b.(*ast.Ident))
	case *ast.StarExpr:
		return StarExprEqual(a, b.(*ast.StarExpr))
	case *ast.StructType:
		return StructEqual(a, b.(*ast.StructType))
	case *ast.InterfaceType:
		return InterfaceEqual(a, b.(*ast.InterfaceType))
	case *ast.ValueSpec:
		return ValueEqual(a, b.(*ast.ValueSpec))
	case *ast.ArrayType:
		return ArrayTypeEqual(a, b.(*ast.ArrayType))
	case *ast.SelectorExpr:
		return SelectorExprEqual(a, b.(*ast.SelectorExpr))
	case *ast.FuncDecl:
		return FuncDeclEqual(a, b.(*ast.FuncDecl))
	case *ast.FuncType:
		return FuncTypeEqual(a, b.(*ast.FuncType))
	case *ast.ReturnStmt:
		return ReturnStmtEqual(a, b.(*ast.ReturnStmt))
	case *ast.AssignStmt:
		return AssignStmtEqual(a, b.(*ast.AssignStmt))
	case *ast.CompositeLit:
		return CompositeLitEqual(a, b.(*ast.CompositeLit))
	case *ast.UnaryExpr:
		return UnaryExprEqual(a, b.(*ast.UnaryExpr))
	case *ast.IfStmt:
		return IfStmtEqual(a, b.(*ast.IfStmt))
	case *ast.IndexExpr:
		return IndexExprEqual(a, b.(*ast.IndexExpr))
	case *ast.DeclStmt:
		return DeclStmtEqual(a, b.(*ast.DeclStmt))
	case *ast.GenDecl:
		return GenDeclEqual(a, b.(*ast.GenDecl))
	case *ast.TypeSpec:
		return TypeSpecEqual(a, b.(*ast.TypeSpec))
	case *ast.CallExpr:
		return CallExprEqual(a, b.(*ast.CallExpr))
	case *ast.BinaryExpr:
		return BinaryExprEqual(a, b.(*ast.BinaryExpr))
	case *ast.RangeStmt:
		return RangeStmtEqual(a, b.(*ast.RangeStmt))
	}
	return fmt.Errorf("unknown node type %q", aT)
}

func IdentEqual(a, b *ast.Ident) error {
	if a.Name == b.Name {
		return nil
	}
	return fmt.Errorf("%q != %q", a, b)
}

func BasicLitEqual(a, b *ast.BasicLit) error {
	if a == b {
		return nil
	}
	if a == nil || b == nil {
		return fmt.Errorf("only one literal is nil")
	}
	if a.Kind != b.Kind {
		return fmt.Errorf("literal kind not the same")
	}
	if a.Value != b.Value {
		return fmt.Errorf("%q != %q", a.Value, b.Value)
	}
	return nil
}

func StarExprEqual(a, b *ast.StarExpr) error {
	return NodeEqual(a.X, b.X)
}

func FieldListEqual(a, b *ast.FieldList) error {
	aFields, err := FieldListToMap(a)
	if err != nil {
		return err
	}
	bFields, err := FieldListToMap(b)
	if err != nil {
		return err
	}

	aAdditionalFields := []string{}
	bAdditionalFields := []string{}

	for name, aField := range aFields {
		bField, exists := bFields[name]
		if !exists {
			aAdditionalFields = append(aAdditionalFields, name)
			continue
		}
		if err := BasicLitEqual(aField.Tag, bField.Tag); err != nil {
			return fmt.Errorf("tag of field %v is not the same", name)
		}
		if err := NodeEqual(aField.Type, bField.Type); err != nil {
			return fmt.Errorf("%v: %v", name, err)
		}
		delete(bFields, name)
	}

	for name := range bFields {
		bAdditionalFields = append(bAdditionalFields, name)
	}

	if len(aAdditionalFields) > 0 || len(bAdditionalFields) > 0 {
		return ErrAdditionalFields{A: aAdditionalFields, B: bAdditionalFields}
	}

	return nil
}

func ValueEqual(a, b *ast.ValueSpec) error {
	if err := NodeEqual(a.Type, b.Type); err != nil {
		return fmt.Errorf("%v different type: %w", a.Names[0].Name, err)
	}
	if len(a.Names) != len(b.Names) {
		return fmt.Errorf("different name count")
	}
	for i, name := range a.Names {
		if err := IdentEqual(name, b.Names[i]); err != nil {
			return fmt.Errorf("mismatch name: %w", err)
		}
	}
	if err := ExprSliceEqual(a.Values, b.Values); err != nil {
		return fmt.Errorf("values not equal: %w", err)
	}
	return nil
}

func InterfaceEqual(a, b *ast.InterfaceType) error {
	if a.Incomplete || b.Incomplete {
		return fmt.Errorf("incomplete interface")
	}
	return FieldListEqual(a.Methods, b.Methods)
}

func ArrayTypeEqual(a, b *ast.ArrayType) error {
	if err := NodeEqual(a.Len, b.Len); err != nil {
		return fmt.Errorf("len expression not the same: %w", err)
	}
	if err := NodeEqual(a.Elt, b.Elt); err != nil {
		return fmt.Errorf("element type not the same: %v", err)
	}
	return nil
}

func StructEqual(a, b *ast.StructType) error {
	if a.Incomplete || b.Incomplete {
		return fmt.Errorf("incomplete struct")
	}
	return FieldListEqual(a.Fields, b.Fields)
}

func SelectorExprEqual(a, b *ast.SelectorExpr) error {
	if err := IdentEqual(a.Sel, b.Sel); err != nil {
		return fmt.Errorf("selector name not equal: %w", err)
	}
	if err := NodeEqual(a.X, b.X); err != nil {
		return fmt.Errorf("selector source not equal: %w", err)
	}
	return nil
}

func FuncDeclEqual(a, b *ast.FuncDecl) error {
	if err := IdentEqual(a.Name, b.Name); err != nil {
		return fmt.Errorf("name not equal: %w", err)
	}
	if err := FieldListEqual(a.Recv, b.Recv); err != nil {
		return fmt.Errorf("receiver not equal: %w", err)
	}
	if err := FuncTypeEqual(a.Type, b.Type); err != nil {
		return fmt.Errorf("type not equal: %w", err)
	}
	if err := BlockStmtEqual(a.Body, b.Body); err != nil {
		return fmt.Errorf("body not equal: %w", err)
	}
	return nil
}

func FuncTypeEqual(a, b *ast.FuncType) error {
	if err := FieldListEqual(a.Params, b.Params); err != nil {
		return fmt.Errorf("params are not equal: %w", err)
	}
	if err := FieldListEqual(a.Results, b.Results); err != nil {
		return fmt.Errorf("results are not equal: %w", err)
	}
	if err := FieldListEqual(a.TypeParams, b.TypeParams); err != nil {
		return fmt.Errorf("type params are not equal: %w", err)
	}
	return nil
}

func BlockStmtEqual(a, b *ast.BlockStmt) error {
	if len(a.List) != len(b.List) {
		return fmt.Errorf("different stmt count")
	}
	for i, aS := range a.List {
		bS := b.List[i]
		if err := NodeEqual(aS, bS); err != nil {
			return fmt.Errorf("stmt not equal: %w", err)
		}
	}
	return nil
}

func ReturnStmtEqual(a, b *ast.ReturnStmt) error {
	if err := ExprSliceEqual(a.Results, b.Results); err != nil {
		return fmt.Errorf("results not equal: %w", err)
	}
	return nil
}

func AssignStmtEqual(a, b *ast.AssignStmt) error {
	if err := ExprSliceEqual(a.Lhs, b.Lhs); err != nil {
		return fmt.Errorf("left hand not equal: %w", err)
	}
	if err := ExprSliceEqual(a.Rhs, b.Rhs); err != nil {
		return fmt.Errorf("right hand not equal: %w", err)
	}
	return nil
}

func ExprSliceEqual(a, b []ast.Expr) error {
	if len(a) != len(b) {
		return fmt.Errorf("different expression count")
	}
	for i, aE := range a {
		bE := b[i]
		if err := NodeEqual(aE, bE); err != nil {
			return fmt.Errorf("expr %v not equal: %w", i, err)
		}
	}
	return nil
}

func CompositeLitEqual(a, b *ast.CompositeLit) error {
	if a.Incomplete || b.Incomplete {
		return fmt.Errorf("compositeLit is incomplete")
	}
	if err := NodeEqual(a.Type, b.Type); err != nil {
		return fmt.Errorf("type not equal: %w", err)
	}
	if err := ExprSliceEqual(a.Elts, b.Elts); err != nil {
		return fmt.Errorf("values not equal: %w", err)
	}
	return nil
}

func UnaryExprEqual(a, b *ast.UnaryExpr) error {
	if a.Op != b.Op {
		return fmt.Errorf("different operator")
	}
	if err := NodeEqual(a.X, b.X); err != nil {
		return fmt.Errorf("different operand: %w", err)
	}
	return nil
}

func IfStmtEqual(a, b *ast.IfStmt) error {
	if err := NodeEqual(a.Init, b.Init); err != nil {
		return fmt.Errorf("init not equal: %w", err)
	}
	if err := NodeEqual(a.Cond, b.Cond); err != nil {
		return fmt.Errorf("cond not equal: %w", err)
	}
	if err := BlockStmtEqual(a.Body, b.Body); err != nil {
		return fmt.Errorf("body not equal: %w", err)
	}
	if err := NodeEqual(a.Else, b.Else); err != nil {
		return fmt.Errorf("else not equal: %w", err)
	}
	return nil
}

func IndexExprEqual(a, b *ast.IndexExpr) error {
	if err := NodeEqual(a.Index, b.Index); err != nil {
		return fmt.Errorf("index not equal: %w", err)
	}
	if err := NodeEqual(a.X, b.X); err != nil {
		return fmt.Errorf("operand not equal: %w", err)
	}
	return nil
}

func DeclStmtEqual(a, b *ast.DeclStmt) error {
	if err := NodeEqual(a.Decl, b.Decl); err != nil {
		return fmt.Errorf("decl not equal: %w", err)
	}
	return nil
}

func GenDeclEqual(a, b *ast.GenDecl) error {
	if a.Tok != b.Tok {
		return fmt.Errorf("tok not equal")
	}
	if len(a.Specs) != len(b.Specs) {
		return fmt.Errorf("different spec count")
	}
	for i, aS := range a.Specs {
		bS := b.Specs[i]
		if err := NodeEqual(aS, bS); err != nil {
			return fmt.Errorf("spec not equal: %w", err)
		}
	}
	return nil
}

func TypeSpecEqual(a, b *ast.TypeSpec) error {
	if err := IdentEqual(a.Name, b.Name); err != nil {
		return fmt.Errorf("name not equal: %w", err)
	}
	if err := NodeEqual(a.Type, b.Type); err != nil {
		return fmt.Errorf("type not equal: %w", err)
	}
	if err := FieldListEqual(a.TypeParams, b.TypeParams); err != nil {
		return fmt.Errorf("type params not equal: %w", err)
	}
	return nil
}

func CallExprEqual(a, b *ast.CallExpr) error {
	if err := NodeEqual(a.Fun, b.Fun); err != nil {
		return fmt.Errorf("fun not equal: %w", err)
	}
	if err := ExprSliceEqual(a.Args, b.Args); err != nil {
		return fmt.Errorf("args not equal: %w", err)
	}
	return nil
}

func BinaryExprEqual(a, b *ast.BinaryExpr) error {
	if a.Op != b.Op {
		return fmt.Errorf("different operand")
	}
	if err := NodeEqual(a.X, b.X); err != nil {
		return fmt.Errorf("left operand not equal: %w", err)
	}
	if err := NodeEqual(a.Y, b.Y); err != nil {
		return fmt.Errorf("right operand not equal: %w", err)
	}
	return nil
}

func RangeStmtEqual(a, b *ast.RangeStmt) error {
	if a.Tok != b.Tok {
		return fmt.Errorf("different tok")
	}
	if err := NodeEqual(a.Key, a.Key); err != nil {
		return fmt.Errorf("key not equal: %w", err)
	}
	if err := NodeEqual(a.Value, b.Value); err != nil {
		return fmt.Errorf("value not equal: %w", err)
	}
	if err := NodeEqual(a.X, b.X); err != nil {
		return fmt.Errorf("range value not equal: %w", err)
	}
	if err := BlockStmtEqual(a.Body, b.Body); err != nil {
		return fmt.Errorf("body not equal: %w", err)
	}
	return nil
}
