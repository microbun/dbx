package dbx

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	rs, err := db.SQLExec("insert into accounts(nickname) values(?)", "OKHome")
	if err != nil {
		t.Fatal(err)
	}
	id, err := rs.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	rs, err = db.SQLExec("insert into accounts(nickname) values(?)", "OKHome1")
	if err != nil {
		t.Fatal(err)
	}
	id, err = rs.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.rawDB.Begin()
	tx.Commit()
	t.Logf("id=%d", id)

	account := &Account{Status: 1}
	db.StructMustInsert(account)
	t.Logf("id=%v", account.ID)
}

func TestDB_Exec_Select(t *testing.T) {
	db, err := Open("sqlite3", "etsme_dev:123456@tcp(47.103.136.234:3306)/etsme_dev?parseTime=True&loc=Local")
	rs1, err := db.SQLExec("select * from accounts")
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
	rs := db.StructMustInsert(&UserNodeRecord{NodeID: "test"})
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
	r := &UserNodeRecord{NodeID: "test"}
	r.ID = 1
	rs, err := db.StructUpdate(r)
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
	err = db.SQLFind(&account, "select * from accounts")
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
	err = db.SQLFind(&nameArr, "select id from accounts")
	if err != nil {
		t.Fatalf("%v", err)
	}

	for _, name := range nameArr {
		fmt.Printf("name:%v\n", name)
	}

	var accountSingle = Account{}
	err = db.SQLFirst(&accountSingle, "select * from accounts where id = 33")
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Printf("account Single:%v\n", accountSingle)

	var nickname = []sql.NullString{}
	err = db.SQLFind(&nickname, "select nickname from accounts")
	if err != nil {
		t.Fatalf("%v", err)
	}
	for _, name := range nickname {
		// nickname
		fmt.Printf("json:%v\n", name.String)
	}

}

func TestDB_MustFind(t *testing.T) {
	db, err := Open("sqlite3", "file:locked.sqlite")
	if err != nil {
		t.Fatalf("open rawDB:%v", err)
	}
	var account = make([]*Account, 0)
	db.DQLMustFind(&account, "select * from accounts order by id desc", DQLArgument{
		"id_list": []int{44, 45},
	})
	for _, a := range account {
		b, err := json.Marshal(a)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("json:%v\n", string(b))
	}
}

func ExampleOpen() {
	db, err := Open("sqlite3", "file:locked.sqlite?cache=shared")
	if err != nil {
		panic(err)
	}
	a := `create table accounts(
		id                       integer primary key autoincrement,
		created_by               varchar(64),
		status                   integer default 0,
		created_at               datetime,
		updated_at               datetime,
		uid                      long,
		nickname                 varchar(24),
		avatar                   varchar(512),
		network_id               varchar(128)
	)`
	rs, err := db.SQLExec(a)
	if err != nil {
		panic(err)
	}

	_, err = db.StructInsert(&Account{
		Record:   Record{},
		NickName: sql.NullString{},
		Status:   0,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", rs)
}

type Account struct {
	Record
	ID        int            `dbx:"column:id" `
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

type Record struct {
	ID        int64     `json:"id,omitempty" dbx:"column:id,primary_key" `
	CreatedAt time.Time `json:"created_at,omitempty" dbx:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" dbx:"column:updated_at"`
}

//UserNodeDeviceStatus is enum
type UserNodeDeviceStatus string

//UserNodeRecord is Database Model of user node
type UserNodeRecord struct {
	Record
	UID       int64              `dbx:"column:uid"`
	NodeType  UserNodeDeviceType `dbx:"column:node_type"`
	NodeID    string             `dbx:"column:node_id"`
	NodeIP    string             `dbx:"column:node_ip"`
	CreatedBy int64              `dbx:"column:created_by" `
	UpdatedBy int64              `dbx:"column:updated_by" `
}

func (a UserNodeRecord) TableName() string {
	return "user_node"
}
