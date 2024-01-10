package controllers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"telegrambot/progress/models"
)

var arrowsKeys = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔼", "left"),
		tgbotapi.NewInlineKeyboardButtonData("🔽", "right"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад в меню", "back"),
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

func GetBeerList(bot *tgbotapi.BotAPI, UserID int64) []models.Beer {
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

func DisplayBeer(bot *tgbotapi.BotAPI, UserID int64, beer *models.Beer) {
	beer_description := fmt.Sprintf("%s от %s \nСтиль: %s\nABV: %.2f Rate: %.2f\n%s\n%d₽",
		beer.Name, beer.Brewery,
		beer.Style, beer.ABV,
		beer.Rate, beer.Brief, beer.Price)
	photo := tgbotapi.NewPhoto(UserID, tgbotapi.FilePath(beer.ImagePath))
	photo.Caption = beer_description
	photo.ReplyMarkup = arrowsKeys

	var editMsg tgbotapi.EditMessageMediaConfig
	editMsg.Media = tgbotapi.FilePath(beer.ImagePath)
	editMsg.MessageID = MsgID
	if _, err := bot.Send(photo); err != nil {
		panic(err)
	}
}
