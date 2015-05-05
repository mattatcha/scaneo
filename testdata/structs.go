package testdata

import "database/sql"

type Fizz struct{}

func (f Fizz) FizzFunc() int {
	return 0
}

type Buzz struct {
	A bool
}

type Foo struct {
	A byte
	B complex128
	C complex64
	D float32
	E float64
	F int
	G int16
	H int32
	I int64
	J int8
	K rune
	L uint
	M uint16
	N uint32
	O uint64
	P uint8
	Q uintptr
}

type Bar struct {
	A sql.NullBool
	B sql.NullFloat64
	C sql.NullInt64
	D sql.NullString
}

type Herp struct {
	a int
	b int
	C int
}

func someFunc() bool {
	return false
}

type Derp struct {
	A string
	b string
	c string
}
