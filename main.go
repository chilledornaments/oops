package main

import (
	"fmt"
	"log"
	"net/http"
)

func createSecret(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == "GET" {
		http.ServeFile(w, r, "static/create.html")
	} else if r.Method == "POST" {
		http.ServeFile(w, r, "static/created.html")
	}

}

func main() {
	fmt.Println("Starting the OTP web server")

	http.HandleFunc("/create", createSecret)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
