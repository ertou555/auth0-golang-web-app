package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"01-Login/platform/authenticator"
	"01-Login/platform/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := router.New(auth)

	log.Print("Server listening on http://localhost:3000/")
	//生成证书
	//openssl genrsa -out key.pem 2048
	//openssl req -new -x509 -key key.pem -out cert.pem -days 3650
	if err := http.ListenAndServeTLS("0.0.0.0:3000", "cert.pem", "key.pem", rtr); err != nil {
	//if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
