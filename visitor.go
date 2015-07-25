package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"log"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/go/types"

	"bitbucket.org/pkg/inflect"
)

const (
	FlagID     = "id"
	FlagPK     = "primary-key"
	FlagInline = "inline"
	FlagIgnore = "-"

	FlagPrefix = "prefix:"

	TableNameSuffix  = "TableName"
	Go2SQLFileSuffix = "go2sql"
)

var cliFlags = struct {
	debug bool
}{}

func init() {
	log.SetFlags(log.Lshortfile)
}

type Option struct {
	Functions []string
}

type Function struct {
	Name     string
	Template template.Template
}

type Table struct {
	Package string

	struc       *ast.StructType
	Name        string
	ColName     string
	SQLName     string
	RefName     string
	ColRefName  string
	IDColumn    *Column
	Columns     []*Column
	PrimaryKeys []*Column

	BelongsTo   []*Table
	HasOnes     []*Table
	HasManys    []*Table
	ManyToManys []*Table

	HasCustomSQLName bool

	w bytes.Buffer
}

type Column struct {
	f            *ast.Field
	Name         string
	SQLName      string
	IsPrimaryKey bool
	Relationship Relationship

	Type      string
	TableType string
	IsPointer bool

	flags []string

	IsTable bool
	Table   *Table
}

func (c *Column) ExpIsZero() string {
	var empty string
	switch c.Type {
	case "bool":
		empty = "false"
	case "uint8", "uint16", "uint32", "uint64",
		"int8", "int16", "int32", "int64",
		"float32", "float64",
		"complex64", "complex128",
		"int", "uint", "uintptr",
		"rune":
		empty = "0"
	case "string":
		empty = `""`
	case "byte":
		empty = "0"
	// case "rune":
	// 	empty = `''`
	case "error":
		empty = `nil`
	}
	return fmt.Sprintf("%s.%s == %s", c.Table.RefName, c.Name, empty)
}

type Relationship int

const (
	RelationshipNone Relationship = iota
	RelationshipBelongsTo
	RelationshipHasOne
	RelationshipHasMany
	RelationshipManyToMany
)

func (r Relationship) String() string {
	switch r {
	case RelationshipNone:
		return "none"
	case RelationshipBelongsTo:
		return "belongs-to"
	case RelationshipHasOne:
		return "has-one"
	case RelationshipHasMany:
		return "has-many"
	case RelationshipManyToMany:
		return "many-to-many"
	}

	return ""
}

func (t *Table) HasColumn(c string) bool {
	for _, cl := range t.Columns {
		if cl.Name == c {
			return true
		}
	}
	return false
}

func (t *Table) GetColumn(c string) *Column {
	for _, cl := range t.Columns {
		if cl.Name == c {
			return cl
		}
	}
	return nil
}

type visitor struct {
	typ string

	consts map[string]string
	tables map[string]*Table
}

func newVisitor() *visitor {
	var v visitor
	v.tables = map[string]*Table{}
	v.consts = map[string]string{}
	return &v
}

func (v *visitor) Visit(n ast.Node) (w ast.Visitor) {
	switch node := n.(type) {
	case *ast.Ident:
		v.typ = node.Name
	case *ast.ValueSpec:
		for i, name := range node.Names {
			// || !strings.HasSuffix(name.String(), "TableName")
			if name.Obj.Kind != ast.Con || len(node.Values) <= i {
				continue
			}
			bl, ok := node.Values[i].(*ast.BasicLit)
			if !ok {
				continue
			}

			v.consts[name.String()] = bl.Value
		}
	case *ast.StructType:
		var table Table
		table.struc = node
		table.Name = v.typ
		table.ColName = inflect.Pluralize(table.Name)
		table.RefName = strings.ToLower(v.typ[:1])
		table.ColRefName = inflect.Pluralize(table.RefName)
		table.SQLName = inflect.Pluralize(toSnake(table.Name)) + TableNameSuffix
	listLoop:
		for _, f := range node.Fields.List {
			var flags []string
			if f.Tag != nil {
				tag, err := strconv.Unquote(f.Tag.Value)
				if err != nil {
					log.Printf("failed to unquote tag in %s.%s: %s\n", v.typ, f.Names[0].Name, err)
				}
				flags = strings.Split(reflect.StructTag(tag).Get("go2sql"), ",")
			}
			for _, n := range f.Names {
				var column Column
				column.Name = n.Name
				column.f = f

				if len(flags) > 0 && flags[0] != "" {
					if flags[0] == FlagIgnore {
						continue listLoop
					}
					column.SQLName = flags[0]
				} else {
					column.SQLName = toSnake(column.Name)
				}
				column.flags = flags
				_, column.IsPointer = f.Type.(*ast.StarExpr)
				if contains(flags, FlagID) {
					table.IDColumn = &column
				}
				if contains(flags, FlagPK) {
					column.IsPrimaryKey = true
					table.PrimaryKeys = append(table.PrimaryKeys, &column)
				}

				column.Type = types.ExprString(f.Type)
				column.IsTable, column.TableType, column.Relationship = isTable(f.Type)
				if column.IsTable && contains(flags, FlagInline) {
					column.IsTable = false
				}
				table.Columns = append(table.Columns, &column)
			}
		}
		v.tables[table.Name] = &table
	}
	return v
}

