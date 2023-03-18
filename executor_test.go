package dbx

import "context"



func ExampleComplexExec_NamedExecContext() {
	_, err = db.NamedExecContext(context.Background(), "update account set name=:name", map[string]interface{}{
		"name": "Lucy",
	})
}

func ExampleComplexExecutor_NamedExec() {
	_, err = db.NamedExec("update account set name=:name", map[string]interface{}{
		"name": "Lucy",
	})
}




func ExampleComplexExec_NamedFind() {
	account := &Account{}
	err := db.NamedQuery(account, "select * from account where id in (:id)", map[string]interface{}{
		"id": []int{1, 2, 3},
	})
	if err != nil {
		panic(err)
	}
}

