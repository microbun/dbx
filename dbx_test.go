package dbx

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql" //justifying
	_ "github.com/mattn/go-sqlite3"    //justifying
)

func Test_open(t *testing.T) {
	ExampleOpen()
}

func TestDB_Exec_Insert(t *testing.T) {
	db, err := Open("sqlite3", "file:locked.sqlite")
	if err != nil {
		t.Fatalf("open db err:%v", err)
	}

	rs, err := db.Exec("insert into accounts(nickname) values(?)", "OKHome")
	if err != nil {
		t.Fatal(err)
	}
	id, err := rs.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.rawDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	err = tx.Commit()
	t.Logf("id=%d", id)

	account := &Account{Status: 1}
	db.MustInsert(account)
	t.Logf("id=%v", account.ID)
}

func TestDB_Exec_Select(t *testing.T) {
	rs1, err := db.Exec("select * from accounts")
	if err != nil {
		t.Fatal(err)
	}
	id1, err := rs1.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("id=%d", id1)
}

func TestDB_Insert(t *testing.T) {
	db, err := Open("mysql", "etsme_dev:123456@tcp(47.103.136.234:3306)/etsme_dev?parseTime=True&loc=Local")
	if err != nil {
		t.Fatalf("open rawDB:%v", err)
	}
	rs := db.MustInsert(&Account{NickName: sql.NullString{String: "test"}})
	lid, err := rs.LastInsertId()
	if err != nil {
		t.Fatalf("last id:%v", err)
	}
	fmt.Printf("rs:%d", lid)
}

func TestDB_Update(t *testing.T) {
	db, err := Open("mysql", "etsme_dev:123456@tcp(47.103.136.234:3306)/etsme_dev?parseTime=True&loc=Local")

	if err != nil {
		t.Fatalf("open rawDB:%v", err)
	}
	r := &Account{NickName: sql.NullString{String: "test"}}
	r.ID = 1
	rs, err := db.Update(r)
	if err != nil {
		t.Fatalf("insert:%v", err)
	}
	lid, err := rs.LastInsertId()
	if err != nil {
		t.Fatalf("last id:%v", err)
	}
	fmt.Printf("rs:%d", lid)
}

func TestDB_Query(t *testing.T) {
	db, err := Open("sqlite3", "file:locked.sqlite?cache=shared")
	if err != nil {
		t.Fatalf("open rawDB:%v", err)
	}
	var account = make([]*Account, 0)
	err = db.Query(&account, "select * from accounts")
	if err != nil {
		t.Fatalf("%v", err)
	}

	for _, a := range account {
		b, err := json.Marshal(a)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("json:%v\n", string(b))
	}

	var nameArr []int
	err = db.Query(&nameArr, "select id from accounts")
	if err != nil {
		t.Fatalf("%v", err)
	}

	for _, name := range nameArr {
		fmt.Printf("name:%v\n", name)
	}

	var accountSingle = Account{}
	err = db.Get(&accountSingle, "select * from accounts where id = 33")
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Printf("account Single:%v\n", accountSingle)

	var nickname = []sql.NullString{}
	err = db.Query(&nickname, "select nickname from accounts")
	if err != nil {
		t.Fatalf("%v", err)
	}
	for _, name := range nickname {
		// nickname
		fmt.Printf("json:%v\n", name.String)
	}

}

func TestExecutor_MustQuery(t *testing.T) {
	var examples = make([]*ExampleRecord, 0)
	now := time.Now()
	//nowPtr:=&now
	//nowPtrPtr:=&nowPtr
	db.MustQuery(&examples, "select * from example where `datetime`<? order by id desc", now)
	//db.MustNamedQuery(&account, "select * from example where `datetime`<:date order by id desc", map[string]interface{}{
	//	"date": now,
	//})
}

