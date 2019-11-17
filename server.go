package main

import (
	"encoding/json"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/joho/godotenv"
	db "github.com/mitchya1/oops/src/db"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

type successJSON struct {
	URL string `json:"url"`
}

var expirationInSeconds int64

func createSecret(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		templateBox, err := rice.FindBox("templates")

		templateString, err := templateBox.String("create.html.tmpl")

		if err != nil {
			log.Println("Unable to find secret create template")
			log.Fatal(err)
		}

		data := createTemplateData{
			CreateEndpoint: fmt.Sprintf("%s/%s", os.Getenv("SITE_URL"), "create"),
		}

		tmplMessage, err := template.New("create").Parse(templateString)

		if err != nil {
			log.Println("Unable to parse secret create template")
			log.Fatal(err)
		}

		tmplMessage.Execute(w, data)

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
			expiration := n + expirationInSeconds
			uuid, err := db.AddSecret(s.Secret, expiration)

			if err != nil {
				log.Println("Error inserting secret into DB")
				log.Println(err)
				w.Write([]byte("Error creating secret"))
			} else {
				log.Println("ID:", uuid)
				b := fmt.Sprintf("%s/%s/%s", os.Getenv("SITE_URL"), "secret", uuid)

				j := successJSON{
					URL: b,
				}
				msg, e := json.Marshal(j)
				if e != nil {
					log.Println("Error creating AJAX success JSON")
					w.Write([]byte("Error creating JSON"))
				}
				w.Write([]byte(msg))

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

// cssFiles serves, well, CSS files
func cssFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	data, err := ioutil.ReadFile(string(path))
	if err != nil {
		log.Println("Error loading CSS file")
		log.Println("Tried to load", path)
		w.Write([]byte("Error loading css file"))
	} else {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(data)

	}
}

func main() {
	log.Println("Starting the OOPS (OOPS One-time Password Sharing) web server")

	err := godotenv.Load(os.Getenv("OOPS_ENV_FILE"))
	log.Println("Using this env file:", os.Getenv("OOPS_ENV_FILE"))
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")

	}

	cssBox := rice.MustFindBox("css")

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
