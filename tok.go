package main

import (
	"bytes"
	"log"
)

type structToken struct {
	Name   string
	Fields []string
	Types  []string
}

func (tok structToken) String() string {
	if len(tok.Fields) != len(tok.Types) {
		log.Println("len(tok.Fields) != len(tok.Types)")
		log.Println("something went wrong with the parsing...")
		log.Println("continuing anyway")
	}

	var buf bytes.Buffer
	buf.WriteString(tok.Name)
	buf.WriteString("\n")

	for i, _ := range tok.Fields {
		buf.WriteString("    ")
		buf.WriteString(tok.Fields[i])
		buf.WriteString(" ")
		buf.WriteString(tok.Types[i])
		buf.WriteString("\n")
	}

	return buf.String()
}
