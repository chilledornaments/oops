package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	internal "github.com/mitchya1/oops/internal"

	rice "github.com/GeertJohan/go.rice"
)

func cssFiles(w http.ResponseWriter, r *http.Request) {
	//path := r.URL.Path[1:]
	data, err := ioutil.ReadFile(string("css/themes.css"))
	if err != nil {
		log.Println("Error loading CSS file")
		log.Println("Tried to load css/themes.css")
		w.Write([]byte("Error loading css file"))
	} else {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(data)
		return
	}
}

func createSecret(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		templateBox, err := rice.FindBox("../templates")

		templateString, err := templateBox.String("create.html.tmpl")

		if err != nil {
			log.Println("Unable to find secret create template")
			log.Fatal(err)
		}

		data := createTemplateData{
			CreateEndpoint: fmt.Sprintf("\"%s/%s\"", os.Getenv("SITE_URL"), "create"),
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
			expiration := n + int64(expirationInSeconds)
			if useDynamo {
				uuid, err = internal.AddDynamoSecret(s.Secret, expiration)
			} else {
				uuid, err = internal.AddSqliteSecret(s.Secret, expiration)
			}

			if err != nil {
				log.Println("Error inserting secret into DB")
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error creating secret"))
				return
			} else {
				log.Println("ID:", uuid)
				b := fmt.Sprintf("%s/%s/%s", os.Getenv("SITE_URL"), "secret", uuid)

				j := successJSON{
					URL: b,
				}
				msg, e := json.Marshal(j)
				if e != nil {
					log.Println("Error creating AJAX success JSON")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Error creating JSON"))
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(msg))
				return
			}
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

}

func showSecret(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		if strings.Contains(r.UserAgent(), "Slack") {
			log.Println("Ignored Slack link expansion")
			w.Write([]byte("Hello Slack"))
			return
		} else {

			id := strings.TrimPrefix(r.URL.Path, "/secret/")
			if useDynamo {
				secret, err = internal.ReturnDynamoSecret(id)
			} else {
				secret, err = internal.ReturnSqliteSecret(id)
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(secret + "\n"))
				return
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(secret + "\n"))
				return
			}

		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}
