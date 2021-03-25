package reflectx

import (
	"errors"
	"reflect"
)

type valueType int

const (
	structType valueType = 1
	basicType  valueType = 2
	sliceType  valueType = 3
	otherType  valueType = 4
)

type Column struct {
	Name  string
	Value interface{}
}

var convertType = map[reflect.Kind]valueType{
	reflect.Invalid:       otherType,
	reflect.Bool:          basicType,
	reflect.Int:           basicType,
	reflect.Int8:          basicType,
	reflect.Int16:         basicType,
	reflect.Int32:         basicType,
	reflect.Int64:         basicType,
	reflect.Uint:          basicType,
	reflect.Uint8:         basicType,
	reflect.Uint16:        basicType,
	reflect.Uint32:        basicType,
	reflect.Uint64:        basicType,
	reflect.Uintptr:       otherType,
	reflect.Float32:       basicType,
	reflect.Float64:       basicType,
	reflect.Complex64:     otherType,
	reflect.Complex128:    otherType,
	reflect.Array:         sliceType,
	reflect.Chan:          otherType,
	reflect.Func:          otherType,
	reflect.Interface:     otherType,
	reflect.Map:           otherType,
	reflect.Ptr:           otherType,
	reflect.Slice:         sliceType,
	reflect.String:        basicType,
	reflect.Struct:        structType,
	reflect.UnsafePointer: otherType,
}

var mapType = map[string]valueType{
	"database/sql.NullBool":    basicType,
	"database/sql.NullFloat64": basicType,
	"database/sql.NullInt32":   basicType,
	"database/sql.NullInt64":   basicType,
	"database/sql.NullTime":    basicType,
	"database/sql.NullString":  basicType,
	"time.Time":                basicType,
}

func ReflectProperty(v reflect.Value, mapping map[string]Property) {
	direct := reflect.Indirect(v)
	count := direct.NumField()
	dt := direct.Type()
	for i := 0; i < count; i++ {
		fv := direct.Field(i)
		ft := dt.Field(i)
		tag := newDbxTag(ft.Tag.Get("dbx"))
		if !IsBasicType(ft.Type) {
			ReflectProperty(fv, mapping)
		}
		if tag.Column != "" {
			mapping[tag.Column] = Property{
				InterValue: fv.Addr().Interface(),
				Value:      &fv,
				Tag:        tag,
			}
		}
	}
}

func indirectPtr(dest interface{}) (reflect.Value, error) {
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return value, errors.New("argument not a ptr")
	}
	return reflect.Indirect(value), nil
}

func IsBasicValue(value reflect.Value) bool {
	fullName := value.Type().PkgPath() + "." + value.Type().Name()
	return convertType[value.Kind()] == basicType || mapType[fullName] == basicType
}

// func isBasicValue(value reflect.InterValue) bool {
// 	fullname := value.Type().PkgPath() + "." + value.Type().Name()
// 	return convertType[value.Kind()] == basicType || mapType[fullname] == basicType
// }

func IsStructValue(value reflect.Value) bool {
	return convertType[value.Kind()] == structType
}

func IsSliceType(value reflect.Type) bool {
	return convertType[value.Kind()] == sliceType
}

func IsBasicType(t reflect.Type) bool {
	fullName := t.PkgPath() + "." + t.Name()
	return convertType[t.Kind()] == basicType || mapType[fullName] == basicType
}

func IsStructType(t reflect.Type) bool {
	return convertType[t.Kind()] == structType
}
