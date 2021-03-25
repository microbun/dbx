package dbx

import (
	"database/sql"
	"testing"
)

func TestTx_Insert(t *testing.T) {
	db, err := Open("sqlite3", "file:locked.sqlite?cache=shared")
	if err != nil {
		t.Fatalf("open rawDB:%v", err)
	}
	account := Account{}
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	// context.
	nn := sql.NullString{}
	nn.Scan("abc")
	account.NickName = nn
	_, err = tx.StructInsert(&account)
	if err != nil {
		t.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	// rawTx.SQLFirst(&account, "select * from accounts")
}