func contains(flags []string, f string) bool {
	for _, fl := range flags {
		if fl == f {
			return true
		}
	}
	return false
}

func isTable(expr ast.Expr) (ok bool, table string, rel Relationship) {
	switch typ := expr.(type) {
	case *ast.Ident:
		// log.Println(typ.Name, typ.Obj)
		if typ.Obj == nil {
			return
		}
		table = typ.Name
		switch decl := typ.Obj.Decl.(type) {
		case *ast.TypeSpec:
			ok, _, rel = isTable(decl.Type)
			return
		case ast.Expr:
			return isTable(decl)
		case *ast.StructType:
			log.Printf("--> %+v\n", typ.Name)
			ok, rel = true, RelationshipHasOne
			return
		default:
			log.Printf("unknown %s", decl)
		}
	case *ast.StarExpr:
		if _, is := typ.X.(*ast.StarExpr); is {
			if cliFlags.debug {
				log.Printf("pointer of pointer is not supported: %s\n", types.ExprString(expr))
			}
			ok = false
			return
		}
		return isTable(typ.X)
	case *ast.SelectorExpr:
		// return isTable(typ.X)
	case *ast.ArrayType:
		rel = RelationshipHasMany
		ok, table, _ = isTable(typ.Elt)
	case *ast.StructType:
		ok = true
		rel = RelationshipHasOne
	case *ast.SliceExpr:
		rel = RelationshipHasMany
		ok, table, _ = isTable(typ.X)
		// case *ast.MapType:
		// case *ast.FuncType:
	}
	return
}

func (v *visitor) analyze() {
	for _, host := range v.tables {
		if name, ok := v.consts[host.Name+TableNameSuffix]; ok {
			host.HasCustomSQLName = true
			host.SQLName = name
		}

		for _, hostc := range host.Columns {
			if !hostc.IsTable {
				continue
			}
			guest := v.tables[hostc.TableType]
			if guest == nil {
				log.Printf("can't found struct %s\n", hostc.TableType)
				continue
			}

			hostc.Table = guest
			if hostc.Relationship == RelationshipHasOne {
				// reanalyze if it's a valid belongs-to
				belongsTo := true
				for _, pk := range guest.PrimaryKeys {
					// TODO: custome primary key naming
					belongsTo = belongsTo && host.HasColumn(hostc.Name+pk.Name)
				}
				if belongsTo {
					hostc.Relationship = RelationshipBelongsTo
					continue
				}

				// reanalyze if it's a valid has-one
				hasOne := true
				for _, pk := range host.PrimaryKeys {
					// TODO: custome primary key naming
					hasOne = hasOne && guest.HasColumn(host.Name+pk.Name)
				}
				if !hasOne {
					hostc.Relationship = RelationshipNone
				}
			} else if hostc.Relationship == RelationshipHasMany {
				// reanalyze if it's valid many-to-many
				if v.hasXrefTable(host, guest) {
					hostc.Relationship = RelationshipManyToMany // TODO: validate if pk keys are consistent
					continue
				}

				hasMany := true
				for _, pk := range host.PrimaryKeys {
					// TODO: custome primary key naming
					hasMany = hasMany && guest.HasColumn(host.Name+pk.Name)
				}
				if !hasMany {
					hostc.Relationship = RelationshipNone
				}
			}
		}
	}
}

