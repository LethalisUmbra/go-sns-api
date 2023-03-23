package utils

import (
	"database/sql"
)

var DB *sql.DB

func OpenDB() error {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/sns-test?parseTime=true")
	return err
}
