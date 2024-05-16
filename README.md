# ORM library for Golang

Open Database

```go
db, err := dbx.Open("sqlite3", "file:locked.sqlite")
```

Create Table

```go
sql := `create table accounts(
    id                       integer primary key autoincrement,
    status                   integer default 0,
    nickname                 varchar(24) not null,
    avatar                   varchar(512),
)`
rs, err := db.SQLExec(a)
if err != nil {
    panic(err)
}
```

Insert

```go
//Define Struct
type Account struct {
    ID        int            `dbx:"column:id" `
    NickName  string         `dbx:"column:nickname" `
    Status    int            `dbx:"column:status" `
    Avatar    sql.NullString `dbx:"column:avatar"`
}

func (a Account) TableName() string {
    return "accounts"
}

_, err = db.StructInsert(&Account{
    NickName: "jack",
    Avatar: sql.NullString{},
    Status:   0,
})
```

Query

```go
var account = make([]*Account, 0)
err = db.SQLFind(&account, "select * from accounts")

var nameArr []string
err = db.SQLFind(&nameArr, "select name from accounts")

var accountSingle = Account{}
err = db.SQLFirst(&accountSingle, "select * from accounts where id = ?", 33)

```

# Generate