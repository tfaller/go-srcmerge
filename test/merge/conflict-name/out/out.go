package out

import "fmt"

type Struct struct{ FirstName, LastName string }

func Hello(s Struct) {
	fmt.Printf("%v - %v", s.FirstName, s.LastName)
}

var (
	Foo  = "Bar"
	File = "0.go"
)

type Struct1 struct{ LastName *string }

func Hello1(s Struct1) {
	fmt.Print(s.LastName)
}

var (
	File1 = "1.go"
)
