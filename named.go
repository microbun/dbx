package dbx

import (
	"fmt"
	"git.basebit.me/enigma/dbx/reflectx"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

var reg = regexp.MustCompile(":\\w+")

func namedCompile(query string, p map[string]interface{}) (string, []interface{}, error) {
	args := make([]interface{}, 0)
	matched := reg.FindAllString(query, -1)
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
			if rt.Elem().String()=="uint8"{
				args = append(args, value)
				placeholders[parameter] = "?"
				continue
			}
			if rt.Elem().String()!="[]uint8" && rt.Elem().String()!="[]int8" && !reflectx.IsBasicType(rt.Elem()) {
				return "", nil, fmt.Errorf("unsupport args type at %s", parameter)
			}
			repeat := make([]string, 0)
			if rv.Len() == 0 {
				return "", nil, fmt.Errorf("`%s` len must be greater than 0", parameter)
			}
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
		placeholder := placeholders[key]
		query = strings.ReplaceAll(query, key, placeholder)
	}

	return query, args, nil

}

