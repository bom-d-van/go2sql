package main

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/bom-d-van/goutils/printutils"

	"testing"
)

func TestVisitor(t *testing.T) {
	var testFile = `package model

import (
	"html/template"
	"time"

	"bitbucket.org/pkg/inflect"
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

	Rule inflect.Rule

	// TODO: support array
	Keywords []*Keyword

	HTML template.HTML

	Teachers []Teacher
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
}

type TeacherLanguageXref struct {
	LanguageID uint
	TeacherID  uint
}
`

	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "model.go", testFile, parser.ParseComments|parser.DeclarationErrors|parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}

	v := newVisitor()
	ast.Walk(v, file)

	v.analyze()

	// log.Printf("--> %+v\n", v)

	language := v.tables["Language"]
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
		t.Errorf("Language.Tag.Relationship = %s; want %s", got, want)
	}
	teachers := language.GetColumn("Teachers")
	if got, want := teachers.Relationship, RelationshipManyToMany; got != want {
		t.Errorf("Language.Tag.Relationship = %s; want %s", got, want)
	}

	language.Package = "model"
	if err := fileHeader.Execute(&language.w, language); err != nil {
		t.Error(err)
	}
	if err := findTmpl.Execute(&language.w, language); err != nil {
		t.Error(err)
	}
	// log.Println(language.w.String())
	printutils.PrettyPrint(language.Columns)
}
