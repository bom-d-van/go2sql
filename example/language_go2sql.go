package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bom-d-van/go2sql/go2sql"
)

const (
	// LanguageTableName = "languages"

	LanguageColumnID         = "id"
	LanguageColumnName       = "name"
	LanguageColumnWordsCount = "words_stat"
	LanguageColumnAuthor     = "author"
	LanguageColumnKeywords   = "keywords"
	LanguageColumnTeachers   = "teachers"

	// TODO
)

var (
	LanguageAllRelatedTables = []string{LanguageColumnAuthor, LanguageColumnKeywords, LanguageColumnTeachers}
)

type Languages []*Language

func FirstLanguage(optsx ...go2sql.QueryOption) (l *Language, err error)  { return }
func FirstLanguages(optsx ...go2sql.QueryOption) (l Languages, err error) { return }
func LastLanguage(optsx ...go2sql.QueryOption) (l *Language, err error)   { return }
func LastLanguages(optsx ...go2sql.QueryOption) (l Languages, err error)  { return }

func FindLanguage(optsx ...go2sql.QueryOption) (l *Language, err error) {
	opts := go2sql.QueryOptions(optsx)
	db, ok := opts.GetDB()
	if !ok {
		db = go2sql.DefaultDB.DB
	}
	if db == nil {
		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
		return
	}
	l = &Language{}
	var fields []interface{}
	var columns []string
	if sel, ok := opts.GetSelect(); ok {
		columns = []string(sel)
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
	} else {
		columns = []string{"id", "name", "words_stat", "field1", "field2", "field3", "field4", "field5", "field6", "field7"}
		fields = []interface{}{&l.ID, &l.Name, &l.WordsCount, &l.Field1, &l.Field2, &l.Field3, &l.Field4, &l.Field5, &l.Field6, &l.Field7}
	}

	var sql go2sql.SQL
	if opt, ok := opts.GetSQL(); ok && opt.Full {
		sql = opt
	} else {
		sql.SQL = fmt.Sprintf("select %s from languages %s", strings.Join(columns, ","), opt.SQL)
		sql.Args = opt.Args
	}

	err = db.QueryRow(sql.SQL, sql.Args...).Scan(fields...)
	if err != nil {
		return
	}

	if tables, ok := opts.GetTables(); ok {
		for _, table := range tables {
			switch table.Name {
			case LanguageColumnAuthor:
				err = l.FetchAuthor(go2sql.DB(db), table.Tables)
			case LanguageColumnKeywords:
				err = l.FetchKeywords(go2sql.DB(db), table.Tables)
			case LanguageColumnTeachers:
				err = l.FetchTeachers(go2sql.DB(db), table.Tables)
			}
			if err != nil {
				return
			}
		}
	}
	return
}

