package dbx

import (
	"errors"
	"fmt"
	"github.com/microbun/dbx/reflectx"
	"reflect"
	"sort"
)

type SQLGenerator interface {
	UpdateSQL(value interface{}) (query string, args []interface{}, err error)
	InsertSQL(value interface{}) (autoIncrement *reflect.Value, query string, args []interface{}, err error)
}

//var tableInterfaceType = reflect.TypeOf(Table).Elem()

//ReflectTable
func ReflectTable(value interface{}) (tableName string, props reflectx.Properties, err error) {
	v := reflect.ValueOf(value)
	direct := reflect.Indirect(v)
	if !reflectx.IsStructValue(direct) {
		return "", nil, fmt.Errorf("value not a struct")
	}
	if tv, ok := value.(Table); ok {
		tableName = tv.TableName()
	} else {
		return "", nil, errors.New("not implement generator.Table interface")
	}

	propertiesMap := map[string]reflectx.Property{}
	reflectx.ReflectProperty(direct, propertiesMap)
	size := len(propertiesMap)
	props = reflectx.NewProperties(size)
	i := 0
	for _, prop := range propertiesMap {
		props[i] = prop
		i++
	}

	sort.Sort(props)
	return tableName, props, nil
}

type CommonSQLGenerator struct {
	AutoUpdated bool
}

func NewCommonSQLGenerator() *CommonSQLGenerator {
	return &CommonSQLGenerator{}
}

func (CommonSQLGenerator) UpdateSQL(value interface{}) (query string, args []interface{}, err error) {
	table, props, err := ReflectTable(value)
	if err != nil {
		return "", nil, err
	}
	n := len(props)
	if n <= 0 {
		return "", nil, fmt.Errorf("not found update columns")
	}
	var values []interface{}
	columns := ""
	primaryKey := ""
	var primaryKeyValue interface{}
	for _, prop := range props {
		if prop.Tag.PrimaryKey {
			primaryKey = prop.Tag.Column + "=?"
			primaryKeyValue = prop.InterValue
		} else {
			columns += ", " + prop.Tag.Column + "=?"
			values = append(values, prop.InterValue)
		}
	}
	sql := fmt.Sprintf("update %s set %s where %s ", table, columns[2:], primaryKey)
	values = append(values, primaryKeyValue)
	return sql, values, nil
}

func (CommonSQLGenerator) InsertSQL(value interface{}) (autoIncrement *reflect.Value, query string, args []interface{}, err error) {
	table, props, err := ReflectTable(value)
	if err != nil {
		return nil, "", nil, err
	}
	n := len(props)
	if n <= 0 {
		return nil, "", nil, fmt.Errorf("not found insert columns")
	}
	var values []interface{}
	columns := ""
	strArg := ""
	for _, prop := range props {
		if prop.Tag.AutoIncrement {
			autoIncrement = prop.Value
		}
		if !prop.Tag.PrimaryKey {
			columns += ", " + prop.Tag.Column
			strArg += ",?"
			values = append(values, prop.InterValue)
		}
	}
	query = fmt.Sprintf("insert into %s(%s) values(%s)", table, columns[2:], strArg[1:])
	return autoIncrement, query, values, err
}

type SQLiteGenerator struct {
	AutoUpdated bool
}

func NewSQLiteSQLiteGenerator() *SQLiteGenerator {
	return &SQLiteGenerator{AutoUpdated: true}
}

