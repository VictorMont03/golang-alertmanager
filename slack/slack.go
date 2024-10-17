package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

type ErrorMessage struct {
	Error string `json:"error"`
}

type SuccessMessage struct {
	Message string `json:"message"`
}

type Message struct {
	ChannelID string `json:"channelId"`
	Text      string `json:"text"`
}

func SendSlack(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("SLACK_AUTH_TOKEN")

	if token == "" {
		log.Fatal("SLACK_AUTH_TOKEN wasnt configured")
	}

	var errorMessage ErrorMessage
	message := Message{}

	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		errorMessage.Error = "Error on decodify JSON"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)
		return
	}

	client := slack.New(token, slack.OptionDebug(true))

	attachment := slack.Attachment{
		Color:   "danger",
		Pretext: "Alerta de desenvolvimento!",
		Text:    message.Text,
	}

	_, timestamp, err := client.PostMessage(
		message.ChannelID,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Message successfully sent to channel %s at %s", message.ChannelID, timestamp)

	var successMessage SuccessMessage

	successMessage.Message = "Slack message sent with success!"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successMessage)
}
