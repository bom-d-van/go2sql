package model

import (
	"github.com/bom-d-van/go2sql/go2sql"

	"database/sql"
)

type Keywords []*Keyword

func (t *Keywords) IsEmptyRow() (is bool)                           { return }
func (t *Keywords) IsNewRow() (is bool)                             { return }
func (t *Keywords) Insert(optsx ...go2sql.InsertOption) (err error) { return }
func (t *Keywords) Update(optsx ...go2sql.UpdateOption) (err error) { return }
func (t *Keywords) Delete(optsx ...go2sql.DeleteOption) (err error) { return }

func (k *Keyword) Insert(optsx ...go2sql.InsertOption) (r sql.Result, err error) {
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

func FindKeywords(opts ...go2sql.QueryOption) (k []*Keyword, err error) {
	return
}

func FindKeyword(opts ...go2sql.QueryOption) (k *Keyword, err error) {
	// rows, err := db.Query(`select id, name, type, language_id from keywords where language_id = ?`, l.ID)
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

	// // l.Keywords = []Keyword{}
	// for rows.Next() {
	// 	var k Keyword
	// 	if err = rows.Scan(&k.ID, &k.Name, &k.Type, &k.LanguageID); err != nil {
	// 		return
	// 	}
	// 	l.Keywords = append(l.Keywords, k)
	// }
	return
}

func (k *Keyword) Update(db *sql.DB) (r sql.Result, err error) { return }
func (k *Keyword) Delete(db *sql.DB) (r sql.Result, err error) { return }

func (k *Keyword) IsEmptyRow() bool {
	empty := Keyword{}
	return *k == empty
}
