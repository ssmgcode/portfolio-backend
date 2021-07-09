package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gopkg.in/gomail.v2"
)

type Form struct {
	Name    string
	Email   string
	Subject string
	Message string
}

func parseBodyRequestToFormStruct(r *http.Request) (*Form, error) {
	var form *Form
	if r.Body == nil {
		return nil, errors.New("Please send name, email, subject and message.")
	}
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		return nil, errors.New("Please send name, email, subject and message.")
	}
	return form, nil
}

func sendEmailHandler(rw http.ResponseWriter, r *http.Request) {
	form, err := parseBodyRequestToFormStruct(r)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
	message := gomail.NewMessage()
	message.SetHeaders(map[string][]string{
		"From":    {""},
		"To":      {form.Email},
		"Subject": {form.Subject},
	})
	message.SetBody("text/plain", form.Message)

	// Settings for SMTP server.
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "", "")

	// This is only needed when SSL/TLS certificate is not valid in server.
	// In production this should be set to false.
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err = dialer.DialAndSend(message); err != nil {
		http.Error(rw, err.Error(), 500)
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", sendEmailHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
