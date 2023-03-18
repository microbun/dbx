package reflectx

import (
	"reflect"
	"strings"
)

//Property contains the value and tags of a field
type Property struct {
	InterValue interface{}
	Value      *reflect.Value
	Tag        *Tag
}

//Properties is a property array
type Properties []Property

//NewProperties return a Properties
func NewProperties(n int) Properties    { return make([]Property, n) }
func (p Properties) Len() int           { return len(p) }
func (p Properties) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Properties) Less(i, j int) bool { return p[i].Tag.Column < p[j].Tag.Column }

//Values return all values of Properties
func (p Properties) Values() []interface{} {
	values := make([]interface{}, len(p))
	for i, v := range p {
		values[i] = v.InterValue
	}
	return values
}

//Tag contains all info of `dbx` tag
type Tag struct {
	Column        string
	PrimaryKey    bool
	Insert        string
	Update        string
	AutoIncrement bool
}

func newDbxTag(tag string) *Tag {
	t := &Tag{}
	props := strings.Split(tag, ",")
	for _, prop := range props {
		if strings.ContainsAny(prop, ":") {
			splitIdx := strings.Index(prop, ":")
			propName := strings.TrimSpace(prop[:splitIdx])
			if propName == "column" {
				t.Column = strings.TrimSpace(prop[splitIdx+1:])
			}
			if propName == "update" {
				t.Update = strings.TrimSpace(prop[splitIdx+1:])
			}
			if propName == "insert" {
				t.Insert = strings.TrimSpace(prop[splitIdx+1:])
			}
		} else {
			if strings.TrimSpace(prop) == "primary_key" {
				t.PrimaryKey = true
			}
			if strings.TrimSpace(prop) == "auto_increment" {
				t.AutoIncrement = true
			}
		}
	}
	return t
}
