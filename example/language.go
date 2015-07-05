package model

import "time"

type Language struct {
	ID         uint `sqlgen:"auto,primary-key"`
	Name       string
	WordsCount uint

	Ignored string `sqlgen:"-"`

	// Info
	// Origin Origin `sqlgen:"inline"`

	AuthorID uint
	Author   Person

	Keywords []Keyword

	Teachers []Teacher
	// LanguagesTeachers []LanguageTeacher
}

type Info struct {
	CreatedAt   time.Time
	Description string
}

type Keyword struct {
	ID   uint `sqlgen:"auto,primary-key"`
	Name string
	Type string

	LanguageID uint
}

type Person struct {
	ID    uint `sqlgen:"auto,primary-key"`
	Name  string
	Email string

	// TODO: DeepSave
	// AnotherTableHere
}

type Teacher struct {
	ID   uint `sqlgen:"auto,primary-key"`
	Name string
	Age  uint
}

type LanguageTeacher struct {
	LanguageID uint
	TeacherID  uint
}
