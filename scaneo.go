package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

const (
	usageText = `SCANEO
    Generate Go code to convert database rows into arbitrary structs.

USAGE
    scaneo [options] paths...

OPTIONS
    -o, -output
        Set the name of the generated file. Default is scans.go.

    -p, -package
        Set the package name for the generated file. Default is current
        directory name.

    -u, -unexport
        Generate unexported functions. Default is export all.

    -w, -whitelist
        Only include structs specified in case-sensitive, comma-delimited
        string.

    -v, -version
        Print version and exit.

    -h, -help
        Print help and exit.

EXAMPLES
    tables.go is a file that contains one or more struct declarations.

    Generate scan functions based on structs in tables.go.
        scaneo tables.go

    Generate scan functions and name the output file funcs.go
        scaneo -o funcs.go tables.go

    Generate scans.go with unexported functions.
        scaneo -u tables.go

    Generate scans.go with only struct Post and struct user.
        scaneo -w "Post,user" tables.go

NOTES
    Struct field names don't have to match database column names at all.
    However, the order of the types must match.

    Integrate this with go generate by adding this line to the top of your
    tables.go file.
        //go:generate scaneo $GOFILE
`
)

type fieldToken struct {
	Name string
	Type string
}

type structToken struct {
	Name   string
	Fields []fieldToken
}

var (
	scansTmpl = template.Must(template.New("scans").Parse(scansText))

	outFilename = flag.String("o", "scans.go", "")
	packName    = flag.String("p", "current directory", "")
	unexport    = flag.Bool("u", false, "")
	whitelist   = flag.String("w", "", "")
	version     = flag.Bool("v", false, "")
	help        = flag.Bool("h", false, "")
)

func init() {
	flag.StringVar(outFilename, "output", "scans.go", "")
	flag.StringVar(packName, "package", "current directory", "")
	flag.BoolVar(unexport, "unexport", false, "")
	flag.StringVar(whitelist, "whitelist", "", "")
	flag.BoolVar(version, "version", false, "")
	flag.BoolVar(help, "help", false, "")

	// if err, send usage to stderr
	// if not err, send usage to stdout
	// that way people can: scaneo -h | less

	flag.Usage = func() {
		// an error happened, send to stderr
		log.Println(usageText)
	}
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *help {
		// not an error, send to stdout
		fmt.Println(usageText)
		return
	}

	if *version {
		fmt.Println("scaneo version 1.1.1")
		return
	}

	inputPaths := flag.Args()
	if len(inputPaths) == 0 {
		log.Println("missing input paths")
		log.Println(usageText)
		os.Exit(1)
	}

	if *packName == "current directory" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalln("couldn't get working directory:", err)
		}

		*packName = filepath.Base(wd)
	}

	files, err := filenames(inputPaths)
	if err != nil {
		log.Fatalln("couldn't get filenames:", err)
	}

	wmap := make(map[string]struct{})
	if *whitelist != "" {
		wSplits := strings.Split(*whitelist, ",")
		for _, s := range wSplits {
			wmap[s] = struct{}{}
		}
	}

	structToks := make([]structToken, 0, 8)
	for _, file := range files {
		toks, err := parseCode(file, wmap)
		if err != nil {
			log.Println(`"syntax error" - parser probably`)
			log.Fatal(err)
		}

		structToks = append(structToks, toks...)
	}

	if len(structToks) < 1 {
		log.Println("heads up! no structs found")
		log.Println("aborting")
		os.Exit(1)
	}

	fout, err := os.Create(*outFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()

	if err := genFile(fout, *packName, *unexport, structToks); err != nil {
		log.Fatalln("couldn't generate file:", err)
	}
}

func filenames(paths []string) ([]string, error) {
	files := make([]string, 0, 8)

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			filepath.Walk(path, func(startDir string, subInfo os.FileInfo, _ error) error {
				if subInfo.IsDir() {
					return nil
				} else if subInfo.Name()[0] == '.' {
					return nil
				}

				files = append(files, startDir)
				return nil
			})

			continue
		}

		files = append(files, path)
	}

	// file list needs to be deduped in case of
	// scaneo foo.go foo.go
	// scaneo somedir/ somedir/foo.go
	// and because it wouldn't make sense to generate code for the same struct twice
	fileMap := make(map[string]struct{})
	for _, f := range files {
		fileMap[f] = struct{}{}
	}

	deduped := make([]string, 0, len(fileMap))
	for f := range fileMap {
		deduped = append(deduped, f)
	}

	return deduped, nil
}

