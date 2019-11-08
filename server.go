package main

import (
	"fmt"
	db "github.com/mitchya1/onetimepass/src/db"
	"log"
	"net/http"
)

func createSecret(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == "GET" {

		http.ServeFile(w, r, "static/create.html")
	} else if r.Method == "POST" {
		_ = db.AddSecret("test")
		http.ServeFile(w, r, "static/created.html")

	} else {
		w.Write([]byte("Method not allowed"))
	}

}

func main() {
	fmt.Println("Starting the OTP web server")

	http.HandleFunc("/create", createSecret)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
