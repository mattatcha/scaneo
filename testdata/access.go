package testdata

type Exported struct {
	A int
	B int
}

type unexported struct {
	a int
	b int
}

type ExAndUn struct {
	a int
	b int
}

type unAndEx struct {
	A int
	B int
}
