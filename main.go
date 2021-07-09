package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type Form struct {
	Name    string
	Email   string
	Subject string
	Message string
}

// Use GoDotEnv package to load/read the .env file.
func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files\n")
	}
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
	fromEmail := os.Getenv("FROM_MAIL")
	fromPassword := os.Getenv("FROM_PASSWORD")
	form, err := parseBodyRequestToFormStruct(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	message := gomail.NewMessage()
	message.SetHeaders(map[string][]string{
		"From":    {fromEmail},
		"To":      {form.Email},
		"Subject": {form.Subject},
	})
	message.SetBody("text/plain", form.Message)

	// Settings for SMTP server.
	dialer := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, fromPassword)

	// This is only needed when SSL/TLS certificate is not valid in server.
	// In production this should be set to false.
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err = dialer.DialAndSend(message); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	http.Error(rw, "Email sent successfully", http.StatusOK)
}

func main() {

	loadEnvVariables()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/send-mail", sendEmailHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
