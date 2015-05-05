package main

import "testing"

func TestParseCode(t *testing.T) {
	files := []string{
		"testdata/structs.go",
	}

	structCnt := []int{
		6,
	}

	fieldsCnt := [][]int{
		[]int{
			0,
			1,
			17,
			4,
			3,
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
