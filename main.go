package main

import (
	"alertmanager/email"
	"alertmanager/slack"
	"alertmanager/sms"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	f, err := os.OpenFile("alertmanager.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	errEnv := godotenv.Load()

	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Starting alert service... ... ...")

	router := mux.NewRouter()

	router.HandleFunc("/email", email.SendEmail).Methods("POST")
	router.HandleFunc("/slack", slack.SendSlack).Methods("POST")
	router.HandleFunc("/sms", sms.SendSMS).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
