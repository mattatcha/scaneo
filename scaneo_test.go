package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"testing"
)

func TestFilenames(t *testing.T) {
	expFiles := 4
	paths := []string{"testdata/", "testdata/access.go"}

	files, err := filenames(paths)
	if err != nil {
		t.Error(err)
	}

	if len(files) != expFiles {
		t.Error("actual files found differs from expected")
		t.Errorf("expected: %d; found: %d\n", expFiles, len(files))
	}
}

func TestParseCode(t *testing.T) {
	files := []string{
		"testdata/access.go",
		"testdata/declarations.go",
		"testdata/methods.go",
		"testdata/types.go",
	}

	structCnt := []int{
		4,
		5,
		1,
		6,
	}

	fieldsCnt := [][]int{
		[]int{
			2,
			2,
			2,
			2,
		},
		[]int{
			2,
			2,
			2,
			6,
			2,
		},
		[]int{
			8,
		},
		[]int{
			1,
			17,
			1,
			3,
			3,
			5,
		},
	}

	noFilter := make(map[string]struct{})
	for i, f := range files {
		toks, err := parseCode(f, noFilter)
		if err != nil {
			t.Error(err)
			continue
		}

		if len(toks) != structCnt[i] {
			t.Error("file:", f)
			t.Error("actual token count differs from expected count")
			t.Errorf("%d != %d", len(toks), structCnt[i])
			t.SkipNow()
		}

		for j, tok := range toks {
			if len(tok.Fields) != len(tok.Types) {
				t.Error("file:", f)
				t.Error("field names and field types don't align")
				t.Errorf("%d != %d", len(tok.Fields), len(tok.Types))
				t.SkipNow()
			}

			if len(tok.Fields) != fieldsCnt[i][j] {
				t.Error("file:", f)
				t.Errorf("struct %d", i)
				t.Errorf("field %d", j)
				t.Error("actual struct field count and expect count differ")
				t.Errorf("%d != %d", len(tok.Fields), fieldsCnt[i][j])
				t.Errorf("%+v", tok)
				t.SkipNow()
			}
		}
	}
}

func TestWhitelist(t *testing.T) {
	whitelist := map[string]struct{}{
		"Exported":   struct{}{},
		"unexported": struct{}{},
	}

	toks, err := parseCode("testdata/access.go", whitelist)
	if err != nil {
		t.Error(err)
		t.SkipNow()
	}

	if len(toks) != 2 {
		t.Error("unexpected struct tokens")
		t.Errorf("expected: %d; found: %d\n", 2, len(toks))
	}
}

func TestGenFile(t *testing.T) {
	toks := []structToken{
		structToken{
			Name:   "lorem",
			Fields: []string{"a", "b", "C"},
			Types:  []string{"int", "int", "int"},
		},
		structToken{
			Name:   "Sit",
			Fields: []string{"A", "b", "c"},
			Types:  []string{"string", "string", "string"},
		},
	}
	expectedScanNames := []string{
		"scanLorem",
		"scanSit",
		"scanLorems",
		"scanSits",
	}

	fout, err := ioutil.TempFile(os.TempDir(), "scaneo-test-")
	if err != nil {
		t.Error(err)
		t.SkipNow()
	}
	defer os.Remove(fout.Name()) // comment this line to examin generated code
	defer fout.Close()

	if err := genFile(fout, "testing", true, toks); err != nil {
		t.Error(err)
		t.SkipNow()
	}

	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, fout.Name(), nil, 0)
	if err != nil {
		t.Error(err)
		t.SkipNow()
	}

	scanFuncs := make([]string, 0, len(toks))
	for _, dec := range astf.Decls {
		funcDecl, isFuncDecl := dec.(*ast.FuncDecl)
		if !isFuncDecl {
			continue
		}

		scanFuncs = append(scanFuncs, funcDecl.Name.String())
	}

	if len(toks)*2 != len(scanFuncs) {
		t.Error("unexpected number of scan functions found")
		t.Errorf("expected: %d; found: %d\n", len(toks)*2, len(scanFuncs))
		t.SkipNow()
	}

	for i := range expectedScanNames {
		if expectedScanNames[i] != scanFuncs[i] {
			t.Error("unexpected scan function found")
			t.Errorf("expected: %s; found: %s\n", expectedScanNames[i], scanFuncs[i])
		}
	}
}
