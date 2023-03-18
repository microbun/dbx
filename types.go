package dbx

type Column interface {
	Take(v interface{}) (interface{},error)
	Put(v interface{}) (interface{},error)
}