package main

const (
	scanStructFunc = `{{range .Tokens}}func {{$.Access}}can{{.Name}}(r *sql.Row) ({{.Name}}, error) {
	var s {{.Name}}

	if err := r.Scan({{range .Fields}}
		&s.{{.}},{{end}}
	); err != nil {
		return {{.Name}}{}, err
	}

	return s, nil
}

{{end}}`

	scanStructsFunc = `{{range .Tokens}}func {{$.Access}}can{{.Name}}s(rs *sql.Rows) ([]{{.Name}}, error) {
	structs := make([]{{.Name}}, 0, 16)

	var err error
	for rs.Next() {
		var s {{.Name}}

		if err = rs.Scan({{range .Fields}}
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