func (v *visitor) hasXrefTable(host, guest *Table) bool {
	if _, ok := v.tables[host.Name+guest.Name+"Xref"]; ok {
		return true
	}
	if _, ok := v.tables[guest.Name+host.Name+"Xref"]; ok {
		return true
	}
	return false
}

var fileHeader = template.Must(template.New("").Parse(`package {{.Package}}

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bom-d-van/go2sql/go2sql"
)

{{$table := .}}

const (
	{{if .HasCustomSQLName}}{{.Name}}TableName = "{{.SQLName}}"{{end}}
	{{range .Columns}}
	{{$table.Name}}Column{{.Name}} = "{{.SQLName}}"{{end}}
)

type {{.ColName}} []*{{.Name}}`))

func (t *Table) NoTableColumns() (cs []*Column) {
	for _, c := range t.Columns {
		if c.IsTable {
			continue
		}
		cs = append(cs, c)
	}
	return
}

func (t *Table) TableColumns() (cs []*Column) {
	for _, c := range t.Columns {
		if !c.IsTable {
			continue
		}
		cs = append(cs, c)
	}
	return
}

func (t *Table) ColumnNamesString(cs []*Column, typ string) (s string) {
	var strs []string
	for _, c := range cs {
		switch typ {
		case "sql":
			strs = append(strs, strconv.Quote(c.SQLName))
		case "go":
			strs = append(strs, fmt.Sprintf("&%s.%s", t.RefName, c.Name))
		}
	}
	return strings.Join(strs, ",")
}

var findTmpl = template.Must(template.New("").Parse(`
{{$table := .}}

func Find{{.Name}}(db *sql.DB, opts ...go2sql.QueryOption) ({{.RefName}} *{{.Name}}, err error) {
	{{.RefName}} = &{{.Name}}{}
	var fields []interface{}
	var columns []string
	if sel, ok := go2sql.GetSelectQueryOption(opts); ok {
		columns = []string(sel)
		for _, c := range columns {
			switch c {
			{{range .NoTableColumns}}
			case {{$table.Name}}Column{{.Name}}:
				fields = append(fields, &{{$table.RefName}}.{{.Name}}){{end}}
			default:
				err = fmt.Errorf("go2sql: unknown column %s", c)
				return
			}
		}
	} else {
		columns = []string{{"{"}}{{.ColumnNamesString .NoTableColumns "sql"}}{{"}"}}
		fields = []interface{}{{"{"}}{{.ColumnNamesString .NoTableColumns "go"}}{{"}"}}
	}

	var query go2sql.SQLQuery
	if q, ok := go2sql.GetFullSQLQueryOption(opts); ok {
		query = go2sql.SQLQuery(q)
	} else {
		query.Query = "select " + strings.Join(columns, ",") + " from languages"
		if q, ok := go2sql.GetPartialSQLQueryOption(opts); ok {
			query.Query += " " + q.Query
			query.Args = q.Args
		}
	}

	err = db.QueryRow(query.Query, query.Args...).Scan(fields...)
	if err != nil {
		return
	}

	if fs, ok := go2sql.GetFetchsQueryOption(opts); ok {
		for _, f := range fs {
			switch f.Name {
			{{range .TableColumns}}
			case {{$table.Name}}Column{{.Name}}:
				err = {{$table.RefName}}.Fetch{{.Name}}(db, f.Fetchs){{end}}
			default:
				err = fmt.Errorf("unknown fetch column: %s", f.Name)
			}
			if err != nil {
				return
			}
		}
	}
	return
}`))

var findsTmpl = template.Must(template.New("").Parse(`
{{$table := .}}

func Find{{.ColName}}(db *sql.DB, opts ...go2sql.QueryOption) ({{.ColRefName}} {{.ColName}}, err error) {
	var columns []string
	if sel, ok := go2sql.GetSelectQueryOption(opts); ok {
		columns = []string(sel)
	} else {
		columns = []string{{"{"}}{{.ColumnNamesString .NoTableColumns "sql"}}{{"}"}}
	}

	var query go2sql.SQLQuery
	if q, ok := go2sql.GetFullSQLQueryOption(opts); ok {
		query = go2sql.SQLQuery(q)
	} else {
		query.Query = "select " + strings.Join(columns, ",") + " from languages"
		if q, ok := go2sql.GetPartialSQLQueryOption(opts); ok {
			query.Query += " " + q.Query
			query.Args = q.Args
		}
	}

	rows, err := db.Query(query.Query, query.Args...)
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

	return
}`))
