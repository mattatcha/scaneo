package testdata

type boolean struct {
	a bool
}

type numerics struct {
	a uint8
	b uint16
	c uint32
	d uint64
	e int8
	f int16
	g int32
	h int64
	i float32
	j float64
	k complex64
	l complex128
	m byte
	n rune
	o uint
	p int
	q uintptr
}

type str struct {
	a string
}

type slices struct {
	a []bool
	b []int
	c []string
}

type structs struct {
	a time.Time
	b pq.NullTime
	c sql.NullString
}

type pointers struct {
	a *bool
	b *int
	c *string
	d *[]byte
	e *time.Time
}
