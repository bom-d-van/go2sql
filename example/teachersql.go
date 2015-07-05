package model

import "database/sql"

func (t *Teacher) Insert(db *sql.DB) (r sql.Result, err error) {
	// re, err := db.Exec(`INSERT INTO teachers (name, age) VALUES(?, ?)`, t.Name, t.Age)
	// if err != nil {
	// 	return r, err
	// }

	// var id int64
	// if id, err = re.LastInsertId(); err != nil {
	// 	return r, err
	// }
	// l.Teachers[i].ID = uint(id)
	return
}

func (t *Teacher) Update(db *sql.DB) (r sql.Result, err error) {
	return
}

func (t *Teacher) IsEmptyRow() bool {
	empty := Teacher{}
	return *t == empty
}
