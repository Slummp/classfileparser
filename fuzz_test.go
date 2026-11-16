package classfileparser

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func FuzzOpen(f *testing.F) {
	f.Add(minimalClass().bytes())
	f.Add(fullClass().bytes())
	_ = filepath.WalkDir(fixturesRoot, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(path, ".class") {
			if data, err := os.ReadFile(path); err == nil {
				f.Add(data)
			}
		}
		return nil
	})

	f.Fuzz(func(t *testing.T, data []byte) {
		cf, err := Open(bytes.NewReader(data))
		if err == nil {
			_, err = cf.GetConstantPool()
		}
		if err == nil {
			_, err = cf.GetClassFile()
		}
		if err != nil && !isTruncation(err) && !errors.Is(err, ErrMalformed) {
			t.Fatalf("unclassified error: %v", err)
		}
	})
}
