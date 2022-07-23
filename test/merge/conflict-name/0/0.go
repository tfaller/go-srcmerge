package conflict

import "fmt"

type Struct struct {
	FirstName, LastName string
}

func Hello(s Struct) {
	fmt.Printf("%v - %v", s.FirstName, s.LastName)
}
