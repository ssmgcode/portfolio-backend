package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	fmt.Println(form)
}

func main() {
	http.HandleFunc("/", sendEmailHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
