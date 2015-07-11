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
	WordsCount uint `go2sql:"word_stat"`

	// WordsCount *uint
	// HTML template.HTML

	Ignored string `go2sql:"-"`

	// Info
	// Origin Origin `go2sql:",inline"`
	// Name Type     `go2sql:"name2"`

	AuthorID **uint
	Author   Person

	Embed struct {
		Name string
	}

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
}

type LanguageTeacher struct {
	LanguageID uint
	TeacherID  uint
}
