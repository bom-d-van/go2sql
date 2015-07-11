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
	PartialSQLQueryOption SQLQuery
	FullSQLQueryOption    SQLQuery
	SelectQueryOption     []string

	SQLQuery struct {
		Query string
		Args  []interface{}
	}
)

var (
	InsertOptionDeep InsertOption = insertOption{}

	DeleteOptionDeep DeleteOption = deleteOption{}

	UpdateOptionDeep UpdateOption = updateOption{}
)

func (PartialSQLQueryOption) QueryOption() {}
func (FullSQLQueryOption) QueryOption()    {}
func (SelectQueryOption) QueryOption()     {}

type insertOption struct{}
type deleteOption struct{}
type updateOption struct{}

func (insertOption) InsertOption() {}
func (deleteOption) DeleteOption() {}
func (updateOption) UpdateOption() {}

func HasInsertOption(opts []InsertOption, opt InsertOption) bool {
	for _, o := range opts {
		if o == opt {
			return true
		}
	}
	return false
}

func HasDeleteOption(opts []DeleteOption, opt DeleteOption) bool {
	for _, o := range opts {
		if o == opt {
			return true
		}
	}
	return false
}

func HasUpdateOption(opts []UpdateOption, opt UpdateOption) bool {
	for _, o := range opts {
		if o == opt {
			return true
		}
	}
	return false
}

func GetPartialSQLQueryOption(opts []QueryOption) (sql PartialSQLQueryOption, ok bool) {
	for _, o := range opts {
		if sql, ok = o.(PartialSQLQueryOption); ok {
			break
		}
	}
	return
}

func GetFullSQLQueryOption(opts []QueryOption) (sql FullSQLQueryOption, ok bool) {
	for _, o := range opts {
		if sql, ok = o.(FullSQLQueryOption); ok {
			break
		}
	}
	return
}

func GetSelectQueryOption(opts []QueryOption) (sel SelectQueryOption, ok bool) {
	for _, o := range opts {
		if sel, ok = o.(SelectQueryOption); ok {
			break
		}
	}
	return
}
