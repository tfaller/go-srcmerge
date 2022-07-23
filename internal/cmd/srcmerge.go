package cmd

import (
	"fmt"
	"log"

	"github.com/tfaller/go-srcmerge/pkg"
)

func Merge(srcFilesNames []string, srcRefactorName []string, outFile, packageName string) error {

	if len(srcFilesNames) == 0 {
		return fmt.Errorf("no source file specified")
	}

	if len(srcFilesNames) != len(srcRefactorName) {
		return fmt.Errorf("for each source file must be refactor name set")
	}

	merger := pkg.NewMerger(packageName)

	for i, srcFile := range srcFilesNames {
		ast, err := pkg.LoadAstFile(srcFile)
		if err != nil {
			log.Fatal(err)
		}
		err = merger.Merge(ast, srcRefactorName[i])
		if err != nil {
			return err
		}
	}

	return pkg.WriteAstFile(outFile, &merger.File)
}
