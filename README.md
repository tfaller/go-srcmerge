# go-srcmerge
This tool can merge multiple Go source files into one. 
In the process it removes all duplicates. 
If there is a name conflict, the names get refactored.

## Example

File a.go
```go
package a

var Hello = "World"

type Foo []string
```

File b.go
```go
package b

var Hello = "World"

type Foo []int
```
Execute the following command
```
srcmerge -f a.go -r A -f b.go -r B -p out -o out.go
```

Resulting out.go
```go
package out

var Hello = "World"

type Foo []string
type FooB []int
```