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
	Status      string   `json:"Status"`
}

type Error struct {
	Error string `json:"Error"`
}

type allNews []News

/*Database*/
var news = allNews{
	{
		ID:          "1",
		Title:       "Silver Price Skyrocks After Reddit Attack",
		Description: "Lorum Ipsum Potum",
		Tag:         []string{"High Frequency", "Risky Trade"},
		Status:      "Publish",
	},
	{
		ID:          "2",
		Title:       "Doge Coin Predicted To Skyrock Again Next Year",
		Description: "Lorum Ipsum Potum",
		Tag:         []string{"Long Trade", "Risky Trade"},
		Status:      "Draft",
	},
}

/*Service*/
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page!")
}

func createNews(w http.ResponseWriter, r *http.Request) {
	var newNews News
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "News ID, Title, Description, and Tag Required")
	}

	json.Unmarshal(reqBody, &newNews)
	newNews.Status = "Publish"
	// Add the newly created event to the array of events
	news = append(news, newNews)

	// Return the 201 created status code
	w.WriteHeader(http.StatusCreated)
	// Return the newly created event
	json.NewEncoder(w).Encode(newNews)
}

func getOneNews(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	newsID := mux.Vars(r)["id"]

	// Get the details from an existing event
	// Use the blank identifier to avoid creating a value that will not be used
	for _, singleNews := range news {
		if singleNews.ID == newsID {
			if singleNews.Status == "Publish" {
				json.NewEncoder(w).Encode(singleNews)
			} else {
				w.WriteHeader(http.StatusNoContent)
				var customError Error
				customError.Error = "The News you are searching is either not published or doesn't exist"
				json.NewEncoder(w).Encode(customError)
			}
		}
	}
}

func getAllNews(w http.ResponseWriter, r *http.Request) {
	var publishedNews = allNews{}

	for _, singleNews := range news {
		if singleNews.Status == "Publish" {
			publishedNews = append(publishedNews, singleNews)
		}
	}
	json.NewEncoder(w).Encode(publishedNews)
}

func updateNews(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	newsID := mux.Vars(r)["id"]
	var updatedNews News
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the News title or description or tags or status in order to update")
	}

	json.Unmarshal(reqBody, &updatedNews)

	for i, singleNews := range news {
		if singleNews.ID == newsID {
			if updatedNews.Title != "" {
				singleNews.Title = updatedNews.Title
			}
			if updatedNews.Description != "" {
				singleNews.Description = updatedNews.Description
			}
			if updatedNews.Tag != nil {
				singleNews.Tag = updatedNews.Tag
			}
			if updatedNews.Status != "" {
				singleNews.Status = updatedNews.Status
			}
			news[i] = singleNews
			json.NewEncoder(w).Encode(singleNews)
		}
	}
}

func deleteNews(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	eventID := mux.Vars(r)["id"]

	// Get the details from an existing event
	for i, singleNews := range news {
		if singleNews.ID == eventID {
			news[i].Status = "Deleted"
			fmt.Fprintf(w, "The event with ID %v has been marked for deletion", eventID)
		}
	}
}

/*Controller*/
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/news", createNews).Methods("POST")
	router.HandleFunc("/news/all", getAllNews).Methods("GET")
	router.HandleFunc("/news/{id}", getOneNews).Methods("GET")
	router.HandleFunc("/news/{id}", updateNews).Methods("PATCH")
	router.HandleFunc("/news/{id}", deleteNews).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
