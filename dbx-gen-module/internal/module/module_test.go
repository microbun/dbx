package module

import (
	"encoding/hex"
	"git.basebit.me/enigma/dbx"
	"git.basebit.me/enigma/dbx/escape"
	_ "github.com/go-sql-driver/mysql" //justifying
	"testing"
)

var db *dbx.DB

func init() {
	var err error
	db, err = dbx.Open("mysql", "root:basebitxdp@tcp(172.18.0.210:32600)/enigma2_accountx?parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
}
func Test_Query(t *testing.T) {
	var rs []*OrgLangRecord
	idstr := []string{
		"65C20768434A442AB32827858773909D",
		"CDAA56C31C194DFDB6FA820E420C8C3A",
		"7974D06932804EAF8C02B94A93A86963",
		"E4CD832F4A8E4344A2B2FFEF7029E689",
		"9C6CF15587964C49B23488B31D4BAB76",
		"55CC41E69C6C4EC9AD82B5337DD0B9EB",
		"31DD74C38AE24EC4A771F80DCF96D149",
		"5BA62E3E901C4DCE90438B818165063A",
	}
	var idBytes [][]byte
	var b []byte
	for _, s := range idstr {
		b, err := hex.DecodeString(s)
		if err != nil {
			t.Fatal(err)
		}
		idBytes = append(idBytes, b)
	}

	db.MustNamedQuery(&rs, "select * from org_lang where id=:id or id in (:id)", map[string]interface{}{
		"id":  b,
		"ids": idBytes,
	})

	user := &UserRecord{}
	exists := db.MustGet(user, "select * from user where status = 1")
	if !exists {
		t.Logf("user exists:%v", false)
		return
	}
	t.Logf("user exists:%v", true)
}

func Test_Update(t *testing.T) {
	db, err := dbx.Open("mysql", "root:basebitxdp@tcp(172.18.0.210:32600)/enigma2_accountx?parseTime=True&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	//uuid.
	_, err = db.Update(&OrgRecord{
		ID:   []byte{0x1, 0x3, 0x5},
		Name: "",
	}, "name")
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuery(t *testing.T) {
	var rs []*OrgLangRecord
	//keyword:="\\%Data"
	keyword := "'"
	err := db.NamedQuery(&rs, "select * from org_lang where name like concat('%',:name,'%')", map[string]interface{}{
		"name":escape.Like(keyword),
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range rs {
		t.Logf("r:%v", i.PettyJSON())
	}

}


type Timeline struct {
	AvgProgress float32 `dbx:"column:avg_progress" json:"avg_progress,omitempty" `
	MinStart    string `dbx:"column:min_start" json:"min_start,omitempty" `
	MaxEnd      string `dbx:"column:max_end" json:"max_end,omitempty" `
}

func TestQueryWorkflow(t *testing.T)  {
	db, err := dbx.Open("mysql", "root:basebitxdp@tcp(172.18.0.210:32600)/enigma2_workflowx?parseTime=True&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	r:=&Timeline{}
	err = db.Get(r,"select AVG(progress) as avg_progress,min( if(UNIX_TIMESTAMP(started_at)>0,started_at,null)) as min_start,max(if(UNIX_TIMESTAMP(end_at)>0,end_at,null)) as max_end from workflow where project_id='2a56f81e-9e87-4ebf-8a5a-6c729bd8fdef' and deleted_at is null")
	if err!=nil{
		t.Fatal(err)
	}

	t.Logf("%v", r)

}
