package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"strconv"
	"strings"
	"text/template"
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
	VarName     string
	ColVarName  string
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
	field        *types.Var
	Name         string
	SQLName      string
	IsPrimaryKey bool
	Relationship Relationship

	Type      string
	TableType string
	IsPointer bool

	flags []string

	IsTable   bool
	Table     *Table
	TypeTable *Table

	parser *Parser
}

func (c *Column) ExpIsZero() string {
	ftype := c.field.Type()

typeSwitch:
	switch typ := ftype.(type) {
	case *types.Named:
		ftype = typ.Underlying()
		goto typeSwitch
	case *types.Basic:
		info := typ.Info()
		if info&types.IsBoolean != 0 {
			return fmt.Sprintf("!%s.%s", c.Table.RefName, c.Name)
		} else if info&types.IsNumeric != 0 {
			log.Println(c.Name)
			// printutils.PrettyPrint(c)
			return fmt.Sprintf("%s.%s == 0", c.Table.RefName, c.Name)
		} else if info&types.IsString != 0 {
			return fmt.Sprintf(`%s.%s == ""`, c.Table.RefName, c.Name)
		}
	case *types.Pointer, *types.Slice:
		return fmt.Sprintf("%s.%s == nil", c.Table.RefName, c.Name)
	case *types.Struct:
		// TODO

	// case *types.Slice:
	// 	return fmt.Sprintf("%s.%s == nil", c.Table.RefName, c.Name)
	case *types.Array:
	}

	return ""
}

func (t *Table) ExpIsZero() string {
	var cs []string
	for _, c := range t.Columns {
		exp := c.ExpIsZero()
		if exp == "" {
			continue
		}
		cs = append(cs, exp)
	}
	return strings.Join(cs, " &&\n")
}

func (t *Table) ExpIsNewRow() string {
	var cs []string
	for _, c := range t.PrimaryKeys {
		exp := c.ExpIsZero()
		if exp == "" {
			continue
		}
		cs = append(cs, exp)
	}
	return strings.Join(cs, " &&\n")
}

func (t *Table) ExpPrimaryKeyValues() string {
	var exps []string
	for _, pk := range t.PrimaryKeys {
		exps = append(exps, t.RefName+"."+pk.Name)
	}
	return strings.Join(exps, ",")
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

func contains(flags []string, f string) bool {
	for _, fl := range flags {
		if fl == f {
			return true
		}
	}
	return false
}

// var fileHeader = template.Must(template.New("").Parse(``))

func (t *Table) NoTableColumns() (cs []*Column) {
	for _, c := range t.Columns {
		if c.IsTable {
			continue
		}
		cs = append(cs, c)
	}
	return
}

func (t *Table) TableColumns(typs ...string) (cs []*Column) {
	var typ string
	if len(typs) > 0 {
		typ = typs[0]
	}
	for _, c := range t.Columns {
		if !c.IsTable {
			continue
		}
		if typ == "has" {
			if c.Relationship == RelationshipNone || c.Relationship == RelationshipBelongsTo {
				continue
			}
		} else if typ == "belongs" {
			if c.Relationship != RelationshipBelongsTo {
				continue
			}
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
		case "*go":
			strs = append(strs, fmt.Sprintf("&%s.%s", t.RefName, c.Name))
		case "go":
			strs = append(strs, fmt.Sprintf("%s.%s", t.RefName, c.Name))
		}
	}
	return strings.Join(strs, ",")
}

func (c *Column) ExpIDValue() string {
	if types.AssignableTo(c.field.Type(), types.Typ[types.Uint]) {
		return "id"
	}

	return fmt.Sprintf("%s(id)", types.TypeString(c.field.Type(), c.parser.Qualifier))
}

// func (c *Column) TableVarName() string {
// 	if !c.IsTable || c.Relationship == RelationshipBelongsTo || c.Relationship == RelationshipHasOne {
// 		return camelCase(c.Name)
// 	}

// 	return inflect.Pluralize(camelCase(c.Name))
// }

func (t *Table) ExpSQLWhere() string {
	var exps []string
	for _, pk := range t.PrimaryKeys {
		exps = append(exps, pk.SQLName+"= ?")
	}
	if len(exps) == 1 {
		return exps[0]
	}

	return fmt.Sprintf("(%s)", strings.Join(exps, " and "))
}

func (c *Column) ExpMany2ManySQLColumns() string {
	var exps []string
	for _, pk := range c.Table.PrimaryKeys {
		exps = append(exps, c.Table.SQLName+"_"+pk.SQLName)
	}
	for _, pk := range c.TypeTable.PrimaryKeys {
		exps = append(exps, c.TypeTable.SQLName+"_"+pk.SQLName)
	}

	return strings.Join(exps, ", ")
}

func (c *Column) ExpMany2ManySQLValues() string {
	var exps []string
	for range c.Table.PrimaryKeys {
		exps = append(exps, "?")
	}
	for range c.TypeTable.PrimaryKeys {
		exps = append(exps, "?")
	}

	return strings.Join(exps, ", ")
}

func (c *Column) ExpMany2ManyFields() string {
	var exps []string
	for _, pk := range c.Table.PrimaryKeys {
		exps = append(exps, fmt.Sprintf("%s.%s%s", c.TypeTable.VarName, c.Table.Name, pk.SQLName))
	}
	for _, pk := range c.TypeTable.PrimaryKeys {
		exps = append(exps, fmt.Sprintf("%s.%s", c.TypeTable.VarName, pk.Name))
	}

	return strings.Join(exps, ", ")
}

var tmpl = template.Must(template.New("tmpl.go").Funcs(template.FuncMap{
	"const_relationship_none":         func() Relationship { return RelationshipNone },
	"const_relationship_belongs_to":   func() Relationship { return RelationshipBelongsTo },
	"const_relationship_has_one":      func() Relationship { return RelationshipHasOne },
	"const_relationship_has_many":     func() Relationship { return RelationshipHasMany },
	"const_relationship_many_to_many": func() Relationship { return RelationshipManyToMany },
}).Parse(rawTmpl))
