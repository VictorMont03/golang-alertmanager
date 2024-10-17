package sms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type ErrorMessage struct {
	Error string `json:"error"`
}

type SuccessMessage struct {
	Message string `json:"message"`
}

type Message struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func SendSMS(w http.ResponseWriter, r *http.Request) {
	endpoint := "https://rest.nexmo.com/sms/json"

	data := url.Values{}

	api_key := os.Getenv("NEXMO_API_KEY")

	if api_key == "" {
		panic("NEXMO_API_KEY was not configured!")
	}

	api_secret := os.Getenv("NEXMO_API_SECRET")

	if api_secret == "" {
		panic("NEXMO_API_SECRET was not configured!")
	}

	message := Message{}
	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		errorMessage := ErrorMessage{Error: "Error on decodify JSON"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	data.Set("api_key", api_key)
	data.Set("api_secret", api_secret)
	data.Set("to", message.Phone)
	data.Set("from", "AlertMgr")
	data.Set("text", message.Message)
	data.Set("type", "unicode")

	client := &http.Client{}

	fmt.Println(strings.NewReader(data.Encode()))

	r, smsErr := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))

	if smsErr != nil {
		errorMessage := ErrorMessage{Error: "Error on create request to send SMS"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)

	if err != nil {
		errorMessage := ErrorMessage{Error: "Error on send SMS"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)

	if err != nil {
		errorMessage := ErrorMessage{Error: "Error on send SMS"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	successMessage := SuccessMessage{}
	successMessage.Message = "SMS sent with success!"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successMessage)
}
