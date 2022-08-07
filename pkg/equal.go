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
		return fmt.Errorf("%v different type: %v", a.Names[0].Name, err)
	}
	for i, name := range a.Names {
		if err := IdentEqual(name, b.Names[i]); err != nil {
			return fmt.Errorf("mismatch name: %v", err)
		}
		if err := NodeEqual(a.Values[i], b.Values[i]); err != nil {
			return fmt.Errorf("%v different value: %v", name, err)
		}
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
