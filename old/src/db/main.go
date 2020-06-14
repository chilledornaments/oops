package onetimepass

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
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
	/*
		SQLITE3 uses the AUTOINCREMENT attribute
		MySQL uses the AUTO_INCREMENT attribute
		Create a local var to store our DB init statement in so that we can switch between either type
	*/
	var initStatement string

	err = godotenv.Load(os.Getenv("OOPS_ENV_FILE"))
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")

	}

	dbType := os.Getenv("DB_DRIVER")

	switch dbType {
	case "sqlite3":
		fmt.Println("Using SQLITE3 database at", os.Getenv("DB_PATH"))
		database, err = sql.Open("sqlite3", os.Getenv("DB_PATH"))
		initStatement = "CREATE TABLE IF NOT EXISTS otp (id INTEGER PRIMARY KEY AUTOINCREMENT, secret TEXT, expiration INT, uuid TEXT)"
	case "mysql":
		sqlString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
		database, err = sql.Open("mysql", sqlString)
		initStatement = "CREATE TABLE IF NOT EXISTS otp (id INTEGER PRIMARY KEY AUTO_INCREMENT, secret TEXT, expiration INT, uuid TEXT)"
	}

	if err != nil {
		panic("Unable to connect to database")
	}

	err = database.Ping()

	if err != nil {
		fmt.Println("Error pinging DB")
		log.Fatal(err)
	}

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

// AddSecret adds a secret to the database and returns the secrets UUID
func AddSecret(s string, exp int64) (string, error) {

	stmt, err := database.Prepare("INSERT INTO otp (secret, expiration, uuid) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Error preparing query")
		return "", err
	}

	b := make([]byte, 16)
	rand.Read(b)
	u := hex.EncodeToString(b)

	_, err = stmt.Exec(s, exp, u)

	if err != nil {
		fmt.Println("Error executing query")
		return "", err
	}

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

// ReturnSecret returns the value of a secret if that secret is not more than an hour old and has not been viewed before
func ReturnSecret(uuid string) (string, error) {

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

func main() {}
