package dbx

import "context"

var db, err = Open("sqlite3", "file:locked.sqlite")

func ExampleComplexExec_DQLExecContext() {
	_, err = db.DQLExecContext(context.Background(), "update account set name=:name", DQLArgument{
		"name": "Lucy",
	})
}

func ExampleComplexExecutor_DQLExec() {
	_, err = db.DQLExec("update account set name=:name", DQLArgument{
		"name": "Lucy",
	})
}

func ExampleComplexExec_DQLFirst() {
	err = db.DQLFirst(nil, "update account set name=:name", DQLArgument{
		"name": "Lucy",
	})
}

func ExampleComplexExec_DQLFind() {
	account := &Account{}
	err := db.DQLFind(account, "select * from account where id in (:id)", DQLArgument{
		"id": []int{1, 2, 3},
	})
	if err != nil {
		panic(err)
	}
}
