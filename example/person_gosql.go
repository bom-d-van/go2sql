package model

import "github.com/bom-d-van/go2sql/go2sql"

type People []*Person

func (t *Person) IsEmptyRow() (is bool)                           { return }
func (t *Person) IsNewRow() (is bool)                             { return }
func (t *People) Insert(optsx ...go2sql.InsertOption) (err error) { return }
func (t *People) Update(optsx ...go2sql.UpdateOption) (err error) { return }
func (t *People) Delete(optsx ...go2sql.DeleteOption) (err error) { return }

func FindPeople(optsx ...go2sql.QueryOption) (ps People, err error) { return }
func FindPerson(optsx ...go2sql.QueryOption) (p *Person, err error) { return }

func (p *Person) Insert(optsx ...go2sql.InsertOption) (err error) {
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

func (p *Person) Update(optsx ...go2sql.UpdateOption) (err error) { return }
func (p *Person) Delete(optsx ...go2sql.DeleteOption) (err error) { return }

// func (l *Person) IsEmptyRow() bool {
// 	empty := Person{}
// 	return *l == empty
// }
