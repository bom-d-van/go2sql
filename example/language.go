package model

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
	ID         uint `go2sql:",id,primary-key"`
	Name       string
	WordsCount uint `go2sql:"word_stat" db:"words_count"`

	// WordsCount *uint
	// HTML template.HTML

	Ignored string `go2sql:"-"`

	Field1 string
	Field2 string
	Field3 string
	Field4 string
	Field5 string
	Field6 string
	Field7 string

	// Info
	// Origin Origin `go2sql:",inline"`
	// Name Type     `go2sql:"name2"`

	AuthorID uint
	Author   *Person

	Embed struct {
		Name string
	}

	MyString MyString

	Rule inflect.Rule

	// TODO: support array
	Keywords []*Keyword

	HTML template.HTML // TODO: convert it to string(HTML) when scanning

	Teachers []*Teacher
	// LanguagesTeachers []LanguageTeacher
	TeacherID uint
}

type Info struct {
	CreatedAt   time.Time
	Description string
}

type Keyword struct {
	ID   uint `go2sql:",id,primary-key"`
	Name string
	Type string

	LanguageID uint
}

type Person struct {
	ID    uint `go2sql:",id,primary-key"`
	Name  string
	Email string

	// TODO: DeepSave
	// AnotherTableHere
}

type Teacher struct {
	ID   uint `go2sql:",id,primary-key"`
	Name string
	Age  uint

	LanguageID uint
}

// type LanguageTeacherXref struct {
// 	LanguageID uint
// 	TeacherID  uint
// }
