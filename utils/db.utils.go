package utils

import (
	"database/sql"
	"net/url"
)

var DB *sql.DB

func OpenDB() error {
	password := "fVQdk;_p}$.y1(@"
	encodedPassword := url.QueryEscape(password)
	dataSourceName := "root:" + encodedPassword + "@tcp(34.176.30.163)/sns-test"
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	return err
}
