package repositories

import (
	"github.com/Exqrch/BTHT/model"
)

var tags = []model.TopicTag{
	{
		TID:          "1",
		TopicTagName: "Risky Trade",
		Status:       "OK",
	},
}

type TagRepositoryInterface interface {
	Create(newTag model.TopicTag) model.TopicTag
	Delete(tagID string) model.TopicTag
	GetById(tagID string) model.TopicTag
	GetOKTag() []model.TopicTag
	GetAllTag() []model.TopicTag
	Update(tagID string, uTopicTagName string, uStatus string) model.TopicTag
}

type TagRepositoryImpl struct{}

func (t TagRepositoryImpl) Create(newTag model.TopicTag) model.TopicTag {
	tags = append(tags, newTag)
	return newTag
}

func (t TagRepositoryImpl) Delete(tagID string) model.TopicTag {
	var returnTag model.TopicTag
	for i, tag := range tags {
		if tag.TID == tagID {
			tag.Status = "Deleted"
			tags[i] = tag
			returnTag = tag
			break
		}
	}
	return returnTag
}

func (t TagRepositoryImpl) GetById(tagID string) model.TopicTag {
	var returnTag model.TopicTag
	for _, tag := range tags {
		if tag.TID == tagID {
			returnTag = tag
			break
		}
	}
	return returnTag
}

func (t TagRepositoryImpl) GetOKTag() []model.TopicTag {
	var filteredTag []model.TopicTag
	for _, tag := range tags {
		if tag.Status == "OK" {
			filteredTag = append(filteredTag, tag)
		}
	}
	return filteredTag
}

func (t TagRepositoryImpl) GetAllTag() []model.TopicTag {
	return tags
}

func (t TagRepositoryImpl) Update(tagID string, uTopicTagName string, uStatus string) model.TopicTag {
	var uTag model.TopicTag
	var index int

	for i, tag := range tags {
		if tag.TID == tagID {
			index = i
			uTag = tag
			break
		}
	}

	if uTopicTagName != "" {
		uTag.TopicTagName = uTopicTagName
	}
	if uStatus != "" {
		uTag.Status = uStatus
	}

	tags[index] = uTag
	return uTag
}
