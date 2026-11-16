package classfileparser

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCorpus(t *testing.T) {
	root := os.Getenv("CLASSFILEPARSER_CORPUS")
	if root == "" {
		t.Skip("set CLASSFILEPARSER_CORPUS to a directory of .class files to run this test")
	}

	count := 0
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".class") {
			return err
		}
		count++

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		cf, err := Open(bytes.NewReader(data))
		if err == nil {
			_, err = cf.GetClassFile()
		}
		if err != nil {
			t.Errorf("%s: %v", path, err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk %s: %v", root, err)
	}
	if count == 0 {
		t.Fatalf("no .class file found in %s", root)
	}
	t.Logf("%d classes parsed", count)
}
