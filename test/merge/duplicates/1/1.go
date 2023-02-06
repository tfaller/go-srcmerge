package duplicates

import "io"

const Hello = "World"

var Foo = "Bar"

type Struct struct{}

type ImportFieldType struct {
	Reader io.Reader
}

func NewImportFieldType(reader io.Reader) ImportFieldType {
	return ImportFieldType{
		Reader: reader,
	}
}

func (i *ImportFieldType) GetReader() io.Reader {
	return i.Reader
}
