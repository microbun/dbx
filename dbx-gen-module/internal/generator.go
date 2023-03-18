package internal

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql" //justifying
	"github.com/microbun/dbx"
)

var types = map[string]string{
	"TINYINT":    "int8",
	"SMALLINT":   "int16",
	"MEDIUMINT":  "int32",
	"INTEGER":    "int32",
	"INT":        "int32",
	"BIGINT":     "int64",
	"DECIMAL":    "float64",
	"NUMERIC":    "float64",
	"FLOAT":      "float32",
	"DOUBLE":     "float64",
	"BIT":        "[]byte",
	"DATE":       "time.Time",
	"DATETIME":   "time.Time",
	"TIMESTAMP":  "time.Time",
	"CHAR":       "string",
	"VARCHAR":    "string",
	"BINARY":     "[]byte",
	"VARBINARY":  "[]byte",
	"BLOB":       "[]byte",
	"MEDIUMBLOB": "[]byte",
	"LONGTEXT":   "string",
	"MEDIUMTEXT": "string",
	"TEXT":       "string",
	"ENUM":       "string",
	"TINYINT(1)": "bool",

	// "SET":   "",
}

type Column struct {
	TableName     string `dbx:"column:table_name" `
	ColumnName    string `dbx:"column:column_name" `
	ColumnComment string `dbx:"column:column_comment" `
	DataType      string `dbx:"column:data_type"`
	Nullable      string `dbx:"column:is_nullable"`
	ColumnKey     string `dbx:"column:column_key"`
	ColumnType    string `dbx:"column:column_type"`
	Extra         string `dbx:"column:extra"`
}

type Table struct {
	StructName string
	TableName  string
	Columns    []Column
}

func (t Table) RecordName() string {
	return t.StructName + "Record"
}

func (t Table) ColumnString() string {
	var s []string
	for _, column := range t.Columns {
		s = append(s, "\""+column.ColumnName+"\"")
	}
	return strings.Join(s, ",")
}

func (c Column) Name() string {
	return toCamelInitCase(c.ColumnName, true)
}

func (c Column) Type() string {
	pointer := ""
	if c.Nullable == "YES" {
		pointer += "*"
	}
	typeName := types[strings.ToUpper(c.DataType)]
	if strings.ToUpper(c.ColumnType) == "TINYINT(1)" {
		typeName = "bool"
	}
	if typeName == "[]byte" {
		pointer = ""
	}
	return pointer + typeName
}

func (c Column) Tag() string {
	omitempty := ""
	dbxTag := "dbx:\"column:" + c.ColumnName
	if c.Nullable == "YES" {
		omitempty = ",omitempty"
	}
	if c.ColumnKey == "PRI" {
		dbxTag += ",primary_key"
	}
	if strings.Contains(c.Extra, "auto_increment") {
		dbxTag += ",auto_increment"
	}
	if c.ColumnName == "created_at" {
		dbxTag += ",insert:current_timestamp"
	}
	if c.ColumnName == "updated_at" {
		dbxTag += ",insert:current_timestamp,update:current_timestamp"
	}

	dbxTag += "\""
	return fmt.Sprintf("`%v json:\"%v%v\" `", dbxTag, c.ColumnName, omitempty)
}

//go:embed module.go.tmpl
var text string

func Run() error {
	return generate()
}

type Context struct {
	Tables  map[string]*Table
	Package string
}

func (c Context) Imports() []string {
	ps := map[string]interface{}{}
	for _, table := range c.Tables {
		for _, column := range table.Columns {
			if column.Type() == "time.Time" || column.Type() == "*time.Time" {
				ps["time"] = nil
			}
		}
	}
	var packages []string
	for k := range ps {
		packages = append(packages, k)
	}
	return packages
}

func generate() error {
	if strings.ToLower(Options.Driver) != "mysql" {
		return errors.New("unsupport database")
	}

	db, err := dbx.Open(Options.Driver, Options.DataSourceName)
	if err != nil {
		return err
	}
	db.Options().Logger = nil
	var columns []*Column
	err = db.Query(&columns, `
	select table_name,column_name,column_comment,data_type,is_nullable,column_key,extra,column_type
	from information_schema.columns
	where TABLE_SCHEMA=?
	`, Options.Schema)
	if err != nil {
		return err
	}

	tables := map[string]*Table{}

	for _, column := range columns {
		name := column.TableName
		if v, ok := tables[name]; ok {
			v.Columns = append(v.Columns, *column)
		} else {
			tables[name] = &Table{
				StructName: toCamelInitCase(column.TableName, true),
				TableName:  column.TableName,
				Columns:    []Column{*column},
			}
		}
	}

	tmpl := template.New("template")
	tmpl, err = tmpl.Parse(string(text))
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(make([]byte, 0))

	err = tmpl.Execute(buf, Context{
		Tables:  tables,
		Package: Options.Package,
	})
	if err != nil {
		return err
	}

	dir := path.Dir(Options.Output)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(Options.Output, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	_, err = out.Write(src)
	if err != nil {
		return err
	}
	return nil
}
