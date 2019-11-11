package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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

type successJSON struct {
	URL string `json:"url"`
}

func createSecret(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("templates/create.html.tmpl"))
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
	fmt.Println("Starting the OOPS (OOPS One-time Password Sharing) web server")

	err := godotenv.Load(os.Getenv("OOPS_ENV_FILE"))
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")

	}

	log.Println("SITE_URL is", os.Getenv("SITE_URL"))
	log.Println("WEB_SERVER_PORT is", os.Getenv("WEB_SERVER_PORT"))
	http.HandleFunc("/css/", cssFiles)
	http.HandleFunc("/", createSecret)
	http.HandleFunc("/create", createSecret)
	http.HandleFunc("/secret/", showSecret)
	portString := fmt.Sprintf(":%s", os.Getenv("WEB_SERVER_PORT"))
	log.Fatal(http.ListenAndServe(portString, nil))
}
