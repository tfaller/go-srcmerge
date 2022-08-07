package conflict

import "fmt"

type Struct struct {
	LastName *string
}

func Hello(s Struct) {
	fmt.Print(s.LastName)
}

var (
	Foo  = "Bar"
	File = "1.go"
)