func (SQLiteGenerator) UpdateSQL(value interface{}) (query string, args []interface{}, err error) {
	table, props, err := ReflectTable(value)
	if err != nil {
		return "", nil, err
	}
	n := len(props)
	if n <= 0 {
		return "", nil, fmt.Errorf("not found update columns")
	}
	var values []interface{}
	columns := ""
	primaryKey := ""
	var primaryKeyValue interface{}
	for _, prop := range props {
		if prop.Tag.PrimaryKey {
			primaryKey = prop.Tag.Column + "=?"
			primaryKeyValue = prop.InterValue
		} else {
			if prop.Tag.Column == "updated_at" {
				columns += ", " + prop.Tag.Column + "=current_timestamp"
			} else if prop.Tag.Column == "created_at" {
				continue
			} else {
				columns += ", " + prop.Tag.Column + "=?"
				values = append(values, prop.InterValue)
			}
		}
	}
	sql := fmt.Sprintf("update %s set %s where %s ", table, columns[2:], primaryKey)
	values = append(values, primaryKeyValue)
	return sql, values, nil
}

func (SQLiteGenerator) InsertSQL(value interface{}) (autoIncrement *reflect.Value, query string, args []interface{}, err error) {
	table, props, err := ReflectTable(value)
	if err != nil {
		return nil, "", nil, err
	}
	n := len(props)
	if n <= 0 {
		return nil, "", nil, fmt.Errorf("not found insert columns")
	}
	var values []interface{}
	columns := ""
	strArg := ""
	for _, prop := range props {
		if prop.Tag.AutoIncrement {
			autoIncrement = prop.Value
		}
		if !prop.Tag.PrimaryKey {
			columns += ", " + prop.Tag.Column

			if prop.Tag.Column == "updated_at" || prop.Tag.Column == "created_at" {
				strArg += ",current_timestamp"
			} else {
				strArg += ",?"
				values = append(values, prop.InterValue)
			}
		}
	}
	query = fmt.Sprintf("insert into %s(%s) values(%s)", table, columns[2:], strArg[1:])
	return autoIncrement, query, values, err
}

type MySQLGenerator struct {
	AutoUpdated bool
}

func NewMySQLSQLiteGenerator() *MySQLGenerator {
	return &MySQLGenerator{AutoUpdated: true}
}

func (MySQLGenerator) UpdateSQL(value interface{}) (query string, args []interface{}, err error) {
	table, props, err := ReflectTable(value)
	if err != nil {
		return "", nil, err
	}
	n := len(props)
	if n <= 0 {
		return "", nil, fmt.Errorf("not found update columns")
	}
	var values []interface{}
	columns := ""
	primaryKey := ""
	var primaryKeyValue interface{}
	for _, prop := range props {
		if prop.Tag.PrimaryKey {
			primaryKey = prop.Tag.Column + "=?"
			primaryKeyValue = prop.InterValue
		} else {
			if prop.Tag.Column == "updated_at" {
				columns += ", " + prop.Tag.Column + "=now()"
			} else if prop.Tag.Column == "created_at" {
				continue
			} else {
				columns += ", " + prop.Tag.Column + "=?"
				values = append(values, prop.InterValue)
			}

		}

	}
	sql := fmt.Sprintf("update %s set %s where %s ", table, columns[2:], primaryKey)
	values = append(values, primaryKeyValue)
	return sql, values, nil
}

func (MySQLGenerator) InsertSQL(value interface{}) (autoIncrement *reflect.Value, query string, args []interface{}, err error) {
	table, props, err := ReflectTable(value)
	if err != nil {
		return nil, "", nil, err
	}
	n := len(props)
	if n <= 0 {
		return nil, "", nil, fmt.Errorf("not found insert columns")
	}
	var values []interface{}
	columns := ""
	strArg := ""
	for _, prop := range props {
		if prop.Tag.AutoIncrement {
			autoIncrement = prop.Value
		}
		if !prop.Tag.PrimaryKey {
			columns += ", " + prop.Tag.Column

			if prop.Tag.Column == "updated_at" || prop.Tag.Column == "created_at" {
				strArg += ",now()"
			} else {
				strArg += ",?"
				values = append(values, prop.InterValue)
			}
		}
	}
	query = fmt.Sprintf("insert into %s(%s) values(%s)", table, columns[2:], strArg[1:])
	return autoIncrement, query, values, err
}
