package go2sql

import "database/sql"

type (
	InsertOption interface {
		InsertOption()
	}
	DeleteOption interface {
		DeleteOption()
	}
	UpdateOption interface {
		UpdateOption()
	}
	QueryOption interface {
		QueryOption()
	}

	Selects []string

	Tables []Table
	Table  struct {
		Name    string
		Tables  Tables
		Columns []string // TODO
	}

	SQL struct {
		SQL  string
		Args []interface{}
		Full bool
	}

	sqldb struct{ *sql.DB }
)

// var (
// 	// QueryOptionDeep  QueryOption  = queryOption{}
// 	InsertOptionDeep InsertOption = insertOption{}
// 	DeleteOptionDeep DeleteOption = deleteOption{}
// 	UpdateOptionDeep UpdateOption = updateOption{}
// )

var (
	DefaultDB *sqldb
)

func (Selects) QueryOption()  {}
func (Selects) UpdateOption() {}

// type insertOption struct{}
// type deleteOption struct{}
// type updateOption struct{}

// func (insertOption) InsertOption() {}
// func (deleteOption) DeleteOption() {}
// func (updateOption) UpdateOption() {}

// func NewTables(name string, tables Tables) Tables {}
func (Tables) InsertOption() {}
func (Tables) DeleteOption() {}
func (Tables) UpdateOption() {}
func (Tables) QueryOption()  {}
func (ts Tables) Get(t string) (Table, bool) {
	for _, ti := range ts {
		if ti.Name == t {
			return ti, true
		}
	}

	return Table{}, false
}

func (SQL) InsertOption() {}
func (SQL) DeleteOption() {}
func (SQL) UpdateOption() {}
func (SQL) QueryOption()  {}

func SetDefaultDB(db *sql.DB) { DefaultDB = &sqldb{db} }
func DB(db *sql.DB) sqldb     { return sqldb{db} }
func (sqldb) InsertOption()   {}
func (sqldb) DeleteOption()   {}
func (sqldb) UpdateOption()   {}
func (sqldb) QueryOption()    {}

func NewSQL(sql string, args ...interface{}) SQL {
	return SQL{SQL: sql, Args: args}
}

func NewFullSQL(sql string, args ...interface{}) SQL {
	s := NewSQL(sql, args...)
	s.Full = true
	return s
}

type InsertOptions []InsertOption
type DeleteOptions []DeleteOption
type UpdateOptions []UpdateOption
type QueryOptions []QueryOption

func (opts InsertOptions) HasOption(opt InsertOption) bool {
	for _, o := range opts {
		if o == opt {
			return true
		}
	}
	return false
}

func (opts DeleteOptions) HasOption(opt DeleteOption) bool {
	for _, o := range opts {
		if o == opt {
			return true
		}
	}
	return false
}

func (opts UpdateOptions) HasOption(opt UpdateOption) bool {
	for _, o := range opts {
		if o == opt {
			return true
		}
	}
	return false
}

func (opts QueryOptions) HasOption(opt QueryOption) bool {
	for _, o := range opts {
		if o == opt {
			return true
		}
	}
	return false
}

func (opts InsertOptions) GetTables() (ts Tables, ok bool) {
	for _, o := range opts {
		if ts, ok = o.(Tables); ok {
			return
		}
	}
	return
}

func (opts DeleteOptions) GetTables() (ts Tables, ok bool) {
	for _, o := range opts {
		if ts, ok = o.(Tables); ok {
			return
		}
	}
	return
}

func (opts UpdateOptions) GetTables() (ts Tables, ok bool) {
	for _, o := range opts {
		if ts, ok = o.(Tables); ok {
			return
		}
	}
	return
}

func (opts UpdateOptions) GetUpdateTables() (ts Tables, ok bool) {
	for _, o := range opts {
		if ts, ok = o.(Tables); ok {
			return
		}
	}
	return
}

func (opts QueryOptions) GetTables() (ts Tables, ok bool) {
	for _, o := range opts {
		if ts, ok = o.(Tables); ok {
			return
		}
	}
	return
}

func (opts InsertOptions) GetSQL() (sql SQL, ok bool) {
	for _, o := range opts {
		if sql, ok = o.(SQL); ok {
			return
		}
	}
	return
}

func (opts DeleteOptions) GetSQL() (sql SQL, ok bool) {
	for _, o := range opts {
		if sql, ok = o.(SQL); ok {
			return
		}
	}
	return
}

func (opts UpdateOptions) GetSQL() (sql SQL, ok bool) {
	for _, o := range opts {
		if sql, ok = o.(SQL); ok {
			return
		}
	}
	return
}

func (opts QueryOptions) GetSQL() (sql SQL, ok bool) {
	for _, o := range opts {
		if sql, ok = o.(SQL); ok {
			return
		}
	}
	return
}

func (opts InsertOptions) GetDB() (*sql.DB, bool) {
	for _, o := range opts {
		if db, ok := o.(sqldb); ok {
			return db.DB, true
		}
	}
	return nil, false
}

func (opts DeleteOptions) GetDB() (*sql.DB, bool) {
	for _, o := range opts {
		if db, ok := o.(sqldb); ok {
			return db.DB, true
		}
	}
	return nil, false
}

func (opts UpdateOptions) GetDB() (*sql.DB, bool) {
	for _, o := range opts {
		if db, ok := o.(sqldb); ok {
			return db.DB, true
		}
	}
	return nil, false
}

func (opts QueryOptions) GetDB() (*sql.DB, bool) {
	for _, o := range opts {
		if db, ok := o.(sqldb); ok {
			return db.DB, true
		}
	}
	return nil, false
}

func (opts QueryOptions) GetSelect() (sel Selects, ok bool) {
	for _, o := range opts {
		if sel, ok = o.(Selects); ok {
			break
		}
	}
	return
}

func (opts UpdateOptions) GetSelect() (sel Selects, ok bool) {
	for _, o := range opts {
		if sel, ok = o.(Selects); ok {
			break
		}
	}
	return
}
