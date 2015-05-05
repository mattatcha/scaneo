package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

const (
	srcFile = "testdata/structs.go"
)

type structData struct {
	Name  string
	Field []string
	Type  []string
}

func main() {
	sdatas := make([]structData, 0, 8)

	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, srcFile, nil, 0)
	if err != nil {
		log.Println("'syntax error' - parser probably")
		log.Fatal(err)
	}

	// ast.Print(fset, astf)
	for _, dec := range astf.Decls {
		sdata := structData{
			Field: make([]string, 0, 8),
			Type:  make([]string, 0, 8),
		}

		decType, isGenDec := dec.(*ast.GenDecl)
		if !isGenDec {
			continue
		}

		for _, spec := range decType.Specs {
			sp, isTypeSpec := spec.(*ast.TypeSpec)
			if !isTypeSpec {
				continue
			}

			sdata.Name = sp.Name.Name

			st, isStructType := sp.Type.(*ast.StructType)
			if !isStructType {
				continue
			}

			//  List: []*ast.Field (len = 2)
			for _, fl := range st.Fields.List {
				for _, ident := range fl.Names {
					sdata.Field = append(sdata.Field, ident.Name)
				}

				switch tp := fl.Type.(type) {
				case *ast.Ident:
					sdata.Type = append(sdata.Type, tp.Name)
				case *ast.SelectorExpr:
					ident, isIdent := tp.X.(*ast.Ident)
					if !isIdent {
						continue
					}

					sdata.Type = append(sdata.Type,
						fmt.Sprintf("%s.%s", ident.Name, tp.Sel.Name))
				}

			}

			sdatas = append(sdatas, sdata)
		}
	}

	for _, sd := range sdatas {
		fmt.Println(sd.Name)

		if len(sd.Field) != len(sd.Type) {
			log.Println(sd)
			log.Fatal("Ahhh!!")
		}

		for i, _ := range sd.Field {
			fmt.Println("   ", sd.Field[i], sd.Type[i])
		}
	}
}
