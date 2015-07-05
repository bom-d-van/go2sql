package model

import "database/sql"

func (k *Keyword) Insert(db *sql.DB) (r sql.Result, err error) {
	// re, err := db.Exec(`INSERT INTO keywords (name, type, language_id) VALUES(?, ?, ?)`, k.Name, k.Type, l.ID)
	// if err != nil {
	// 	return r, err
	// }

	// var id int64
	// if id, err = re.LastInsertId(); err != nil {
	// 	return r, err
	// }
	// l.Keywords[i].ID = uint(id)
	return
}

func (k *Keyword) Update(db *sql.DB) (r sql.Result, err error) {
	return
}

func (k *Keyword) IsEmptyRow() bool {
	empty := Keyword{}
	return *k == empty
}
