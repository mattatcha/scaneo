package testdata

import "database/sql"

type Ipsum struct {
	A sql.NullBool
	B sql.NullFloat64
	C sql.NullInt64
	D sql.NullString
}
