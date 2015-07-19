package model

import "github.com/bom-d-van/go2sql/go2sql"

type Teachers []*Teacher

func FindTeachers(opts ...go2sql.QueryOption) (ts Teachers, err error) {
	return
}

func (t *Teacher) IsEmptyRow() (is bool)                            { return }
func (t *Teacher) IsNewRow() (is bool)                              { return }
func (t *Teachers) Insert(optsx ...go2sql.InsertOption) (err error) { return }
func (t *Teachers) Update(optsx ...go2sql.UpdateOption) (err error) { return }
func (t *Teachers) Delete(optsx ...go2sql.DeleteOption) (err error) { return }

func (t *Teacher) Insert(optsx ...go2sql.InsertOption) (err error) {
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

func (t *Teacher) Update(optsx ...go2sql.UpdateOption) (err error) { return }
func (t *Teacher) Delete(optsx ...go2sql.DeleteOption) (err error) { return }

// func (t *Teacher) IsEmptyRow() bool {
// 	empty := Teacher{}
// 	return *t == empty
// }
