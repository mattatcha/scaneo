package main

import (
	"path/filepath"
	"testing"
)

func TestFilenames(t *testing.T) {
	expectedFiles := 4
	inputPaths := []string{"testdata/", "testdata/access.go"}
	expectedFilenames := map[string]struct{}{
		"access.go":       struct{}{},
		"declarations.go": struct{}{},
		"methods.go":      struct{}{},
		"types.go":        struct{}{},
	}

	files, err := filenames(inputPaths)
	if err != nil {
		t.Error(err)
		t.SkipNow()
	}

	if expectedFiles != len(files) {
		t.Error("unexpected file count")
		t.Errorf("expected: %d; found: %d\n", expectedFiles, len(files))
		t.SkipNow()
	}

	for _, fp := range files {
		fn := filepath.Base(fp)
		if _, exists := expectedFilenames[fn]; !exists {
			t.Error("unexpected filename")
			t.Errorf("expected: %v\n", expectedFilenames)
			t.Errorf("found: %d\n", fn)
			t.SkipNow()
		}
	}
}

func TestWhitelist(t *testing.T) {
	whitelist := map[string]struct{}{
		"Exported":   struct{}{},
		"unexported": struct{}{},
	}
	expectedToks := len(whitelist)

	toks, err := parseCode("testdata/access.go", whitelist)
	if err != nil {
		t.Error(err)
		t.SkipNow()
	}

	if expectedToks != len(toks) {
		t.Error("unexpected struct tokens length")
		t.Errorf("expected: %d; found: %d\n", expectedToks, len(toks))
	}
}

func TestParseCode(t *testing.T) {
}

func TestGenFile(t *testing.T) {
}
