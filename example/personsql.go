package model

import "database/sql"

func (p *Person) Insert(db *sql.DB) (r sql.Result, err error) {
	// r, err := db.Exec(`INSERT INTO people (name, email) VALUES(?, ?)`, l.Author.Name, l.Author.Email)
	// if err != nil {
	// 	return nil, err
	// }

	// var id int64
	// if id, err = r.LastInsertId(); err != nil {
	// 	return r, err
	// }
	// l.Author.ID = uint(id)
	return
}

func (p *Person) Update(db *sql.DB) (r sql.Result, err error) { return }
func (p *Person) Delete(db *sql.DB) (r sql.Result, err error) { return }

func (l *Person) IsEmptyRow() bool {
	empty := Person{}
	return *l == empty
}
