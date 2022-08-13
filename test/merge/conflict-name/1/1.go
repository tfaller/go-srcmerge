package conflict

import "fmt"

type Struct struct {
	LastName *string
	Struct   *Struct
}

func Hello(s Struct) {
	s = *s.Struct
	fmt.Print(s.LastName)
}

var (
	Foo  = "Bar"
	File = "1.go"
)
