package onetimepass

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var (
	database *sql.DB
	err      error
)

func main() {
	database, err = sql.Open("sqlite3", "./otp.sql")

	if err != nil {
		panic("Unable to create database")
	}

	defer database.Close()

	err = database.Ping()
	if err != nil {
		fmt.Println(err)
	}

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS otp (id INTEGER PRIMARY KEY, secret TEXT)")

	r, err := statement.Exec()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r)
}

func AddSecret(s string) error {
	stmt, err := database.Prepare("INSERT INTO otp (1, ?)")

	if err != nil {
		fmt.Println(err)
	}

	r, err := stmt.Exec(s)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
	return nil
}
