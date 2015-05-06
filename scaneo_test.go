package main

import "testing"

func TestFilenames(t *testing.T) {
	expFiles := 7
	dirs := []string{"testdata/"}

	files, err := filenames(dirs)
	if err != nil {
		t.Error(err)
	}

	if len(files) != expFiles {
		t.Error("actual files found differs from expected")
		t.Errorf("%d != %d", len(files), expFiles)
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
