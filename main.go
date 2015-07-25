package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bom-d-van/goutils/printutils"
)

func main() {
	var pkgStr string
	flag.StringVar(&pkgStr, "pkg", "", "import path")
	flag.Parse()
	args := flag.Args()

	if pkgStr == "" {
		// TODO: use pwd as default import path
		exitf("need to specify package import path")
	}
	if len(args) == 0 {
		exitf("need to specify types")
	}

	// log.Printf("--> %+v\n", typ)

	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, filepath.Join(os.Getenv("GOPATH"), "src", pkgStr), ignoreGo2SQLFiles, parser.ParseComments|parser.DeclarationErrors|parser.AllErrors)
	if err != nil {
		exitf(err.Error())
	}

	var pkg *ast.Package
	for _, pkg = range pkgs {
	}
	// v := visitor{typ: args[0]}
	v := newVisitor()
	ast.Walk(v, pkg)
	printutils.PrettyPrint(v.tables["Language"].GetColumn("Author"))
	// log.Printf("--> %+v\n", v.consts)
	// printutils.PrettyPrint(v.tables)

	// for _, f := range v.struc.Fields.List {
	// 	switch typ := f.Type.(type) {
	// 	case *ast.Ident:
	// 		log.Printf("--> %+v\n", typ.Name)
	// 	case *ast.ArrayType:
	// 		log.Printf("--> %+v\n", typ.Elt)
	// 	}
	// }

	// pkg, err := types.Check(filepath.Join(os.Getenv("GOPATH"), "src", pkgStr), fs, files)
	// if err != nil {
	// 	exitf(err.Error())
	// }
	// typ := args[0]
	// log.Printf("--> %+v\n", typ)
	// // o := pkg.Scope().Lookup(typ)
	// log.Printf("--> %+v\n", pkg.Complete())
	// log.Printf("--> %+v\n", pkg.Scope().Names())
}

func ignoreGo2SQLFiles(fi os.FileInfo) bool {
	if fi.IsDir() {
		return true
	}
	if strings.HasSuffix(fi.Name(), Go2SQLFileSuffix+".go") || strings.HasSuffix(fi.Name(), Go2SQLFileSuffix+"_test.go") {
		return false
	}
	return true
}

func exitf(fmt string, args ...interface{}) {
	log.Printf("need to specify package import path"+"\n", args...)
	os.Exit(1)
}
