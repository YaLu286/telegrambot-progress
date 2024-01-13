package controllers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegrambot/progress/keyboards"
	"telegrambot/progress/models"
)

func FindAllBeer() []models.Beer {
	var beer_list []models.Beer
	models.DB.Find(&beer_list)
	return beer_list
}

func FindBeer(favorite_breweries []string, favorite_styles []string) []models.Beer {
	var beer_list []models.Beer
	if len(favorite_breweries) == 0 {
		models.DB.Where("style IN ?", favorite_styles).Find(&beer_list)
	} else if len(favorite_styles) == 0 {
		models.DB.Where("brewery IN ?", favorite_breweries).Find(&beer_list)
	} else {
		models.DB.Where("brewery IN ? AND style IN ?", favorite_breweries, favorite_styles).Find(&beer_list)
	}
	return beer_list
}

func GetBeerList(UserID int64) []models.Beer {

	session := &models.UserSession{UserID: UserID}
	session.LoadInfo()

	var beers []models.Beer

	if len(session.Breweries) == 0 && len(session.Styles) == 0 {
		beers = FindAllBeer()
	} else {
		beers = FindBeer(session.Breweries, session.Styles)
	}
	return beers
}

func DisplayBeer(bot *tgbotapi.BotAPI, UserID int64, beer *models.Beer, callerID int) {
	beer_description := fmt.Sprintf("%s от %s \nСтиль: %s\nABV: %.2f Rate: %.2f\n%s\n%d₽",
		beer.Name, beer.Brewery,
		beer.Style, beer.ABV,
		beer.Rate, beer.Brief, beer.Price)

	var beerImage tgbotapi.InputMediaPhoto
	beerImage.Media = tgbotapi.FilePath(beer.ImagePath)
	beerImage.Caption = beer_description
	beerImage.Type = "photo"
	editMsg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      UserID,
			MessageID:   callerID,
			ReplyMarkup: &keyboards.ArrowsKeys,
		},
		Media: beerImage,
	}
	if _, err := bot.Request(editMsg); err != nil {
		panic(err)
	}
}

// func NextPage(UserID int64, CallerMsgID int, BeerMap map[int64][]models.Beer) (bool, int){
// 	ctx := context.Background()
// 	next_page, _ := strconv.Atoi(models.RedisClient.HGet(ctx, fmt.Sprint(UserID), "page").Val())
// 	next_page++
// 	if next_page < len(BeerMap[UserID]) {
// 		DisplayBeer(bot, UserID, &BeerMap[UserID][next_page], CallerMsgID)
// 		models.RedisClient.HSet(ctx, fmt.Sprint(UserID), "page", next_page)
// 	}
// 	return false, next_page
// }
