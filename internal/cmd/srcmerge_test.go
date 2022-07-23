package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"testing"
)

const TestCaseBasePath = "../../test/merge"

func Test(t *testing.T) {
	tests, err := os.ReadDir(TestCaseBasePath)
	if err != nil {
		t.Fatalf("Failed to read test cases: %v", err)
	}

	for _, test := range tests {
		t.Run(test.Name(), testMerge)
	}
}

func testMerge(t *testing.T) {
	t.Parallel()

	testBasePath := path.Join(TestCaseBasePath, path.Base(t.Name()))
	testDirEntries, err := os.ReadDir(testBasePath)
	if err != nil {
		log.Fatal(err)
	}

	srcFiles := []string{}
	refactorNames := []string{}

	for _, entry := range testDirEntries {

		if entry.Name() == "out" {
			// out contains no test data ... it is the test result
			continue
		}

		entryPath := path.Join(testBasePath, entry.Name())

		if !entry.IsDir() {
			srcFiles = append(srcFiles, entryPath)
			refactorNames = append(refactorNames, fmt.Sprint(len(srcFiles)))
		} else {
			subEntries, err := os.ReadDir(entryPath)
			if err != nil {
				log.Fatal(err)
			}
			for _, subEntry := range subEntries {
				if subEntry.IsDir() {
					t.Fatalf("expected no sub dir %v", subEntry.Name())
				}
				srcFiles = append(srcFiles, path.Join(entryPath, subEntry.Name()))
				refactorNames = append(refactorNames, strings.ToTitle(entry.Name()))
			}
		}
	}

	err = Merge(srcFiles, refactorNames, path.Join(testBasePath, "out", "out.go"), "out")
	if err != nil {
		t.Error(err)
	}
}
