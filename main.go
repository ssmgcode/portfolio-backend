package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mailgun/mailgun-go/v3"
	"log"
	"net/http"
	"os"
	"time"
)

type Form struct {
	Name    string
	Email   string
	Subject string
	Message string
}

func sendInternalServerError(err error, rw http.ResponseWriter) {
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		log.Panic(err)
	}
}

func parseBodyRequestToFormStruct(r *http.Request) (*Form, error) {
	var form *Form
	if r.Body == nil {
		return nil, errors.New("Please send name, email, subject and message.")
	}
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		return nil, errors.New("Please send name, mail, subject and message.")
	}
	return form, nil
}

func sendEmailHandler(rw http.ResponseWriter, r *http.Request) {
	mailgunDomain := os.Getenv("MAILGUN_DOMAIN")
	mailgunApiKey := os.Getenv("MAILGUN_API_KEY")

	form, err := parseBodyRequestToFormStruct(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	messageText := "New job proposal:\n\n" + form.Message + "\n\nContact sender: " + form.Email
	messageSender := form.Name + " " + "<" + form.Email + ">"
	mg := mailgun.NewMailgun(mailgunDomain, mailgunApiKey)
	message := mg.NewMessage(messageSender, form.Subject, messageText, "SSMG Code <ssmg.sg@gmail.com>")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err = mg.Send(ctx, message)
	sendInternalServerError(err, rw)

	http.Error(rw, "Email sent successfully", http.StatusOK)
	fmt.Println("Email sent successfully")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/send-mail", sendEmailHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
