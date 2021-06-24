package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Exqrch/BTHT/model"
	"github.com/gorilla/mux"
)

type allNews []model.News

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
	var newNews model.News
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "News ID, Title, Description, and Tag Required")
	}

	json.Unmarshal(reqBody, &newNews)
	newNews.Status = "Publish"
	// Add the newly created news to the array of news
	news = append(news, newNews)

	// Return the 201 created status code
	w.WriteHeader(http.StatusCreated)
	// Return the newly created news
	json.NewEncoder(w).Encode(newNews)
}

func getOneNews(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	newsID := mux.Vars(r)["id"]

	// Get the details from an existing news
	// Use the blank identifier to avoid creating a value that will not be used
	for _, singleNews := range news {
		if singleNews.ID == newsID {
			if singleNews.Status == "Publish" {
				json.NewEncoder(w).Encode(singleNews)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		}
	}
}

func getAllNews(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(news)
}

func updateNews(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	newsID := mux.Vars(r)["id"]
	var updatedNews model.News
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
	newsID := mux.Vars(r)["id"]

	// Get the details from an existing news
	for i, singleNews := range news {
		if singleNews.ID == newsID {
			news[i].Status = "Deleted"
			fmt.Fprintf(w, "The news with ID %v has been marked for deletion", newsID)
		}
	}
}

func filterNewsByStatus(w http.ResponseWriter, r *http.Request) {
	status := mux.Vars(r)["s"]

	var filteredNews = allNews{}
	for _, singleNews := range news {
		if singleNews.Status == status {
			filteredNews = append(filteredNews, singleNews)
		}
	}
	json.NewEncoder(w).Encode(filteredNews)
}

func filterNewsByTag(w http.ResponseWriter, r *http.Request) {
	// Convert r.Body into a readable formart
	tagQuery := r.URL.Query().Get("filter")

	var filteredNews = allNews{}

	tagFilter := strings.Split(tagQuery, ",")

	for _, singleNews := range news {
		if hasTag(singleNews, tagFilter) {
			filteredNews = append(filteredNews, singleNews)
		}
	}

	if len(filteredNews) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	// Return the newly created news
	json.NewEncoder(w).Encode(filteredNews)
}

/*Helper Function*/
func hasTag(singleNews model.News, tagFilter []string) bool {
	for _, tag := range tagFilter {
		if !foundIn(strings.TrimSpace(tag), singleNews.Tag) {
			return false
		}
	}
	return true
}

func foundIn(s1 string, sArray []string) bool {
	for _, s := range sArray {
		if s == s1 {
			return true
		}
	}
	return false
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
	router.HandleFunc("/news", filterNewsByTag).Queries("filter", "{filter}").Methods("GET")
	router.HandleFunc("/news/status/{s}", filterNewsByStatus).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