func parseCode(srcFile string, wlist map[string]struct{}) ([]structToken, error) {
	structToks := make([]structToken, 0, 8)

	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, srcFile, nil, 0)
	if err != nil {
		return nil, err
	}

	var filter bool
	if len(wlist) > 0 {
		filter = true
	}

	//ast.Print(fset, astf)
	for _, decl := range astf.Decls {
		genDecl, isGeneralDeclaration := decl.(*ast.GenDecl)
		if !isGeneralDeclaration {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, isTypeDeclaration := spec.(*ast.TypeSpec)
			if !isTypeDeclaration {
				continue
			}

			structType, isStructTypeDeclaration := typeSpec.Type.(*ast.StructType)
			if !isStructTypeDeclaration {
				continue
			}

			// found a struct in the source code!

			var structTok structToken

			// filter logic
			if structName := typeSpec.Name.Name; !filter {
				// no filter, collect everything
				structTok.Name = structName
			} else if _, exists := wlist[structName]; filter && !exists {
				// if structName not in whitelist, continue
				continue
			} else if filter && exists {
				// structName exists in whitelist
				structTok.Name = structName
			}

			structTok.Fields = make([]fieldToken, 0, len(structType.Fields.List))

			// iterate through struct fields (1 line at a time)
			for _, fieldLine := range structType.Fields.List {
				fieldToks := make([]fieldToken, len(fieldLine.Names))

				// get field name (or names because multiple vars can be declared in 1 line)
				for i, fieldName := range fieldLine.Names {
					fieldToks[i].Name = parseIdent(fieldName)
				}

				var fieldType string

				// get field type
				switch typeToken := fieldLine.Type.(type) {
				case *ast.Ident:
					// simple types, e.g. bool, int
					fieldType = parseIdent(typeToken)
				case *ast.SelectorExpr:
					// struct fields, e.g. time.Time, sql.NullString
					if fieldType = parseSelector(typeToken); fieldType == "" {
						continue
					}
				case *ast.ArrayType:
					// arrays
					if fieldType = parseArray(typeToken); fieldType == "" {
						continue
					}
				case *ast.StarExpr:
					// pointers
				}

				// apply type to all variables declared in this line
				for i := range fieldToks {
					fieldToks[i].Type = fieldType
				}

				structTok.Fields = append(structTok.Fields, fieldToks...)
			}

			structToks = append(structToks, structTok)
		}
	}

	return structToks, nil
}

func parseIdent(fieldType *ast.Ident) string {
	// return like byte, string, int
	return fieldType.Name
}

func parseSelector(fieldType *ast.SelectorExpr) string {
	// return like time.Time, sql.NullString
	ident, isIdent := fieldType.X.(*ast.Ident)
	if !isIdent {
		return ""
	}

	return fmt.Sprintf("%s.%s", parseIdent(ident), fieldType.Sel.Name)
}

func parseArray(fieldType *ast.ArrayType) string {
	// return like []byte, []time.Time, []*byte, []*sql.NullString
	switch arrayType := fieldType.Elt.(type) {
	case *ast.Ident:
		return fmt.Sprintf("[]%s", parseIdent(arrayType))
	case *ast.SelectorExpr:
		sel := parseSelector(arrayType)
		if sel == "" {
			return ""
		}
		return fmt.Sprintf("[]%s", sel)
	}

	return ""
}

func parseStar(fieldType *ast.StarExpr) string {
	// return like *bool, *time.Time, *[]byte, and other array stuff
	return ""
}

func genFile(fout *os.File, pkg string, unexport bool, toks []structToken) error {
	data := struct {
		PackageName string
		Tokens      []structToken
		Access      string
	}{
		PackageName: pkg,
		Access:      "S",
	}

	// always capitalize the first letter of struct types so camel case is
	// correct, e.g. if struct is foo, then generate scanFoo or ScanFoo, but
	// never scanfoo or Scanfoo
	for i := range toks {
		toks[i].Name = string(unicode.ToTitle(rune(toks[i].Name[0]))) +
			toks[i].Name[1:]
	}

	data.Tokens = toks

	if unexport {
		data.Access = "s"
	}

	if err := scansTmpl.Execute(fout, data); err != nil {
		return err
	}

	return nil
}
