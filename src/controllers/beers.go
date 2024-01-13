package controllers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	// "strconv"
	"telegrambot/progress/models"
)

var arrowsKeys = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔼", "left"),
		tgbotapi.NewInlineKeyboardButtonData("🔽", "right"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад в меню", "backToMenu"),
	),
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
	ctx := context.Background()
	var favorite_breweries []string
	var favorite_styles []string
	favorite_breweries = strings.Split(models.RedisClient.HGetAll(ctx, fmt.Sprint(UserID)).Val()["brewery"], ",")
	favorite_styles = strings.Split(models.RedisClient.HGetAll(ctx, fmt.Sprint(UserID)).Val()["style"], ",")
	favorite_breweries = RemoveStrFromArray(favorite_breweries, "")
	favorite_styles = RemoveStrFromArray(favorite_styles, "")

	var beers []models.Beer

	if len(favorite_breweries) == 0 && len(favorite_styles) == 0 {
		beers = FindAllBeer()
	} else {
		beers = FindBeer(favorite_breweries, favorite_styles)
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
			ReplyMarkup: &arrowsKeys,
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
