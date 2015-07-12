package go2sql

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
		Name   string
		Tables Tables
	}

	SQL struct {
		SQL  string
		Args []interface{}
		Full bool
	}
)

// var (
// 	// QueryOptionDeep  QueryOption  = queryOption{}
// 	InsertOptionDeep InsertOption = insertOption{}
// 	DeleteOptionDeep DeleteOption = deleteOption{}
// 	UpdateOptionDeep UpdateOption = updateOption{}
// )

func (Selects) QueryOption() {}

// type insertOption struct{}
// type deleteOption struct{}
// type updateOption struct{}

// func (insertOption) InsertOption() {}
// func (deleteOption) DeleteOption() {}
// func (updateOption) UpdateOption() {}

func (Tables) InsertOption() {}
func (Tables) DeleteOption() {}
func (Tables) UpdateOption() {}
func (Tables) QueryOption()  {}

func (SQL) InsertOption() {}
func (SQL) DeleteOption() {}
func (SQL) UpdateOption() {}
func (SQL) QueryOption()  {}

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

func (opts UpdateOptions) GetUpdateSQL() (sql SQL, ok bool) {
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

func (opts QueryOptions) GetSelect() (sel Selects, ok bool) {
	for _, o := range opts {
		if sel, ok = o.(Selects); ok {
			break
		}
	}
	return
}
