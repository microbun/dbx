package main

import (
	"flag"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" //justifying
	_ "github.com/mattn/go-sqlite3"    //justifying
	"github.com/microbun/dbx/dbx-gen/internal"
)

type Avatar struct {
	Id        []byte     `dbx:"column:id"`
	Avatar    []byte     `dbx:"column:avatar"`
	CreatedAt time.Time  `dbx:"column:created_at"`
	UpdatedAt time.Time  `dbx:"column:updated_at"`
	DeletedAt *time.Time `dbx:"column:deleted_at"`
}

func main() {
	flag.StringVar(&internal.Options.Package, "p", "module", "Golang package name")
	flag.StringVar(&internal.Options.Output, "o", "module.gen.go", "Write output to a `file`")
	flag.StringVar(&internal.Options.Driver, "driver", "mysql", "Database driver name")
	flag.StringVar(&internal.Options.DataSourceName, "uri", "", "Data source name")
	flag.StringVar(&internal.Options.Schema, "schema", "", "Database schema")
	flag.Parse()
	if !flag.Parsed() {
		flag.PrintDefaults()
		return
	}
	err := internal.Run()
	if err != nil {
		fmt.Println(err)
	}
}
