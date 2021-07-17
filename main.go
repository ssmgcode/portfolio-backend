package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
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
	myEmail := os.Getenv("FROM_MAIL")
	myPassword := os.Getenv("FROM_PASSWORD")

	form, err := parseBodyRequestToFormStruct(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	from := mail.Address{Name: form.Name, Address: form.Email}
	to := mail.Address{Name: "SSMG Code", Address: myEmail}

	headers := map[string]string{
		"From":         from.String(),
		"To":           to.String(),
		"Subject":      form.Subject,
		"Content-Type": `text/html; charset="UTF-8"`,
	}

	var message string
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n"

	t, err := template.ParseFiles("template.html")
	sendInternalServerError(err, rw)
	buf := new(bytes.Buffer)
	err = t.Execute(buf, form)
	sendInternalServerError(err, rw)
	message += buf.String() + "\r\n"

	host := "smtp.gmail.com"
	auth := smtp.PlainAuth("", myEmail, myPassword, host)

	err = smtp.SendMail(host+":465", auth, myEmail, []string{to.Address}, []byte(message))
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
