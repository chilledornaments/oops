package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	db "github.com/mitchya1/oops/src/db"
)

type newSecret struct {
	Secret string `json:"secret"`
}

type createTemplateData struct {
	CreateEndpoint string
}

type successTemplateData struct {
	URL string
}

func createSecret(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("static/create.html.tmpl"))
		data := createTemplateData{
			CreateEndpoint: fmt.Sprintf("%s/%s", os.Getenv("SITE_URL"), "create"),
		}
		e := tmpl.Execute(w, data)
		if e != nil {
			w.Write([]byte("Could not render template"))
		}

	} else if r.Method == "POST" {

		s := newSecret{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&s)
		if err != nil {
			log.Println("Error decoding JSON request to create new secret")
			log.Println(err)
			w.Write([]byte("Error reading incoming JSON"))
		} else {
			n := time.Now().Unix()
			expiration := n + 3600
			uuid, err := db.AddSecret(s.Secret, expiration)

			if err != nil {
				log.Println("Error inserting secret into DB")
				log.Println(err)
				w.Write([]byte("Error creating secret"))
			} else {
				log.Println("ID:", uuid)
				b := fmt.Sprintf("Secret URL: %s/%s/%s\n", os.Getenv("SITE_URL"), "secret", uuid)
				tmpl := template.Must(template.ParseFiles("static/created.html.tmpl"))
				data := successTemplateData{
					URL: b,
				}
				e := tmpl.Execute(w, data)
				if e != nil {
					w.Write([]byte("Could not render template"))
				}
			}
		}

	} else {
		w.Write([]byte("Method not allowed"))
	}

}

func showSecret(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		id := strings.TrimPrefix(r.URL.Path, "/secret/")
		secret, err := db.ReturnSecret(id)

		if err != nil {
			w.Write([]byte(secret + "\n"))
		} else {
			w.Write([]byte(secret + "\n"))
		}

	} else {
		w.Write([]byte("Method not allowed"))
	}
}

func main() {
	fmt.Println("Starting the OOPS (OOPS One-time Password Sharing) web server")

	err := godotenv.Load(os.Getenv("OTP_ENV_FILE"))
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")

	}

	log.Println("SITE_URL is", os.Getenv("SITE_URL"))
	log.Println("WEB_SERVER_PORT is", os.Getenv("WEB_SERVER_PORT"))

	http.HandleFunc("/", createSecret)
	http.HandleFunc("/create", createSecret)
	http.HandleFunc("/secret/", showSecret)
	portString := fmt.Sprintf(":%s", os.Getenv("WEB_SERVER_PORT"))
	log.Fatal(http.ListenAndServe(portString, nil))
}
