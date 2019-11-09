package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	db "github.com/mitchya1/onetimepass/src/db"
)

type newSecret struct {
	Secret string `json:"secret"`
}

func createSecret(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == "GET" {
		http.ServeFile(w, r, "static/create.html")
	} else if r.Method == "POST" {
		s := newSecret{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&s)
		if err != nil {
			log.Println("Error decoding JSON request to create new secret")
			w.Write([]byte("Error reading incoming JSON"))
		}
		n := time.Now().Unix()
		expiration := n + 3600
		id, err := db.AddSecret(s.Secret, expiration)

		if err != nil {
			log.Println("Error inserting secret into DB")
			log.Println(err)
			w.Write([]byte("Error creating secret"))
		} else {
			log.Println("ID:", id)
			b := fmt.Sprintf("Secret URL: %s/%d\n", "http://localhost:8081/secret", id)
			w.Write([]byte(b))
			//http.ServeFile(w, r, "static/created.html")
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
	fmt.Println("Starting the OTP web server")
	//db.Init()

	http.HandleFunc("/create", createSecret)
	http.HandleFunc("/secret/", showSecret)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
