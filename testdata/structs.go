package testdata

import "database/sql"

type Fizz struct {
	Foo  int
	Herp string
}

type Buzz struct {
	Bar  bool
	Derp uint8
	Meep sql.NullString
}
