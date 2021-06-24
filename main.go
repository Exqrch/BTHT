package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/*Model*/
type News struct {
	ID          string   `json:"ID"`
	Title       string   `json:"Title"`
	Description string   `json:"Description"`
	Tag         []string `json:"Tag"`
}

type allNews []News

/*Database*/
var news = allNews{
	{
		ID:          "1",
		Title:       "Silver Price Skyrocks After Reddit Attack",
		Description: "Lorum Ipsum Potum",
		Tag:         []string{"High Frequency", "Risky Trade"},
	},
	{
		ID:          "2",
		Title:       "Doge Coin Predicted To Skyrock Again Next Year",
		Description: "Lorum Ipsum Potum",
		Tag:         []string{"Long Trade", "Risky Trade"},
	},
}

/*Service*/
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newNews News
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "News ID, Title, Description, and Tag Required")
	}

	json.Unmarshal(reqBody, &newNews)

	// Add the newly created event to the array of events
	news = append(news, newNews)

	// Return the 201 created status code
	w.WriteHeader(http.StatusCreated)
	// Return the newly created event
	json.NewEncoder(w).Encode(newNews)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	newsID := mux.Vars(r)["id"]

	// Get the details from an existing event
	// Use the blank identifier to avoid creating a value that will not be used
	for _, singleNews := range news {
		if singleNews.ID == newsID {
			json.NewEncoder(w).Encode(singleNews)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(news)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	newsID := mux.Vars(r)["id"]
	var updatedNews News
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the News title and description only in order to update")
	}

	json.Unmarshal(reqBody, &updatedNews)

	for i, singleNews := range news {
		if singleNews.ID == newsID {
			singleNews.Title = updatedNews.Title
			singleNews.Description = updatedNews.Description
			news[i] = singleNews
			json.NewEncoder(w).Encode(singleNews)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	eventID := mux.Vars(r)["id"]

	// Get the details from an existing event
	// Use the blank identifier to avoid creating a value that will not be used
	for i, singleNews := range news {
		if singleNews.ID == eventID {
			news = append(news[:i], news[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

/*Controller*/
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/news", createEvent).Methods("POST")
	router.HandleFunc("/news/all", getAllEvents).Methods("GET")
	router.HandleFunc("/news/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/news/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/news/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
