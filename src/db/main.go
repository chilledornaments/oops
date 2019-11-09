package onetimepass

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	database *sql.DB
	//err      error
)

func init() {

	var err error

	database, err = sql.Open("sqlite3", "./otp.sql")

	if err != nil {
		panic("Unable to create database")
	}

	//defer database.Close()

	err = database.Ping()

	if err != nil {
		fmt.Println("Error pinging DB")
		fmt.Println(err)
	}

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS otp (id INTEGER PRIMARY KEY, secret TEXT, expiration INT)")

	_, err = statement.Exec()

	if err != nil {
		fmt.Println(err)
	}

}

func main() {

	fmt.Println("Hello from the database")

}

func AddSecret(s string, exp int64) error {

	stmt, err := database.Prepare("INSERT INTO otp (secret, expiration) VALUES (?, ?)")
	if err != nil {
		fmt.Println("Error preparing query")
		return err
	}

	_, err = stmt.Exec(s, exp)

	if err != nil {
		fmt.Println("Error executing query")
		return err
	}

	return nil
}
