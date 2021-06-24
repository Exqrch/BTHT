package repositories

import (
	"github.com/Exqrch/BTHT/model"
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

type NewsRepositoryInterface interface {
	Create(newNews model.News) model.News
	Delete(newsID string) model.News
	GetById(newsID string) model.News
	GetAllNews() []model.News
	Update(newsID string, uTitle string, uDescription string, uTag []string, uStatus string) model.News
}

type NewsRepositoryImpl struct{}

func (n NewsRepositoryImpl) Create(newNews model.News) model.News {
	news = append(news, newNews)
	return newNews
}

func (n NewsRepositoryImpl) Delete(newsID string) model.News {
	emptyNews := model.News{}

	for i, singleNews := range news {
		if singleNews.ID == newsID {
			singleNews.Status = "Deleted"
			news[i] = singleNews
			return singleNews
		}
	}

	return emptyNews
}

func (n NewsRepositoryImpl) GetById(newsID string) model.News {
	emptyNews := model.News{}

	for _, singleNews := range news {
		if singleNews.ID == newsID {
			return singleNews
		}
	}

	return emptyNews
}

func (n NewsRepositoryImpl) GetAllNews() []model.News {
	return news
}

func (n NewsRepositoryImpl) Update(newsID string, uTitle string, uDescription string, uTag []string, uStatus string) model.News {
	var uNews model.News
	var index int

	for i, singleNews := range news {
		if singleNews.ID == newsID {
			uNews = singleNews
			index = i
			break
		}
	}

	if uTitle != "" {
		uNews.Title = uTitle
	}
	if uDescription != "" {
		uNews.Description = uDescription
	}
	if uTag != nil {
		uNews.Tag = uTag
	}
	if uStatus != "" {
		uNews.Status = uStatus
	}
	news[index] = uNews
	return uNews

}
