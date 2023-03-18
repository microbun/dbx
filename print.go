package dbx

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type nullable interface {
	Value() (driver.Value, error)
}

func printSQL(query string, args []interface{}, opts *Options) {
	if opts.Logger != nil {
		opts.Logger.Printf(formatSQL(query, args, opts))
	}
}

func formatSQL(query string, args []interface{}, opts *Options) string {
	compile := query
	for _, arg := range args {
		var value string
		rv := reflect.ValueOf(arg)
		kind := rv.Kind()
		for kind == reflect.Ptr {
			rv = reflect.Indirect(rv)
			kind = rv.Kind()
		}
		if !rv.IsValid() {
			value = "null"
		} else {
			typeName := rv.Type().String()
			switch typeName {
			case "string":
				{
					value = "'" + strings.ReplaceAll(rv.String(), "'", "\\'") + "'"
				}
			case "int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64":
				{
					value = fmt.Sprintf("%v", rv.Int())
				}
			case "float32", "float64":
				{
					value = fmt.Sprintf("%v", rv.Float())
				}
			case "time.Time":
				{
					iv := rv.Interface()
					t := iv.(time.Time)
					t = t.In(opts.Location)
					value = fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04.05.999"))
				}
			case "[]int8", "[]uint8":
				{
					if len(rv.Bytes()) == 0 {
						value = "null"
					} else {
						value = "0x" + hex.EncodeToString(rv.Bytes())
					}
				}
			case "sql.NullString":
				{
					iv := rv.Interface()
					v := iv.(nullable)
					v0, err := v.Value()
					if err != nil {
						fmt.Printf("value err:%v", err)
					}
					if v0 == nil {
						value = "null"
					} else {
						value = fmt.Sprintf("%v", v0)
						value = "'" + strings.ReplaceAll(value, "'", "\\'") + "'"
					}
				}
			case "sql.NullTime":
				{
					iv := rv.Interface()
					v := iv.(nullable)
					v0, err := v.Value()
					if err != nil {
						fmt.Printf("value err:%v", err)
					}
					if v0 == nil {
						value = "null"
					} else {
						t := v0.(time.Time)
						t = t.In(opts.Location)
						value = fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04.05.999"))
					}
				}
			case "sql.NullInt64", "sql.NullInt32", "sql.NullFloat64", "sql.NullBool":
				{
					iv := rv.Interface()
					v := iv.(nullable)
					v0, err := v.Value()
					if err != nil {
						fmt.Printf("value err:%v", err)
					}
					if v0 == nil {
						value = "null"
					} else {
						value = fmt.Sprintf("%v", v)
					}
				}

			default:
				value = fmt.Sprintf("'%v'", rv.Interface())
			}
		}
		compile = strings.Replace(compile, "?", value, 1)
	}
	return compile
}
