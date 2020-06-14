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

func main() {

	err := godotenv.Load(os.Getenv("OOPS_ENV_FILE"))
	log.Println("Using this env file:", os.Getenv("OOPS_ENV_FILE"))
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")

	}

	cssBox := rice.MustFindBox("../css")

	if os.Getenv("LINK_EXPIRATION_TIME") == "" {
		log.Println("LINK_EXPIRATION_TIME not set in env file. Using 3600")
		expirationInSeconds = 3600
	} else {

		convertedString, err := strconv.ParseInt(os.Getenv("LINK_EXPIRATION_TIME"), 10, 64)

		if err != nil {
			panic("Unable to convert string to int")
		}
		log.Println("Expiration is", convertedString)
		expirationInSeconds = convertedString
	}

	log.Println("SITE_URL is", os.Getenv("SITE_URL"))
	log.Println("WEB_SERVER_PORT is", os.Getenv("WEB_SERVER_PORT"))

	cssFileServer := http.StripPrefix("/css/", http.FileServer(cssBox.HTTPBox()))
	http.Handle("/css/", cssFileServer)
	http.HandleFunc("/", createSecret)
	http.HandleFunc("/create", createSecret)
	http.HandleFunc("/secret/", showSecret)
	portString := fmt.Sprintf(":%s", os.Getenv("WEB_SERVER_PORT"))
	tls, err := strconv.ParseBool(os.Getenv("SERVE_TLS"))

	internal.TableName = os.Getenv("DYNAMO_TABLE_NAME")

	if err != nil {
		log.Println("Unable to decide if we should server TLS. Double check your SERVE_TLS variable")
		log.Fatal(err)
	}
	if tls {
		log.Println("Serving TLS")
		log.Fatal(http.ListenAndServeTLS(portString, os.Getenv("TLS_CERTIFICATE"), os.Getenv("TLS_KEY"), nil))
	} else {

		log.Fatal(http.ListenAndServe(portString, nil))
	}

}
