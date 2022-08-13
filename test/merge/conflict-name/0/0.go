package conflict

import "fmt"

type Struct struct {
	FirstName, LastName string
	Struct              *Struct
}

func Hello(s Struct) {
	s = *s.Struct
	fmt.Printf("%v - %v", s.FirstName, s.LastName)
}

var (
	Foo  = "Bar"
	File = "0.go"
)
