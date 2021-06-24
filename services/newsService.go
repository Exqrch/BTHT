package service

import (
	"strings"

	"github.com/Exqrch/BTHT/model"
)

type NewsServiceInterface interface {
	FilterByTags(newsList []model.News, tagList []string) []model.News
	FilterByStatus(newsList []model.News, status string) []model.News
}

type NewsServiceImpl struct{}

func (n NewsServiceImpl) FilterByTags(newsList []model.News, tagList []string) []model.News {
	var filteredNews []model.News

	for _, singleNews := range newsList {
		if hasTag(singleNews, tagList) {
			filteredNews = append(filteredNews, singleNews)
		}
	}

	return filteredNews

}

func (n NewsServiceImpl) FilterByStatus(newsList []model.News, status string) []model.News {
	var filteredNews []model.News

	for _, singleNews := range newsList {
		if singleNews.Status == status {
			filteredNews = append(filteredNews, singleNews)
		}
	}

	return filteredNews
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