func Example_MustNamedQuery() {
	var examples = make([]*ExampleRecord, 0)
	now := time.Now()
	db.MustNamedQuery(&examples, "select * from example where `datetime`<? order by id desc", map[string]interface{}{
		"date": now,
	})
}
func ExampleExecutor_MustNamedQuery() {
	var examples = make([]*ExampleRecord, 0)
	now := time.Now()
	db.MustNamedQuery(&examples, "select * from example where `datetime`<? order by id desc", map[string]interface{}{
		"date": now,
	})
}

func TestExecutor_MustNamedQuery(t *testing.T) {
	ExampleExecutor_MustNamedQuery()
}
func TestDB_MustFind(t *testing.T) {
	var account = make([]*ExampleRecord, 0)
	now := time.Now()
	//nowPtr:=&now
	//nowPtrPtr:=&nowPtr
	db.MustQuery(&account, "select * from example where `datetime`<? order by id desc", now)
	//db.MustNamedQuery(&account, "select * from example where `datetime`<:date order by id desc", map[string]interface{}{
	//	"date": now,
	//})
	for _, a := range account {
		b, err := json.Marshal(a)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("json:%v\n", string(b))
	}
}

type Account struct {
	ID        int64          `json:"id,omitempty" dbx:"column:id,primary_key" `
	CreatedAt sql.NullTime   `json:"created_at,omitempty" dbx:"column:created_at,insert:current_timestamp"`
	UpdatedAt sql.NullTime   `json:"updated_at,omitempty" dbx:"column:updated_at,update:current_timestamp,insert:current_timestamp"`
	UID       sql.NullInt64  `dbx:"column:uid" `
	NickName  sql.NullString `dbx:"column:nickname" `
	Status    int            `dbx:"column:status" `
	CreatedBy sql.NullString `dbx:"column:created_by"`
	Avatar    sql.NullString `dbx:"column:avatar"`
	NetworkID sql.NullInt64  `dbx:"column:network_id"`
}

func (a Account) TableName() string {
	return "accounts"
}

//UserNodeDeviceType is enum
type UserNodeDeviceType string

const (
	//UserNodeDeviceTypeIOS is iOS device of UserNodeRecord
	UserNodeDeviceTypeIOS UserNodeDeviceType = "iOS"
	//UserNodeDeviceTypeAndroid is Android device of UserNodeRecord
	UserNodeDeviceTypeAndroid UserNodeDeviceType = "Android"
)

func TestDB_MySQLConn(t *testing.T) {
	db, err := Open("mysql", "basebit:123@tcp(localhost:3306)/basebit?parseTime=True&loc=Local")
	if err != nil {
		t.Fatal(err)
	}
	wg := sync.WaitGroup{}
	rand.Seed(time.Now().Unix())
	printDB := func() {
		wg.Add(1)
		for {
			account := &ExampleRecord{}
			time.Sleep(time.Duration(rand.Int31n(3000)) * time.Millisecond)
			tx, err := db.Begin()
			if err != nil {
				t.Logf("open tx err:%v\n", err)
			}
			err = tx.NamedGet(account, "select * from example ", map[string]interface{}{
				"name": []string{"bang", "xxx"},
			})
			if err != nil {
				t.Logf("select err:%v\n", err)
			}
			// t.Logf("result name=%v ,wait 1s", name)
			time.Sleep(time.Duration(rand.Int31n(2000)) * time.Millisecond)
			tx.Commit()
		}
		// wg.Done()
	}
	for i := 0; i < 10; i++ {
		go printDB()
	}
	for {
		stats := db.RawDB().Stats()
		t.Logf("maxConnections:%v idle:%v inuse:%v wait:%v", stats.OpenConnections, stats.Idle, stats.InUse, stats.WaitCount)
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func TestDefaultExecutor_Insert(t *testing.T) {
	e := &ExampleRecord{
		Datetime:  time.Now(),
		Timestamp: time.Now(),
		Enum: "yes",
		Bit: []byte{0x1},
	}
	_, err = db.Insert(e)
	if err != nil {
		t.Fatal(err)
	}
}

func ExampleOpen() {
	db, err = Open("mysql", "root:basebitxdp@tcp(172.18.0.210:32600)/enigma2_accountx?parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
}
