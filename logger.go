package dbx

import "log"

type Logger interface {
	Printf(sql string)
}

type defaultLogger struct {
}

func (defaultLogger) Printf(sql string) {
	log.Printf("SQL=>%v", sql)
}

var logger = defaultLogger{}
