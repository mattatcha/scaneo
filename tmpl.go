package main

const (
	scanStructFunc = `{{range .}}func Scan{{.StructName}}(r *sql.Row) ({{.StructName}}, error) {
	var s {{.StructName}}

	if err := r.Scan({{range .FieldName}}
		&s.{{.}},{{end}}
	); err != nil {
		return {{.StructName}}{}, err
	}

	return s, nil
}

{{end}}`

	scanStructsFunc = `{{range .}}func Scan{{.StructName}}s(rs *sql.Rows) ([]{{.StructName}}, error) {
	structs := make([]{{.StructName}}, 0, 16)

	var err error
	for rs.Next() {
		var s {{.StructName}}

		if err = rs.Scan({{range .FieldName}}
			&s.{{.}},{{end}}
		); err != nil {
			return nil, err
		}

		structs = append(structs, s)
	}

	return structs, nil
}

{{end}}`
)
