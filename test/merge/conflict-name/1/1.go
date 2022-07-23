package conflict

import "fmt"

type Struct struct {
	Name string
}

func Hello(s Struct) {
	fmt.Print(s.Name)
}
