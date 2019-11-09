package onetimepass

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	database *sql.DB
)

type secretRecord struct {
	ID         int
	Secret     string
	Expiration time.Time
}

func init() {

	var err error

	database, err = sql.Open("sqlite3", "./otp.sql")

	if err != nil {
		panic("Unable to create database")
	}

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

func AddSecret(s string, exp int64) (int, error) {

	stmt, err := database.Prepare("INSERT INTO otp (secret, expiration) VALUES (?, ?)")
	if err != nil {
		fmt.Println("Error preparing query")
		return 0, err
	}

	r, err := stmt.Exec(s, exp)

	if err != nil {
		fmt.Println("Error executing query")
		return 0, err
	}

	id, err := r.LastInsertId()

	if err != nil {
		fmt.Println("Error retrieving last inserted ID")
		return 0, err
	}

	return int(id), nil
}

func deleteSecret(id string) error {
	stmt, err := database.Prepare("DELETE FROM otp WHERE id=?")
	if err != nil {
		fmt.Println("Error preparing deleteSecret query")
		return err
	}

	_, err = stmt.Exec(id)

	if err != nil {
		fmt.Println("Error executing query")
		return err
	}

	return nil
}

func ReturnSecret(id string) (string, error) {

	var secret string
	var expiration int64

	rows := database.QueryRow("SELECT secret, expiration FROM otp WHERE id=?", id)
	/*
		if err != nil {
			fmt.Println("Error querying database to get secret")
			log.Println(err)
			return "Internal error", err
		}
	*/

	switch err := rows.Scan(&secret, &expiration); err {
	case sql.ErrNoRows:
		return "Secret not found", errors.New("empty")
	case nil:
		if time.Now().Unix() > expiration {
			e := errors.New("expired")
			_ = deleteSecret(id)

			return "Secret expired", e

		}

		_ = deleteSecret(id)

		return secret, nil

	default:
		fmt.Println(err)
		return "Internal error", err
	}

}
