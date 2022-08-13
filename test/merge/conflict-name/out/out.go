package out

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

type Struct1 struct {
	LastName *string
	Struct   *Struct1
}

func Hello1(s Struct1) {
	s = *s.Struct
	fmt.Print(s.LastName)
}

var (
	File1 = "1.go"
)
