package model

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	LanguageTableName = "languages"

	LanguageColumnID         = "id"
	LanguageColumnName       = "name"
	LanguageColumnWordsCount = "words_count"
)

// TODO: use FindLanguageSQL
func FindLanguage(db *sql.DB) (l *Language, err error) {
	return FindLanguageSQL(db, "")
}

func FindLanguageSQL(db *sql.DB, sql string) (l *Language, err error) {
	l = &Language{}
	err = db.QueryRow("select id, name, words_count from languages "+sql).Scan(&l.ID, &l.Name, &l.WordsCount)
	return
}

func FindLanguages(db *sql.DB) (ls []*Language, err error) {
	return FindLanguagesSQL(db, "")
}

func FindLanguagesSQL(db *sql.DB, sql string) (ls []*Language, err error) {
	rows, err := db.Query("select * from languages " + sql)
	if err != nil {
		return
	}

	defer func() {
		if er := rows.Close(); er != nil {
			if err != nil {
				log.Println(er)
			} else {
				err = er
			}
		}
	}()

	for rows.Next() {
		var l Language
		err = rows.Scan(&l.ID, &l.Name, &l.WordsCount)
		if err != nil {
			return
		}
		ls = append(ls, &l)
	}

	return
}

// TODO: use SelectLanguageSQL
func SelectLanguage(db *sql.DB, columns []string) (l *Language, err error) {
	return SelectLanguageSQL(db, columns, "")
}

func SelectLanguageSQL(db *sql.DB, columns []string, sql string) (l *Language, err error) {
	l = &Language{}
	var fields []interface{}
	for _, c := range columns {
		switch c {
		case LanguageColumnID:
			fields = append(fields, &l.ID)
		case LanguageColumnName:
			fields = append(fields, &l.Name)
		case LanguageColumnWordsCount:
			fields = append(fields, &l.WordsCount)
		default:
			err = fmt.Errorf("go2sql: unknown column %s", c)
			return
		}
	}
	err = db.QueryRow("select " + strings.Join(columns, ",") + " from languages " + sql).Scan(fields...)
	return
}

func SelectLanguages(db *sql.DB, columns []string) (ls []*Language, err error) {
	return SelectLanguagesSQL(db, columns, "")
}

func SelectLanguagesSQL(db *sql.DB, columns []string, sql string) (ls []*Language, err error) {
	rows, err := db.Query("select " + strings.Join(columns, ",") + " from languages " + sql)
	if err != nil {
		return
	}

	defer func() {
		if er := rows.Close(); er != nil {
			if err != nil {
				log.Println(er)
			} else {
				err = er
			}
		}
	}()

	for rows.Next() {
		var l Language
		var fields []interface{}
		for _, c := range columns {
			switch c {
			case LanguageColumnID:
				fields = append(fields, &l.ID)
			case LanguageColumnName:
				fields = append(fields, &l.Name)
			case LanguageColumnWordsCount:
				fields = append(fields, &l.WordsCount)
			default:
				err = fmt.Errorf("go2sql: unknown column %s", c)
				return
			}
		}

		err = rows.Scan(fields...)
		if err != nil {
			return
		}
		ls = append(ls, &l)
	}

	return
}

func (l *Language) IsEmptyRow() bool {
	if l == nil {
		return true
	}

	return l.ID == 0 &&
		l.Name == "" &&
		l.WordsCount == 0 &&
		l.AuthorID == 0 &&
		l.Author.IsEmptyRow() &&
		len(l.Keywords) == 0 &&
		len(l.Teachers) == 0
}

func (l *Language) Insert(db *sql.DB) (r sql.Result, err error) {
	if !l.Author.IsEmptyRow() {
		if _, err = l.Author.Insert(db); err != nil {
			return
		}
		l.AuthorID = l.Author.ID
	}

	r, err = db.Exec(`INSERT INTO languages
			(name, words-count, author_id)
		VALUES(?, ?, ?)`, l.Name, l.WordsCount, l.AuthorID)
	if err != nil {
		return
	}
	id, err := r.LastInsertId()
	if err != nil {
		return
	}
	l.ID = uint(id)

	if len(l.Keywords) > 0 {
		for i, k := range l.Keywords {
			l.Keywords[i].LanguageID = l.ID
			if k.ID > 0 {
				if r, err = db.Exec("UPDATE keywords SET language_id = ? WHERE id = ?", l.ID, k.ID); err != nil {
					return
				}
			} else if _, err = l.Keywords[i].Insert(db); err != nil {
				return
			}
		}
	}

	if len(l.Teachers) > 0 {
		for i, t := range l.Teachers {
			if t.ID <= 0 {
				if _, err = l.Teachers[i].Insert(db); err != nil {
					return
				}
			}
			if _, err = db.Exec("INSERT INTO languages_teachers_xref (language_id, teacher_id) VALUES (?, ?)", l.ID, l.Teachers[i].ID); err != nil {
				return
			}
		}
	}

	return
}

// TODO: deep update
func (l *Language) Update(db *sql.DB) (r sql.Result, err error) {
	if !l.Author.IsEmptyRow() {
		if _, err = l.Author.Update(db); err != nil {
			return
		}
		l.AuthorID = l.Author.ID
	}

	if l.ID == 0 {
		return l.Insert(db)
	}

	r, err = db.Exec(`UPDATE languages SET name = ?, words_count = ? WHERE id = ?`, l.Name, l.WordsCount, l.ID)
	if err != nil {
		return
	}

	if len(l.Keywords) > 0 {
		for i := range l.Keywords {
			if _, err = l.Keywords[i].Update(db); err != nil {
				return
			}
		}
	}

	if len(l.Teachers) > 0 {
		for i := range l.Teachers {
			if _, err = l.Teachers[i].Update(db); err != nil {
				return
			}
		}
	}

	return
}

