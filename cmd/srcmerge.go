package main

import (
	"flag"
	"log"

	"github.com/tfaller/go-srcmerge/internal/cmd"
	"github.com/tfaller/go-srcmerge/pkg/sliceflag"
)

func main() {
	srcFilesNames := sliceflag.StringSliceFlag{}
	flag.Var(&srcFilesNames, "f", "go source file (can be set multiple time)")

	srcRefactorName := sliceflag.StringSliceFlag{}
	flag.Var(&srcRefactorName, "r", "refactor name for a given source file (can be set multiple time)")

	packageName := flag.String("p", "merged", "package name")
	outFile := flag.String("o", "", "out file")
	flag.Parse()

	err := cmd.Merge(srcFilesNames, srcRefactorName, *outFile, *packageName)
	if err != nil {
		log.Fatal(err)
	}
}
