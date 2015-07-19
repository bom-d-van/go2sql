package model

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/bom-d-van/go2sql/go2sql"

	_ "github.com/go-sql-driver/mysql"

	"testing"
	"time"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@/go2sql_example")
	if err != nil {
		panic(err)
	}

	go2sql.SetDefaultDB(db)
}

func resetDB() {
	must(db.Exec(`
		DROP TABLE IF EXISTS languages;
	`))

	must(db.Exec(`
		Create TABLE languages (
			id int NOT NULL AUTO_INCREMENT,
			name TEXT,
			words_stat int,
			field1 varchar(255) not null default 'text',
			field2 varchar(255) not null default 'text',
			field3 varchar(255) not null default 'text',
			field4 varchar(255) not null default 'text',
			field5 varchar(255) not null default 'text',
			field6 varchar(255) not null default 'text',
			field7 varchar(255) not null default 'text',
			PRIMARY KEY (id)
		);
	`))
	rand.Seed(time.Now().UnixNano())
	for i := 1; i < 100; i++ {
		must(db.Exec(fmt.Sprintf("INSERT INTO languages (name, words_stat) VALUES ('Mr. Tester', %d);", i)))
	}
}

func must(r sql.Result, err error) {
	if err != nil {
		panic(err)
	}
}

func TestFindLanguage(t *testing.T) {
	resetDB()

	var l *Language
	var err error

	// normal
	l, err = FindLanguage()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := l.Name, "Mr. Tester"; got != want {
		t.Errorf("l.Name = %s; want %s", got, want)
	}

	// sql patial
	l, err = FindLanguage(go2sql.NewSQL("order by id desc"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := l.ID, uint(99); got != want {
		t.Errorf("l.Name = %d; want %d", got, want)
	}

	// db specification
	go2sql.SetDefaultDB(nil)
	_, err = FindLanguage()
	if err == nil {
		t.Error("expect error for nil db")
	}
	l, err = FindLanguage(go2sql.DB(db))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := l.Name, "Mr. Tester"; got != want {
		t.Errorf("l.Name = %s; want %s", got, want)
	}
	go2sql.SetDefaultDB(db)

	// selects
	l, err = FindLanguage(go2sql.Selects{"id", "words_stat"})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := l.Name, ""; got != want {
		t.Errorf("l.Name = %s; want %s", got, want)
	}
	if got, want := l.WordsCount, uint(1); got != want {
		t.Errorf("l.Name = %d; want %d", got, want)
	}
}

func TestFindLanguages(t *testing.T) {
	resetDB()

	var ls Languages
	var err error

	// normal
	ls, err = FindLanguages()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(ls), 99; got != want {
		t.Error("len(ls) = %d; want %d", got, want)
	}
	if got, want := ls[0].Name, "Mr. Tester"; got != want {
		t.Errorf("ls[0].Name = %s; want %s", got, want)
	}

	// sql patial
	ls, err = FindLanguages(go2sql.NewSQL("order by id desc limit 10"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(ls), 10; got != want {
		t.Error("len(ls) = %d; want %d", got, want)
	}
	if got, want := ls[0].ID, uint(99); got != want {
		t.Errorf("ls[0].Name = %d; want %d", got, want)
	}

	// db specification
	go2sql.SetDefaultDB(nil)
	_, err = FindLanguages()
	if err == nil {
		t.Error("expect error for nil db")
	}
	ls, err = FindLanguages(go2sql.DB(db))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := ls[0].Name, "Mr. Tester"; got != want {
		t.Errorf("ls[0].Name = %s; want %s", got, want)
	}
	go2sql.SetDefaultDB(db)

	// selects
	ls, err = FindLanguages(go2sql.Selects{"id", "words_stat"})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := ls[0].Name, ""; got != want {
		t.Errorf("ls[0].Name = %s; want %s", got, want)
	}
	if got, want := ls[0].WordsCount, uint(1); got != want {
		t.Errorf("ls[0].Name = %d; want %d", got, want)
	}
}

func TestInsertLanguage(t *testing.T) {
	resetDB()

	var Language Language
}
