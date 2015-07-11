package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/bom-d-van/goutils/printutils"
)

const (
	FlagID     = "id"
	FlagPK     = "primary-key"
	FlagInline = "inline"
	FlagIgnore = "-"

	FlagPrefix = "prefix:"

	Go2SQLFileSuffix = "go2sql"
)

var cliFlags = struct {
	debug bool
}{}

func init() {
	log.SetFlags(log.Lshortfile)
}

type Option struct {
	Functions []string
}

type Function struct {
	Name     string
	Template template.Template
}

type Table struct {
	Name        string
	SQLName     string
	RefName     string
	IDColumn    *Column
	Columns     []*Column
	PrimaryKeys []*Column

	BelongsTo   []*Table
	HasOnes     []*Table
	HasManys    []*Table
	ManyToManys []*Table
}

type Column struct {
	f            *ast.Field
	Name         string
	SQLName      string
	IsPrimaryKey bool
	Relationship Relationship

	Type      string
	IsPointer bool

	flags []string

	IsTable bool
	Table   *Table
}

type Relationship int

const (
	RelationshipBelongsTo Relationship = iota
	RelationshipHasOne
	RelationshipHasMany
	RelationshipManyToMany
)

func (t *Table) HasColumn(c string) bool {
	for _, cl := range t.Columns {
		if cl.Name == c {
			return true
		}
	}
	return false
}

func (t *Table) GetColumn(c string) *Column {
	for _, cl := range t.Columns {
		if cl.Name == c {
			return cl
		}
	}
	return nil
}

type visitor struct {
	// typ        string
	// foundIdent bool
	// done       bool
	// struc      *ast.TableType
	typ    string
	consts map[string]string
	tables map[string]*Table
}

func newVisitor() *visitor {
	var v visitor
	v.tables = map[string]*Table{}
	v.consts = map[string]string{}
	return &v
}

func (v *visitor) Visit(n ast.Node) (w ast.Visitor) {
	switch node := n.(type) {
	case *ast.Ident:
		v.typ = node.Name
	case *ast.ValueSpec:
		for i, name := range node.Names {
			// || !strings.HasSuffix(name.String(), "TableName")
			if name.Obj.Kind != ast.Con || len(node.Values) <= i {
				continue
			}
			bl, ok := node.Values[i].(*ast.BasicLit)
			if !ok {
				continue
			}

			v.consts[name.String()] = bl.Value
		}
	case *ast.StructType:
		var table Table
		table.Name = v.typ
		table.RefName = strings.ToLower(v.typ[:1])
	listLoop:
		for _, f := range node.Fields.List {
			var flags []string
			if f.Tag != nil {
				tag, err := strconv.Unquote(f.Tag.Value)
				if err != nil {
					log.Printf("failed to unquote tag in %s.%s: %s\n", v.typ, f.Names[0].Name, err)
				}
				flags = strings.Split(reflect.StructTag(tag).Get("go2sql"), ",")
			}
			for _, n := range f.Names {
				var column Column
				column.Name = n.Name
				column.f = f

				if len(flags) > 0 && flags[0] != "" {
					if flags[0] == FlagIgnore {
						continue listLoop
					}
					column.SQLName = flags[0]
				} else {
					column.SQLName = toSnake(column.Name)
				}
				column.flags = flags
				_, column.IsPointer = f.Type.(*ast.StarExpr)
				if contains(flags, FlagID) {
					table.IDColumn = &column
				}
				if contains(flags, FlagPK) {
					column.IsPrimaryKey = true
					table.PrimaryKeys = append(table.PrimaryKeys, &column)
				}

				column.Type = types.ExprString(f.Type)
				column.IsTable, column.Relationship = isTable(f.Type)
				if column.IsTable && contains(flags, FlagInline) {
					column.IsTable = false
				}
				table.Columns = append(table.Columns, &column)
			}
		}
		v.tables[table.Name] = &table
	}
	return v
}

func contains(flags []string, f string) bool {
	for _, fl := range flags {
		if fl == f {
			return true
		}
	}
	return false
}

func isTable(expr ast.Expr) (ok bool, rel Relationship) {
	switch typ := expr.(type) {
	case *ast.Ident:
		log.Println(typ.Name, typ.Obj)
		if typ.Obj == nil {
			return
		}
		log.Printf("--> %+v\n", typ.Name)
		switch decl := typ.Obj.Decl.(type) {
		case *ast.TypeSpec:
			return isTable(decl.Type)
		case ast.Expr:
			return isTable(decl)
		case *ast.StructType:
			log.Printf("--> %+v\n", typ.Name)
			return true, RelationshipHasOne
		default:
			log.Printf("unknown %s", decl)
		}
	case *ast.StarExpr:
		if _, is := typ.X.(*ast.StarExpr); is {
			if cliFlags.debug {
				log.Printf("pointer of pointer is not supported: %s\n", types.ExprString(expr))
			}
			ok = false
			return
		}
		return isTable(typ.X)
	case *ast.SelectorExpr:
		// return isTable(typ.X)
	case *ast.ArrayType:
		rel = RelationshipHasMany
		ok, _ = isTable(typ.Elt)
	case *ast.StructType:
		return true, RelationshipHasOne
	case *ast.SliceExpr:
		rel = RelationshipHasMany
		ok, _ = isTable(typ.X)
		// case *ast.MapType:
		// case *ast.FuncType:
	}
	return
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
