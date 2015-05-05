package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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

func main() {
	toks, err := parseCode("testdata/structs.go")
	if err != nil {
		log.Println(`"syntax error" - parser probably`)
		log.Fatal(err)
	}

	for _, t := range toks {
		fmt.Println(t)
	}
}

func parseCode(srcFile string) ([]structToken, error) {
	structToks := make([]structToken, 0, 8)

	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, srcFile, nil, 0)
	if err != nil {
		return nil, err
	}

	// ast.Print(fset, astf)
	for _, dec := range astf.Decls {
		structTok := structToken{
			Fields: make([]string, 0, 8),
			Types:  make([]string, 0, 8),
		}

		genDec, isGenDec := dec.(*ast.GenDecl)
		if !isGenDec {
			continue
		}

		for _, spec := range genDec.Specs {
			typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
			if !isTypeSpec {
				continue
			}

			structTok.Name = typeSpec.Name.Name

			structType, isStructType := typeSpec.Type.(*ast.StructType)
			if !isStructType {
				continue
			}

			for _, field := range structType.Fields.List {
				for _, ident := range field.Names {
					structTok.Fields = append(structTok.Fields, ident.Name)
				}

				switch fieldType := field.Type.(type) {
				case *ast.Ident:
					structTok.Types = append(structTok.Types, fieldType.Name)
				case *ast.SelectorExpr:
					ident, isIdent := fieldType.X.(*ast.Ident)
					if !isIdent {
						continue
					}

					structTok.Types = append(structTok.Types,
						fmt.Sprint(ident.Name, ".", fieldType.Sel.Name))
				}

			}

			structToks = append(structToks, structTok)
		}
	}

	return structToks, nil
}
