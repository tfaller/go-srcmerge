package duplicates

import "io"

const Hello = "World"

var Foo = "Bar"

type Struct struct{}

type ImportFieldType struct {
	Reader io.Reader
}

func (i *ImportFieldType) GetReader() io.Reader {
	return i.Reader
}
