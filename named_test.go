package dbx

import (
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_compile(t *testing.T) {
	cq, ca, err := namedCompile(`update conversation_list set status=:status, updated_at = datetime('now')
	where conversation_id=:conversationID and created_by IN (:uid) or status!=:status and status!=:status`, map[string]interface{}{
		"status":         500,
		"conversationID": "abc",
		"uid":            []int64{1, 2, 3, 4},
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("SQL:%v args:%v", cq, ca)
}
