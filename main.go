package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/beeker1121/mailchimp-go/lists/members"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	mailchimp "github.com/beeker1121/mailchimp-go"
)

type ContactForm struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Company   string `json:"company"`
	Message   string `json:"message"`
}
type Response struct {
	Text   string `json:"text"`
	Status int    `json:"status"`
}

func main() {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	r := mux.NewRouter()
	r.HandleFunc("/contact", handleContact).Methods("POST")
	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func handleContact(w http.ResponseWriter, r *http.Request) {

	contactForm := ContactForm{}
	err := mailchimp.SetKey("5bcab47e4600071c2820deb7a46b4fd7-us17")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Encountered error while reading request body")
		return
	}

	err = json.Unmarshal(body, &contactForm)
	if err != nil {
		fmt.Println("Encountered error while unmarshalling JSON")
		return
	}

	mergeFields := map[string]interface{}{
		"FNAME":   contactForm.FirstName,
		"LNAME":   contactForm.LastName,
		"COMPANY": contactForm.Company,
		"MESSAGE": contactForm.Message,
	}

	params := &members.NewParams{}
	params.EmailAddress = contactForm.Email
	params.Status = members.StatusSubscribed
	params.MergeFields = mergeFields

	_, err = members.New("4935d5e586", params)
	if err != nil {
		fmt.Println("Encountered error adding member to mailchimp")
	}

	respondJson("true", http.StatusOK, w)
	return

}

func respondJson(text string, status int, w http.ResponseWriter) {

	response := Response{Text: text, Status: status}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResponse)

}
