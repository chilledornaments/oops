package onetimepass

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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

	err = godotenv.Load(os.Getenv("OTP_ENV_FILE"))
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")

	}

	dbType := os.Getenv("DB_DRIVER")

	switch dbType {
	case "sqlite3":
		database, err = sql.Open("sqlite3", os.Getenv("DB_PATH"))
	case "mysql":
		sqlString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
		database, err = sql.Open("mysql", sqlString)
	}

	if err != nil {
		panic("Unable to connect to database")
	}

	err = database.Ping()

	if err != nil {
		fmt.Println("Error pinging DB")
		log.Fatal(err)
	}

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS otp (id INTEGER PRIMARY KEY AUTO_INCREMENT, secret TEXT, expiration INT, uuid TEXT)")

	_, err = statement.Exec()

	if err != nil {
		fmt.Println(err)
	}

}

func main() {

	fmt.Println("Hello from the database")

}

func AddSecret(s string, exp int64) (string, error) {

	stmt, err := database.Prepare("INSERT INTO otp (secret, expiration, uuid) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing query")
		return "", err
	}

	b := make([]byte, 16) //equals 8 charachters
	rand.Read(b)
	u := hex.EncodeToString(b)

	_, err = stmt.Exec(s, exp, u)

	if err != nil {
		fmt.Println("Error executing query")
		return "", err
	}

	/*
		id, err := r.LastInsertId()

		if err != nil {
			fmt.Println("Error retrieving last inserted ID")
			return "", err
		}
	*/

	return u, nil
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

func ReturnSecret(uuid string) (string, error) {

	var secret string
	var expiration int64

	rows := database.QueryRow("SELECT secret, expiration FROM otp WHERE uuid=?", uuid)
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
			_ = deleteSecret(uuid)

			return "Secret expired", e

		}

		_ = deleteSecret(uuid)

		return secret, nil

	default:
		fmt.Println(err)
		return "Internal error", err
	}

}
