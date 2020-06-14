package internal

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func init() {

	/*
		SQLITE3 uses the AUTOINCREMENT attribute
		MySQL uses the AUTO_INCREMENT attribute
		Create a local var to store our DB init statement in so that we can switch between either type
	*/
	var err error

	godotenv.Load(os.Getenv("OOPS_ENV_FILE"))

	dbType := os.Getenv("DB_DRIVER")

	switch dbType {
	case "sqlite3":
		log.Println("Using SQLITE3 database at", os.Getenv("DB_PATH"))
		database, err = sql.Open("sqlite3", os.Getenv("DB_PATH"))
		if err != nil {
			log.Println(err.Error())
			panic("Unable to open SQLite database. Check your DB_PATH")
		}
		initStatement = "CREATE TABLE IF NOT EXISTS otp (id INTEGER PRIMARY KEY AUTOINCREMENT, secret TEXT, expiration INT, uuid TEXT)"
		initDB()
	}

}

func initDB() {

	var err error

	err = database.Ping()

	if err != nil {
		fmt.Println("Error pinging DB")
		log.Fatal(err)
	}

	log.Println("Initializing database")

	statement, err := database.Prepare(initStatement)

	if err != nil {
		log.Println("Error preparing DB init statement")
		log.Fatal(err)
	}

	_, err = statement.Exec()

	if err != nil {
		log.Println("Error initializing the database")
		log.Fatal(err)
	}

}

func AddSqliteSecret(secret string, exp int64) (string, error) {

	stmt, err := database.Prepare("INSERT INTO otp (secret, expiration, uuid) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing query")
		return "", err
	}

	b := make([]byte, 16)
	rand.Read(b)
	u := hex.EncodeToString(b)

	_, err = stmt.Exec(secret, exp, u)

	if err != nil {
		fmt.Println("Error executing query")
		return "", err
	}

	return u, nil
}

func ReturnSqliteSecret(uuid string) (string, error) {
	var secret string
	var expiration int64

	rows := database.QueryRow("SELECT secret, expiration FROM otp WHERE uuid=?", uuid)

	switch err := rows.Scan(&secret, &expiration); err {
	case sql.ErrNoRows:
		return "Secret not found", errors.New("empty")
	case nil:
		if time.Now().Unix() > expiration {
			e := errors.New("expired")
			_ = deleteSecret(uuid)

			return "Secret expired", e

		}

		err := deleteSecret(uuid)

		if err != nil {
			log.Println("Error deleting secret:", uuid)
			log.Println(err)
		}

		return secret, nil

	default:
		fmt.Println(err)
		return "Internal error", err
	}
}

func deleteSecret(uuid string) error {

	stmt, err := database.Prepare("DELETE FROM otp WHERE uuid=?")
	if err != nil {
		fmt.Println("Error preparing deleteSecret query")
		return err
	}

	_, err = stmt.Exec(uuid)

	if err != nil {
		fmt.Println("Error executing query")
		return err
	}

	return nil

}
