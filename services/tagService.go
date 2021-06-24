package service

import (
	"github.com/Exqrch/BTHT/model"
)

type TagServiceInterface interface {
	GetOKTag(tagList []model.TopicTag) []model.TopicTag
}

type TagServiceImpl struct{}

func (t TagServiceImpl) GetOKTag(tagList []model.TopicTag) []model.TopicTag {
	var filteredTag []model.TopicTag
	for _, tag := range tagList {
		if tag.Status == "OK" {
			filteredTag = append(filteredTag, tag)
		}
	}
	return filteredTag
}
