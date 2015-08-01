package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"

	"testing"
)

var funcsf string

func init() {
	flag.StringVar(&funcsf, "funcs", "", "func flags")
	flag.Parse()
}

func TestVisitor(t *testing.T) {
	var testFile = `package model

import (
	"html/template"
	"time"

	// "bitbucket.org/pkg/inflect"
)

const (
	LanguageTableName = "languages"
)

type MyString string

type Language struct {
	ID         uint ` + "`go2sql:\",id,primary-key\"`" + `
	Name       string
	WordsCount uint ` + "`go2sql:\"word_stat\"`" + `

	// WordsCount *uint
	// HTML template.HTML

	Ignored string ` + "`go2sql:\"-\"`" + `

	// Info
	// Origin Origin ` + "`go2sql:\",inline\"`" + `
	// Name Type     ` + "`go2sql:\"name2\"`" + `

	AuthorID **uint
	Author   Person

	Tag *Keyword
	// Maintainer Person

	// Embed struct {
	// 	Name string
	// }

	MyString MyString

	// Rule inflect.Rule

	// TODO: support array
	Keywords []*Keyword

	HTML template.HTML

	// Teachers   []Teacher
	// TeacherIDs []uint
	// LanguagesTeachers []LanguageTeacher
}

type Info struct {
	CreatedAt   time.Time
	Description string
}

type Keyword struct {
	ID   uint ` + "`go2sql:\",id,primary-key\"`" + `
	Name string
	Type string

	LanguageID uint
}

type Person struct {
	ID    uint ` + "`go2sql:\",id,primary-key\"`" + `
	Name  string
	Email string

	// LanguageID

	// TODO: DeepSave
	// AnotherTableHere
}

type Teacher struct {
	ID   uint ` + "`go2sql:\",id,primary-key\"`" + `
	Name string
	Age  uint

	Languages []*Language
	LanguageID uint
}

type TeacherLanguageXref struct {
	LanguageID uint
	TeacherIDs []uint
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "model.go", testFile, parser.ParseComments|parser.DeclarationErrors|parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}

	info := types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
		InitOrder:  make([]*types.Initializer, 0),
	}

	var conf types.Config
	// log.Println(f.Name.Name)
	conf.Importer = importer.For("gc", nil)
	pkg, err := conf.Check(file.Name.Name, fset, []*ast.File{file}, &info)
	if err != nil {
		panic(err)
	}

	// analyzeConfig(info)
	// var p Parser
	// p.Package = pkg.Name()
	p := NewParser(pkg.Name())
	if err := p.Parse(info); err != nil {
		t.Error(err)
	}

	// v := newVisitor()
	// ast.Walk(v, file)

	// v.analyze()

	// // log.Printf("--> %+v\n", v)

	language := p.Tables["Language"]
	author := language.GetColumn("Author")
	if got, want := author.Relationship, RelationshipBelongsTo; got != want {
		t.Errorf("Language.Author.Relationship = %s; want %s", got, want)
	}
	tag := language.GetColumn("Tag")
	if got, want := tag.Relationship, RelationshipHasOne; got != want {
		t.Errorf("Language.Tag.Relationship = %s; want %s", got, want)
	}
	keywords := language.GetColumn("Keywords")
	if got, want := keywords.Relationship, RelationshipHasMany; got != want {
		t.Errorf("Language.Keywords.Relationship = %s; want %s", got, want)
	}
	// teachers := language.GetColumn("Teachers")
	// if got, want := teachers.Relationship, RelationshipManyToMany; got != want {
	// 	t.Errorf("Language.Teachers.Relationship = %s; want %s", got, want)
	// }

	language.Package = "model"
	// if err := fileHeader.Execute(&language.w, language); err != nil {
	// 	t.Error(err)
	// }
	// if err := findTmpl.Execute(&language.w, language); err != nil {
	// 	t.Error(err)
	// }
	funcs := []string{
		"header",
	}
	if funcsf != "" {
		funcs = append(funcs, strings.Split(funcsf, ",")...)
	} else {
		funcs = append(
			funcs,
			"is_empty_row", "is_new_row",
			"find", "find_many",
			"insert", "insert_many",
			"update", "update_many",
			"delete", "delete_many",
		)
	}
	for _, name := range funcs {
		if err := tmpl.ExecuteTemplate(&language.w, name, language); err != nil {
			printWithLineNo(rawTmpl)
			t.Fatal(err)
		}
	}

	src, err := format.Source(language.w.Bytes())
	if err != nil {
		printWithLineNo(language.w.String())
		t.Error(err)
	}
	printWithLineNo(string(src))
	// printutils.PrettyPrint(language.Columns)
}

func printWithLineNo(src string) {
	codes := strings.Split(src, "\n")
	for i, line := range codes {
		fmt.Printf("%3d: %s\n", i+1, line)
	}
}
