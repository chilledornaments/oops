package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	rice "github.com/GeertJohan/go.rice"
	"github.com/joho/godotenv"

	internal "github.com/mitchya1/oops/internal"
)

type config struct {
	File        string
	DynamoTable string
	SqlitePath  string
	UseSQL      bool
	UseDynamo   bool
	Expiration  int64
	URL         string
	Port        string
	TLS         bool
	Certificate string
	Key         string
}

func loadConfig() {
	err := godotenv.Load(os.Getenv("OOPS_ENV_FILE"))
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")
	}

	conf = config{}

	conf.File = os.Getenv("OOPS_ENV_FILE")
	log.Println("Using this env file:", os.Getenv("OOPS_ENV_FILE"))

	if os.Getenv("DB_DRIVER") == "sqlite3" {
		conf.UseSQL = true
		conf.UseDynamo = false
	} else if os.Getenv("DB_DRIVER") == "dynamo" {
		conf.UseDynamo = true
		conf.UseSQL = false
		conf.DynamoTable = os.Getenv("DYNAMO_TABLE_NAME")
		internal.TableName = os.Getenv("DYNAMO_TABLE_NAME")
	} else {
		log.Fatal("Unable to determine database driver. Must be either sqlite3 or dynamo. Please check the README")
	}

	if os.Getenv("LINK_EXPIRATION_TIME") == "" {
		log.Println("LINK_EXPIRATION_TIME not set in env file. Using 3600")
		conf.Expiration = 3600
	} else {
		convertedString, err := strconv.ParseInt(os.Getenv("LINK_EXPIRATION_TIME"), 10, 64)

		if err != nil {
			panic("Unable to convert string to int")
		}

		log.Println("Expiration is", convertedString)
		conf.Expiration = convertedString

	}

	log.Println("SITE_URL is", os.Getenv("SITE_URL"))
	log.Println("WEB_SERVER_PORT is", os.Getenv("WEB_SERVER_PORT"))

	conf.URL = os.Getenv("SITE_URL")
	conf.Port = os.Getenv("WEB_SERVER_PORT")

	tls, err := strconv.ParseBool(os.Getenv("SERVE_TLS"))
	conf.TLS = tls

	if err != nil {
		log.Println("Unable to decide if we should server TLS. Double check your SERVE_TLS variable")
		log.Fatal(err)
	}

	conf.Certificate = os.Getenv("TLS_CERTIFICATE")
	conf.Key = os.Getenv("TLS_KEY")
}

func getRiceBox(p string) *rice.Box {
	cssBox := rice.MustFindBox(p)
	return cssBox
}

func main() {

	loadConfig()

	cssFileServer := http.StripPrefix("/css/", http.FileServer(getRiceBox("../css").HTTPBox()))
	http.Handle("/css/", cssFileServer)
	http.HandleFunc("/", createSecret)
	http.HandleFunc("/create", createSecret)
	http.HandleFunc("/secret/", showSecret)

	if conf.TLS {
		log.Println("Serving TLS")
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%s", conf.Port), conf.Certificate, conf.Key, nil))
	} else {

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", conf.Port), nil))
	}

}
