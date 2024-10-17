package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"os"
)

type ErrorMessage struct {
	Error string `json:"error"`
}

type Message struct {
	Subject      string   `json:"subject"`
	Server       string   `json:"server"`
	Time         string   `json:"time"`
	To           []string `json:"to"`
	ErrorMessage string   `json:"errorMessage"`
}

func SendEmail(w http.ResponseWriter, r *http.Request) {
	var errorMessage ErrorMessage

	from := os.Getenv("SMTP_USERNAME")

	if from == "" {
		errorMessage.Error = "SMTP_USERNAME was not configured."
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	password := os.Getenv("SMTP_PASSWORD")

	if password == "" {
		errorMessage.Error = "SMTP_PASSWORD was not configured."
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	smtpHost := "smtp.mailersend.net"
	smtpPort := "587"

	message := Message{}

	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		errorMessage.Error = "Error on decodify JSON"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("./email/template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", message.Subject, mimeHeaders)))

	t.Execute(&body, struct {
		Server  string
		Error   string
		Horario string
	}{
		Server:  message.Server,
		Error:   message.ErrorMessage,
		Horario: message.Time,
	})

	errSmtp := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, message.To, body.Bytes())

	if errSmtp != nil {
		fmt.Println(errSmtp)
		errorMessage.Error = fmt.Sprintf("Error while sending email to user: %v", errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	fmt.Println("Email enviado com sucesso!")
}
