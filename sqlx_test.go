package dbx

import (
	"testing"
	"time"
)

type Acc struct {
	FirstName string `db:"id"`
	LastName  string `db:"name"`
	S         string `db:"xx"`
}

func TestSQLX(t *testing.T) {
	//db, err := sqlx.Open("mysql", "basebit:123@tcp(localhost:3306)/basebit?parseTime=True&loc=Local")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//acc := &Acc{}
	//err = db.Get(acc, "select id,name from accounts limit 1")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Logf("acc:%v", acc)

	time.Now().Format(time.RFC3339)

}
