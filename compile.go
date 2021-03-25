package dbx

import (
	"database/sql/driver"
	"fmt"
	"github.com/microbun/dbx/reflectx"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
)

var reg = regexp.MustCompile(":\\w+")

type DQLArgument map[string]interface{}

func DSLCompile(dbxQuery string, p DQLArgument) (query string, args []interface{}, err error) {
	query = dbxQuery
	args = make([]interface{}, 0)
	matched := reg.FindAllString(dbxQuery, -1)
	placeholders := map[string]string{}
	for _, parameter := range matched {
		value, ok := p[parameter[1:]]
		if !ok {
			return "", nil, fmt.Errorf("`%s` not found ", parameter)
		}
		rt := reflect.TypeOf(value)
		rv := reflect.ValueOf(value)

		if reflectx.IsBasicType(rt) {
			args = append(args, value)
			placeholders[parameter] = "?"
		} else if reflectx.IsSliceType(rt) {
			if !reflectx.IsBasicType(rt.Elem()) {
				return "", nil, fmt.Errorf("unsupport args type at %s", parameter)
			}
			repeat := make([]string, 0)
			for j := 0; j < rv.Len(); j++ {
				args = append(args, rv.Index(j).Interface())
				repeat = append(repeat, "?")
			}
			placeholders[parameter] = strings.Join(repeat, ",")
		} else if reflectx.IsStructType(rt) {
			if len(p) > 1 {
				return "", nil, fmt.Errorf("unsupport args type at %s", parameter)
			} else {
				//解析结构体的参数
				return "", nil, fmt.Errorf("unsupport key type %s", parameter)
			}
		} else {
			return "", nil, fmt.Errorf("unsupport args type at %s", parameter)
		}
	}

	params := make([]string, 0)
	for p := range placeholders {
		params = append(params, p)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(params)))
	for _, key := range params {
		query = strings.ReplaceAll(query, key, placeholders[key])
	}

	return query, args, nil

}

type nullValue interface {
	Value() (driver.Value, error)
}

func printSQL(query string, args []interface{}, out io.Writer) {
	if out != nil {
		now := time.Now()
		s := fmt.Sprintf(" %s [SQL]:%s\n", now.Format("2006/1/2 15:04:05"), formatSQL(query, args))
		out.Write([]byte(s))
	}
}

func formatSQL(query string, args []interface{}) string {
	compile := query
	for _, arg := range args {
		var value string
		switch v := arg.(type) {
		case nullValue:
			{
				v, err := v.Value()
				if err != nil {
					fmt.Printf("value err:%v", err)
				}
				if v == nil {
					value = "null"
				} else {
					value = fmt.Sprintf("%v", v)
				}
			}
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64, int, uint:
			value = fmt.Sprintf("%v", v)
		case *uint8, *uint16, *uint32, *uint64, *int8, *int16, *int32, *int64, *float32, *float64, *int, *uint:
			iv := reflect.Indirect(reflect.ValueOf(v))
			value = fmt.Sprintf("%v", iv)
		default:
			{
				iv := reflect.Indirect(reflect.ValueOf(v))
				value = fmt.Sprintf("'%v'", iv)
			}
		}
		compile = strings.Replace(compile, "?", value, 1)
	}
	return compile
}