func FindLanguages(optsx ...go2sql.QueryOption) (ls Languages, err error) {
	opts := go2sql.QueryOptions(optsx)
	db, ok := opts.GetDB()
	if !ok {
		db = go2sql.DefaultDB.DB
	}
	if db == nil {
		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
		return
	}
	var columns []string
	if sel, ok := opts.GetSelect(); ok {
		columns = []string(sel)
	} else {
		columns = []string{"id", "name", "words_stat", "field1", "field2", "field3", "field4", "field5", "field6", "field7"}
	}

	var sql go2sql.SQL
	if opt, ok := opts.GetSQL(); ok && opt.Full {
		sql = opt
	} else {
		sql.SQL = fmt.Sprintf("select %s from languages %s", strings.Join(columns, ","), opt.SQL)
		sql.Args = opt.Args
	}

	rows, err := db.Query(sql.SQL, sql.Args...)
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
			case "field1":
				fields = append(fields, &l.Field1)
			case "field2":
				fields = append(fields, &l.Field2)
			case "field3":
				fields = append(fields, &l.Field3)
			case "field4":
				fields = append(fields, &l.Field4)
			case "field5":
				fields = append(fields, &l.Field5)
			case "field6":
				fields = append(fields, &l.Field6)
			case "field7":
				fields = append(fields, &l.Field7)
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

	if tables, ok := opts.GetTables(); ok {
		for _, table := range tables {
			switch table.Name {
			case LanguageColumnAuthor:
				err = ls.FetchAuthor(go2sql.DB(db), table.Tables)
			case LanguageColumnKeywords:
				err = ls.FetchKeywords(go2sql.DB(db), table.Tables)
			case LanguageColumnTeachers:
				err = ls.FetchTeachers(go2sql.DB(db), table.Tables)
			}
			if err != nil {
				return
			}
		}
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

// At least one of primary keys is not zero value
func (l *Language) IsNewRow() bool {
	if l == nil {
		return true
	}
	return l.ID == 0
}

// func init() {
// 	var l Language
// 	var db *sql.DB
// 	FindLanguage(go2sql.DB(db), go2sql.NewSQL("limit 1 ordered by id dsc"))
// }

func (l *Language) Duplicate(optsx ...go2sql.InsertOption) (nl *Language, err error) {
	return
}

func (l *Language) ZeroPrimaryKeys() {
}

func (ls *Languages) Insert(optsx ...go2sql.InsertOption) (err error) {
	if len(*ls) == 0 {
		return
	}

	opts := go2sql.InsertOptions(optsx)
	db, ok := opts.GetDB()
	if !ok {
		db = go2sql.DefaultDB.DB
	}
	if db == nil {
		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
		return
	}

	tables, _ := opts.GetTables()

	for _, table := range tables {
		switch table.Name {
		case LanguageColumnAuthor:
			var authors People
			for _, l := range *ls {
				if l.Author != nil {
					authors = append(authors, l.Author)
				}
			}
			if err = authors.Insert(); err != nil {
				return
			}
		default:
			err = fmt.Errorf("go2sql: unknown table %s", table)
			return
		}
	}

	// TODO: support selects
	stmt, err := db.Prepare(`INSERT INTO languages (name, words_stat) VALUES(?, ?)`)
	if err != nil {
		return
	}
	defer func() {
		er := stmt.Close()
		if err == nil {
			err = er
		} else {
			log.Println(er)
		}
	}()
	for _, l := range *ls {
		var r sql.Result
		r, err = stmt.Exec(l.Name, l.WordsCount)
		if err != nil {
			return
		}
		var id int64
		id, err = r.LastInsertId()
		if err != nil {
			return
		}
		l.ID = uint(id)
	}

	for _, table := range tables {
		switch table.Name {
		case LanguageColumnKeywords:
			var keywords Keywords
			for _, l := range *ls {
				for _, keyword := range l.Keywords {
					keyword.LanguageID = l.ID
					keywords = append(keywords, keyword)
				}
			}
			if err = keywords.Update(go2sql.DB(db), table.Tables); err != nil {
				return
			}
		case LanguageColumnTeachers:
			var teachers Teachers
			for _, l := range *ls {
				for _, teacher := range l.Teachers {
					teacher.LanguageID = l.ID
					teachers = append(teachers, teacher)
				}
			}
			if err = teachers.Update(go2sql.DB(db), table.Tables); err != nil {
				return
			}
			for _, teacher := range teachers {
				if _, err = db.Exec("INSERT INTO languages_teachers_xref (language_id, teacher_id) VALUES (?, ?)", teacher.LanguageID, teacher.ID); err != nil {
					return
				}
			}
			// teachers := Teachers(l.Teachers)
			// if err = teachers.Update(go2sql.DB(db), table.Tables); err != nil {
			// 	return
			// }
			// for index := range l.Teachers {
			// 	if _, err = db.Exec("INSERT INTO languages_teachers_xref (language_id, teacher_id) VALUES (?, ?)", l.ID, l.Teachers[index].ID); err != nil {
			// 		return
			// 	}
			// }
		default:
			err = fmt.Errorf("go2sql: unknown column %s", table)
			return
		}
	}

	return
}

func (l *Language) Insert(optsx ...go2sql.InsertOption) (err error) {
	if !l.IsNewRow() {
		return
	}

	opts := go2sql.InsertOptions(optsx)
	db, ok := opts.GetDB()
	if !ok {
		db = go2sql.DefaultDB.DB
	}
	if db == nil {
		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
		return
	}

	tables, _ := opts.GetTables()

	for _, table := range tables {
		switch table.Name {
		case LanguageColumnAuthor:
			if l.Author.IsEmptyRow() {
				continue
			}
			if err = l.Author.Insert(go2sql.DB(db), table.Tables); err != nil {
				return
			}
			l.AuthorID = l.Author.ID
		default:
			err = fmt.Errorf("go2sql: unknown column %s", table)
			return
		}
	}

	// TODO: support selects
	r, err := db.Exec(`INSERT INTO languages (name, words_stat) VALUES(?, ?)`, l.Name, l.WordsCount)
	if err != nil {
		return
	}
	id, err := r.LastInsertId()
	if err != nil {
		return
	}
	l.ID = uint(id)

	for _, table := range tables {
		switch table.Name {
		case LanguageColumnKeywords:
			for index := range l.Keywords {
				l.Keywords[index].LanguageID = l.ID
				// if err := l.Keywords[index].Update(db); err != nil {
				// 	return
				// }
			}
			keywords := Keywords(l.Keywords)
			if err = keywords.Update(go2sql.DB(db), table.Tables); err != nil {
				return
			}
		case LanguageColumnTeachers:
			teachers := Teachers(l.Teachers)
			if err = teachers.Update(go2sql.DB(db), table.Tables); err != nil {
				return
			}
			for index := range l.Teachers {
				// if l.Teachers[index].ID <= 0 {
				// 	if _, err = l.Teachers[index].Insert(db); err != nil {
				// 		return
				// 	}
				// }
				if _, err = db.Exec("INSERT INTO languages_teachers_xref (language_id, teacher_id) VALUES (?, ?)", l.ID, l.Teachers[index].ID); err != nil {
					return
				}
			}
		default:
			err = fmt.Errorf("go2sql: unknown column %s", table)
			return
		}
	}

	return
}

func (l *Language) Update(optsx ...go2sql.UpdateOption) (err error) {
	if l == nil {
		return
	}

	opts := go2sql.UpdateOptions(optsx)
	db, ok := opts.GetDB()
	if !ok {
		db = go2sql.DefaultDB.DB
	}
	if db == nil {
		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
		return
	}

	tables, _ := opts.GetTables()

	for _, table := range tables {
		switch table.Name {
		case LanguageColumnAuthor:
			if l.Author.IsNewRow() {
				continue
			}
			if err = l.Author.Update(go2sql.DB(db), table.Tables); err != nil {
				return
			}
			l.AuthorID = l.Author.ID
		default:
			err = fmt.Errorf("go2sql: unknown column %s", table)
			return
		}
	}

	if l.IsNewRow() {
		err = l.Insert(go2sql.DB(db))
	} else {
		_, err = db.Exec(`UPDATE languages SET name = ?, words_stat = ? WHERE id = ?`, l.Name, l.WordsCount, l.ID)
	}
	if err != nil {
		return
	}

	for _, table := range tables {
		switch table.Name {
		case LanguageColumnKeywords:
			for index := range l.Keywords {
				l.Keywords[index].LanguageID = l.ID
			}
			keywords := Keywords(l.Keywords)
			if err = keywords.Update(go2sql.DB(db), table.Tables); err != nil {
				return
			}
		case LanguageColumnTeachers:
			teachers := Teachers(l.Teachers)
			if err = teachers.Update(go2sql.DB(db), table.Tables); err != nil {
				return
			}
		default:
			err = fmt.Errorf("go2sql: unknown column %s", table)
			return
		}
	}

	return
}

func (ls *Languages) Update(optsx ...go2sql.UpdateOption) (err error) {
	if l == nil {
		return
	}

	opts := go2sql.UpdateOptions(optsx)
	db, ok := opts.GetDB()
	if !ok {
		db = go2sql.DefaultDB.DB
	}
	if db == nil {
		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
		return
	}

	tables, _ := opts.GetTables()

	for _, table := range tables {
		switch table.Name {
		case LanguageColumnAuthor:
			if l.Author.IsNewRow() {
				continue
			}
			if err = l.Author.Update(go2sql.DB(db), table.Tables); err != nil {
				return
			}
			l.AuthorID = l.Author.ID
		default:
			err = fmt.Errorf("go2sql: unknown column %s", table)
			return
		}
	}

	if l.IsNewRow() {
		err = l.Insert(go2sql.DB(db))
	} else {
		_, err = db.Exec(`UPDATE languages SET name = ?, words_stat = ? WHERE id = ?`, l.Name, l.WordsCount, l.ID)
	}
	if err != nil {
		return
	}

	for _, table := range tables {
		switch table.Name {
		case LanguageColumnKeywords:
			for index := range l.Keywords {
				l.Keywords[index].LanguageID = l.ID
			}
			keywords := Keywords(l.Keywords)
			if err = keywords.Update(go2sql.DB(db), table.Tables); err != nil {
				return
			}
		case LanguageColumnTeachers:
			teachers := Teachers(l.Teachers)
			if err = teachers.Update(go2sql.DB(db), table.Tables); err != nil {
				return
			}
		default:
			err = fmt.Errorf("go2sql: unknown column %s", table)
			return
		}
	}

	return
}

func (l *Language) UpdateColumns(optsx ...go2sql.UpdateOption) (err error) {
	opts := go2sql.UpdateOptions(optsx)
	db, ok := opts.GetDB()
	if !ok {
		db = go2sql.DefaultDB.DB
	}
	if db == nil {
		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
		return
	}

	if l.IsNewRow() {
		// return l.Insert(db)
		return errors.New("can't not update a new row")
	}

	var columns []string
	if sel, ok := opts.GetSelect(); ok {
		columns = []string(sel)
	} else {
		return errors.New("empty select/columns")
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

	_, err = db.Exec(fmt.Sprintf(`UPDATE languages SET %s where id = ?`, strings.Join(updates, ","), l.ID), args...)
	return
}

func (l *Language) Delete(optsx ...go2sql.DeleteOption) (err error) {
	if l == nil || l.IsEmptyRow() {
		return
	}

	opts := go2sql.DeleteOptions(optsx)
	db, ok := opts.GetDB()
	if !ok {
		db = go2sql.DefaultDB.DB
	}
	if db == nil {
		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
		return
	}

	_, err = db.Exec(`DELETE FROM languages WHERE id = ?`, l.ID)
	if err != nil {
		return
	}
	tables, _ := opts.GetTables()
	for _, table := range tables {
		switch table.Name {
		case LanguageColumnAuthor:
			err = l.Author.Delete(go2sql.DB(db), table.Tables)
		case LanguageColumnKeywords:
			keywords := Keywords(l.Keywords)
			err = keywords.Delete(go2sql.DB(db), table.Tables)
		case LanguageColumnTeachers:
			teachers := Teachers(l.Teachers)
			err = teachers.Delete(go2sql.DB(db), table.Tables)
		default:
			err = fmt.Errorf("go2sql: unknown column %s", table)
		}
		if err != nil {
			return
		}
	}

	return
}

func (ls *Languages) Delete(optsx ...go2sql.DeleteOption) (err error) {
	opts := go2sql.DeleteOptions(optsx)
	db, ok := opts.GetDB()
	if !ok {
		db = go2sql.DefaultDB.DB
	}
	if db == nil {
		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
		return
	}

	var ids []interface{}
	var idsCons []string
	for _, l := range *ls {
		ids = append(ids, l.ID)
		idsCons = append(idsCons, "id = ?")
	}
	_, err = db.Exec(`DELETE FROM languages WHERE `+strings.Join(idsCons, " or "), ids...)
	if err != nil {
		return
	}
	tables, _ := opts.GetTables()
	for _, table := range tables {
		switch table.Name {
		case LanguageColumnAuthor:
			var people People
			for _, l := range *ls {
				people = append(people, l.Author)
			}
			err = people.Delete(go2sql.DB(db))
		case LanguageColumnKeywords:
			var keywords Keywords
			for _, l := range *ls {
				keywords = append(keywords, l.Keywords...)
			}
			err = keywords.Delete(go2sql.DB(db))
		case LanguageColumnTeachers:
			var teachers Teachers
			for _, l := range *ls {
				teachers = append(teachers, l.Teachers...)
			}
			err = teachers.Delete(go2sql.DB(db))
		default:
			err = fmt.Errorf("go2sql: unknown column %s", table)
		}
		if err != nil {
			return
		}
	}

	return
}

func (l *Language) FetchKeywords(opts ...go2sql.QueryOption) (err error) {
	opts = append(opts, go2sql.NewSQL("where language_id = ?", l.ID))
	l.Keywords, err = FindKeywords(opts...)

	return
}

func (ls *Languages) FetchKeywords(optsx ...go2sql.QueryOption) (err error) {
	opts := go2sql.QueryOptions(optsx)
	db, ok := opts.GetDB()
	if !ok {
		db = go2sql.DefaultDB.DB
	}
	if db == nil {
		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
		return
	}

	var ids []string
	for _, l := range *ls {
		ids = append(ids, strconv.Itoa(int(l.ID)))
	}
	// rows, err := db.Query(`SELECT id, name, type, language_id FROM keywords WHERE language_id IN (?)`, strings.Join(ids, ","))
	// if err != nil {
	// 	return
	// }

	// defer func() {
	// 	if er := rows.Close(); er != nil {
	// 		if err != nil {
	// 			log.Println(er)
	// 		} else {
	// 			err = er
	// 		}
	// 	}
	// }()

	// var keywords []Keyword
	// for rows.Next() {
	// 	var k Keyword
	// 	if err = rows.Scan(&k.ID, &k.Name, &k.Type, &k.LanguageID); err != nil {
	// 		return
	// 	}
	// 	keywords = append(keywords, k)
	// }
	keywords, err := FindKeywords(go2sql.DB(db), go2sql.NewSQL("WHERE language_id IN (?)", strings.Join(ids, ",")))
	if err != nil {
		return
	}

	for i, l := range *ls {
		for _, keyword := range keywords {
			if keyword.LanguageID == l.ID {
				(*ls)[i].Keywords = append((*ls)[i].Keywords, keyword)
			}
		}
	}

	return
}

func (l *Language) FetchAuthor(opts ...go2sql.QueryOption) (err error) {
	// err = db.QueryRow("select id, name, email from people where id = ?", l.AuthorID).Scan(&l.Author.ID, &l.Author.Name, &l.Author.Email)
	opts = append(opts, go2sql.NewSQL("where language_id = ?", l.ID))
	l.Author, err = FindPerson(opts...)
	return
}

// func (ls *Languages) FetchAuthor(opts ...go2sql.QueryOption) (err error) {
// 	// err = db.QueryRow("select id, name, email from people where id = ?", l.AuthorID).Scan(&l.Author.ID, &l.Author.Name, &l.Author.Email)
// 	var ids []string
// 	for _, l := range ls {
// 		ids = append(ids, strconv.Itoa(l.ID))
// 	}
// 	// TODO: composite primary keys support
// 	opts = append(opts, go2sql.NewSQL("where language_id in (?)", strings.Join(ids, ",")))
// 	l.Author, err = FindPersons(opts...)
// 	return
// }

func (ls *Languages) FetchAuthor(optsx ...go2sql.QueryOption) (err error) {
	if len(*ls) == 0 {
		return
	}
	opts := go2sql.QueryOptions(optsx)
	db, ok := opts.GetDB()
	if !ok {
		db = go2sql.DefaultDB.DB
	}
	if db == nil {
		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
		return
	}

	var ids []string
	for _, l := range *ls {
		ids = append(ids, strconv.Itoa(int(l.AuthorID)))
	}
	// rows, err := db.Query(`SELECT id, name, email FROM people WHERE id IN (?)`, strings.Join(ids, ","))
	// if err != nil {
	// 	return
	// }

	// defer func() {
	// 	if er := rows.Close(); er != nil {
	// 		if err != nil {
	// 			log.Println(er)
	// 		} else {
	// 			err = er
	// 		}
	// 	}
	// }()

	// var people []*Person
	// for rows.Next() {
	// 	var person Person
	// 	if err = rows.Scan(&person.ID, &person.Name, &person.Email); err != nil {
	// 		return
	// 	}
	// 	people = append(people, &person)
	// }
	people, err := FindPeople(go2sql.DB(db), go2sql.NewSQL("WHERE id IN (?)", strings.Join(ids, ", ")))
	if err != nil {
		return
	}

	for index, l := range *ls {
		for _, person := range people {
			if person.ID == l.AuthorID {
				(*ls)[index].Author = person
			}
		}
	}

	return
}

func (l *Language) FetchTeachers(opts ...go2sql.QueryOption) (err error) {
	// rows, err := db.Query(`SELECT teachers.id, teachers.name, teachers.age FROM teachers
	// 	INNER JOIN languages_teachers_xref
	// 	ON teachers.id = languages_teachers_xref.teacher_id
	// 	WHERE languages_teachers_xref.language_id = ?`, l.ID)
	// if err != nil {
	// 	return
	// }

	// defer func() {
	// 	if er := rows.Close(); er != nil {
	// 		if err != nil {
	// 			log.Println(er)
	// 		} else {
	// 			err = er
	// 		}
	// 	}
	// }()

	// for rows.Next() {
	// 	var t Teacher
	// 	if err = rows.Scan(&t.ID, &t.Name, &t.Age); err != nil {
	// 		return
	// 	}
	// 	l.Teachers = append(l.Teachers, t)
	// }

	opts = append(opts, go2sql.NewSQL("where language_id = ?", l.ID))
	teachers, err := FindTeachers(opts...)
	l.Teachers = []*Teacher(teachers)

	return
}

func (ls *Languages) FetchTeachers(optsx ...go2sql.QueryOption) (err error) {
	if len(*ls) == 0 {
		return
	}
	opts := go2sql.QueryOptions(optsx)
	// db, ok := opts.GetDB()
	// if !ok {
	// 	if go2sql.DefaultDB == nil {
	// 		err = errors.New("should specify *sql.DB by go2sql.DB or init go2sql.DefaultDB")
	// 		return
	// 	}
	// 	db = go2sql.DefaultDB.DB
	// }

	var ids []string
	for _, l := range *ls {
		ids = append(ids, strconv.Itoa(int(l.ID)))
	}
	// rows, err := db.Query(`SELECT
	// 		teachers.id,
	// 		teachers.name,
	// 		teachers.age,
	// 		languages_teachers_xref.language_id
	// 	FROM teachers
	// 	INNER JOIN languages_teachers_xref
	// 	ON teachers.id = languages_teachers_xref.teacher_id
	// 	WHERE languages_teachers_xref.language_id in (?)`, strings.Join(ids, ","))
	// if err != nil {
	// 	return
	// }

	// defer func() {
	// 	if er := rows.Close(); er != nil {
	// 		if err != nil {
	// 			log.Println(er)
	// 		} else {
	// 			err = er
	// 		}
	// 	}
	// }()

	// type teacher struct {
	// 	Teacher
	// 	languageID uint
	// }

	// var ts []teacher
	// for rows.Next() {
	// 	var t teacher
	// 	if err = rows.Scan(&t.ID, &t.Name, &t.Age, &t.languageID); err != nil {
	// 		return
	// 	}
	// 	ts = append(ts, t)
	// }

	// // TODO: support where override?
	// optsx = append(optsx, go2sql.NewSQL(`SELECT
	// 		teachers.id,
	// 		teachers.name,
	// 		teachers.age,
	// 	FROM teachers
	// 	INNER JOIN languages_teachers_xref
	// 	ON teachers.id = languages_teachers_xref.teacher_id
	// 	WHERE languages_teachers_xref.language_id in (?)`, strings.Join(ids, ",")))
	// teachers, err := FindTeachers(optsx...)
	// if err != nil {
	// 	return
	// }

	opts = append(opts, go2sql.NewSQL(`INNER JOIN languages_teachers_xref
	ON teachers.id = languages_teachers_xref.teacher_id
	WHERE languages_teachers_xref.language_id in (?)`, strings.Join(ids, ",")))
	teachers, err := FindTeachers(opts...)
	if err != nil {
		return
	}

	for i, l := range *ls {
		for _, t := range teachers {
			if t.LanguageID == l.ID {
				(*ls)[i].Teachers = append((*ls)[i].Teachers, t)
			}
		}
	}

	return
}
