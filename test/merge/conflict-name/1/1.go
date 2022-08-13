package conflict

import "fmt"

type Struct struct {
	LastName *string
	Struct   *Struct
}

func Hello(s Struct) {
	type NewType Struct

	var data struct {
		NewType
	}

	*s.Struct = Struct(*data.NewType.Struct)
	fmt.Print(s.LastName)
}

var (
	Foo  = "Bar"
	File = "1.go"
)
