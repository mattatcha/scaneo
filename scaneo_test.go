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
	expFiles := 7
	paths := []string{"testdata/", "testdata/adipiscing.go"}

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
		"testdata/adipiscing.go",
		"testdata/amet.go",
		"testdata/consectetur.go",
		"testdata/dolor.go",
		"testdata/ipsum.go",
		"testdata/lorem.go",
		"testdata/sit.go",
	}

	structCnt := []int{
		1,
		3,
		1,
		1,
		1,
		1,
		1,
	}

	fieldsCnt := [][]int{
		[]int{
			5,
		},
		[]int{
			0,
			1,
			1,
		},
		[]int{
			1,
		},
		[]int{
			17,
		},
		[]int{
			4,
		},
		[]int{
			3,
		},
		[]int{
			3,
		},
	}

	for i, f := range files {
		toks, err := parseCode(f)
		if err != nil {
			t.Error(err)
			continue
		}

		if len(toks) != structCnt[i] {
			t.Error("actual token count differs from expected count")
			t.Errorf("%d != %d", len(toks), structCnt[i])
		}

		for j, tok := range toks {
			if len(tok.Fields) != len(tok.Types) {
				t.Error("field names and field types don't align")
				t.Errorf("%d != %d", len(tok.Fields), len(tok.Types))
			}

			if len(tok.Fields) != fieldsCnt[i][j] {
				t.Errorf("struct %d", j+1)
				t.Error("actual struct field count and expect count differ")
				t.Errorf("%d != %d", len(tok.Fields), fieldsCnt[i][j])
			}
		}
	}
}

func TestWriteCode(t *testing.T) {
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

	fout, err := ioutil.TempFile(os.TempDir(), "scaneo-test-")
	if err != nil {
		t.Error(err)
		t.SkipNow()
	}
	defer os.Remove(fout.Name())
	defer fout.Close()

	if err := writeCode(fout, "testing", true, toks); err != nil {
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
	}
}
