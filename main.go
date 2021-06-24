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
	service "github.com/Exqrch/BTHT/services"
	"github.com/gorilla/mux"
)

type allNews []model.News

/*Repository -- No SQL, using hard coded data*/
var newsRepository repositories.NewsRepositoryInterface = new(repositories.NewsRepositoryImpl)

/*Service*/
var newsService service.NewsServiceInterface = new(service.NewsServiceImpl)

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
	var filteredNews = newsService.FilterByStatus(newsRepository.GetAllNews(), status)

	json.NewEncoder(w).Encode(filteredNews)
}

func filterNewsByTag(w http.ResponseWriter, r *http.Request) {
	// Convert r.Body into a readable formart
	tagQuery := r.URL.Query().Get("filter")
	tagFilter := strings.Split(tagQuery, ",")

	var filteredNews = newsService.FilterByTags(newsRepository.GetAllNews(), tagFilter)

	if len(filteredNews) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	// Return the newly created news
	json.NewEncoder(w).Encode(filteredNews)
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
