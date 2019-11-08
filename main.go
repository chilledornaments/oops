package main

import (
	"fmt"
	"log"
	"net/http"
)

func createSecret(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/create.html")
}

func main() {
	fmt.Println("Starting the OTP web server")

	http.HandleFunc("/create", createSecret)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
