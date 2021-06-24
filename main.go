package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Exqrch/BTHT/model"
	"github.com/Exqrch/BTHT/repositories"
	"github.com/gorilla/mux"
)

type allNews []model.News

/*Repository -- No SQL, using hard coded data*/
var newsRepository repositories.NewsRepositoryInterface = new(repositories.NewsRepositoryImpl)

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
	newsRepository.Create(newNews)

	// Return the 201 created status code
	w.WriteHeader(http.StatusCreated)
	// Return the newly created news
	json.NewEncoder(w).Encode(newNews)
}

func getOneNews(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	newsID := mux.Vars(r)["id"]
	news := newsRepository.GetById(newsID)
	json.NewEncoder(w).Encode(news)
}

func getAllNews(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(newsRepository.GetAllNews())
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

	var afterUpdate model.News = newsRepository.Update(newsID, updatedNews.Title, updatedNews.Description, updatedNews.Tag, updatedNews.Status)

	json.NewEncoder(w).Encode(afterUpdate)
}

func deleteNews(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	var returnNews model.News
	newsID := mux.Vars(r)["id"]
	returnNews = newsRepository.Delete(newsID)
	json.NewEncoder(w).Encode(returnNews)
	fmt.Fprintf(w, "The news with ID %v has been marked for deletion", newsID)
}

func filterNewsByStatus(w http.ResponseWriter, r *http.Request) {
	status := mux.Vars(r)["s"]

	var filteredNews = allNews{}
	for _, singleNews := range newsRepository.GetAllNews() {
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

	for _, singleNews := range newsRepository.GetAllNews() {
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
