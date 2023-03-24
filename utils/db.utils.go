package utils

import (
	"database/sql"
)

var DB *sql.DB

func OpenDB() error {
	var err error
	DB, err = sql.Open("mysql", "root:fVQdk;_p}$.y1(@,@tcp(test-sns-api:southamerica-west1:test-sns)/sns-test?parseTime=true")
	return err
}
