package dbx

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/microbun/dbx/reflectx"
	"reflect"
)

type mode int

const (
	one   mode = 2
	array mode = 3
)

func toDest(rows *sql.Rows, dest interface{}, m mode) error {
	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return errors.New("argument not a ptr")
	}
	direct := reflect.Indirect(value)
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	switch m {
	case array:
		{
			return toArray(direct, columns, rows)
		}
	case one:
		{
			return toOne(direct, columns, rows)
		}
	default:
		{
			return errors.New("unknown scan m")
		}
	}
}

func toOne(dest reflect.Value, columns []string, rows *sql.Rows) error {
	if reflectx.IsBasicValue(dest) {
		if len(columns) == 1 {
			for rows.Next() {
				return rows.Scan(dest.Addr().Interface())
			}
		} else {
			return fmt.Errorf("multi columns not scan to a basic type")
		}
	} else if reflectx.IsStructValue(dest) {
		properties := reflectx.NewProperties(len(columns))
		for rows.Next() {
			err := traversal(dest, properties, columns)
			if err != nil {
				return err
			}
			return rows.Scan(properties.Values()...)
		}
	} else {
		return fmt.Errorf("argument not a struct or basic")
	}
	return nil
}

func toArray(dest reflect.Value, columns []string, rows *sql.Rows) error {
	kind := dest.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		return fmt.Errorf("argument not a array or slice")
	}
	valueType := dest.Type().Elem()
	isPtr := valueType.Kind() == reflect.Ptr
	if isPtr {
		valueType = valueType.Elem()
	}
	if reflectx.IsBasicType(valueType) {
		if len(columns) == 1 {
			for rows.Next() {
				pv := reflect.New(valueType)
				dv := reflect.Indirect(pv)
				err := rows.Scan(pv.Interface())
				if err != nil {
					return err
				}
				if isPtr {
					dest.Set(reflect.Append(dest, pv))
				} else {
					dest.Set(reflect.Append(dest, dv))
				}
			}
		} else {
			return fmt.Errorf("dest slice not a struct or basic")
		}
	} else if reflectx.IsStructType(valueType) {

		for rows.Next() {
			pv := reflect.New(valueType)
			dv := reflect.Indirect(pv)
			properties := reflectx.NewProperties(len(columns))
			err := traversal(dv, properties, columns)
			if err != nil {
				return err
			}
			err = rows.Scan(properties.Values()...)
			if err != nil {
				return err
			}
			if isPtr {
				dest.Set(reflect.Append(dest, pv))
			} else {
				dest.Set(reflect.Append(dest, dv))
			}
		}
	} else {
		return fmt.Errorf("unknown array type")
	}
	return nil
}

func traversal(v reflect.Value, props reflectx.Properties, columns []string) error {
	direct := reflect.Indirect(v)
	sv := map[string]reflectx.Property{}
	reflectx.ReflectProperty(v, sv)
	for i, name := range columns {
		prop, ok := sv[name]
		if ok {
			props[i] = prop
		} else {
			return fmt.Errorf("missing field `%s` in %s.%s", name, direct.Type().PkgPath(), direct.Type().Name())
		}
	}
	return nil
}
