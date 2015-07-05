package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

type visitor struct {
	typ        string
	foundIdent bool
	done       bool
	struc      *ast.StructType
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.Ident:
		if n.Name == v.typ {
			// log.Printf("--> %+v\n", n)
			v.foundIdent = true
		}
	case *ast.StructType:
		if !v.foundIdent {
			break
		}
		v.struc = n
		// for _, f := range n.Fields.List {
		// }
		v.done = true
		return nil
	}
	if v.done {
		return nil
	}
	return v
}

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
	pkgs, err := parser.ParseDir(fs, filepath.Join(os.Getenv("GOPATH"), "src", pkgStr), nil, parser.ParseComments|parser.DeclarationErrors|parser.AllErrors)
	if err != nil {
		exitf(err.Error())
	}

	var pkg *ast.Package
	for _, pkg = range pkgs {
	}
	v := visitor{typ: args[0]}
	ast.Walk(&v, pkg)

	for _, f := range v.struc.Fields.List {
		switch typ := f.Type.(type) {
		case *ast.Ident:
			log.Printf("--> %+v\n", typ.Name)
		case *ast.ArrayType:
			log.Printf("--> %+v\n", typ.Elt)
		}
	}

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

func exitf(fmt string, args ...interface{}) {
	log.Printf("need to specify package import path"+"\n", args...)
	os.Exit(1)
}
