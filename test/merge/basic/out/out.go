package out

import "fmt"

type Struct struct{ Name string }

func (s Struct) Hello() {
	fmt.Print("Hello %w", s.Name)
}

type Interface interface{ Hello() }
