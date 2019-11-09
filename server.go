package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/mitchya1/onetimepass/src/db"
)

type newSecret struct {
	Secret string `json:"secret"`
}

func secrets(w http.ResponseWriter, r *http.Request) {
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
		err = db.AddSecret(s.Secret, expiration)
		if err != nil {
			log.Println("Error inserting secret into DB")
			log.Println(err)
			w.Write([]byte("Error creating secret"))
		} else {
			http.ServeFile(w, r, "static/created.html")
		}

	} else {
		w.Write([]byte("Method not allowed"))
	}

}

func main() {
	fmt.Println("Starting the OTP web server")
	//db.Init()

	http.HandleFunc("/create", secrets)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