func (l *Language) UpdateColumns(db *sql.DB, columns []string) (r sql.Result, err error) {
	if l.ID == 0 {
		return l.Insert(db)
	}

	var args []interface{}
	var updates []string
	for _, c := range columns {
		updates = append(updates, c+" = ?")
		switch c {
		case LanguageColumnID:
			args = append(args, &l.ID)
		case LanguageColumnName:
			args = append(args, &l.Name)
		case LanguageColumnWordsCount:
			args = append(args, &l.WordsCount)
		default:
			err = fmt.Errorf("go2sql: unknown column %s", c)
			return
		}
	}

	r, err = db.Exec(`UPDATE languages SET `+strings.Join(updates, ","), args...)
	return
}

// TODO: deep delete
func (l *Language) Delete(db *sql.DB) (r sql.Result, err error) {
	r, err = db.Exec(`DELETE FROM languages WHERE id = ?`, l.ID)
	return
}

func (l *Language) LoadKeywords(db *sql.DB) (err error) {
	rows, err := db.Query(`select id, name, type, language_id from keywords where language_id = ?`, l.ID)
	if err != nil {
		return
	}

	defer func() {
		if er := rows.Close(); er != nil {
			if err != nil {
				log.Println(er)
			} else {
				err = er
			}
		}
	}()

	// l.Keywords = []Keyword{}
	for rows.Next() {
		var k Keyword
		if err = rows.Scan(&k.ID, &k.Name, &k.Type, &k.LanguageID); err != nil {
			return
		}
		l.Keywords = append(l.Keywords, k)
	}

	return
}

func LoadLanguagesKeywords(db *sql.DB, ls []*Language) (err error) {
	var ids []string
	for _, l := range ls {
		ids = append(ids, strconv.Itoa(int(l.ID)))
	}
	rows, err := db.Query(`SELECT id, name, type, language_id FROM keywords WHERE language_id IN (?)`, strings.Join(ids, ","))
	if err != nil {
		return
	}

	defer func() {
		if er := rows.Close(); er != nil {
			if err != nil {
				log.Println(er)
			} else {
				err = er
			}
		}
	}()

	var keywords []Keyword
	for rows.Next() {
		var k Keyword
		if err = rows.Scan(&k.ID, &k.Name, &k.Type, &k.LanguageID); err != nil {
			return
		}
		keywords = append(keywords, k)
	}

	for i, l := range ls {
		for _, k := range keywords {
			if k.LanguageID == l.ID {
				ls[i].Keywords = append(ls[i].Keywords, k)
			}
		}
	}

	return
}

func (l *Language) LoadAuthor(db *sql.DB) (err error) {
	err = db.QueryRow("select id, name, email from people where id = ?", l.AuthorID).Scan(&l.Author.ID, &l.Author.Name, &l.Author.Email)
	return
}

func LoadLanguagesAuthor(db *sql.DB, ls []*Language) (err error) {
	var ids []string
	for _, l := range ls {
		ids = append(ids, strconv.Itoa(int(l.AuthorID)))
	}
	rows, err := db.Query(`SELECT id, name, email FROM people WHERE id IN (?)`, strings.Join(ids, ","))
	if err != nil {
		return
	}

	defer func() {
		if er := rows.Close(); er != nil {
			if err != nil {
				log.Println(er)
			} else {
				err = er
			}
		}
	}()

	var ps []Person
	for rows.Next() {
		var p Person
		if err = rows.Scan(&p.ID, &p.Name, &p.Email); err != nil {
			return
		}
		ps = append(ps, p)
	}

	for i, l := range ls {
		for _, p := range ps {
			if p.ID == l.AuthorID {
				ls[i].Author = p
			}
		}
	}

	return
}

func (l *Language) LoadTeachers(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT teachers.id, teachers.name, teachers.age FROM teachers
		INNER JOIN languages_teachers_xref
		ON teachers.id = languages_teachers_xref.teacher_id
		WHERE languages_teachers_xref.language_id = ?`, l.ID)
	if err != nil {
		return
	}

	defer func() {
		if er := rows.Close(); er != nil {
			if err != nil {
				log.Println(er)
			} else {
				err = er
			}
		}
	}()

	for rows.Next() {
		var t Teacher
		if err = rows.Scan(&t.ID, &t.Name, &t.Age); err != nil {
			return
		}
		l.Teachers = append(l.Teachers, t)
	}

	return
}

func LoadLanguagesTeachers(db *sql.DB, ls []*Language) (err error) {
	var ids []string
	for _, l := range ls {
		ids = append(ids, strconv.Itoa(int(l.ID)))
	}
	rows, err := db.Query(`SELECT
			teachers.id,
			teachers.name,
			teachers.age,
			languages_teachers_xref.language_id
		FROM teachers
		INNER JOIN languages_teachers_xref
		ON teachers.id = languages_teachers_xref.teacher_id
		WHERE languages_teachers_xref.language_id in (?)`, strings.Join(ids, ","))
	if err != nil {
		return
	}

	defer func() {
		if er := rows.Close(); er != nil {
			if err != nil {
				log.Println(er)
			} else {
				err = er
			}
		}
	}()

	type teacher struct {
		Teacher
		languageID uint
	}

	var ts []teacher
	for rows.Next() {
		var t teacher
		if err = rows.Scan(&t.ID, &t.Name, &t.Age, &t.languageID); err != nil {
			return
		}
		ts = append(ts, t)
	}

	for i, l := range ls {
		for _, t := range ts {
			if t.languageID == l.ID {
				ls[i].Teachers = append(ls[i].Teachers, t.Teacher)
			}
		}
	}

	return
}
