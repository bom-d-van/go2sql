package main

import (
	"go/ast"
	"go/types"
	"log"
	"reflect"
	"strconv"
	"strings"

	"bitbucket.org/pkg/inflect"
)

type Parser struct {
	Package string
	Consts  map[string]string
	Tables  map[string]*Table
}

func NewParser(pkg string) *Parser {
	return &Parser{
		Package: pkg,
		Consts:  make(map[string]string),
		Tables:  make(map[string]*Table),
	}
}

func (p *Parser) Parse(info types.Info) (err error) {
	for ident, obj := range info.Defs {
		if ident.Obj == nil {
			continue
		}
		switch ident.Obj.Kind {
		case ast.Con:
			p.Consts[obj.Name()], err = strconv.Unquote(info.Types[ident.Obj.Decl.(*ast.ValueSpec).Values[0]].Value.String())
			if err != nil {
				return
			}
		case ast.Typ:
			struc, ok := obj.Type().Underlying().(*types.Struct)
			if !ok {
				continue
			}

			var table Table
			// table.struc = node
			table.Name = ident.Name
			table.RefName = strings.ToLower(ident.Name[:1])
			table.VarName = strings.ToLower(ident.Name[:1]) + ident.Name[1:]
			table.ColName = inflect.Pluralize(table.Name)
			table.ColRefName = inflect.Pluralize(table.RefName)
			table.ColVarName = inflect.Pluralize(table.VarName)
			table.SQLName = inflect.Pluralize(toSnake(table.Name)) + TableNameSuffix

			for i := 0; i < struc.NumFields(); i++ {
				field := struc.Field(i)
				flags := strings.Split(reflect.StructTag(struc.Tag(i)).Get("go2sql"), ",")

				var column Column
				column.Name = field.Name()
				column.field = field
				column.flags = flags
				column.Table = &table
				column.parser = p

				if len(flags) > 0 && flags[0] != "" {
					if flags[0] == FlagIgnore {
						continue
					}
					column.SQLName = flags[0]
				} else {
					column.SQLName = toSnake(column.Name)
				}

				_, column.IsPointer = field.Type().(*types.Pointer)
				if contains(flags, FlagID) {
					table.IDColumn = &column
				}
				if contains(flags, FlagPK) {
					column.IsPrimaryKey = true
					table.PrimaryKeys = append(table.PrimaryKeys, &column)
				}

				column.Type = types.TypeString(field.Type(), p.Qualifier)
				column.IsTable, column.TableType, column.Relationship = p.IsTable(field.Type())
				if column.IsTable && contains(flags, FlagInline) {
					column.IsTable = false
				}
				table.Columns = append(table.Columns, &column)
			}
			p.Tables[table.Name] = &table
		}
	}

	for _, host := range p.Tables {
		if name, ok := p.Consts[host.Name+TableNameSuffix]; ok {
			host.HasCustomSQLName = true
			host.SQLName = name
		}

		for _, hostc := range host.Columns {
			if !hostc.IsTable {
				continue
			}
			guest := p.Tables[hostc.TableType]
			if guest == nil {
				log.Printf("can't found struct %s\n", hostc.TableType)
				continue
			}

			hostc.TypeTable = guest
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
				hasMany := true
				for _, pk := range host.PrimaryKeys {
					// TODO: custome primary key naming
					hasMany = hasMany && guest.HasColumn(host.Name+pk.Name)
				}
				if !hasMany {
					hostc.Relationship = RelationshipNone
					continue
				}

				many2Many := true
				for _, pk := range guest.PrimaryKeys {
					// TODO: custome primary key naming
					many2Many = many2Many && host.HasColumn(guest.Name+pk.Name)
				}
				if many2Many {
					hostc.Relationship = RelationshipManyToMany
				}
			}
		}
	}
	return
}

func (p *Parser) Qualifier(pkg *types.Package) string {
	if pkg.Name() == p.Package {
		return ""
	}

	return pkg.Name()
}

func (p *Parser) IsTable(typ types.Type) (bool, string, Relationship) {
	// log.Printf("--> %s %T\n", types.TypeString(typ, p.Qualifier), typ)
	switch utyp := typ.(type) {
	case *types.Named:
		is, table, rel := p.IsTable(utyp.Underlying())
		if _, ok := utyp.Underlying().(*types.Pointer); !ok {
			table = types.TypeString(utyp, p.Qualifier)
		}
		return is, table, rel
	case *types.Pointer:
		return p.IsTable(utyp.Elem())
	case *types.Struct:
		return true, types.TypeString(utyp, p.Qualifier), RelationshipHasOne
	case *types.Slice:
		is, table, _ := p.IsTable(utyp.Elem())
		return is, table, RelationshipHasMany
	case *types.Array:
		is, table, _ := p.IsTable(utyp.Elem())
		return is, table, RelationshipHasMany
	}

	return false, "", 0
}
