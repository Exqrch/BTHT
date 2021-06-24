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

/*Repository -- No SQL, using hard coded data*/
var newsRepository repositories.NewsRepositoryInterface = new(repositories.NewsRepositoryImpl)
var tagRepository repositories.TagRepositoryInterface = new(repositories.TagRepositoryImpl)

/*Service*/
var newsService service.NewsServiceInterface = new(service.NewsServiceImpl)
var tagService service.TagServiceInterface = new(service.TagServiceImpl)

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

func createTag(w http.ResponseWriter, r *http.Request) {
	var newTag model.TopicTag
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "News ID, Title, Description, and Tag Required")
	}
	json.Unmarshal(reqBody, &newTag)

	tagRepository.Create(newTag)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTag)
}

func getAllTag(w http.ResponseWriter, r *http.Request) {
	var tagList []model.TopicTag = tagRepository.GetAllTag()
	json.NewEncoder(w).Encode(tagList)
}

func getOneTag(w http.ResponseWriter, r *http.Request) {
	tagID := mux.Vars(r)["id"]
	tag := tagRepository.GetById(tagID)
	json.NewEncoder(w).Encode(tag)
}

func updateTag(w http.ResponseWriter, r *http.Request) {
	tagID := mux.Vars(r)["id"]
	var updatedTag model.TopicTag
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the Tag name or status in order to update")
	}

	json.Unmarshal(reqBody, &updatedTag)

	var afterUpdate model.TopicTag = tagRepository.Update(tagID, updatedTag.TopicTagName, updatedTag.Status)

	json.NewEncoder(w).Encode(afterUpdate)
}

func deleteTag(w http.ResponseWriter, r *http.Request) {
	var returnTag model.TopicTag
	tagID := mux.Vars(r)["id"]
	returnTag = tagRepository.Delete(tagID)
	json.NewEncoder(w).Encode(returnTag)
	fmt.Fprintf(w, "The tag with ID %v has been marked for deletion", tagID)
}

func getOKTag(w http.ResponseWriter, r *http.Request) {
	var filteredTag []model.TopicTag = tagService.GetOKTag(tagRepository.GetAllTag())
	json.NewEncoder(w).Encode(filteredTag)
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

	router.HandleFunc("/tag", createTag).Methods("POST")
	router.HandleFunc("/tag/all", getAllTag).Methods("GET")
	router.HandleFunc("/tag/ok", getOKTag).Methods("GET")
	router.HandleFunc("/tag/{id}", getOneTag).Methods("GET")
	router.HandleFunc("/tag/{id}", updateTag).Methods("PATCH")
	router.HandleFunc("/tag/{id}", deleteTag).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
